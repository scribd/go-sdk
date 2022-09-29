//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Test.Run

type Test mg.Namespace

// Runs the tests.
func (Test) Run() error {
	goCmd := mg.GoCmd()
	env := map[string]string{
		"APP_ENV": "test",
	}
	return sh.RunWithV(env, goCmd, "test", "./...", "-race", "-count=1", "-v")
}

// Generates proto test files.
func (Test) GenerateProto() error {
	return sh.RunV(
		"protoc",
		"--proto_path=./pkg/testing/testproto",
		"--go_out=./pkg/testing/testproto",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=./pkg/testing/testproto",
		"--go-grpc_opt=paths=source_relative",
		"sdktest.proto")
}

type Fmt mg.Namespace

// Runs gofmt.
func (Fmt) Run() error {
	goCmd := mg.GoCmd()
	return sh.RunV(goCmd, "fmt", "./...")
}

// Checks the code formatting.
func (Fmt) Check() error {
	goModule, err := sh.Output(mg.GoCmd(), "list", "-m")
	if err != nil {
		return err
	}

	// flag -local string: put imports beginning with this string after 3rd-party packages; comma-separated list.
	// flag -e: report all errors (not just the first 10 on different lines).
	// flag -d: display diffs instead of rewriting files.
	output, err := sh.Output("goimports", "-local="+goModule, "-e", "-d", ".")
	if err != nil {
		return err
	}

	if output != "" {
		return fmt.Errorf("source code is not formatted:\n\n%s", output)
	}

	return nil
}

// Runs the linter.
func (Fmt) Lint() error {
	lintCmd := "golangci-lint"
	return sh.RunV(lintCmd, "run")
}
