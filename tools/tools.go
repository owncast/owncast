//go:build tools

// Package tools manages development tool dependencies.
// This file is used to track tool dependencies separately from the main module,
// allowing `go mod vendor` to work without vendoring tool dependencies.
//
// Use the Makefile to install and run tools:
//
//	make install-tools  # Install all tools to ./bin
//	make install-hooks  # Install git hooks
//	make lint           # Run golangci-lint
//	make fmt            # Format code with gofumpt
//	make sqlc           # Generate database models
//	make api-generate   # Generate API code
//
// See `make help` for all available targets.
package tools

import (
	_ "github.com/evilmartians/lefthook"
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	_ "mvdan.cc/gofumpt"
)
