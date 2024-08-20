//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build.Binary

var Aliases = map[string]interface{}{
	"bb": Build.Binary,
	"bc": Build.Container,
	"t":  Test,
	"l":  Lint,
}

type Build mg.Namespace

func (Build) Container() error {
	fmt.Println("Building container...")
	if err := sh.RunV("podman", "build", "-t", "ghcr.io/shikachuu/metrics-tool", "-f", "build/Dockerfile", "."); err != nil {
		return fmt.Errorf("failed to build container: %v", err)
	}

	return nil
}

// BuildBinary builds the binary for the application using the main.go file
func (Build) Binary() error {
	mg.Deps(Test)
	fmt.Println("Building...")

	if err := sh.Rm("bin/metrics-tool"); err != nil {
		return err
	}

	_, err := sh.Exec(
		map[string]string{"GOOS": "linux", "CGO_ENABLED": "0"},
		os.Stdout, os.Stderr,
		mg.GoCmd(), "build", "-o", "bin/metrics-tool", `-ldflags=-w`, "cmd/main.go",
	)

	return err
}

// Lint command runs fmt, vet and golangci-lint on the codebase
func Lint() error {
	fmt.Println("Running go fmt and go vet...")
	fmtCmd := exec.Command(mg.GoCmd(), "fmt", "./...")
	if err := fmtCmd.Run(); err != nil {
		return fmt.Errorf("go fmt failed: %v", err)
	}

	vetCmd := exec.Command(mg.GoCmd(), "vet", "./...")
	if err := vetCmd.Run(); err != nil {
		return fmt.Errorf("go vet failed: %v", err)
	}

	fmt.Println("Running golangci-lint...")

	lintCmd := exec.Command(mg.GoCmd(), "run", "github.com/golangci/golangci-lint/cmd/golangci-lint", "run", "./...")

	return lintCmd.Run()
}

// Test runs the unit tests for the given codebase
func Test() error {
	mg.Deps(Lint)
	return sh.RunV(mg.GoCmd(), "test", "-cover", "./...")
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command(mg.GoCmd(), "mod", "tidy")
	return cmd.Run()
}
