// +build mage

/*
	Copyright (c) 2020 - 2021 Digit Game Studios Ltd. - All Rights Reserved
	Copyright (c) 2020 - 2021 Scopely Inc. - All Rights Reserved

	Reproduction of this material is strictly forbidden unless prior written
	permission is obtained from Scopely Inc. and/or Digit Game Studios Ltd,
*/

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	color "github.com/alecthomas/colour"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

//
// Required constants and variables
//

const (
	packageNamePrefix = "hq.0xa1.red/axdx/scheduler/cmd/"
)

const (
	Check string = "\u2713"
	Cross string = "\u2717"
)

var ldflags = "-s -w -X main.commitHash=$COMMIT_HASH -X main.buildtime=$BUILD_TIME -X main.tag=$VERSION_TAG"

var Default = Release
var binaries = []string{"scheduler"}
var targetOS = []string{"linux", "darwin"}
var checksumFormats = []string{"sha256", "md5"}

//
// Magefile targets
//

// Clean removes any previous compilation products
func Clean() error {

	if err := sh.Run("go", "clean", "./..."); err != nil {
		return err
	}

	if err := sh.Rm("./scheduler"); err != nil {
		return err
	}

	if err := sh.Rm("./target"); err != nil {
		return err
	}

	if err := sh.Rm("./cover"); err != nil {
		return err
	}

	return nil
}

// Test run tests on project
func Test() error {

	env := map[string]string{"CGO_ENABLED": "0"}
	return runWithV(env, "go", "test", "-v", "-count=1", "./...")
}

// TestSilent run tests on silent mode
func TestSilent() error {

	env := map[string]string{"CGO_ENABLED": "0"}
	return runWithV(env, "go", "test", "./...")
}

// Coverage run tests and outputs a HTML coverage report
func Coverage() error {

	if err := sh.Run("mkdir", "-p", "./cover/"); err != nil {
		return fmt.Errorf("could not create coverage directory: %w", err)
	}

	env := map[string]string{"CGO_ENABLED": "0"}
	commitHash := getCommitHash()
	log.Printf("commit hash: %s", commitHash)
	err := runWithV(env, "go", "test", "-coverprofile", fmt.Sprintf("./cover/%s.out", commitHash), "./...")
	if err != nil {
		return fmt.Errorf("could not run coverage test: %w", err)
	}

	if err := sh.Run("go", "tool", "cover", "-html", fmt.Sprintf("./cover/%s.out", commitHash), "-o", fmt.Sprintf("./cover/%s.html", commitHash)); err != nil {
		return fmt.Errorf("could not convert coverage report into HTML format: %w", err)
	}

	return nil
}

// Cover run tests and opens a browser to its HTML coverage report
func Cover() error {

	if err := sh.Run("mkdir", "-p", "./cover/"); err != nil {
		return fmt.Errorf("could not create coverage directory: %w", err)
	}

	env := map[string]string{"CGO_ENABLED": "0"}
	commitHash := getCommitHash()
	err := runWithV(env, "go", "test", "-coverprofile", fmt.Sprintf("./cover/%s.out", commitHash), "./...")
	if err != nil {
		return fmt.Errorf("could not run coverage test: %w", err)
	}

	if err := sh.Run("go", "tool", "cover", "-html", fmt.Sprintf("./cover/%s.out", commitHash)); err != nil {
		return fmt.Errorf("could not convert coverage report into HTML format: %w", err)
	}

	return nil
}

// Benchmark runs any benchmarks present in the code and shows an output with the results
func Benchmark() error {

	env := map[string]string{"CGO_ENABLED": "0"}
	err := runWithV(env, "go", "test", "-run", "none", "-bench=.", "-benchmem", "-benchtime=5s", "./...")
	if err != nil {
		return fmt.Errorf("could not execute project benchmarks: %w", err)
	}

	return nil
}

// Linter runs the configured battery of linters over the code
func Linter() error {

	fmt.Println("running golang-ci-lint...")
	golangciPath, lookupErr := lookupLinter()
	if lookupErr == nil {
		err := sh.RunV(golangciPath, "run", "-v", "--config", "./.golangci.yml", "./...")
		if err != nil {
			return fmt.Errorf("golangci-lint returned an error: %w", err)
		}
	}

	return lookupErr
}

// EssenceD builds essenced binary
func Scheduler() error {

	for _, binary := range binaries {
		buildTag := "netgo"
		err := runWithV(flagEnv(), "go", "build", "-ldflags", ldflags, "-tags", buildTag, "-o", binary, packageNamePrefix+binary)
		if err != nil {
			return err
		}
	}

	return nil
}

