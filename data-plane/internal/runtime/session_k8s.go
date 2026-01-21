package runtime

import (
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubernetesSessionRuntime struct {
	Client       kubernetes.Interface
	Config       *rest.Config
	Namespace    string
	RuntimeClass string
	Image        string
	PythonImage  string
	NodeImage    string
	Env          string
	AgentAddr    string
	AgentAuthMode  string
	AgentAuthToken string
}

func (r KubernetesSessionRuntime) StartSession(ctx context.Context, sessionID string, policyID string, workspaceRef string, runtime string) (SessionRoute, error) {
	_ = policyID
	_ = workspaceRef
	if r.Client == nil {
		return SessionRoute{}, errors.New("missing kubernetes client")
	}
	if sessionID == "" {
		return SessionRoute{}, errors.New("missing session id")
	}
	if r.Namespace == "" {
		r.Namespace = "default"
	}
	image := imageForRuntime(runtime, r.PythonImage, r.NodeImage, r.Image)
	podName := fmt.Sprintf("session-%s", sessionID)
	envVars := []corev1.EnvVar{}
	if r.Env != "" {
		envVars = append(envVars, corev1.EnvVar{Name: "ENV", Value: r.Env})
	}
	if r.AgentAddr != "" {
		envVars = append(envVars, corev1.EnvVar{Name: "SESSION_AGENT_ADDR", Value: r.AgentAddr})
	}
	if r.AgentAuthMode == "bypass" {
		envVars = append(envVars, corev1.EnvVar{Name: "SESSION_AGENT_AUTH_BYPASS", Value: "true"})
	} else {
		envVars = append(envVars, corev1.EnvVar{Name: "SESSION_AGENT_AUTH_BYPASS", Value: "false"})
		if r.AgentAuthToken != "" {
			envVars = append(envVars, corev1.EnvVar{Name: "SESSION_AGENT_AUTH_TOKEN", Value: r.AgentAuthToken})
		}
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: r.Namespace,
			Labels: map[string]string{
				"app":        "sandbox-session",
				"session_id": sessionID,
			},
		},
		Spec: corev1.PodSpec{
			RuntimeClassName: runtimeClassName(r.RuntimeClass),
			Containers: []corev1.Container{
				{
					Name:    "session",
					Image:   image,
					Env:     envVars,
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	if _, err := r.Client.CoreV1().Pods(r.Namespace).Create(ctx, pod, metav1.CreateOptions{}); err != nil {
		return SessionRoute{}, err
	}
	endpoint := r.buildAgentEndpoint(ctx, podName)
	return SessionRoute{
		RuntimeID: podName,
		Endpoint:  endpoint,
		Token:     r.AgentAuthToken,
		AuthMode:  r.AgentAuthMode,
	}, nil
}

func (r KubernetesSessionRuntime) RunStep(ctx context.Context, runtimeID string, command string) (StepOutput, error) {
	_ = ctx
	_ = runtimeID
	_ = command
	return StepOutput{}, errors.New("kubernetes exec is disabled; use session-agent endpoint")
}

func (r KubernetesSessionRuntime) TerminateSession(ctx context.Context, runtimeID string) error {
	if r.Client == nil {
		return errors.New("missing kubernetes client")
	}
	if runtimeID == "" {
		return errors.New("missing runtime id")
	}
	namespace := r.Namespace
	if namespace == "" {
		namespace = "default"
	}
	grace := int64(5)
	return r.Client.CoreV1().Pods(namespace).Delete(ctx, runtimeID, metav1.DeleteOptions{
		GracePeriodSeconds: &grace,
	})
}

func runtimeClassName(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func imageForRuntime(runtime string, pythonImage string, nodeImage string, fallback string) string {
	switch runtime {
	case "python":
		if pythonImage != "" {
			return pythonImage
		}
	case "node":
		if nodeImage != "" {
			return nodeImage
		}
	}
	if fallback != "" {
		return fallback
	}
	switch runtime {
	case "python":
		return "python:3.12-slim"
	case "node":
		return "node:20-alpine"
	default:
		return "busybox:1.36"
	}
}

func (r KubernetesSessionRuntime) buildAgentEndpoint(ctx context.Context, podName string) string {
	namespace := r.Namespace
	if namespace == "" {
		namespace = "default"
	}
	if r.Client != nil {
		pod, err := r.Client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
		if err == nil && pod.Status.PodIP != "" {
			return "http://" + pod.Status.PodIP + ":9000"
		}
	}
	return fmt.Sprintf("http://%s.%s.pod:9000", podName, namespace)
}
