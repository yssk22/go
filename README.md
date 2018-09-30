# yssk22 Go

[![CircleCI](https://circleci.com/gh/yssk22/go/tree/master.svg?style=svg)](https://circleci.com/gh/yssk22/go/tree/master)
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/yssk22/go)

A set of go packages and tools used in yssk22.net internal tools.

## Usage

```bash
$ go get github.com/yssk22/go
```

## Code structure

- `x/*` packages provide extended functionality of default packages.
- `tools/cmd/*` packages provide command line tools to go project development.
  - `tools/cmd/enum` to generate enum functions from type aliases.

See full doc at [godoc.org](https://godoc.org/github.com/yssk22/go).
