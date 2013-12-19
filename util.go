// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"errors"
)

const (
	BytesPerInt     = 4
	SynchByteLength = 7
	NormByteLength  = 8
)

func byteint(buf []byte, base uint) (i int32, err error) {
	if len(buf) != BytesPerInt {
		err = errors.New("byte integer: invalid []byte length")
		return
	}

	for _, b := range buf {
		if base < NormByteLength && b >= (1<<base) {
			err = errors.New("byte integer: exceed max bit")
			return
		}

		i = (i << base) | int32(b)
	}

	return
}

func synchint(buf []byte) (i int32, err error) {
	i, err = byteint(buf, SynchByteLength)
	return
}

func normint(buf []byte) (i int32, err error) {
	i, err = byteint(buf, NormByteLength)
	return
}
