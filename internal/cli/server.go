package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewServerCmd(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Manage MCP servers",
		Long: `Commands for managing MCP server deployments.

For building images from source, use 'server build'.
For pushing images, use 'registry push'.`,
	}

	cmd.AddCommand(newServerListCmd(logger))
	cmd.AddCommand(newServerGetCmd(logger))
	cmd.AddCommand(newServerCreateCmd(logger))
	cmd.AddCommand(newServerDeleteCmd(logger))
	cmd.AddCommand(newServerLogsCmd(logger))
	cmd.AddCommand(newServerStatusCmd(logger))
	cmd.AddCommand(newServerBuildCmd(logger))

	return cmd
}

func newServerBuildCmd(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build MCP server images (push via `registry push`)",
	}

	// Only expose image build here; pushing is handled by `registry push`.
	cmd.AddCommand(newBuildImageCmd(logger))

	return cmd
}

func newServerListCmd(logger *zap.Logger) *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List MCP servers",
		Long:  "List all MCP server deployments",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listServers(logger, namespace)
		},
	}

	cmd.Flags().StringVar(&namespace, "namespace", "mcp-servers", "Namespace to list servers from")

	return cmd
}

func newServerGetCmd(logger *zap.Logger) *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "get [name]",
		Short: "Get MCP server details",
		Long:  "Get detailed information about an MCP server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getServer(logger, args[0], namespace)
		},
	}

	cmd.Flags().StringVar(&namespace, "namespace", "mcp-servers", "Namespace")

	return cmd
}

func newServerCreateCmd(logger *zap.Logger) *cobra.Command {
	var namespace string
	var image string
	var imageTag string
	var file string

	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create an MCP server",
		Long:  "Create a new MCP server deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file != "" {
				return createServerFromFile(logger, file)
			}
			return createServer(logger, args[0], namespace, image, imageTag)
		},
	}

	cmd.Flags().StringVar(&namespace, "namespace", "mcp-servers", "Namespace")
	cmd.Flags().StringVar(&image, "image", "", "Container image")
	cmd.Flags().StringVar(&imageTag, "tag", "latest", "Image tag")
	cmd.Flags().StringVar(&file, "file", "", "YAML file with server spec")

	return cmd
}

func newServerDeleteCmd(logger *zap.Logger) *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete an MCP server",
		Long:  "Delete an MCP server deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteServer(logger, args[0], namespace)
		},
	}

	cmd.Flags().StringVar(&namespace, "namespace", "mcp-servers", "Namespace")

	return cmd
}

func newServerLogsCmd(logger *zap.Logger) *cobra.Command {
	var namespace string
	var follow bool

	cmd := &cobra.Command{
		Use:   "logs [name]",
		Short: "View server logs",
		Long:  "View logs from an MCP server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return viewServerLogs(logger, args[0], namespace, follow)
		},
	}

	cmd.Flags().StringVar(&namespace, "namespace", "mcp-servers", "Namespace")
	cmd.Flags().BoolVar(&follow, "follow", false, "Follow log output")

	return cmd
}

func newServerStatusCmd(logger *zap.Logger) *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show MCP server runtime status (pods, images, pull secrets)",
		Long:  "List MCPServer resources with their Deployment/pod status, image, and pull secrets.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return serverStatus(logger, namespace)
		},
	}

	cmd.Flags().StringVar(&namespace, "namespace", "mcp-servers", "Namespace to inspect")

	return cmd
}

