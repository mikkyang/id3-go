// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"errors"
	"fmt"
)

const (
	FrameHeaderSize = 10
)

// FrameType holds frame id metadata and constructor method
// A set number of these are created in the version specific files
type FrameType struct {
	id          string
	description string
	constructor func(FrameHead, []byte) Framer
}

// Framer provides a generic interface for frames
// This is the default type returned when creating frames
type Framer interface {
	Id() string
	Size() int
	StatusFlags() byte
	FormatFlags() byte
	String() string
	Bytes() []byte
}

// FrameHead represents the header of each frame
// Additional metadata is kept through the embedded frame type
// These do not usually need to be manually created
type FrameHead struct {
	FrameType
	statusFlags byte
	formatFlags byte
	size        int32
}

func (h FrameHead) Id() string {
	return h.id
}

func (h FrameHead) Size() int {
	return int(h.size)
}

func (h FrameHead) StatusFlags() byte {
	return h.statusFlags
}

func (h FrameHead) FormatFlags() byte {
	return h.formatFlags
}

// DataFrame is the default frame for binary data
type DataFrame struct {
	FrameHead
	data []byte
}

func NewDataFrame(head FrameHead, data []byte) Framer {
	return &DataFrame{head, data}
}

func (f DataFrame) Data() []byte {
	return f.data
}

func (f *DataFrame) SetData(b []byte) {
	f.size += int32(len(b)) - f.size
	f.data = b
}

func (f DataFrame) String() string {
	return "<binary data>"
}

func (f DataFrame) Bytes() []byte {
	return f.data
}

// TextFramer represents frames that contain encoded text
type TextFramer interface {
	Framer
	Encoding() string
	SetEncoding(string) error
	Text() string
	SetText(string) error
}

// TextFrame represents frames that contain encoded text
type TextFrame struct {
	FrameHead
	encoding byte
	text     string
}

func NewTextFrame(head FrameHead, data []byte) Framer {
	var err error
	f := &TextFrame{FrameHead: head}
	rd := newReader(data)

	if f.encoding, err = rd.readByte(); err != nil {
		return nil
	}

	if f.text, err = rd.readRestString(f.encoding); err != nil {
		return nil
	}

	return f
}

func (f TextFrame) Encoding() string {
	return encodingForIndex(f.encoding)
}

func (f *TextFrame) SetEncoding(encoding string) error {
	i := byte(indexForEncoding(encoding))
	if i < 0 {
		return errors.New("encoding: invalid encoding")
	}

	f.encoding = i
	return nil
}

func (f TextFrame) Text() string {
	return f.text
}

func (f *TextFrame) SetText(text string) error {
	diff, err := encodedDiff(f.encoding, text, f.text)
	if err != nil {
		return err
	}

	f.size += int32(diff)
	f.text = text
	return nil
}

func (f TextFrame) String() string {
	return f.text
}

func (f TextFrame) Bytes() []byte {
	var err error
	bytes := make([]byte, f.Size())
	wr := newWriter(bytes)

	if err = wr.writeByte(f.encoding); err != nil {
		return bytes
	}

	if err = wr.writeString(f.text, f.encoding); err != nil {
		return bytes
	}

	return bytes
}

type DescTextFrame struct {
	FrameHead
	TextFrame
	description string
}

// DescTextFrame represents frames that contain encoded text and descriptions
func NewDescTextFrame(head FrameHead, data []byte) Framer {
	var err error
	f := &DescTextFrame{FrameHead: head}
	rd := newReader(data)

	if f.encoding, err = rd.readByte(); err != nil {
		return nil
	}

	if f.description, err = rd.readNullTermString(f.encoding); err != nil {
		return nil
	}

	if f.text, err = rd.readRestString(f.encoding); err != nil {
		return nil
	}

	return f
}

func (f DescTextFrame) Description() string {
	return f.description
}

func (f *DescTextFrame) SetDescription(description string) error {
	diff, err := encodedDiff(f.encoding, description, f.description)
	if err != nil {
		return err
	}

	f.size += int32(diff)
	f.description = description
	return nil
}

func (f DescTextFrame) String() string {
	return fmt.Sprintf("%s: %s", f.description, f.text)
}

