#!/bin/bash
# dev-setup.sh - Setup development environment for mcp-runtime
# This script installs dev dependencies (controller-gen, kustomize) and generates code/manifests.
# End users don't need this - they use pre-generated manifests.
#
# Usage:
#   ./hack/dev-setup.sh [command]
#
# Commands:
#   install    - Install dev tools (controller-gen, kustomize)
#   generate   - Generate CRD manifests and DeepCopy methods
#   format     - Format code with go fmt
#   validate   - Validate code with go vet
#   all        - Run all steps (default)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${PROJECT_ROOT}"

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        echo "❌ Error: Go is not installed. Please install Go to continue."
        echo "   Visit: https://golang.org/dl/"
        exit 1
    fi
}

# Install dev tools
install_tools() {
    echo "📦 Installing dev tools..."
    echo ""
    
    # Install controller-gen
    echo "📦 Installing controller-gen..."
    if [ ! -f "bin/controller-gen" ]; then
        echo "   Downloading controller-gen..."
        make -f Makefile.operator controller-gen
        echo "   ✓ controller-gen installed to bin/controller-gen"
    else
        echo "   ✓ controller-gen already installed"
    fi
    echo ""
    
    # Install kustomize
    echo "📦 Installing kustomize..."
    if [ ! -f "bin/kustomize" ]; then
        echo "   Downloading kustomize..."
        make -f Makefile.operator kustomize
        echo "   ✓ kustomize installed to bin/kustomize"
    else
        echo "   ✓ kustomize already installed"
    fi
    echo ""
    
    echo "✅ Tools installation complete!"
}

# Generate manifests and code
generate_code() {
    echo "📝 Generating code and manifests..."
    echo ""
    
    # Generate CRD manifests
    echo "📝 Generating CRD manifests..."
    make -f Makefile.operator manifests
    echo "   ✓ CRD manifests generated in config/crd/bases/"
    echo ""
    
    # Generate DeepCopy methods
    echo "📝 Generating DeepCopy methods..."
    make -f Makefile.operator generate
    echo "   ✓ DeepCopy methods generated in api/v1alpha1/zz_generated.deepcopy.go"
    echo ""
    
    echo "✅ Code generation complete!"
}

# Format code
format_code() {
    echo "🎨 Formatting code..."
    make -f Makefile.operator fmt
    echo "   ✓ Code formatted"
    echo ""
    echo "✅ Formatting complete!"
}

# Validate code
validate_code() {
    echo "🔍 Validating code..."
    make -f Makefile.operator vet
    echo "   ✓ Code validated"
    echo ""
    echo "✅ Validation complete!"
}

# Show usage
show_usage() {
    cat << EOF
🔧 MCP Runtime Development Setup
==================================

Usage: ./hack/dev-setup.sh [command]

Commands:
  install    Install dev tools (controller-gen, kustomize)
  generate   Generate CRD manifests and DeepCopy methods
  format     Format code with go fmt
  validate   Validate code with go vet
  all        Run all steps (default)

Examples:
  ./hack/dev-setup.sh install     # Install tools only
  ./hack/dev-setup.sh generate    # Generate manifests only
  ./hack/dev-setup.sh format      # Format code only
  ./hack/dev-setup.sh validate    # Validate code only
  ./hack/dev-setup.sh all         # Run everything
  ./hack/dev-setup.sh             # Run everything (default)

Note: End users don't need these tools - they use pre-generated manifests.
EOF
}

# Main execution
COMMAND="${1:-all}"

case "$COMMAND" in
    install)
        check_go
        echo "✓ Go is installed: $(go version)"
        echo ""
        install_tools
        ;;
    generate)
        check_go
        echo "✓ Go is installed: $(go version)"
        echo ""
        generate_code
        ;;
    format)
        check_go
        format_code
        ;;
    validate)
        check_go
        validate_code
        ;;
    all)
        check_go
        echo "🔧 MCP Runtime Development Setup"
        echo "=================================="
        echo ""
        echo "✓ Go is installed: $(go version)"
        echo ""
        
        install_tools
        generate_code
        format_code
        validate_code
        
        echo "✅ Development setup complete!"
        echo ""
        echo "Next steps:"
        echo "  - Make changes to api/v1alpha1/mcpserver_types.go"
        echo "  - Run './hack/dev-setup.sh generate' to regenerate manifests"
        echo "  - Or use: make -f Makefile.operator manifests generate"
        echo ""
        echo "Tools installed:"
        echo "  - bin/controller-gen (for generating CRDs and RBAC)"
        echo "  - bin/kustomize (for building Kubernetes manifests)"
        ;;
    help|--help|-h)
        show_usage
        ;;
    *)
        echo "❌ Unknown command: $COMMAND"
        echo ""
        show_usage
        exit 1
        ;;
esac