func listServers(logger *zap.Logger, namespace string) error {
	cmd := exec.Command("kubectl", "get", "mcpserver", "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func getServer(logger *zap.Logger, name, namespace string) error {
	cmd := exec.Command("kubectl", "get", "mcpserver", name, "-n", namespace, "-o", "yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func createServer(logger *zap.Logger, name, namespace, image, imageTag string) error {
	if image == "" {
		return fmt.Errorf("image is required")
	}

	logger.Info("Creating MCP server", zap.String("name", name), zap.String("image", image))

	// Create a basic MCPServer manifest
	manifest := fmt.Sprintf(`apiVersion: mcp.agent-hellboy.io/v1alpha1
kind: MCPServer
metadata:
  name: %s
  namespace: %s
spec:
  image: %s
  imageTag: %s
  replicas: 1
  port: 8088
  servicePort: 80
  ingressPath: /%s
`, name, namespace, image, imageTag, name)

	tmpFile := fmt.Sprintf("/tmp/mcpserver-%s.yaml", name)
	if err := os.WriteFile(tmpFile, []byte(manifest), 0644); err != nil {
		return fmt.Errorf("failed to create manifest: %w", err)
	}
	defer os.Remove(tmpFile)

	cmd := exec.Command("kubectl", "apply", "-f", tmpFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func createServerFromFile(logger *zap.Logger, file string) error {
	cmd := exec.Command("kubectl", "apply", "-f", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func deleteServer(logger *zap.Logger, name, namespace string) error {
	logger.Info("Deleting MCP server", zap.String("name", name))

	cmd := exec.Command("kubectl", "delete", "mcpserver", name, "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func viewServerLogs(logger *zap.Logger, name, namespace string, follow bool) error {
	args := []string{"logs", "-l", "app=" + name, "-n", namespace}
	if follow {
		args = append(args, "-f")
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func serverStatus(logger *zap.Logger, namespace string) error {
	fmt.Printf("Namespace: %s\n", namespace)
	fmt.Println("MCPServers:")
	getServers := exec.Command("kubectl", "get", "mcpserver", "-n", namespace, "-o", "wide")
	getServers.Stdout = os.Stdout
	getServers.Stderr = os.Stderr
	if err := getServers.Run(); err != nil {
		return err
	}

	fmt.Println("\nDeployments:")
	getDeploy := exec.Command("kubectl", "get", "deploy", "-n", namespace, "-l", "app", "-o", "wide")
	getDeploy.Stdout = os.Stdout
	getDeploy.Stderr = os.Stderr
	if err := getDeploy.Run(); err != nil {
		logger.Warn("Could not list deployments", zap.Error(err))
	}

	fmt.Println("\nPods:")
	getPods := exec.Command("kubectl", "get", "pods", "-n", namespace, "-o", "wide", "--show-labels")
	getPods.Stdout = os.Stdout
	getPods.Stderr = os.Stderr
	if err := getPods.Run(); err != nil {
		logger.Warn("Could not list pods", zap.Error(err))
	}

	// Show image and pull secrets for each MCPServer
	fmt.Println("\nDetails per MCPServer:")
	getServersYaml := exec.Command("kubectl", "get", "mcpserver", "-n", namespace, "-o", "jsonpath={range .items[*]}{.metadata.name} {.spec.image}{\":\"}{.spec.imageTag} {.spec.useProvisionedRegistry} {.spec.registryOverride}{\"\\n\"}{end}")
	out, err := getServersYaml.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 1 {
				name := fields[0]
				image := ""
				if len(fields) >= 2 {
					image = fields[1]
				}
				useProv := ""
				if len(fields) >= 3 {
					useProv = fields[2]
				}
				regOverride := ""
				if len(fields) >= 4 {
					regOverride = fields[3]
				}
				fmt.Printf("  %s\n", name)
				fmt.Printf("    image: %s\n", image)
				if useProv != "" {
					fmt.Printf("    useProvisionedRegistry: %s\n", useProv)
				}
				if regOverride != "" {
					fmt.Printf("    registryOverride: %s\n", regOverride)
				}

				// pull secrets from deployment
				psCmd := exec.Command("kubectl", "get", "deploy", name, "-n", namespace, "-o", "jsonpath={.spec.template.spec.imagePullSecrets[*].name}")
				psOut, psErr := psCmd.Output()
				if psErr == nil && len(psOut) > 0 {
					fmt.Printf("    imagePullSecrets: %s\n", string(psOut))
				}
			}
		}
	}

	return nil
}
