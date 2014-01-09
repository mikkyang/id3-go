// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"errors"
	v1 "github.com/mikkyang/id3-go/v1"
	v2 "github.com/mikkyang/id3-go/v2"
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
	AllFrames() []v2.Framer
	Frames(string) []v2.Framer
	Frame(string) v2.Framer
	DeleteFrames(string) []v2.Framer
	AddFrame(v2.Framer)
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

const (
	fileEndFlag = 2
)

// Opens a new tagged file
func Open(name string) (*File, error) {
	fi, err := os.OpenFile(name, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	file := &File{file: fi}

	if v2Tag := v2.ParseTag(fi); v2Tag != nil {
		file.Tagger = v2Tag
		file.originalSize = v2Tag.Size()
	} else if v1Tag := v1.ParseTag(fi); v1Tag != nil {
		file.Tagger = v1Tag
	} else {
		return nil, errors.New("Open: unknown tag format")
	}

	return file, nil
}

// Saves any edits to the tagged file
func (f *File) Close() error {
	defer f.file.Close()

	switch f.Tagger.(type) {
	case (*v1.Tag):
		f.file.Seek(-v1.TagSize, fileEndFlag)
		f.file.Write(f.Tagger.Bytes())
	case (*v2.Tag):
		if f.Size() > f.originalSize {
			stat, err := f.file.Stat()
			if err != nil {
				return err
			}

			start := f.originalSize + v2.HeaderSize
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
	default:
		return errors.New("Close: unknown tag version")
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