func (f DescTextFrame) Bytes() []byte {
	var err error
	bytes := make([]byte, f.Size())
	wr := newWriter(bytes)

	if err = wr.writeByte(f.encoding); err != nil {
		return bytes
	}

	if err = wr.writeString(f.description, f.encoding); err != nil {
		return bytes
	}

	if err = wr.writeString(f.text, f.encoding); err != nil {
		return bytes
	}

	return bytes
}

// UnsynchTextFrame represents frames that contain unsynchronized text
type UnsynchTextFrame struct {
	FrameHead
	DescTextFrame
	language string
}

func NewUnsynchTextFrame(head FrameHead, data []byte) Framer {
	var err error
	f := &UnsynchTextFrame{FrameHead: head}
	rd := newReader(data)

	if f.encoding, err = rd.readByte(); err != nil {
		return nil
	}

	if f.language, err = rd.readNumBytesString(3); err != nil {
		return nil
	}

	if f.description, err = rd.readNullTermString(f.encoding); err != nil {
		return nil
	}

	if f.text, err = rd.readRestString(f.encoding); err != nil {
		return nil
	}

	return f
}

func (f UnsynchTextFrame) Language() string {
	return f.language
}

func (f *UnsynchTextFrame) SetLanguage(language string) error {
	if len(language) != 3 {
		return errors.New("language: invalid language string")
	}

	f.language = language
	return nil
}

func (f UnsynchTextFrame) String() string {
	return fmt.Sprintf("%s\t%s:\n%s", f.language, f.description, f.text)
}

func (f UnsynchTextFrame) Bytes() []byte {
	var err error
	bytes := make([]byte, f.Size())
	wr := newWriter(bytes)

	if err = wr.writeByte(f.encoding); err != nil {
		return bytes
	}

	if err = wr.writeString(f.language, NativeEncoding); err != nil {
		return bytes
	}

	if err = wr.writeString(f.description, f.encoding); err != nil {
		return bytes
	}

	if err = wr.writeString(f.text, f.encoding); err != nil {
		return bytes
	}

	return bytes
}

// ImageFrame represent frames that have media attached
type ImageFrame struct {
	FrameHead
	DataFrame
	encoding    byte
	mimeType    string
	pictureType byte
	description string
}

func NewImageFrame(head FrameHead, data []byte) Framer {
	var err error
	f := &ImageFrame{FrameHead: head}
	rd := newReader(data)

	if f.encoding, err = rd.readByte(); err != nil {
		return nil
	}

	if f.mimeType, err = rd.readNullTermString(NativeEncoding); err != nil {
		return nil
	}

	if f.pictureType, err = rd.readByte(); err != nil {
		return nil
	}

	if f.description, err = rd.readNullTermString(f.encoding); err != nil {
		return nil
	}

	if f.data, err = rd.readRest(); err != nil {
		return nil
	}

	return f
}

func (f ImageFrame) Encoding() string {
	return encodingForIndex(f.encoding)
}

func (f *ImageFrame) SetEncoding(encoding string) error {
	i := byte(indexForEncoding(encoding))
	if i < 0 {
		return errors.New("encoding: invalid encoding")
	}

	f.encoding = i
	return nil
}

func (f ImageFrame) MIMEType() string {
	return f.mimeType
}

func (f *ImageFrame) SetMIMEType(mimeType string) {
	f.size += int32(len(mimeType)) - f.size
	if mimeType[len(mimeType)-1] != 0 {
		nullTermBytes := append([]byte(mimeType), 0x00)
		f.mimeType = string(nullTermBytes)
		f.size += 1
	} else {
		f.mimeType = mimeType
	}
}

func (f ImageFrame) String() string {
	return fmt.Sprintf("%s\t%s: <binary data>", f.mimeType, f.description)
}

func (f ImageFrame) Bytes() []byte {
	var err error
	bytes := make([]byte, f.Size())
	wr := newWriter(bytes)

	if err = wr.writeByte(f.encoding); err != nil {
		return bytes
	}

	if err = wr.writeString(f.mimeType, NativeEncoding); err != nil {
		return bytes
	}

	if err = wr.writeByte(f.pictureType); err != nil {
		return bytes
	}

	if err = wr.writeString(f.description, f.encoding); err != nil {
		return bytes
	}

	if n, err := wr.write(f.data); n < len(f.data) || err != nil {
		return bytes
	}

	return bytes
}
