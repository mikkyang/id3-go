// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"fmt"
	"io"
)

const (
	HeaderSize = 10
)

// Tag represents an ID3v2 tag
type Tag struct {
	Header
	frames  map[string][]Framer
	padding uint
}

// Creates a new tag
func NewTag(reader io.Reader) *Tag {
	t := &Tag{NewHeader(reader), make(map[string][]Framer), 0}
	if t.Header == nil {
		return nil
	}

	var frame Framer
	size := t.Header.Size()
	for size > 0 {
		switch t.Header.Version() {
		case "2.3.0":
			frame = NewV3Frame(reader)
		default:
			frame = NewV3Frame(reader)
		}

		if frame == nil {
			break
		}

		id := frame.Id()
		t.frames[id] = append(t.frames[id], frame)

		size -= FrameHeaderSize + frame.Size()
	}

	t.padding = uint(size)
	nAdvance := int(t.padding - FrameHeaderSize)
	if n, err := io.ReadFull(reader, make([]byte, nAdvance)); n != nAdvance || err != nil {
		return nil
	}

	return t
}

// Size of the tag
// Recalculated as frames and padding can be changed
func (t Tag) Size() int {
	size := 0
	for _, v := range t.frames {
		for _, f := range v {
			size += FrameHeaderSize + f.Size()
		}
	}

	headerSize := t.Header.Size()
	if padding := headerSize - size; padding < 0 {
		t.padding = 0
		head := t.Header.(Head)
		head.size = int32(size)
		return size
	} else {
		t.padding = uint(padding)
		return headerSize
	}
}

func (t Tag) Bytes() []byte {
	data := make([]byte, t.Size())

	index := 0
	for _, v := range t.frames {
		for _, f := range v {
			size := FrameHeaderSize + f.Size()

			switch t.Header.Version() {
			case "2.3":
				copy(data[index:index+size], V3Bytes(f))
			default:
				copy(data[index:index+size], V3Bytes(f))
			}

			index += size
		}
	}

	return append(t.Header.Bytes(), data...)
}

// All frames
func (t Tag) Frames(id string) []Framer {
	if frames, ok := t.frames[id]; ok && frames != nil {
		return frames
	}

	return nil
}

// First frame
func (t Tag) Frame(id string) Framer {
	if frames := t.Frames(id); frames != nil {
		return frames[0]
	}

	return nil
}

// Delete and return all frames
func (t *Tag) DeleteFrames(id string) []Framer {
	frames := t.Frames(id)
	if frames == nil {
		return nil
	}

	delete(t.frames, id)

	return frames
}

// Add frame
func (t *Tag) AddFrame(frame Framer) {
	id := frame.Id()
	t.frames[id] = append(t.frames[id], frame)
}

// Header represents the useful information contained in the data
type Header interface {
	Version() string
	Size() int
	Bytes() []byte
}

func (t Tag) Title() string {
	return t.textFrameText("TIT2")
}

func (t Tag) Artist() string {
	return t.textFrameText("TPE1")
}

func (t Tag) Album() string {
	return t.textFrameText("TALB")
}

func (t Tag) Year() string {
	return t.textFrameText("TYER")
}

func (t Tag) Genre() string {
	return t.textFrameText("TCON")
}

func (t *Tag) SetTitle(text string) {
	t.setTextFrameText("TIT2", text)
}

func (t *Tag) SetArtist(text string) {
	t.setTextFrameText("TPE1", text)
}

func (t *Tag) SetAlbum(text string) {
	t.setTextFrameText("TALB", text)
}

func (t *Tag) SetGenre(text string) {
	t.setTextFrameText("TCON", text)
}

func (t *Tag) SetYear(text string) {
	t.setTextFrameText("TYER", text)
}

func (t *Tag) textFrame(id string) *TextFrame {
	if frame := t.Frame(id); frame != nil {
		switch frame.(type) {
		case (*TextFrame):
			return frame.(*TextFrame)
		default:
		}
	}

	return nil
}

func (t Tag) textFrameText(id string) string {
	if frame := t.textFrame(id); frame != nil {
		return frame.Text()
	}

	return ""
}

func (t *Tag) setTextFrameText(id, text string) {
	if frame := t.textFrame(id); frame != nil {
		frame.SetText(text)
	}
}

func NewHeader(reader io.Reader) Header {
	data := make([]byte, HeaderSize)
	n, err := io.ReadFull(reader, data)
	if n < HeaderSize || err != nil || string(data[:3]) != "ID3" {
		return nil
	}

	size, err := synchint(data[6:])
	if err != nil {
		return nil
	}

	return Head{
		version:  data[3],
		revision: data[4],
		flags:    data[5],
		size:     size,
	}
}

// Head represents the data of the header of the entire tag
type Head struct {
	version, revision byte
	flags             byte
	size              int32
}

func (h Head) Version() string {
	return fmt.Sprintf("2.%d.%d", h.version, h.revision)
}

func (h Head) Size() int {
	return int(h.size)
}

func (h Head) Bytes() []byte {
	data := make([]byte, HeaderSize)

	copy(data[:3], []byte("ID3"))
	copy(data[6:], synchbytes(h.size))
	data[3] = h.version
	data[4] = h.revision
	data[5] = h.flags

	return data
}
