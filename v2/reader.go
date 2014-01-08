// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package v2

import (
	"io"
)

// reader is a helper for reading frame bytes
type reader struct {
	data  []byte
	index int // current reading index
}

func (r *reader) read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	if r.index >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(b, r.data[r.index:])
	r.index += n
	return
}

func (r *reader) readByte() (b byte, err error) {
	if r.index >= len(r.data) {
		return 0, io.EOF
	}
	b = r.data[r.index]
	r.index++
	return
}

func (r *reader) readNumBytes(n int) ([]byte, error) {
	if n <= 0 {
		return []byte{}, nil
	}
	if r.index+n > len(r.data) {
		return []byte{}, io.EOF
	}

	b := make([]byte, n)
	_, err := r.read(b)

	return b, err
}

func (r *reader) readNumBytesString(n int) (string, error) {
	b, err := r.readNumBytes(n)
	return string(b), err
}

func (r *reader) readRest() ([]byte, error) {
	return r.readNumBytes(len(r.data) - r.index)
}

func (r *reader) readRestString(encoding byte) (string, error) {
	b, err := r.readRest()
	if err != nil {
		return "", err
	}

	return Decoders[encoding].ConvertString(string(b))
}

func (r *reader) readNullTermString(encoding byte) (string, error) {
	b, err := r.readNumBytes(afterNullIndex(r.data[r.index:], encoding))
	if err != nil {
		return "", err
	}

	return Decoders[encoding].ConvertString(string(b))
}

func newReader(b []byte) *reader { return &reader{b, 0} }
