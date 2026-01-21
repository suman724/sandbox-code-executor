package runtime

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"shared/sessionagent"
)

type KubernetesSessionRuntime struct {
	Client        kubernetes.Interface
	Config        *rest.Config
	Namespace     string
	RuntimeClass  string
	Image         string
	PythonImage   string
	NodeImage     string
	Env           string
	AgentAddr     string
	AgentAuthMode string
}

func (r KubernetesSessionRuntime) StartSession(ctx context.Context, sessionID string, policyID string, workspaceRef string, runtime string) (SessionRoute, error) {
	_ = policyID
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
	workspaceRoot := "/workspace"
	workspaceDir := filepath.Join(workspaceRoot, sessionID)
	if workspaceRef != "" {
		workspaceDir = filepath.Join(workspaceRoot, workspaceRef)
	}
	envVars := []corev1.EnvVar{}
	if r.Env != "" {
		envVars = append(envVars, corev1.EnvVar{Name: "ENV", Value: r.Env})
	}
	if r.AgentAddr != "" {
		envVars = append(envVars, corev1.EnvVar{Name: "SESSION_AGENT_ADDR", Value: r.AgentAddr})
	}
	envVars = append(envVars, corev1.EnvVar{Name: "WORKSPACE_ROOT", Value: workspaceRoot})
	if r.AgentAuthMode == "bypass" {
		envVars = append(envVars, corev1.EnvVar{Name: "SESSION_AGENT_AUTH_BYPASS", Value: "true"})
	} else {
		envVars = append(envVars, corev1.EnvVar{Name: "SESSION_AGENT_AUTH_BYPASS", Value: "false"})
	}
	volumeName := "workspace"
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
			Volumes: []corev1.Volume{
				{
					Name: volumeName,
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name:  "session",
					Image: image,
					Env:   envVars,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      volumeName,
							MountPath: workspaceRoot,
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	if _, err := r.Client.CoreV1().Pods(r.Namespace).Create(ctx, pod, metav1.CreateOptions{}); err != nil {
		return SessionRoute{}, err
	}
	cleanupPod := func() {
		namespace := r.Namespace
		if namespace == "" {
			namespace = "default"
		}
		grace := int64(0)
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = r.Client.CoreV1().Pods(namespace).Delete(cleanupCtx, podName, metav1.DeleteOptions{
			GracePeriodSeconds: &grace,
		})
	}
	timeout := durationFromEnv("SESSION_READY_TIMEOUT", 60*time.Second)
	readyCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := r.WaitForPodReady(readyCtx, podName); err != nil {
		cleanupPod()
		return SessionRoute{}, err
	}
	endpoint := r.buildAgentEndpoint(ctx, podName)
	client := NewAgentClient()
	if err := client.WaitForHealth(readyCtx, endpoint, 500*time.Millisecond); err != nil {
		cleanupPod()
		return SessionRoute{}, err
	}
	token := ""
	if r.AgentAuthMode != "bypass" {
		token = generateSessionToken()
	}
	if err := client.RegisterSession(ctx, AgentRoute{
		Endpoint: endpoint,
		Token:    token,
		AuthMode: r.AgentAuthMode,
	}, sessionagent.SessionRegisterRequest{
		SessionID:    sessionID,
		Runtime:      runtime,
		Token:        token,
		WorkspaceDir: workspaceDir,
	}); err != nil {
		cleanupPod()
		return SessionRoute{}, err
	}
	return SessionRoute{
		RuntimeID: podName,
		Runtime:   runtime,
		Endpoint:  endpoint,
		Token:     token,
		AuthMode:  r.AgentAuthMode,
	}, nil
}

func (r KubernetesSessionRuntime) RunStep(ctx context.Context, runtimeID string, command string) (StepOutput, error) {
	_ = ctx
	_ = runtimeID
	_ = command
	return StepOutput{}, ErrRuntimeUnavailable
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

func (r KubernetesSessionRuntime) WaitForPodReady(ctx context.Context, podName string) error {
	namespace := r.Namespace
	if namespace == "" {
		namespace = "default"
	}
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		pod, err := r.Client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
		if err == nil && podReady(pod) {
			return nil
		}
		select {
		case <-ctx.Done():
			if err != nil {
				return err
			}
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func podReady(pod *corev1.Pod) bool {
	if pod == nil {
		return false
	}
	if pod.Status.Phase != corev1.PodRunning {
		return false
	}
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}
