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

type Header interface {
	Version() string
	Size() int
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
