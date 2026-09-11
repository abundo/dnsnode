# Build directory
BUILD_DIR := build

# Install directory
INSTALL_DIR := /usr/bin

# Go build flags (for release)
GO_BUILD_FLAGS := -ldflags="-s -w"

build: dnsnode_cli

dnsnode_cli:
	@mkdir -p $(BUILD_DIR)
	@go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/dnsnode_cli ./cmd/cli

install: build
	install -m 755 $(BUILD_DIR)/dnsnode_cli    $(INSTALL_DIR)