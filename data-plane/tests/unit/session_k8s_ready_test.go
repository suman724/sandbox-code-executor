package unit

import (
	"context"
	"testing"
	"time"

	"data-plane/internal/runtime"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestWaitForPodReady(t *testing.T) {
	client := fake.NewSimpleClientset()
	r := runtime.KubernetesSessionRuntime{
		Client:    client,
		Namespace: "default",
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "session-test",
			Namespace: "default",
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodPending,
		},
	}
	if _, err := client.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{}); err != nil {
		t.Fatalf("create pod: %v", err)
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		pod.Status.Phase = corev1.PodRunning
		pod.Status.Conditions = []corev1.PodCondition{
			{
				Type:   corev1.PodReady,
				Status: corev1.ConditionTrue,
			},
		}
		_, _ = client.CoreV1().Pods("default").UpdateStatus(context.Background(), pod, metav1.UpdateOptions{})
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := r.WaitForPodReady(ctx, "session-test"); err != nil {
		t.Fatalf("wait for ready: %v", err)
	}
}
