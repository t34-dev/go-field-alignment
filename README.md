[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Coverage Status](https://coveralls.io/repos/github/t34-dev/go-pad/badge.svg?branch=main&ver=1724705312)](https://coveralls.io/github/t34-dev/go-pad?branch=main&ver=1724705312)
![Go Version](https://img.shields.io/badge/Go-1.22-blue?logo=go&ver=1724705312)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/t34-dev/go-pad?ver=1724705312)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/t34-dev/go-pad?sort=semver&style=flat&logo=git&logoColor=white&label=Latest%20Version&color=blue&ver=1724705312)


# Go-Pad

Gopad is a powerful tool designed for Golang developers to enhance code readability by performing multi-level field alignment in struct declarations while preserving original metadata.

## Features

- Analyzes struct field alignment and padding in Go source files
- Calculates the size and alignment of each struct and its fields
- Optimizes struct layout by reordering fields for better memory efficiency
- Performs multi-level struct field alignment for improved readability
- Preserves original comments and metadata
- Supports nested structs and complex type hierarchies
- Processes single files or entire directories
- Offers an option to automatically apply optimizations to source files
- Easily integrates with existing Go projects

## Installation

To install `gopad`, make sure you have Go installed on your system, then run:

```shell
go install github.com/t34-dev/go-pad
```

For local installation

```shell
################ Bash
go build -o $GOPATH/bin/gopad       # unix
go build -o $GOPATH/bin/gopad.exe   # window

################ Makefile
make install                          # any system
```

## Usage

```
gopad [options] <file or directory paths>
```

### Options

- `-fix`: Apply fixes to optimize struct layout
- `-help`: Display help information

### Examples

Analyze a single file:
```
gopad main.go
```

Optimize structs in all Go files in the current directory:
```
gopad -fix .
```

Analyze all Go files in a specific directory:
```
gopad /path/to/project
```

## Output

For each struct found in the processed files, `gopad` will output:

- Struct name
- Total size of the struct
- Alignment of the struct
- For each field:
    - Field name
    - Field type
    - Offset within the struct
    - Size of the field
    - Alignment of the field

If the `-fix` option is used, it will also show the optimized layout of the struct.

## License

This project is licensed under the ISC License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.


---

Developed with ❤️ by [T34](https://github.com/t34-dev)