// Race builds the essence binary with race condition track support
func Race() error {

	for _, binary := range binaries {
		buildTag := "netgo"
		flags := flagEnv()
		flags["CGO_ENABLED"] = "1"
		err := runWithV(flags, "go", "build", "-race", "-ldflags", ldflags, "-tags", buildTag, "-o", binary, packageNamePrefix+binary)
		if err != nil {
			return err
		}
	}

	return nil
}

// Install installs the essenced binary
func Install() error {

	// mg.SerialDeps runs the functions passed to it in the order they appear
	env := flagEnv()
	env["CGO_ENABLED"] = "0"
	mg.SerialDeps(Clean, Test)
	return runWithV(env, "go", "install", "-ldflags", ldflags, "-tags", "netgo", packageNamePrefix+binaries[0])
}

// Release compiles all essence binaries and targets after a test run
func Release() error {

	env := flagEnv()
	mg.SerialDeps(Clean, TestSilent)

	for _, binary := range binaries {

		for _, target := range targetOS {

			fmt.Printf("compiling %s for %s, please wait...", binary, target)

			buildTag := "netgo"
			env["GOOS"] = target
			if err := runWithV(
				env, "go", "build", "-ldflags", ldflags,
				"-tags", buildTag, "-o", fmt.Sprintf("./target/%s/%s", target, binary), packageNamePrefix+binary,
			); err != nil {
				color.Printf(" ^1%s^R\n", Cross)
				return fmt.Errorf("failed to compile target %s of binary %s: %w", target, binary, err)
			}

			color.Printf(" ^2%s^R\n", Check)
			fmt.Printf("generating sha256 and md5 checksum files for %s target of %s binary...", target, binary)
			if err := generateCheckSumFiles(binary, target); err != nil {
				color.Printf(" ^1%s^R\n", Cross)
				return fmt.Errorf("failed to generate checksum files for target %s of binary %s: %w", target, binary, err)
			}

			color.Printf(" ^2%s^R\n", Check)
		}
	}

	return nil
}

//
//  Function helpers
//

// getCommitHash returns the last commit hash using git
func getCommitHash() string {

	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}

// getVersionTag returns the latest version tag defaulting to .release_version
// file is present (declared on CI pipelines) and falling back to git describe
func getVersionTag() string {

	versiontag, err := sh.Output("cat", ".release_version")
	if err != nil {
		versiontag, _ = sh.Output("git", "describe", "--tags", "--abbrev=0")
	}

	return versiontag
}

// flagEnv creates and return a map with environment variables to inject to the linker
func flagEnv() map[string]string {

	return map[string]string{
		"COMMIT_HASH": getCommitHash(),
		"BUILD_TIME":  time.Now().Format("2006-01-02T15:04:05Z0700"),
		"VERSION_TAG": getVersionTag(),
		"CGO_ENABLED": "0",
	}
}

// generateCheckSumFiles generates sha256 and md5 sum files for the given file system file
func generateCheckSumFiles(binary, target string) error {

	targetpath := fmt.Sprintf("target/%s/%s", target, binary)
	for _, checksumFormat := range checksumFormats {
		out, err := sh.Output(fmt.Sprintf("%ssum", checksumFormat), targetpath)
		if err != nil {
			return err
		}

		f, openerr := os.OpenFile(fmt.Sprintf("%s.%s", targetpath, checksumFormat), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
		if openerr != nil {
			return openerr
		}

		defer f.Close()
		if _, err := f.Write([]byte(strings.Split(out, " ")[0])); err != nil {
			return err
		}
	}

	return nil
}

// lookupLinter looks up for golangci-lint binary on the system
func lookupLinter() (string, error) {
	return exec.LookPath("golangci-lint")
}

// prepare makes sure the compilation output directories exists on the file system
func prepare() {
	sh.Run("mkdir", "-p", "target/linux target/darwin")
}

// runWithV sends the command's stdout to os.Stdout on real time
func runWithV(env map[string]string, cmd string, args ...string) error {

	_, err := sh.Exec(env, os.Stdout, os.Stderr, cmd, args...)
	return err
}

// getMagefileDirectory returns back the magefile location in the local disk
func getMagefileDirectory() string {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic(fmt.Errorf("could not determine current file path"))
	}

	return filepath.Dir(filename)
}

//
// Init
//

func init() {

	// make sure we use Go 1.11 modules even if the source lives inside GOPATH
	os.Setenv("GO111MODULE", "on")
}
