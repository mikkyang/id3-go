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
	name string
	data []byte
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

	return &File{tag, name, data}, nil
}

func (f *File) Close() {
	fi, err := os.Create(f.name)
	defer fi.Close()
	if err != nil {
		panic(err)
	}

	wr := bufio.NewWriter(fi)
	wr.Write(f.Tag.Bytes())
	wr.Write(f.data)
}
