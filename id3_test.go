// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"bytes"
	v2 "github.com/mikkyang/id3-go/v2"
	"io/ioutil"
	"testing"
)

const (
	testFile = "test.mp3"
)

func TestOpen(t *testing.T) {
	file, err := Open(testFile)
	if err != nil {
		t.Errorf("Open: unable to open file")
	}

	tag, ok := file.Tagger.(*v2.Tag)
	if !ok {
		t.Errorf("Open: incorrect tagger type")
	}

	if s := tag.Artist(); s != "Paloalto\x00" {
		t.Errorf("Open: incorrect artist, %v", s)
	}

	if s := tag.Title(); s != "Nice Life (Feat. Basick)" {
		t.Errorf("Open: incorrect title, %v", s)
	}

	if s := tag.Album(); s != "Chief Life" {
		t.Errorf("Open: incorrect album, %v", s)
	}
}

func TestClose(t *testing.T) {
	before, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf("test file error")
	}

	file, err := Open(testFile)
	if err != nil {
		t.Errorf("Close: unable to open file")
	}
	beforeCutoff := file.originalSize

	file.SetArtist("Paloalto")
	file.SetTitle("Test test test test test test")

	afterCutoff := file.Size()

	if err := file.Close(); err != nil {
		t.Errorf("Close: unable to close file")
	}

	after, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf("Close: unable to reopen file")
	}

	if !bytes.Equal(before[beforeCutoff:], after[afterCutoff:]) {
		t.Errorf("Close: nontag data lost on close")
	}

	if err := ioutil.WriteFile(testFile, before, 0666); err != nil {
		t.Errorf("Close: unable to write original contents to test file")
	}
}
