markdown
### Hexlet tests and linter status:
[![Actions Status](https://github.com/misterDillett/go-project-242/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/misterDillett/go-project-242/actions)

## Hexlet Path Size

A CLI tool to calculate the size of files and directories with support for human-readable output, hidden files, and recursive directory traversal.

## Installation

```bash
go install hexlet-boilerplates/gopackage/cmd/hexlet-path-size@latest
```

or build from source:
```bash
git clone https://github.com/misterDillett/go-project-242
cd go-project-242
make build
```

## Usage
```bash
# Basic usage
./bin/hexlet-path-size file.txt
# 123B	file.txt

# Human-readable format
./bin/hexlet-path-size -H directory/
# 24.0MB	directory/

# Include hidden files
./bin/hexlet-path-size -a directory/
# 27.0MB	directory/

# Recursive directory traversal
./bin/hexlet-path-size -r directory/
# 31.0MB	directory/

# Combine flags
./bin/hexlet-path-size -H -a -r directory/
# 31.0MB	directory/
```

##Flags
-r, --recursive	Recursive size of directories
-H, --human	Human-readable sizes (auto-select unit)
-a, --all	Include hidden files and directories
-h, --help	Show help

##Development
```bash
# Build
make build

# Run linter
make lint

# Run tests
make test

# Format code
make fmt
```
