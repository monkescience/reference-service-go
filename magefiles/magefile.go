//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	appName    = "reference-service"
	binaryName = "service"
	cmdPath    = "./cmd/main.go"
	buildDir   = "./build"
	chartPath  = "./chart"
)

// getVersion returns the version from git or "dev" if unavailable
func getVersion() string {
	version, err := sh.Output("git", "describe", "--tags", "--always", "--dirty")
	if err != nil {
		return "dev"
	}
	return version
}

// Build builds the application binary
func Build() error {
	fmt.Println("Building", binaryName, "...")
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return err
	}
	output := filepath.Join(buildDir, binaryName)
	return sh.RunWith(map[string]string{"CGO_ENABLED": "0"}, "go", "build", "-o", output, cmdPath)
}

// Run runs the application locally
func Run() error {
	fmt.Println("Running", appName, "...")
	return sh.RunV("go", "run", cmdPath)
}

// Generate runs go generate to generate code from OpenAPI specs
func Generate() error {
	fmt.Println("Generating code from OpenAPI specifications...")
	return sh.RunV("go", "generate", "./...")
}

// Test runs tests with race detection
func Test() error {
	fmt.Println("Running tests...")
	return sh.RunV("go", "test", "-v", "-race", "./...")
}

// Coverage runs tests with coverage report
func Coverage() error {
	fmt.Println("Running tests with coverage...")
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return err
	}
	coverageOut := filepath.Join(buildDir, "coverage.out")
	coverageHTML := filepath.Join(buildDir, "coverage.html")

	if err := sh.RunV("go", "test", "-v", "-race", "-coverprofile="+coverageOut, "-covermode=atomic", "./..."); err != nil {
		return err
	}
	if err := sh.RunV("go", "tool", "cover", "-html="+coverageOut, "-o", coverageHTML); err != nil {
		return err
	}
	fmt.Printf("Coverage report generated at %s\n", coverageHTML)
	return nil
}

// Fmt formats Go code
func Fmt() error {
	fmt.Println("Formatting code...")
	return sh.RunV("golangci-lint", "fmt")
}

// Lint runs golangci-lint
func Lint() error {
	fmt.Println("Running linter...")
	return sh.RunV("golangci-lint", "run", "--timeout=5m")
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning build artifacts...")
	if err := sh.Rm(buildDir); err != nil {
		return err
	}
	return sh.Run("go", "clean")
}

// Docker namespace for Docker operations
type Docker mg.Namespace

// Build builds the Docker image
func (Docker) Build() error {
	fmt.Printf("Building Docker image %s:latest...\n", appName)
	version := getVersion()
	return sh.RunV("docker", "build",
		"--build-arg", "VERSION="+version,
		"-t", appName+":latest",
		".")
}

// Run runs the Docker container
func (Docker) Run() error {
	fmt.Println("Running Docker container...")
	return sh.RunV("docker", "run", "--rm", "-p", "8080:8080", appName+":latest")
}

// Helm namespace for Helm operations
type Helm mg.Namespace

// Lint lints the Helm chart
func (Helm) Lint() error {
	fmt.Println("Linting Helm chart...")
	return sh.RunV("helm", "lint", chartPath)
}

// Template renders Helm chart templates
func (Helm) Template() error {
	fmt.Println("Rendering Helm templates...")
	return sh.RunV("helm", "template", appName, chartPath)
}

// Mod namespace for Go module operations
type Mod mg.Namespace

// Tidy tidies Go module dependencies
func (Mod) Tidy() error {
	fmt.Println("Tidying dependencies...")
	return sh.RunV("go", "mod", "tidy")
}
