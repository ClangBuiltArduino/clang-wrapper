# Clang Wrapper

A lightweight wrapper for Clang and Clang++ to work around limitations in the Arduino build system.

## Why did we make this wrapper?

The Arduino build system applies specific compiler flag combinations based on file extensions (`.c`, `.cpp`, `.S`). It determines the appropriate recipe by extracting the file extension and applying predefined flags. However, this approach lacks flexibility when per-file flag customization is required.

### The Problem

One specific issue arises when compiling the `wiring.c` file from the Arduino AVR core. Clang's Link-Time Optimization (LTO) causes breakage in the final program.

### The Workaround

This wrapper introduces a special `--skip-lto` flag, which allows users to specify files where LTO should be disabled. When invoked, the wrapper:
1. Checks if any of the specified files match the current compilation target.
2. If a match is found, appends `-fno-lto` to the compiler arguments before invoking Clang.
3. Otherwise, it forwards all arguments unchanged.

This approach ensures that essential functions remain intact while maintaining the benefits of LTO for other files. Future versions may extend functionality for additional per-file flag manipulations.

### Other functionality: Hanlding glibc and musl hosts

Since we dynamically link the **BFD linker** and the **LLVMgold.so** plugin, we provide separate binaries for **musl** and **glibc** hosts. However, the **Arduino IDE** lacks a mechanism to detect **musl-based** systems.  

To address this, we introduced two new flags:  

- `--bfd-dir` → Specifies the root directory containing **BFD linker** binaries for both **musl** and **glibc**.  
- `--llvmgold-dir` → Specifies the root directory containing **LLVMgold.so** binaries for both **musl** and **glibc**.  

Within these directories, binaries are organized into separate subdirectories for each **libc** variant. The **wrapper** automatically detects the system's **libc** and selects the appropriate binaries accordingly.

## Installation

```bash
make
sudo make install PREFIX=/usr/local
```

### Cross-Compiling

To build for Windows:
```bash
make windows
```

## Usage

```bash
# Disable LTO for specific files
clang-wrapper -O3 -flto --skip-lto=file1.c;file2.c -o output.o -c input.c

# Use as C++ compiler
clang++-wrapper -O3 -flto --skip-lto=file1.cpp;file2.cpp -o output.o -c input.cpp
```

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](https://github.com/ClangBuiltArduino/clang-wrapper/blob/main/LICENSE) file for details.

