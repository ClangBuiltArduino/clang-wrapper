/*
 * Copyright (C) 2025 ClangBuiltArduino. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package wrapper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ClangBuiltArduino/clang-wrapper/internal/utils"
)

type Wrapper struct {
	execName    string
	clangPath   string
	verbose     bool
	bfdDir      string
	llvmGoldDir string
}

var gitSHA string

func New(execName string) *Wrapper {
	// Get the absolute path of the wrapper binary
	absPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving absolute path: %v\n", err)
		os.Exit(1)
	}

	// Get the directory where the wrapper is located
	wrapperDir := filepath.Dir(absPath)

	// Construct the expected Clang binary path
	clangBinary := strings.TrimSuffix(execName, "-wrapper")
	clangPath := filepath.Join(wrapperDir, clangBinary)

	return &Wrapper{
		execName:  execName,
		clangPath: clangPath,
	}
}

func (w *Wrapper) Run(args []string) error {
	for _, arg := range args {
		if arg == "--version" || arg == "--help" {
			fmt.Println("clang-wrapper: A wrapper to workaround Arduino build system limitations.")
			fmt.Println("git-commit:", gitSHA)
			fmt.Println("For more details check: https://github.com/ClangBuiltArduino/clang-wrapper")
			return nil
		}
	}

	skipLTOFiles := make(map[string]bool)
	var newArgs []string
	var targetFile string

	// Parse arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Handle --wrapper-verbose argument
		if arg == "--wrapper-verbose" {
			w.verbose = true
			continue // Do not pass this to clang
		}

		// Handle --skip-lto argument
		if strings.HasPrefix(arg, "--skip-lto=") {
			files := strings.Split(strings.TrimPrefix(arg, "--skip-lto="), ";")
			for _, f := range files {
				skipLTOFiles[filepath.Base(f)] = true
			}
			continue // Do not pass this to clang
		}

		// Handle --llvmgold-dir argument
		if strings.HasPrefix(arg, "--llvmgold-dir=") {
			w.llvmGoldDir = strings.TrimPrefix(arg, "--llvmgold-dir=")
			continue // Do not pass this to clang
		}

		// Handle --bfd-dir argument
		if strings.HasPrefix(arg, "--bfd-dir=") {
			w.bfdDir = strings.TrimPrefix(arg, "--bfd-dir=")
			continue // Do not pass this to clang
		}

		// Track the input file
		if strings.HasSuffix(arg, ".c") || strings.HasSuffix(arg, ".cpp") ||
			strings.HasSuffix(arg, ".cc") || strings.HasSuffix(arg, ".cxx") ||
			strings.HasSuffix(arg, ".S") {
			targetFile = filepath.Base(arg)
		}

		newArgs = append(newArgs, arg)
	}

	// Handle LLVMgold.so symlink creation
	if w.llvmGoldDir != "" && runtime.GOOS == "linux" {
		clangLibDir := filepath.Join(filepath.Dir(w.clangPath), "..", "lib")
		llvmGoldPath := filepath.Join(clangLibDir, "LLVMgold.so")

		if _, err := os.Stat(llvmGoldPath); os.IsNotExist(err) {
			libcType := utils.DetectLibC()
			sourcePath := filepath.Join(w.llvmGoldDir, libcType, "lib", "LLVMgold.so")
			if _, err := os.Stat(sourcePath); err == nil {
				if err := os.Symlink(sourcePath, llvmGoldPath); err != nil {
					fmt.Fprintf(os.Stderr, "Error creating symlink for LLVMgold.so: %v\n", err)
					return err
				}
				if w.verbose {
					fmt.Println("Created symlink:", llvmGoldPath, "->", sourcePath)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Source LLVMgold.so not found at %s\n", sourcePath)
			}
		}
	}

	// If the current file is in skip-lto list, add -fno-lto flag
	if targetFile != "" && skipLTOFiles[targetFile] {
		filteredArgs := []string{}
		for _, arg := range newArgs {
			if arg != "-flto" {
				filteredArgs = append(filteredArgs, arg)
			}
		}
		newArgs = filteredArgs
	}

	// Append --ld-path if --bfd-dir was provided
	if w.bfdDir != "" {
		var ldPath string
		if runtime.GOOS == "windows" {
			ldPath = filepath.Join(w.bfdDir, "bin", "avr-ld.bfd.exe")
		} else {
			libcType := utils.DetectLibC()
			ldPath = filepath.Join(w.bfdDir, libcType, "bin", "avr-ld.bfd")
		}
		newArgs = append(newArgs, fmt.Sprintf("--ld-path=%s", ldPath))
	}

	// Temporarily modify PATH to prevent using system clang
	// Prepend the current directory to PATH
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath) // Restore original PATH after command runs

	// Only print the execution command if verbose is enabled
	if w.verbose {
		fmt.Println("Executing:", w.clangPath, newArgs)
	}

	os.Setenv("PATH", filepath.Dir(os.Args[0])+":"+oldPath)
	cmd := exec.Command(w.clangPath, newArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
