package runtime

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type KubernetesSessionRuntime struct {
	Client       kubernetes.Interface
	Config       *rest.Config
	Namespace    string
	RuntimeClass string
	Image        string
}

func (r KubernetesSessionRuntime) StartSession(ctx context.Context, sessionID string, policyID string, workspaceRef string) (string, error) {
	_ = policyID
	_ = workspaceRef
	if r.Client == nil {
		return "", errors.New("missing kubernetes client")
	}
	if sessionID == "" {
		return "", errors.New("missing session id")
	}
	if r.Namespace == "" {
		r.Namespace = "default"
	}
	image := r.Image
	if image == "" {
		image = "busybox:1.36"
	}
	podName := fmt.Sprintf("session-%s", sessionID)
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
					Command: []string{"sh", "-c", "sleep 3600"},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	if _, err := r.Client.CoreV1().Pods(r.Namespace).Create(ctx, pod, metav1.CreateOptions{}); err != nil {
		return "", err
	}
	return podName, nil
}

func (r KubernetesSessionRuntime) RunStep(ctx context.Context, runtimeID string, command string) (StepOutput, error) {
	if r.Client == nil {
		return StepOutput{}, errors.New("missing kubernetes client")
	}
	if r.Config == nil {
		return StepOutput{}, errors.New("missing kubernetes config")
	}
	if runtimeID == "" {
		return StepOutput{}, errors.New("missing runtime id")
	}
	if command == "" {
		return StepOutput{}, errors.New("missing command")
	}
	namespace := r.Namespace
	if namespace == "" {
		namespace = "default"
	}
	req := r.Client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(runtimeID).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: "session",
			Command:   []string{"sh", "-c", command},
			Stdout:    true,
			Stderr:    true,
		}, scheme.ParameterCodec)
	executor, err := remotecommand.NewSPDYExecutor(r.Config, "POST", req.URL())
	if err != nil {
		return StepOutput{}, err
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	}); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return StepOutput{}, fmt.Errorf("kubernetes exec error: %w: %s", err, msg)
		}
		return StepOutput{}, fmt.Errorf("kubernetes exec error: %w", err)
	}
	return StepOutput{Stdout: stdout.String(), Stderr: stderr.String()}, nil
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
