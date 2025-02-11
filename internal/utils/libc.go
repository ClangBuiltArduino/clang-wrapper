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

package utils

import (
	"os/exec"
	"strings"
)

// detectLibC checks whether the system uses glibc or musl.
func DetectLibC() string {
	// Check getconf first
	cmd := exec.Command("getconf", "GNU_LIBC_VERSION")
	output, err := cmd.CombinedOutput()
	if err == nil && strings.Contains(string(output), "glibc") {
		return "glibc"
	}

	// Check for ldd (glibc always ships with it)
	cmd = exec.Command("ldd", "--version")
	output, err = cmd.CombinedOutput()
	if err == nil && strings.Contains(string(output), "glibc") {
		return "glibc"
	}

	// Check for glibc binary
	cmd = exec.Command("ldconfig", "-p")
	output, err = cmd.CombinedOutput()
	if err == nil && strings.Contains(string(output), "libc.so.6") {
		return "glibc"
	}

	// Check dynamic linker directly for musl
	cmd = exec.Command("strings", "/proc/self/exe")
	output, err = cmd.CombinedOutput()
	if err == nil && strings.Contains(string(output), "musl") {
		return "musl"
	}

	// Default to musl if no other checks confirm glibc
	return "musl"
}
