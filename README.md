jsoncompare
======

[![Build Status](https://travis-ci.org/martinohmann/jsoncompare.svg)](https://travis-ci.org/martinohmann/jsoncompare)
[![Go Report Card](https://goreportcard.com/badge/github.com/martinohmann/jsoncompare)](https://goreportcard.com/report/github.com/martinohmann/jsoncompare)
[![GoDoc](https://godoc.org/github.com/martinohmann/jsoncompare?status.svg)](https://godoc.org/github.com/martinohmann/jsoncompare)

jsoncompare is a helper utility to compare two byte slices containing json with each other. The matching behaviour is configurable. This is mostly useful for assertions in tests where we want to validate that a given haystack contains a needle but we do not mind if the haystack contains additional data not present in needle.

Installation
------------

```sh
go get -u github.com/martinohmann/jsoncompare
```

Usage
-----

See [`example_test.go`](example_test.go) for usage examples.

License
-------

The source code of jsoncompare is released under the MIT License. See the bundled
LICENSE file for details.
