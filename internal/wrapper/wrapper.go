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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Wrapper struct {
	execName string
	clangPath string
}

func New(execName string) *Wrapper {
	return &Wrapper{
		execName: execName,
		// Clang binary should be in the same directory as the wrapper
		clangPath: filepath.Join(filepath.Dir(os.Args[0]), strings.TrimSuffix(execName, "-wrapper")),
	}
}

func (w *Wrapper) Run(args []string) error {
	skipLTOFiles := make(map[string]bool)
	var newArgs []string
	var targetFile string

	// Parse arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]
		
		// Handle --skip-lto argument
		if strings.HasPrefix(arg, "--skip-lto=") {
			files := strings.Split(strings.TrimPrefix(arg, "--skip-lto="), ";")
			//fmt.Println("files:", files)
			for _, f := range files {
				skipLTOFiles[filepath.Base(f)] = true
			}
			continue
		}

		// Track the input file
		if strings.HasSuffix(arg, ".c") || strings.HasSuffix(arg, ".cpp") || 
		   strings.HasSuffix(arg, ".cc") || strings.HasSuffix(arg, ".cxx") ||
		   strings.HasSuffix(arg, ".S") {
			targetFile = filepath.Base(arg)
		}

		newArgs = append(newArgs, arg)
	}

	// If the current file is in skip-lto list, add -fno-lto flag
	if targetFile != "" && skipLTOFiles[targetFile] {
		newArgs = append([]string{"-fno-lto"}, newArgs...)
	}

	// Temporarily modify PATH to prevent using system clang
	// Prepend the current directory to PATH
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath) // Restore original PATH after command runs

	os.Setenv("PATH", filepath.Dir(os.Args[0])+":"+oldPath)
	cmd := exec.Command("./"+w.clangPath, newArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}