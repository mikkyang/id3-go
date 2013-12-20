// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"errors"
	iconv "github.com/djimenez/iconv-go"
)

const (
	BytesPerInt     = 4
	SynchByteLength = 7
	NormByteLength  = 8
	NativeEncoding  = "UTF-8"
)

var (
	EncodingMap = [...]string{
		"ISO-8859-1",
		"UTF-16",
		"UTF-16BE",
		"ISO-8859-1",
	}
	Decoders = make([]*iconv.Converter, len(EncodingMap))
	Encoders = make([]*iconv.Converter, len(EncodingMap))
)

func init() {
	for i, e := range EncodingMap {
		Decoders[i], _ = iconv.NewConverter(e, NativeEncoding)
		Encoders[i], _ = iconv.NewConverter(NativeEncoding, e)
	}
}

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

func encodingForIndex(b byte) string {
	encodingIndex := int(b)
	if encodingIndex < 0 || encodingIndex > len(EncodingMap) {
		encodingIndex = 0
	}

	return EncodingMap[encodingIndex]
}

func indexForEncoding(e string) byte {
	for i, v := range EncodingMap {
		if v == e {
			return byte(i)
		}
	}

	return 0
}
