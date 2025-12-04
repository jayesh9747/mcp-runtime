package cli

import (
	"testing"
)

// TestSplitImage tests the splitImage function (now SplitImage for testing)
// This is in the same package, so it can access unexported functions
func TestSplitImage(t *testing.T) {
	tests := []struct {
		image string
		want  string
		tag   string
	}{
		{"registry.example.com/example-mcp-server:latest", "registry.example.com/example-mcp-server", "latest"},
		{"registry.example.com/example-mcp-server", "registry.example.com/example-mcp-server", ""},
		{"example-mcp-server:latest", "example-mcp-server", "latest"},
		{"example-mcp-server", "example-mcp-server", ""},
	}
	for _, test := range tests {
		image, tag := splitImage(test.image)
		if image != test.want {
			t.Errorf("SplitImage(%q) = %q, want %q", test.image, image, test.want)
		}
		if tag != test.tag {
			t.Errorf("SplitImage(%q) tag = %q, want %q", test.image, tag, test.tag)
		}
	}
}

func TestDropRegistryPrefix(t *testing.T) {
	tests := []struct {
		repo string
		want string
	}{
		{"registry.example.com/example-mcp-server", "example-mcp-server"},
		{"example-mcp-server", "example-mcp-server"},
		{"localhost:5000/my-image", "my-image"},
		{"192.168.1.1:5000/my-image", "my-image"},
		{"my-image", "my-image"},
	}
	for _, test := range tests {
		repo := dropRegistryPrefix(test.repo)
		if repo != test.want {
			t.Errorf("dropRegistryPrefix(%q) = %q, want %q", test.repo, repo, test.want)
		}
	}
}
