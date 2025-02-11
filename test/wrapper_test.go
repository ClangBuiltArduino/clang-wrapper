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

package test

import (
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/ClangBuiltArduino/clang-wrapper/internal/utils"
)

func TestClangWrapper(t *testing.T) {
	// Define tests with input flags and expected output
	libcType := utils.DetectLibC()
	bfdDir := "/home/test/bfd" // Mock directory for test

	expectedLdPath := ""
	if runtime.GOOS == "windows" {
		expectedLdPath = "--ld-path=/home/test/bfd/bin/avr-ld.bfd.exe"
	} else {
		expectedLdPath = "--ld-path=/home/test/bfd/" + libcType + "/bin/avr-ld.bfd"
	}

	tests := []struct {
		name          string
		args          []string
		expectedFlags string
	}{
		{
			name:          "Clang Wrapper: --version",
			args:          []string{"--version"},
			expectedFlags: "",
		},
		{
			name:          "Clang Wrapper: --help",
			args:          []string{"--help"},
			expectedFlags: "",
		},
		{
			name:          "Clang Wrapper: no --skip-lto",
			args:          []string{"test.c"},
			expectedFlags: "test.c",
		},
		{
			name:          "Clang Wrapper: --skip-lto=file1.c test.c",
			args:          []string{"-O3", "-flto", "--skip-lto=file1.c", "test.c"},
			expectedFlags: "test.c",
		},
		{
			name:          "Clang Wrapper: --skip-lto=file1.c file1.c",
			args:          []string{"-O3", "-flto", "--skip-lto=file1.c", "file1.c"},
			expectedFlags: "-O3 file1.c",
		},
		{
			name:          "Clang Wrapper: --skip-lto=file1.c;file2.c test.c",
			args:          []string{"-O3", "-flto", "--skip-lto=file1.c;file2.c", "test.c"},
			expectedFlags: "-O3 -flto test.c",
		},
		{
			name:          "Clang Wrapper: --skip-lto=file1.c;file2.c file1.c",
			args:          []string{"-O3", "-flto", "--skip-lto=file1.c;file2.c", "file1.c"},
			expectedFlags: "-O3 file1.c",
		},
		{
			name:          "Clang Wrapper: --bfd-dir argument",
			args:          []string{"--bfd-dir=" + bfdDir, "test.c"},
			expectedFlags: expectedLdPath,
		},
	}

	// Create a temporary directory for the test
	tmpDir := t.TempDir()

	// Copy the clang-wrapper binary the test directory
	err := copyFile("../clang-wrapper", tmpDir+"/clang-wrapper")
	if err != nil {
		t.Fatalf("Error copying clang-wrapper: %v", err)
	}

	// Create clang++-wrapper symlink and set executable permission
	err = os.Symlink(tmpDir+"/clang-wrapper", tmpDir+"/clang++-wrapper")
	if err != nil {
		t.Fatalf("Error creating symlink for clang++-wrapper: %v", err)
	}
	err = os.Chmod(tmpDir+"/clang++-wrapper", 0755) // Make clang-wrapper executable
	if err != nil {
		t.Fatalf("Error setting permissions on clang++-wrapper: %v", err)
	}

	// Copy mock compiler to the test directory and set executable permission
	err = copyFile("./mock_compiler/mock_compiler", tmpDir+"/mock_compiler")
	if err != nil {
		t.Fatalf("Error copying mock_compiler: %v", err)
	}
	err = os.Chmod(tmpDir+"/mock_compiler", 0755) // Make mock_compiler executable
	if err != nil {
		t.Fatalf("Error setting permissions on mock_compiler: %v", err)
	}

	// Create symlinks to the mock compiler for clang and clang++
	err = os.Symlink(tmpDir+"/mock_compiler", tmpDir+"/clang")
	if err != nil {
		t.Fatalf("Error creating symlink for clang: %v", err)
	}
	err = os.Symlink(tmpDir+"/mock_compiler", tmpDir+"/clang++")
	if err != nil {
		t.Fatalf("Error creating symlink for clang++: %v", err)
	}

	// Change the working directory to the temporary directory for the test
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Error changing directory: %v", err)
	}

	// Run tests for both clang-wrapper and clang++-wrapper
	for _, tt := range tests {
		t.Run(tt.name+" with clang-wrapper", func(t *testing.T) {
			// Prepare the command
			cmd := exec.Command("./clang-wrapper", tt.args...)
			output, err := cmd.CombinedOutput()

			// Check for expected error
			if err != nil && len(tt.expectedFlags) > 0 {
				t.Fatalf("Error executing clang-wrapper: %s", output)
				t.Fatalf("Error executing clang-wrapper: %v", err)
			}

			// Check if the output matches the expected flags
			outputStr := string(output)
			if !strings.Contains(outputStr, tt.expectedFlags) {
				t.Errorf("Expected flags to contain %q, but got %q", tt.expectedFlags, outputStr)
			}
		})

		t.Run(tt.name+" with clang++-wrapper", func(t *testing.T) {
			// Prepare the command for clang++
			cmd := exec.Command("./clang++-wrapper", tt.args...)
			output, err := cmd.CombinedOutput()

			// Check for expected error
			if err != nil && len(tt.expectedFlags) > 0 {
				t.Fatalf("Error executing clang++-wrapper: %v", err)
			}

			// Check if the output matches the expected flags
			outputStr := string(output)
			if !strings.Contains(outputStr, tt.expectedFlags) {
				t.Errorf("Expected flags to contain %q, but got %q", tt.expectedFlags, outputStr)
			}
		})
	}
}

// Helper function to copy file
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
