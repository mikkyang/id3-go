// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"fmt"
)

type Header interface {
	Version() string
	Size() int
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
