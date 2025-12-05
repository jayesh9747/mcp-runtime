package operator

import (
	"testing"
)

func TestRewriteRegistry(t *testing.T) {
	tests := []struct {
		name     string
		image    string
		registry string
		want     string
	}{
		{
			name:     "test-image",
			image:    "test-image",
			registry: "registry.registry.svc.cluster.local:5000",
			want:     "registry.registry.svc.cluster.local:5000/test-image",
		},
	}
	for _, test := range tests {
		got := rewriteRegistry(test.image, test.registry)
		if got != test.want {
			t.Errorf("rewriteRegistry(%q, %q) = %q, want %q", test.image, test.registry, got, test.want)
		}
	}
}
