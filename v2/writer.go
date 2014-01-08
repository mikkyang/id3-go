// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package v2

import (
	"io"
)

// writer is a helper for writing frame bytes
type writer struct {
	data  []byte
	index int // current writing index
}

func (w *writer) write(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	if w.index >= len(w.data) {
		return 0, io.EOF
	}
	n = copy(w.data[w.index:], b)
	w.index += n
	return
}

func (w *writer) writeByte(b byte) (err error) {
	if w.index >= len(w.data) {
		return io.EOF
	}
	w.data[w.index] = b
	w.index++
	return
}

func (w *writer) writeString(s string, encoding byte) (err error) {
	encodedString, err := Encoders[encoding].ConvertString(s)
	if err != nil {
		return err
	}

	_, err = w.write([]byte(encodedString))
	if err != nil {
		return err
	}

	return
}

func newWriter(b []byte) *writer { return &writer{b, 0} }
