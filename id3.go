// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"bufio"
	"io/ioutil"
	"os"
)

type File struct {
	*Tag
	name                        string
	data                        []byte
	Title, Artist, Album, Genre string
}

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
		tag.textFrame("TIT2"),
		tag.textFrame("TPE1"),
		tag.textFrame("TALB"),
		tag.textFrame("TCON"),
	}, nil
}

func (f *File) Close() {
	f.setTextFrame("TIT2", f.Title)
	f.setTextFrame("TPE1", f.Artist)
	f.setTextFrame("TALB", f.Album)
	f.setTextFrame("TCON", f.Genre)

	fi, err := os.Create(f.name)
	defer fi.Close()
	if err != nil {
		panic(err)
	}

	wr := bufio.NewWriter(fi)
	wr.Write(f.Tag.Bytes())
	wr.Write(f.data)
}

func (t Tag) textFrame(id string) string {
	if frames, ok := t.Frames[id]; ok {
		frame := frames[0]
		switch frame.(type) {
		case (*TextFrame):
			return frame.(*TextFrame).Text()
		default:
		}
	}

	return ""
}

func (t *Tag) setTextFrame(id, text string) {
	if frames, ok := t.Frames[id]; ok {
		frame := frames[0]
		switch frame.(type) {
		case (*TextFrame):
			frame.(*TextFrame).SetText(text)
		default:
		}
	}
}
