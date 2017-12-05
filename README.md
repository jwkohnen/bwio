# BWIO: io wrappers that limit bandwidth

[![Apache License v2.0](https://img.shields.io/badge/license-Apache%20License%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0.txt)
[![GoDoc](https://godoc.org/github.com/wjkohnen/bwio?status.svg)](https://godoc.org/github.com/wjkohnen/bwio)
[![Build Status](https://travis-ci.org/wjkohnen/bwio.svg?branch=master)](https://travis-ci.org/wjkohnen/bwio)
[![Go Report](https://goreportcard.com/badge/github.com/wjkohnen/bwio)](https://goreportcard.com/report/github.com/wjkohnen/bwio)
[![codebeat badge](https://codebeat.co/badges/8bf05d3c-f74d-4a86-9cc7-dd3aaaac2213)](https://codebeat.co/projects/github-com-wjkohnen-bwio)
[![codecov](https://codecov.io/gh/wjkohnen/bwio/branch/master/graph/badge.svg)](https://codecov.io/gh/wjkohnen/bwio)

Package bwio provides wrappers for io.Reader, io.Writer, io.Copy and
io.CopyBuffer that limit the throughput to a given bandwidth. The limiter
uses an internal time bucket and hibernates each io operation for short time
periods, whenever the configured bandwidth has been exceeded.

The limiter tries to detect longer stalls and resets the bucket such that
stalls do not cause subsequent high bursts. Usually you should choose small
buffer sizes for low bandwidths and vice versa. The limiter tries to
compensate for high buffer size / bandwidth ratio when detecting stalls, but
this is not well tested.

## Monotonic time

Support for monotonic time before Go version 1.8 was implemented in a branch,
but has been dropped. Use Go version 1.9 or later in order to profit from 
transparent non-monotonic time robustness.

## License

Copyright (c) 2017 Johannes Kohnen <wjkohnen@users.noreply.github.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
