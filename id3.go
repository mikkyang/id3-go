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
	SetTitle(string)
	SetArtist(string)
	SetAlbum(string)
	SetYear(string)
	SetGenre(string)
	Frame(string) Framer
	Frames(string) []Framer
	DeleteFrames(string) []Framer
	AddFrame(Framer)
	Bytes() []byte
	Version() string
}

// File represents the tagged file
type File struct {
	Tagger
	name string
	data []byte
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
		name,
		data,
	}, nil
}

// Saves any edits to the tagged file
func (f *File) Close() {
	fi, err := os.Create(f.name)
	defer fi.Close()
	if err != nil {
		panic(err)
	}

	wr := bufio.NewWriter(fi)
	wr.Write(f.Tagger.Bytes())
	wr.Write(f.data)
}
