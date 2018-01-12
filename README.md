# jsonptrerror - Extends [`encoding/json`](http://golang.org/pkg/encoding/json) with errors reported as JSON Pointer ([RFC 6901](https://tools.ietf.org/html/rfc6901))

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/dolmen-go/jsonptrerror)
[![Travis-CI](https://img.shields.io/travis/dolmen-go/jsonptrerror.svg)](https://travis-ci.org/dolmen-go/jsonptrerror)
[![Codecov](https://img.shields.io/codecov/c/github/dolmen-go/jsonptrerror/master.svg)](https://codecov.io/gh/dolmen-go/jsonptrerror/branch/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/dolmen-go/jsonptrerror)](https://goreportcard.com/report/github.com/dolmen-go/jsonptrerror)

This package is a wrapper around the standard [`encoding/json` package](https://golang.org/pkg/encoding/json)
to improve the reporting of location for unmarshalling errors:
[`UnmarshalTypeError`](https://golang.org/pkg/encoding/json/#UnmarshalTypeError) are enhanced to also include
a JSON Pointer (see RFC6901) indicating the location of the error: see
[`jsonptrerror.UnmarshalTypeError`](https://godoc.org/github.com/dolmen-go/jsonptrerror#UnmarshalTypeError).

## Status

The aim is code coverage of 100%. Use go coverage tools and consider any
code not covered by the testsuite as never tested and full of bugs.

See [latest coverage report](https://codecov.io/gh/dolmen-go/jsonptrerror/src/master/decoder.go) at codecov.io.

[![Coverage for dolmen-go/jsonptr](https://img.shields.io/codecov/c/github/dolmen-go/jsonptr/master.svg)](https://codecov.io/gh/dolmen-go/jsonptr/branch/master)
[![Coverage for dolmen-go/jsonptrerror](https://img.shields.io/codecov/c/github/dolmen-go/jsonptrerror/master.svg)](https://codecov.io/gh/dolmen-go/jsonptrerror/branch/master)


## License

Copyright 2017-2018 Olivier Mengué

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
