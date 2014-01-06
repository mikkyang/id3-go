// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"io"
	"os"
)

// Tagger represents the metadata of a tag
type Tagger interface {
	Title() string
	Artist() string
	Album() string
	Year() string
	Genre() string
	Comments() []string
	SetTitle(string)
	SetArtist(string)
	SetAlbum(string)
	SetYear(string)
	SetGenre(string)
	AllFrames() []Framer
	Frames(string) []Framer
	Frame(string) Framer
	DeleteFrames(string) []Framer
	AddFrame(Framer)
	Bytes() []byte
	Padding() uint
	Size() int
	Version() string
}

// File represents the tagged file
type File struct {
	Tagger
	originalSize int
	file         *os.File
}

// Opens a new tagged file
func Open(name string) (*File, error) {
	fi, err := os.OpenFile(os.Args[1], os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	tag := ParseTag(fi)

	return &File{tag, tag.Size(), fi}, nil
}

// Saves any edits to the tagged file
func (f *File) Close() error {
	defer f.file.Close()

	if f.Size() > f.originalSize {
		stat, err := f.file.Stat()
		if err != nil {
			return err
		}

		start := f.originalSize + HeaderSize
		end := stat.Size()
		offset := f.Tagger.Size() - f.originalSize

		err = f.shiftBytesBack(int64(start), end, int64(offset))
		if err != nil {
			return err
		}
	}

	if _, err := f.file.WriteAt(f.Tagger.Bytes(), 0); err != nil {
		return err
	}

	return nil
}

func (f *File) shiftBytesBack(start, end, offset int64) error {
	wrBuf := make([]byte, offset)
	rdBuf := make([]byte, offset)

	wrOffset := offset
	rdOffset := start

	rn, err := f.file.ReadAt(wrBuf, rdOffset)
	if err != nil && err != io.EOF {
		panic(err)
	}
	rdOffset += int64(rn)

	for {
		if rdOffset >= end {
			break
		}

		n, err := f.file.ReadAt(rdBuf, rdOffset)
		if err != nil && err != io.EOF {
			return err
		}

		if rdOffset+int64(n) > end {
			n = int(end - rdOffset)
		}

		if _, err := f.file.WriteAt(wrBuf[:rn], wrOffset); err != nil {
			return err
		}

		rdOffset += int64(n)
		wrOffset += int64(rn)
		copy(wrBuf, rdBuf)
		rn = n
	}

	if _, err := f.file.WriteAt(wrBuf[:rn], wrOffset); err != nil {
		return err
	}

	return nil
}
