/*
 * Copyright (c) 2017 Johannes Kohnen <wjkohnen@users.noreply.github.com>
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

package bwio

import (
	"io"
	"io/ioutil"
	"testing"
	"time"
)

type BWTestReader struct {
	size  int
	count int
	stall time.Duration
}

func (r *BWTestReader) Read(p []byte) (n int, err error) {
	l := len(p)
	if l < r.size {
		n = l
	} else {
		n = r.size
		err = io.EOF
	}
	if r.count == 1 {
		time.Sleep(r.stall)
	}
	r.count++
	r.size -= n
	return n, err
}

func TestRead(t *testing.T) {
	t.Parallel()

	tr := &BWTestReader{size: 1 << 20, stall: 2 * time.Second}
	br := NewReader(tr, 500<<10)

	start := time.Now()
	n, err := io.Copy(ioutil.Discard, br)
	dur := time.Since(start)
	if err != nil {
		t.Error(err)
	}
	if n != 1<<20 {
		t.Errorf("Want %d bytes, got %d.", 1<<20, n)
	}
	t.Logf("Read %d bytes in %s", n, dur)
	if dur < 3600*time.Millisecond || dur > 4400*time.Millisecond {
		t.Errorf("Took %s, want 4s.", dur)
	}
}

func TestWrite(t *testing.T) {
	t.Parallel()

	tr := &BWTestReader{size: 1 << 20, stall: 2 * time.Second}
	bw := NewWriter(ioutil.Discard, 500<<10)

	start := time.Now()
	n, err := io.Copy(bw, tr)
	dur := time.Since(start)
	if err != nil {
		t.Error(err)
	}
	if n != 1<<20 {
		t.Errorf("Want %d bytes, got %d.", 1<<20, n)
	}
	t.Logf("Wrote %d bytes in %s.", n, dur)
	if dur < 3600*time.Millisecond || dur > 4400*time.Millisecond {
		t.Errorf("Took %s, want 4s.", dur)
	}
}

func TestCopy(t *testing.T) {
	t.Parallel()

	tr := &BWTestReader{size: 1 << 20, stall: 2 * time.Second}

	start := time.Now()
	n, err := Copy(ioutil.Discard, tr, 500<<10)
	dur := time.Since(start)
	if err != nil {
		t.Error(err)
	}
	if n != 1<<20 {
		t.Errorf("Want %d bytes, got %d.", 1<<20, n)
	}
	t.Logf("Copied %d bytes in %s.", n, dur)
	if dur < 3600*time.Millisecond || dur > 4400*time.Millisecond {
		t.Errorf("Took %s, want 4s.", dur)
	}

}
