#!/bin/bash

# install.sh - Installation script for genie CLI tool
# Usage: ./install.sh

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Emojis for better visual appeal
ROCKET="🚀"
SPARKLES="✨"
HAMMER="🔨"
CHECK="✅"
CROSS="❌"
INFO="💡"
PACKAGE="📦"
GEAR="⚙️"

# App info
APP_NAME="genie"
INSTALL_PATH="/usr/local/bin"
BINARY_NAME="genie"

print_header() {
    echo
    echo -e "${PURPLE}${BOLD}================================${NC}"
    echo -e "${WHITE}${BOLD}    ${SPARKLES} Genie CLI Installer ${SPARKLES}    ${NC}"
    echo -e "${PURPLE}${BOLD}================================${NC}"
    echo -e "${CYAN}AI-powered Git commit message generator${NC}"
    echo
}

print_step() {
    echo -e "${BLUE}${BOLD}$1${NC} $2"
}

print_success() {
    echo -e "${GREEN}${CHECK}${NC} $1"
}

print_error() {
    echo -e "${RED}${CROSS}${NC} $1" >&2
}

print_warning() {
    echo -e "${YELLOW}${INFO}${NC} $1"
}

print_info() {
    echo -e "${CYAN}${INFO}${NC} $1"
}

check_requirements() {
    print_step "${GEAR}" "Checking requirements..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go first:"
        echo -e "  ${CYAN}https://golang.org/doc/install${NC}"
        exit 1
    fi
    
    local go_version=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go ${go_version} found"
    
    # Check if git is installed
    if ! command -v git &> /dev/null; then
        print_error "Git is not installed. Please install Git first."
        exit 1
    fi
    
    local git_version=$(git --version | awk '{print $3}')
    print_success "Git ${git_version} found"
    
    # Check if main.go exists
    if [[ ! -f "main.go" ]]; then
        print_error "main.go not found in current directory"
        echo -e "  ${YELLOW}Please run this script from the genie project directory${NC}"
        exit 1
    fi
    
    print_success "All requirements met"
    echo
}

build_binary() {
    print_step "${HAMMER}" "Building ${APP_NAME}..."
    
    # Clean any existing binary
    if [[ -f "${BINARY_NAME}" ]]; then
        rm "${BINARY_NAME}"
    fi
    
    # Build the binary
    if go build -o "${BINARY_NAME}" main.go; then
        print_success "Binary built successfully"
    else
        print_error "Failed to build binary"
        exit 1
    fi
    
    # Check if binary was created and is executable
    if [[ -f "${BINARY_NAME}" && -x "${BINARY_NAME}" ]]; then
        local binary_size=$(ls -lh "${BINARY_NAME}" | awk '{print $5}')
        print_success "Binary size: ${binary_size}"
    else
        print_error "Binary was not created properly"
        exit 1
    fi
    echo
}

install_binary() {
    print_step "${PACKAGE}" "Installing to ${INSTALL_PATH}..."
    
    # Check if install directory exists and is writable
    if [[ ! -d "${INSTALL_PATH}" ]]; then
        print_warning "${INSTALL_PATH} does not exist"
        if sudo mkdir -p "${INSTALL_PATH}"; then
            print_success "Created ${INSTALL_PATH}"
        else
            print_error "Failed to create ${INSTALL_PATH}"
            exit 1
        fi
    fi
    
    # Install the binary
    if sudo cp "${BINARY_NAME}" "${INSTALL_PATH}/${BINARY_NAME}"; then
        print_success "Installed to ${INSTALL_PATH}/${BINARY_NAME}"
    else
        print_error "Failed to install binary"
        exit 1
    fi
    
    # Make sure it's executable
    if sudo chmod +x "${INSTALL_PATH}/${BINARY_NAME}"; then
        print_success "Made binary executable"
    else
        print_error "Failed to make binary executable"
        exit 1
    fi
    
    # Clean up local binary
    rm "${BINARY_NAME}"
    print_success "Cleaned up build artifacts"
    echo
}

verify_installation() {
    print_step "${CHECK}" "Verifying installation..."
    
    # Check if binary is in PATH and executable
    if command -v "${BINARY_NAME}" &> /dev/null; then
        local installed_version=$(${BINARY_NAME} --version 2>/dev/null || echo "unknown version")
        print_success "${BINARY_NAME} is available in PATH"
        print_success "Version: ${installed_version}"
    else
        print_error "${BINARY_NAME} is not available in PATH"
        print_warning "You may need to add ${INSTALL_PATH} to your PATH"
        echo -e "  ${CYAN}Add this to your ~/.bashrc or ~/.zshrc:${NC}"
        echo -e "  ${WHITE}export PATH=\"${INSTALL_PATH}:\$PATH\"${NC}"
        exit 1
    fi
    echo
}

print_setup_instructions() {
    echo -e "${GREEN}${BOLD}================================${NC}"
    echo -e "${WHITE}${BOLD}     ${SPARKLES} Installation Complete! ${SPARKLES}     ${NC}"
    echo -e "${GREEN}${BOLD}================================${NC}"
    echo
    echo -e "${WHITE}${BOLD}Next steps:${NC}"
    echo -e "${YELLOW}1.${NC} Get your Gemini API key:"
    echo -e "   ${CYAN}https://aistudio.google.com/apikey${NC}"
    echo
    echo -e "${YELLOW}2.${NC} Set up your environment variable:"
    echo -e "   ${WHITE}export GOOGLE_AI_TOKEN=your_api_key_here${NC}"
    echo
    echo -e "${YELLOW}3.${NC} Add it to your shell profile for persistence:"
    echo -e "   ${WHITE}echo 'export GOOGLE_AI_TOKEN=your_api_key' >> ~/.bashrc${NC}"
    echo -e "   ${WHITE}# or for zsh: echo 'export GOOGLE_AI_TOKEN=your_api_key' >> ~/.zshrc${NC}"
    echo
    echo -e "${WHITE}${BOLD}Usage:${NC}"
    echo -e "   ${GREEN}genie${NC}                    ${CYAN}# Generate commit message for changes${NC}"
    echo -e "   ${GREEN}genie \"bug fix\"${NC}          ${CYAN}# Generate with context${NC}"
    echo -e "   ${GREEN}genie --help${NC}             ${CYAN}# Show help${NC}"
    echo -e "   ${GREEN}genie --version${NC}          ${CYAN}# Show version${NC}"
    echo
    echo -e "${GREEN}${ROCKET} Happy committing! ${ROCKET}${NC}"
    echo
}

main() {
    print_header
    
    # Check if running as root (not recommended)
    if [[ $EUID -eq 0 ]]; then
        print_warning "Running as root. This is not recommended."
        echo -e "  ${CYAN}Consider running as a regular user (sudo will be used when needed)${NC}"
        echo
    fi
    
    # Run installation steps
    check_requirements
    build_binary
    install_binary
    verify_installation
    print_setup_instructions
}

# Handle Ctrl+C gracefully
trap 'echo -e "\n${RED}${CROSS} Installation cancelled${NC}"; exit 1' INT

# Run main function
main "$@"