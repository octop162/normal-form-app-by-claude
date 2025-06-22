#!/bin/bash

# Install Go static analysis tools for normal-form-app
# Usage: ./scripts/install-go-tools.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

install_tool() {
    local tool_name=$1
    local tool_package=$2
    local version=$3
    
    print_status "Installing $tool_name..."
    
    if command -v $tool_name &> /dev/null; then
        print_warning "$tool_name is already installed"
        return 0
    fi
    
    if [ -n "$version" ]; then
        go install ${tool_package}@${version}
    else
        go install ${tool_package}@latest
    fi
    
    if command -v $tool_name &> /dev/null; then
        print_success "$tool_name installed successfully"
    else
        print_error "Failed to install $tool_name"
        return 1
    fi
}

main() {
    print_status "Installing Go static analysis tools..."
    echo ""
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go first."
        exit 1
    fi
    
    print_status "Go version: $(go version)"
    echo ""
    
    # Install golangci-lint (comprehensive linter)
    print_status "Installing golangci-lint..."
    if ! command -v golangci-lint &> /dev/null; then
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
        print_success "golangci-lint installed"
    else
        print_warning "golangci-lint is already installed"
    fi
    
    # Install individual tools
    install_tool "goimports" "golang.org/x/tools/cmd/goimports"
    install_tool "staticcheck" "honnef.co/go/tools/cmd/staticcheck"
    install_tool "gosec" "github.com/securecodewarrior/gosec/v2/cmd/gosec"
    install_tool "ineffassign" "github.com/gordonklaus/ineffassign"
    install_tool "misspell" "github.com/client9/misspell/cmd/misspell"
    install_tool "gocyclo" "github.com/fzipp/gocyclo/cmd/gocyclo"
    install_tool "goconst" "github.com/jgautheron/goconst/cmd/goconst"
    install_tool "deadcode" "golang.org/x/tools/cmd/deadcode"
    install_tool "godot" "github.com/tetafro/godot/cmd/godot"
    install_tool "gofumpt" "mvdan.cc/gofumpt"
    
    echo ""
    print_success "All Go static analysis tools installed successfully!"
    echo ""
    
    # Verify installations
    print_status "Verifying installations..."
    echo ""
    
    tools=(
        "golangci-lint:golangci-lint version"
        "goimports:goimports --version"
        "staticcheck:staticcheck -version"
        "gosec:gosec -version"
        "ineffassign:ineffassign --version"
        "misspell:misspell -version"
        "gocyclo:gocyclo --version"
        "goconst:goconst --version"
        "deadcode:deadcode --version"
        "godot:godot --version"
        "gofumpt:gofumpt --version"
    )
    
    for tool_info in "${tools[@]}"; do
        IFS=':' read -r tool_name version_cmd <<< "$tool_info"
        if command -v $tool_name &> /dev/null; then
            version_output=$($version_cmd 2>&1 | head -1 || echo "Version command failed")
            print_success "$tool_name: $version_output"
        else
            print_error "$tool_name: Not found in PATH"
        fi
    done
    
    echo ""
    print_status "Installation complete!"
    print_status "You can now run 'make lint' or 'golangci-lint run' to analyze your code."
    echo ""
    
    # Add PATH information
    GOPATH=$(go env GOPATH)
    if [[ ":$PATH:" != *":$GOPATH/bin:"* ]]; then
        print_warning "Note: $GOPATH/bin is not in your PATH."
        print_status "Add the following to your shell profile (.bashrc, .zshrc, etc.):"
        echo "export PATH=\$PATH:\$(go env GOPATH)/bin"
    fi
}

main "$@"