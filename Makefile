# Compiler settings
GO ?= go
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)
PREFIX ?= /usr/local

# Version information
GIT_SHA := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS := -X github.com/ClangBuiltArduino/clang-wrapper/internal/wrapper.gitSHA=$(GIT_SHA)

# Binary names
BINARY := clang-wrapper
BINARY_WINDOWS := $(BINARY).exe

# Install paths
INSTALL_PATH := $(PREFIX)/bin

.PHONY: all
all: build

.PHONY: build
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/clang-wrapper

.PHONY: windows
windows:
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_WINDOWS) ./cmd/clang-wrapper

.PHONY: install
install: build
	install -d $(INSTALL_PATH)
	install -m 755 $(BINARY) $(INSTALL_PATH)/$(BINARY)
	ln -sf $(BINARY) $(INSTALL_PATH)/clang++-wrapper

.PHONY: uninstall
uninstall:
	rm -f $(INSTALL_PATH)/$(BINARY)
	rm -f $(INSTALL_PATH)/clang++-wrapper

.PHONY: clean
clean:
	rm -f $(BINARY) $(BINARY_WINDOWS)
	rm -f test/mock_compiler/mock_compiler

# Additional test targets
.PHONY: test
test: build_mock_compiler
	go test -v ./test
	rm -f test/mock_compiler/mock_compiler

.PHONY: build_mock_compiler
build_mock_compiler:
	go build -o test/mock_compiler/mock_compiler test/mock_compiler/mock_compiler.go

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        - Build for current platform (default)"
	@echo "  build      - Same as 'all'"
	@echo "  windows    - Cross-compile for Windows"
	@echo "  install    - Install binary and symlinks to PREFIX (default: /usr/local)"
	@echo "  clean      - Remove built binaries"
	@echo "  uninstall  - Remove installed binary and symlinks"
	@echo "  test       - Run unit tests"
	@echo "  coverage   - Run tests with coverage report"
	@echo "  check      - Run all tests including integration tests"
	@echo ""
	@echo "Variables:"
	@echo "  PREFIX     - Installation prefix (default: /usr/local)"
	@echo "  GOOS       - Target operating system"
	@echo "  GOARCH     - Target architecture"