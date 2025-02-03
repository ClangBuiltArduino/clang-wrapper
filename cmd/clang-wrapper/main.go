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

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ClangBuiltArduino/clang-wrapper/internal/wrapper"
)

func main() {
	// Get the binary name to determine if we're running as clang or clang++
	execName := filepath.Base(os.Args[0])

	// Print execName
	//fmt.Println("Executable Name:", execName)

	w := wrapper.New(execName)
	if err := w.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}