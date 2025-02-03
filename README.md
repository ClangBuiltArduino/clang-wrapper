# Clang Wrapper

A lightweight wrapper for Clang and Clang++ to work around limitations in the Arduino build system.

## Why did we make this wrapper?

The Arduino build system applies specific compiler flag combinations based on file extensions (`.c`, `.cpp`, `.S`). It determines the appropriate recipe by extracting the file extension and applying predefined flags. However, this approach lacks flexibility when per-file flag customization is required.

### The Problem

One specific issue arises when compiling the `HardwareSerial1.cpp` file from the Arduino AVR core. Clang's Link-Time Optimization (LTO) aggressively removes what it considers unused functions, causing critical functionality to break in the final binary.

### The Solution

This wrapper introduces a special `--skip-lto` flag, which allows users to specify files where LTO should be disabled. When invoked, the wrapper:
1. Checks if any of the specified files match the current compilation target.
2. If a match is found, appends `-fno-lto` to the compiler arguments before invoking Clang.
3. Otherwise, it forwards all arguments unchanged.

This approach ensures that essential functions remain intact while maintaining the benefits of LTO for other files. Future versions may extend functionality for additional per-file flag manipulations.

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

