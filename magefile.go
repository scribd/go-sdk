//+build mage

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
	return sh.RunWith(env, goCmd, "test", "./...", "-count=1", "-v")
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
	goCmd := "gofmt"
	o, e := sh.OutCmd(goCmd)("-l", ".")
	if e != nil {
		return e
	}
	if o != "" {
		fmtRes, e := sh.OutCmd(goCmd)("-d", ".")
		if e != nil {
			return e
		}
		return fmt.Errorf("Go code is not formatted:\n\n%s", fmtRes)
	}
	return nil
}

// Runs the linter.
func (Fmt) Lint() error {
	lintCmd := "golangci-lint"
	return sh.RunV(lintCmd, "run")
}
