//go:build tools
// +build tools

package tools

import (
    _ "github.com/magefile/mage"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
