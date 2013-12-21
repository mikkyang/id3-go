// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"fmt"
	"io"
)

const (
	HeaderSize = 10
)

type Tag struct {
	Header
	Frames map[string][]Framer
}

func NewTag(reader io.Reader) *Tag {
	t := &Tag{NewHeader(reader), make(map[string][]Framer)}
	if t.Header == nil {
		return nil
	}

	var frame Framer
	for size := t.Header.Size(); size > 0; {
		switch t.Header.Version() {
		case "2.3":
			frame = NewV3Frame(reader)
		default:
			frame = NewV3Frame(reader)
		}

		if frame == nil {
			break
		}

		id := frame.Id()
		t.Frames[id] = append(t.Frames[id], frame)

		size -= frame.Size()
	}

	return t
}

func (t Tag) Bytes() []byte {
	data := make([]byte, t.Size())

	index := 0
	for _, v := range t.Frames {
		for _, f := range v {
			size := FrameHeaderSize + f.Size()

			switch t.Header.Version() {
			case "2.3":
				copy(data[index:index+size], V3Bytes(f))
			default:
				copy(data[index:index+size], V3Bytes(f))
			}

			index += size
		}
	}

	return append(t.Header.Bytes(), data...)
}

type Header interface {
	Version() string
	Size() int
	Bytes() []byte
}

func NewHeader(reader io.Reader) Header {
	data := make([]byte, HeaderSize)
	n, err := io.ReadFull(reader, data)
	if n < HeaderSize || err != nil || string(data[:3]) != "ID3" {
		return nil
	}

	size, err := synchint(data[6:])
	if err != nil {
		return nil
	}

	return &Head{
		version:  data[3],
		revision: data[4],
		flags:    data[5],
		size:     size,
	}
}

type Head struct {
	version, revision byte
	flags             byte
	size              int32
}

func (h Head) Version() string {
	return fmt.Sprintf("%d.%d", h.version, h.revision)
}

func (h Head) Size() int {
	return int(h.size)
}

func (h Head) Bytes() []byte {
	data := make([]byte, HeaderSize)

	copy(data[:3], []byte("ID3"))
	copy(data[6:], synchbytes(h.size))
	data[3] = h.version
	data[4] = h.revision
	data[5] = h.flags

	return data
}
