package midec

/*
Based on image/format.go in standard library.

---

Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/

import (
	"bufio"
	"errors"
	"io"
	"sync"
	"sync/atomic"
)

// ErrFormat indicates that detecting encountered an unknown format.
var ErrFormat = errors.New("midec: unknown format")

type format struct {
	name, magic string
	isAnimated  func(io.Reader) (bool, error)
}

var (
	formatsMu     sync.Mutex
	atomicFormats atomic.Value
)

// RegisterFormat registers an image format for use by IsAnimated.
func RegisterFormat(name, magic string, isAnimated func(io.Reader) (bool, error)) {
	formatsMu.Lock()
	formats, _ := atomicFormats.Load().([]format)
	atomicFormats.Store(append(formats, format{name, magic, isAnimated}))
	formatsMu.Unlock()
}

// A reader is an io.Reader that can also peek ahead.
type reader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

// asReader converts an io.Reader to a reader.
func asReader(r io.Reader) reader {
	if rr, ok := r.(reader); ok {
		return rr
	}
	return bufio.NewReader(r)
}

// Match reports whether magic matches b. Magic may contain "?" wildcards.
func match(magic string, b []byte) bool {
	if len(magic) != len(b) {
		return false
	}
	for i, c := range b {
		if magic[i] != c && magic[i] != '?' {
			return false
		}
	}
	return true
}

// Sniff determines the format of r's data.
func sniff(r reader) format {
	formats, _ := atomicFormats.Load().([]format)
	for _, f := range formats {
		b, err := r.Peek(len(f.magic))
		if err == nil && match(f.magic, b) {
			return f
		}
	}
	return format{}
}

// IsAnimated detects whether it is an animated image that has been encoded in a registered format.
func IsAnimated(r io.Reader) (bool, error) {
	rr := asReader(r)
	f := sniff(rr)
	if f.isAnimated == nil {
		return false, ErrFormat
	}
	m, err := f.isAnimated(rr)
	return m, err
}
