# BWIO: io wrappers that limit bandwidth

[![Apache License v2.0](https://img.shields.io/badge/license-Apache%20License%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0.txt)
[![GoDoc](https://godoc.org/github.com/wjkohnen/bwio?status.svg)](https://godoc.org/github.com/wjkohnen/bwio)


Package bwio provides wrappers for io.Reader, io.Writer, io.Copy and
io.CopyBuffer that limit the throughput to a given bandwidth. The limiter
uses an internal time bucket and hibernates each io operation for short time
time periods, whenever the configured bandwidth has been exceeded.

The limiter tries to detect longer stalls and resets the bucket such that
stalls do not cause subsequent high bursts. Usually you should choose small
buffer sizes for low bandwidths and vice versa. The limiter tries to
compensate for high buffer size / bandwidth ratio when detecting stalls, but
this is not well tested.

Licensed under the Apache License, Version 2.0.
