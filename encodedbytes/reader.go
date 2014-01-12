// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package encodedbytes

import (
	"io"
)

// Reader is a helper for Reading frame bytes
type Reader struct {
	data  []byte
	index int // current Reading index
}

func (r *Reader) Read(b []byte) (n int, err error) {
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

func (r *Reader) ReadByte() (b byte, err error) {
	if r.index >= len(r.data) {
		return 0, io.EOF
	}
	b = r.data[r.index]
	r.index++
	return
}

func (r *Reader) ReadNumBytes(n int) ([]byte, error) {
	if n <= 0 {
		return []byte{}, nil
	}
	if r.index+n > len(r.data) {
		return []byte{}, io.EOF
	}

	b := make([]byte, n)
	_, err := r.Read(b)

	return b, err
}

func (r *Reader) ReadNumBytesString(n int) (string, error) {
	b, err := r.ReadNumBytes(n)
	return string(b), err
}

func (r *Reader) ReadRest() ([]byte, error) {
	return r.ReadNumBytes(len(r.data) - r.index)
}

func (r *Reader) ReadRestString(encoding byte) (string, error) {
	b, err := r.ReadRest()
	if err != nil {
		return "", err
	}

	return Decoders[encoding].ConvertString(string(b))
}

func (r *Reader) ReadNullTermString(encoding byte) (string, error) {
	b, err := r.ReadNumBytes(afterNullIndex(r.data[r.index:], encoding))
	if err != nil {
		return "", err
	}

	return Decoders[encoding].ConvertString(string(b))
}

func NewReader(b []byte) *Reader { return &Reader{b, 0} }
