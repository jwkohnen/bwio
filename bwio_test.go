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
	"bytes"
	"errors"
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

func TestIllegalLimit(t *testing.T) {
	t.Parallel()

	testt := []struct {
		name  string
		size  int
		panic bool
	}{
		{"negative", -1, true},
		{"zero", 0, true},
		{"one", 1, false},
	}
	for _, testc := range testt {
		t.Run(testc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if testc.panic && r == nil {
					t.Error("Should have paniced, but didn't.")
				}
				if !testc.panic && r != nil {
					t.Errorf("Shouldn't have paniced, but did: %v", r)
				}
			}()
			r := bytes.NewReader([]byte{0x00})
			lr := NewReader(r, testc.size)
			_, _ = io.Copy(ioutil.Discard, lr)
		})
	}
}

var poison = errors.New("poison")

type pReader struct{}
type pWriter struct{}

func (*pReader) Read(_ []byte) (int, error) {
	return 0, poison
}

func (*pWriter) Write(_ []byte) (int, error) {
	return 0, poison
}

func TestError(t *testing.T) {
	t.Parallel()

	t.Run("read", func(t *testing.T) {
		t.Parallel()
		lr := NewReader(new(pReader), 1)
		_, err := io.Copy(ioutil.Discard, lr)
		if err != poison {
			t.Errorf("Want %v, got %v", poison, err)
		}
	})
	t.Run("write", func(t *testing.T) {
		t.Parallel()
		lw := NewWriter(new(pWriter), 1)
		r := bytes.NewReader([]byte{0x00})
		_, err := io.Copy(lw, r)
		if err != poison {
			t.Errorf("Want %v, got %v", poison, err)
		}
	})
	t.Run("copyR", func(t *testing.T) {
		t.Parallel()
		_, err := Copy(ioutil.Discard, new(pReader), 1)
		if err != poison {
			t.Errorf("Want %v, got %v", poison, err)
		}
	})
	t.Run("copyW", func(t *testing.T) {
		t.Parallel()
		r := bytes.NewReader([]byte{0x00})
		_, err := Copy(new(pWriter), r, 1)
		if err != poison {
			t.Errorf("Want %v, got %v", poison, err)
		}
	})
	t.Run("copyRW", func(t *testing.T) {
		t.Parallel()
		_, err := Copy(new(pWriter), new(pReader), 1)
		if err != poison {
			t.Errorf("Want %v, got %v", poison, err)
		}
	})
}
