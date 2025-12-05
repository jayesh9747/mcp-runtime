package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+kubebuilder:object:generate=true

// MCPServerSpec defines the desired state of MCPServer
type MCPServerSpec struct {
	// Image is the container image for the MCP server
	Image string `json:"image"`

	// ImageTag is the tag of the container image (defaults to "latest")
	ImageTag string `json:"imageTag,omitempty"`

	// RegistryOverride, if set, overrides the registry portion of the image (e.g., registry.example.com)
	RegistryOverride string `json:"registryOverride,omitempty"`

	// UseProvisionedRegistry tells the controller to use the provisioned registry (from operator env) for this server
	UseProvisionedRegistry bool `json:"useProvisionedRegistry,omitempty"`

	// ImagePullSecrets are secrets to use for pulling the image
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`

	// Replicas is the number of desired replicas (defaults to 1)
	Replicas *int32 `json:"replicas,omitempty"`

	// Port is the port the container listens on (defaults to 8088)
	Port int32 `json:"port,omitempty"`

	// ServicePort is the port exposed by the service (defaults to 80)
	ServicePort int32 `json:"servicePort,omitempty"`

	// IngressPath is the path for the ingress route (defaults to /{name}/mcp)
	IngressPath string `json:"ingressPath,omitempty"`

	// IngressHost is the hostname for the ingress (optional; defaults from MCP_DEFAULT_INGRESS_HOST env var if set on the operator)
	IngressHost string `json:"ingressHost,omitempty"`

	// IngressClass is the ingress class to use (e.g., "traefik", "nginx", "istio"). Defaults to "traefik"
	IngressClass string `json:"ingressClass,omitempty"`

	// IngressAnnotations are additional annotations for the ingress controller
	IngressAnnotations map[string]string `json:"ingressAnnotations,omitempty"`

	// Resources defines resource limits and requests
	Resources ResourceRequirements `json:"resources,omitempty"`

	// EnvVars are environment variables to pass to the container
	EnvVars []EnvVar `json:"envVars,omitempty"`
}

//+kubebuilder:object:generate=true

// ResourceRequirements defines resource limits and requests
type ResourceRequirements struct {
	Limits   *ResourceList `json:"limits,omitempty"`
	Requests *ResourceList `json:"requests,omitempty"`
}

//+kubebuilder:object:generate=true

// ResourceList defines CPU and memory resources
type ResourceList struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

//+kubebuilder:object:generate=true

// EnvVar represents an environment variable
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

//+kubebuilder:object:generate=true

// MCPServerStatus defines the observed state of MCPServer
type MCPServerStatus struct {
	// Phase represents the current phase of the MCPServer
	Phase string `json:"phase,omitempty"`

	// Message provides additional information about the status
	Message string `json:"message,omitempty"`

	// Conditions represent the latest available observations
	Conditions []Condition `json:"conditions,omitempty"`

	// DeploymentReady indicates if the deployment is ready
	DeploymentReady bool `json:"deploymentReady,omitempty"`

	// ServiceReady indicates if the service is ready
	ServiceReady bool `json:"serviceReady,omitempty"`

	// IngressReady indicates if the ingress is ready
	IngressReady bool `json:"ingressReady,omitempty"`
}

//+kubebuilder:object:generate=true

// Condition represents a condition status
type Condition struct {
	Type               string      `json:"type"`
	Status             string      `json:"status"`
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
//+kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image"
//+kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.deploymentReady"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// MCPServer is the Schema for the mcpservers API
type MCPServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MCPServerSpec   `json:"spec,omitempty"`
	Status MCPServerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MCPServerList contains a list of MCPServer
type MCPServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MCPServer `json:"items"`
}
