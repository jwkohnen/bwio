/*
 * Copyright (c) 2016 Wolfgang Johannes Kohnen <wjkohnen@users.noreply.github.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package bwio provides wrappers for io.Reader, io.Writer, io.Copy and
// io.CopyBuffer that limit the throughput to a given bandwidth. The limiter
// uses an internal time bucket and hibernates each io operation for short time
// periods, whenever the configured bandwidth has been exceeded.
//
// The limiter tries to detect longer stalls and resets the bucket such that
// stalls do not cause subsequent high bursts. Usually you should choose small
// buffer sizes for low bandwidths and vice versa. The limiter tries to
// compensate for high buffer size / bandwidth ratio when detecting stalls, but
// this is not well tested.
package bwio

import (
	"io"
	"time"
)

type limiter struct {
	bandwidth     int
	start         time.Time
	bucket        int64
	isInitialized bool
}

func (l *limiter) init() {
	if !l.isInitialized {
		if l.bandwidth < 1 {
			panic("bandwidth must be greater than zero")
		}

		l.reset()
		l.isInitialized = true
	}
}

func (l *limiter) reset() {
	l.bucket = 0
	l.start = time.Now()
}

func (l *limiter) limit(n, bufsize int) {
	l.bucket += int64(n)
	bucketAge := time.Now().Sub(l.start)
	penalty := time.Duration(l.bucket)*time.Second/time.Duration(l.bandwidth) - bucketAge

	if penalty > 0 {
		time.Sleep(penalty)
		l.reset()
	} else {
		// Prevent peak after stall. Compensate in case of large buffer
		// and small bandwidth. TODO: The test cases could get more
		// love.
		compensation := time.Duration(bufsize/l.bandwidth) * time.Second
		stallThreshold := time.Second + compensation
		if bucketAge > stallThreshold {
			l.reset()
		}
	}
}

// Reader wraps another reader that maintains a given bandwidth.
type Reader struct {
	lim limiter
	src io.Reader
}

// NewReader returns a new reader that wraps reader r and maintains the
// given bandwidth.
func NewReader(r io.Reader, bandwidth int) *Reader {
	reader := &Reader{
		src: r,
		lim: limiter{bandwidth: bandwidth},
	}
	return reader
}

// Read implements the io.Reader interface, but maintains the given bandwidth.
func (r *Reader) Read(p []byte) (n int, err error) {
	r.lim.init()

	n, err = r.src.Read(p)
	if err != nil {
		// return all err, including io.EOF
		return n, err
	}

	r.lim.limit(n, len(p))

	return n, err
}

// Writer wraps another writer that maintains a given bandwidth.
type Writer struct {
	lim limiter
	dst io.Writer
}

// NewWriter returns a new writer that wraps writer d and maintains the given
// bandwidth.
func NewWriter(d io.Writer, bandwidth int) *Writer {
	writer := &Writer{
		dst: d,
		lim: limiter{bandwidth: bandwidth},
	}
	return writer
}

// Write implements the io.Writer interface, but maintains the given bandwidth.
func (w *Writer) Write(p []byte) (n int, err error) {
	w.lim.init()

	n, err = w.dst.Write(p)
	if err != nil {
		return n, err
	}

	w.lim.limit(n, len(p))

	return n, err
}

// Copy copies the same way io.Copy does, except maintaining the given
// bandwidth. It uses a buffer size of 16 KiBytes.
func Copy(dst io.Writer, src io.Reader, bandwidth int) (written int64, err error) {
	return CopyBuffer(dst, src, bandwidth, nil)
}

// CopyBuffer copies the same way io.CopyBuffer does, except maintaining the
// given bandwidth. If buf is nil, CopyBuffer will create a buffer with size of
// 16 KiBytes.
func CopyBuffer(dst io.Writer, src io.Reader, bandwidth int, buf []byte) (written int64, err error) {
	if len(buf) == 0 {
		buf = make([]byte, 16<<10)
	}
	bwreader := NewReader(src, bandwidth)
	return io.CopyBuffer(dst, bwreader, buf)
}
