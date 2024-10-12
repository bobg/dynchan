# Dynchan - Go channel with a dynamic buffer

[![Go Reference](https://pkg.go.dev/badge/github.com/bobg/dynchan.svg)](https://pkg.go.dev/github.com/bobg/dynchan)
[![Go Report Card](https://goreportcard.com/badge/github.com/bobg/dynchan)](https://goreportcard.com/report/github.com/bobg/dynchan)
[![Tests](https://github.com/bobg/dynchan/actions/workflows/go.yml/badge.svg)](https://github.com/bobg/dynchan/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/bobg/dynchan/badge.svg?branch=main)](https://coveralls.io/github/bobg/dynchan?branch=main)

This is dynchan,
a library that provides a Go-like channel with a dynamic buffer.

The buffer in a normal Go channel (when it has one) has a fixed size and behavior.
The buffer in a dynchan `Chan[T]` resizes dynamically
(so sends never block)
and can also change behavior.
This library includes a buffer with normal FIFO behavior,
plus another that has heap (a.k.a. priority-queue) behavior.
Other behaviors can be added by implementing the `Buffer` interface.
