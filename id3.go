// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"bufio"
	"io/ioutil"
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
	name         string
	data         []byte
}

// Opens a new tagged file
func Open(name string) (*File, error) {
	fi, err := os.Open(name)
	defer fi.Close()
	if err != nil {
		return nil, err
	}

	rd := bufio.NewReader(fi)
	tag := NewTag(rd)
	data, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	return &File{
		tag,
		tag.Size(),
		name,
		data,
	}, nil
}

// Saves any edits to the tagged file
func (f *File) Close() {
	fi, err := os.OpenFile(f.name, os.O_RDWR, 0666)
	defer fi.Close()
	if err != nil {
		panic(err)
	}

	wr := bufio.NewWriter(fi)
	wr.Write(f.Tagger.Bytes())

	if f.Size() > f.originalSize {
		wr.Write(f.data)
	}
}
