package operator_test

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	mcpv1alpha1 "mcp-runtime/api/v1alpha1"
	"mcp-runtime/internal/operator"
)

func TestMCPServerReconciler_ReconcileCreatesResources(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = mcpv1alpha1.AddToScheme(scheme)
	_ = appsv1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
	_ = networkingv1.AddToScheme(scheme)

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithStatusSubresource(&mcpv1alpha1.MCPServer{}).
		Build()

	reconciler := &operator.MCPServerReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	mcpServer := &mcpv1alpha1.MCPServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-server",
			Namespace: "mcp-servers",
		},
		Spec: mcpv1alpha1.MCPServerSpec{
			Image:        "my-registry.com/my-image:v1.0",
			ImageTag:     "v1.0",
			Port:         9000,
			ServicePort:  8080,
			Replicas:     int32Ptr(3),
			IngressPath:  "/custom/path",
			IngressHost:  "example.com",
			IngressClass: "traefik",
			IngressAnnotations: map[string]string{
				"custom": "annotation",
			},
		},
	}

	ctx := context.Background()
	if err := fakeClient.Create(ctx, mcpServer); err != nil {
		t.Fatalf("failed to create MCPServer: %v", err)
	}

	_, err := reconciler.Reconcile(ctx, ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-server",
			Namespace: "mcp-servers",
		},
	})
	if err != nil {
		t.Fatalf("Reconcile returned error: %v", err)
	}

	deployment := &appsv1.Deployment{}
	if err := fakeClient.Get(ctx, types.NamespacedName{Name: "test-server", Namespace: "mcp-servers"}, deployment); err != nil {
		t.Fatalf("expected Deployment to be created: %v", err)
	}
	if *deployment.Spec.Replicas != 3 {
		t.Errorf("expected 3 replicas, got %d", *deployment.Spec.Replicas)
	}
	if deployment.Spec.Template.Spec.Containers[0].Image != "my-registry.com/my-image:v1.0" {
		t.Errorf("expected image my-registry.com/my-image:v1.0, got %s", deployment.Spec.Template.Spec.Containers[0].Image)
	}
	if deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort != 9000 {
		t.Errorf("expected container port 9000, got %d", deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort)
	}

	service := &corev1.Service{}
	if err := fakeClient.Get(ctx, types.NamespacedName{Name: "test-server", Namespace: "mcp-servers"}, service); err != nil {
		t.Fatalf("expected Service to be created: %v", err)
	}
	if service.Spec.Ports[0].Port != 8080 {
		t.Errorf("expected service port 8080, got %d", service.Spec.Ports[0].Port)
	}
	if service.Spec.Ports[0].TargetPort.IntVal != 9000 {
		t.Errorf("expected service target port 9000, got %d", service.Spec.Ports[0].TargetPort.IntVal)
	}

	ingress := &networkingv1.Ingress{}
	if err := fakeClient.Get(ctx, types.NamespacedName{Name: "test-server", Namespace: "mcp-servers"}, ingress); err != nil {
		t.Fatalf("expected Ingress to be created: %v", err)
	}
	if ingress.Spec.Rules[0].HTTP.Paths[0].Path != "/custom/path" {
		t.Errorf("expected ingress path /custom/path, got %s", ingress.Spec.Rules[0].HTTP.Paths[0].Path)
	}
	if ingress.Spec.Rules[0].Host != "example.com" {
		t.Errorf("expected ingress host example.com, got %s", ingress.Spec.Rules[0].Host)
	}
	if ingress.Spec.IngressClassName == nil || *ingress.Spec.IngressClassName != "traefik" {
		t.Errorf("expected ingress class traefik, got %v", ingress.Spec.IngressClassName)
	}
	if ingress.Annotations["custom"] != "annotation" {
		t.Errorf("expected custom annotation to be present, got %v", ingress.Annotations["custom"])
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}
