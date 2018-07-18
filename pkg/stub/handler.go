package stub

import (
	"context"

	"github.com/alaypatel07/tocker-operator/pkg/apis/tocker/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.TockerApp:
		err := sdk.Create(newTockerAppPod(o))
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create tocker pod : %v", err)
			return err
		}
		err = sdk.Create(newTockerAppService(o))
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create tocker service: %v", err)
			return err
		}
	}
	return nil
}

// newTockerAppPod demonstrates how to create a busybox pod
func newTockerAppPod(cr *v1alpha1.TockerApp) *corev1.Pod {
	labels := map[string]string{
		"app": "tocker",
	}
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tocker",
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "TockerApp",
				}),
			},
			Labels: labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "tocker",
					Image: "alaypatel07/tocker",
				},
			},
		},
	}
}

func newTockerAppService(cr *v1alpha1.TockerApp) *corev1.Service {
	labels := map[string]string{
		"app": "tocker",
	}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tocker-service",
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "TockerApp",
				}),
			},
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "tocker",
			},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.Protocol("TCP"),
					TargetPort: intstr.IntOrString{
						StrVal: "8080"},
					Port: 8080,
				},
			},
		},
	}
}
