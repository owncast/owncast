# Development tools are managed in tools/go.mod and installed to ./bin
GOBIN := $(shell pwd)/bin

# Tool binaries
LEFTHOOK := $(GOBIN)/lefthook
GOLANGCI_LINT := $(GOBIN)/golangci-lint
GOFUMPT := $(GOBIN)/gofumpt
SQLC := $(GOBIN)/sqlc
OAPI_CODEGEN := $(GOBIN)/oapi-codegen

.PHONY: install-tools install-hooks lint fmt sqlc api-generate build test clean

## Install all development tools to ./bin
install-tools: $(LEFTHOOK) $(GOLANGCI_LINT) $(GOFUMPT) $(SQLC) $(OAPI_CODEGEN)

$(LEFTHOOK):
	GOBIN=$(GOBIN) go install -C tools github.com/evilmartians/lefthook

$(GOLANGCI_LINT):
	GOBIN=$(GOBIN) go install -C tools github.com/golangci/golangci-lint/v2/cmd/golangci-lint

$(GOFUMPT):
	GOBIN=$(GOBIN) go install -C tools mvdan.cc/gofumpt

$(SQLC):
	GOBIN=$(GOBIN) go install -C tools github.com/sqlc-dev/sqlc/cmd/sqlc

$(OAPI_CODEGEN):
	GOBIN=$(GOBIN) go install -C tools github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen

## Install git hooks using lefthook
install-hooks: $(LEFTHOOK)
	$(LEFTHOOK) install
	@echo "Patching hooks to use ./bin/lefthook..."
	@for hook in .git/hooks/*; do \
		if [ -f "$$hook" ] && grep -q 'call_lefthook' "$$hook" && ! grep -q 'export LEFTHOOK_BIN' "$$hook"; then \
			sed -i.bak 's|call_lefthook|export LEFTHOOK_BIN="$$(git rev-parse --show-toplevel)/bin/lefthook"; call_lefthook|' "$$hook" && rm -f "$$hook.bak"; \
		fi \
	done
	@echo "Hooks installed successfully"

## Run golangci-lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

## Format Go code with gofumpt
fmt: $(GOFUMPT)
	$(GOFUMPT) -l -w .

## Generate database models with sqlc
sqlc: $(SQLC)
	$(SQLC) generate

## Generate API code from OpenAPI spec
api-generate: $(OAPI_CODEGEN)
	./build/gen-api.sh

## Build the application
build:
	go build -o owncast .

## Run tests
test:
	go test ./...

## Clean build artifacts
clean:
	rm -rf bin/
	rm -f owncast

## Show help
help:
	@echo "Available targets:"
	@echo "  install-tools  - Install all development tools to ./bin"
	@echo "  install-hooks  - Install git hooks using lefthook"
	@echo "  lint           - Run golangci-lint"
	@echo "  fmt            - Format Go code with gofumpt"
	@echo "  sqlc           - Generate database models"
	@echo "  api-generate   - Generate API code from OpenAPI spec"
	@echo "  build          - Build the application"
	@echo "  test           - Run tests"
	@echo "  clean          - Remove build artifacts and tools"
