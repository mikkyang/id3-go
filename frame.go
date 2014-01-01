// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"bytes"
	"errors"
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

	f.encoding = data[0]

	if f.text, err = Decoders[f.encoding].ConvertString(string(data[1:])); err != nil {
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
	bytes := make([]byte, f.Size())

	encodedString, err := Encoders[f.encoding].ConvertString(f.text)
	if err != nil {
		return bytes
	}

	bytes[0] = f.encoding
	copy(bytes[1:], []byte(encodedString))

	return bytes
}

type DescTextFrame struct {
	FrameHead
	TextFrame
	description string
}

// DescTextFrame represents frames that contain encoded text and descriptions
func NewDescTextFrame(head FrameHead, data []byte) Framer {
	f := &DescTextFrame{FrameHead: head}

	var err error

	f.encoding = data[0]

	cutoff := 1
	if i := afterNullIndex(data[1:], f.encoding); i < 0 {
		return nil
	} else {
		cutoff += i
	}

	if f.description, err = Decoders[f.encoding].ConvertString(string(data[1:cutoff])); err != nil {
		return nil
	}

	if f.text, err = Decoders[f.encoding].ConvertString(string(data[cutoff:])); err != nil {
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
	bytes := make([]byte, f.Size())

	encodedDescription, err := Encoders[f.encoding].ConvertString(f.description)
	if err != nil {
		return bytes
	}
	encodedText, err := Encoders[f.encoding].ConvertString(f.text)
	if err != nil {
		return bytes
	}

	bytes[0] = f.encoding
	index := 1
	copy(bytes[index:index+len(encodedDescription)], []byte(encodedDescription))
	index += len(encodedDescription)
	copy(bytes[index:index+len(encodedText)], []byte(encodedText))

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

	f.encoding = data[0]
	f.language = string(data[1:4])

	cutoff := 4
	if i := afterNullIndex(data[4:], f.encoding); i < 0 {
		return nil
	} else {
		cutoff += i
	}

	if f.description, err = Decoders[f.encoding].ConvertString(string(data[4:cutoff])); err != nil {
		return nil
	}

	if f.text, err = Decoders[f.encoding].ConvertString(string(data[cutoff:])); err != nil {
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
	bytes := make([]byte, f.Size())

	encodedDescription, err := Encoders[f.encoding].ConvertString(f.description)
	if err != nil {
		return bytes
	}
	encodedText, err := Encoders[f.encoding].ConvertString(f.text)
	if err != nil {
		return bytes
	}

	bytes[0] = f.encoding
	copy(bytes[1:4], []byte(f.language))
	index := 4
	copy(bytes[index:index+len(encodedDescription)], []byte(encodedDescription))
	index += len(encodedDescription)
	copy(bytes[index:index+len(encodedText)], []byte(encodedText))

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
	f := &ImageFrame{FrameHead: head}

	var err error
	encodingIndex := data[0]

	f.encoding = encodingIndex

	buffer := bytes.NewBuffer(data[1:])
	if f.mimeType, err = buffer.ReadString(0); err != nil {
		return nil
	}

	if f.pictureType, err = buffer.ReadByte(); err != nil {
		return nil
	}

	beginIndex := 1 + len(f.mimeType) + 1
	var cutoff int
	if i := afterNullIndex(data[beginIndex:], f.encoding); i < 0 {
		return nil
	} else {
		cutoff = beginIndex + i
	}

	if f.description, err = Decoders[encodingIndex].ConvertString(string(data[beginIndex:cutoff])); err != nil {
		return nil
	}

	f.data = data[cutoff:]

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
	bytes := make([]byte, f.Size())

	encodedDescription, err := Encoders[f.encoding].ConvertString(f.description)
	if err != nil {
		return bytes
	}

	bytes[0] = f.encoding
	index := 1
	copy(bytes[index:index+len(f.mimeType)], []byte(f.mimeType))
	index += len(f.mimeType)
	bytes[index] = f.pictureType
	index += 1
	copy(bytes[index:index+len(encodedDescription)], []byte(encodedDescription))
	index += len(encodedDescription)
	copy(bytes[index:index+len(f.data)], f.data)

	return bytes
}
