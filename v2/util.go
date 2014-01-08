// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package v2

import (
	"bytes"
	"errors"
	iconv "github.com/djimenez/iconv-go"
)

const (
	BytesPerInt     = 4
	SynchByteLength = 7
	NormByteLength  = 8
	NativeEncoding  = 3
	UTF16NullLength = 2
)

var (
	EncodingMap = [...]string{
		"ISO-8859-1",
		"UTF-16",
		"UTF-16BE",
		"UTF-8",
	}
	Decoders = make([]*iconv.Converter, len(EncodingMap))
	Encoders = make([]*iconv.Converter, len(EncodingMap))
)

func init() {
	n := encodingForIndex(NativeEncoding)
	for i, e := range EncodingMap {
		Decoders[i], _ = iconv.NewConverter(e, n)
		Encoders[i], _ = iconv.NewConverter(n, e)
	}
}

func byteint(buf []byte, base uint) (i int32, err error) {
	if len(buf) > BytesPerInt {
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

func intbytes(n int32, base uint) []byte {
	mask := int32(1<<base - 1)
	bytes := make([]byte, BytesPerInt)

	for i, _ := range bytes {
		bytes[len(bytes)-i-1] = byte(n & mask)
		n >>= base
	}

	return bytes
}

func synchbytes(n int32) []byte {
	return intbytes(n, SynchByteLength)
}

func normbytes(n int32) []byte {
	return intbytes(n, NormByteLength)
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

func afterNullIndex(data []byte, encoding byte) int {
	encodingString := encodingForIndex(encoding)

	if encodingString == "UTF-16" || encodingString == "UTF-16BE" {
		limit, byteCount := len(data), UTF16NullLength
		null := bytes.Repeat([]byte{0x0}, byteCount)

		for i, _ := range data[:limit/byteCount] {
			atIndex := byteCount * i
			afterIndex := atIndex + byteCount

			if bytes.Equal(data[atIndex:afterIndex], null) {
				return afterIndex
			}
		}
	} else {
		if index := bytes.IndexByte(data, 0x00); index >= 0 {
			return index + 1
		}
	}

	return -1
}

func encodedDiff(ea byte, a string, eb byte, b string) (int, error) {
	encodedStringA, err := Encoders[ea].ConvertString(a)
	if err != nil {
		return 0, err
	}

	encodedStringB, err := Encoders[eb].ConvertString(b)
	if err != nil {
		return 0, err
	}

	return len(encodedStringA) - len(encodedStringB), nil
}
