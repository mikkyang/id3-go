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

type FrameType struct {
	id          string
	description string
	constructor func(FrameHead, []byte) Framer
}

type Framer interface {
	Id() string
	Size() int
	StatusFlags() byte
	FormatFlags() byte
	String() string
	Bytes() []byte
}

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

type TextFrame struct {
	FrameHead
	encoding string
	text     string
}

func NewTextFrame(head FrameHead, data []byte) Framer {
	var err error
	f := &TextFrame{FrameHead: head}

	encodingIndex := data[0]
	f.encoding = encodingForIndex(data[0])

	if f.text, err = Decoders[encodingIndex].ConvertString(string(data[1:])); err != nil {
		return nil
	}

	return f
}

func (f TextFrame) Encoding() string {
	return f.encoding
}

func (f *TextFrame) SetEncoding(encoding string) error {
	if indexForEncoding(encoding) < 0 {
		return errors.New("encoding: invalid encoding")
	}

	f.encoding = encoding
	return nil
}

func (f TextFrame) Text() string {
	return f.text
}

func (f *TextFrame) SetText(text string) error {
	encodingIndex := indexForEncoding(f.encoding)
	encodedString, err := Encoders[encodingIndex].ConvertString(text)
	if err != nil {
		return err
	}

	f.size += int32(len(encodedString)) - f.size
	f.text = text
	return nil
}

func (f TextFrame) String() string {
	return f.text
}

func (f TextFrame) Bytes() []byte {
	bytes := make([]byte, f.Size())

	encodingIndex := indexForEncoding(f.encoding)
	encodedString, err := Encoders[encodingIndex].ConvertString(f.text)
	if err != nil {
		return bytes
	}

	bytes[0] = encodingIndex
	copy(bytes[1:], []byte(encodedString))

	return bytes
}

type DescTextFrame struct {
	FrameHead
	TextFrame
	description string
}

func NewDescTextFrame(head FrameHead, data []byte) Framer {
	f := &DescTextFrame{FrameHead: head}

	var err error
	encodingIndex := data[0]

	f.encoding = encodingForIndex(encodingIndex)

	cutoff := 1
	if i := afterNullIndex(data[1:], f.encoding); i < 0 {
		return nil
	} else {
		cutoff += i
	}

	if f.description, err = Decoders[encodingIndex].ConvertString(string(data[1:cutoff])); err != nil {
		return nil
	}

	if f.text, err = Decoders[encodingIndex].ConvertString(string(data[cutoff:])); err != nil {
		return nil
	}

	return f
}

func (f DescTextFrame) Description() string {
	return f.description
}

func (f *DescTextFrame) SetDescription(description string) error {
	encodingIndex := indexForEncoding(f.encoding)
	encodedString, err := Encoders[encodingIndex].ConvertString(description)
	if err != nil {
		return err
	}

	f.size += int32(len(encodedString)) - f.size
	f.description = description
	return nil
}

func (f DescTextFrame) Bytes() []byte {
	bytes := make([]byte, f.Size())

	encodingIndex := indexForEncoding(f.encoding)
	encodedDescription, err := Encoders[encodingIndex].ConvertString(f.description)
	if err != nil {
		return bytes
	}
	encodedText, err := Encoders[encodingIndex].ConvertString(f.text)
	if err != nil {
		return bytes
	}

	bytes[0] = encodingIndex
	index := 1
	copy(bytes[index:index+len(encodedDescription)], []byte(encodedDescription))
	index += len(encodedDescription)
	copy(bytes[index:index+len(encodedText)], []byte(encodedText))

	return bytes
}

type UnsynchTextFrame struct {
	FrameHead
	DescTextFrame
	language          string
	contentDescriptor string
}

func NewUnsynchTextFrame(head FrameHead, data []byte) Framer {
	var err error
	f := &UnsynchTextFrame{FrameHead: head}

	encodingIndex := data[0]

	f.encoding = encodingForIndex(encodingIndex)
	f.language = string(data[1:4])

	cutoff := 4
	if i := afterNullIndex(data[4:], f.encoding); i < 0 {
		return nil
	} else {
		cutoff += i
	}

	if f.description, err = Decoders[encodingIndex].ConvertString(string(data[4:cutoff])); err != nil {
		return nil
	}

	if f.text, err = Decoders[encodingIndex].ConvertString(string(data[cutoff:])); err != nil {
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

func (f UnsynchTextFrame) ContentDescriptor() string {
	return f.contentDescriptor
}

func (f *UnsynchTextFrame) SetContentDescriptor(contentDescriptor string) error {
	encodingIndex := indexForEncoding(f.encoding)
	encodedString, err := Encoders[encodingIndex].ConvertString(contentDescriptor)
	if err != nil {
		return err
	}

	f.size += int32(len(encodedString)) - f.size
	f.contentDescriptor = contentDescriptor
	return nil
}

func (f UnsynchTextFrame) Bytes() []byte {
	bytes := make([]byte, f.Size())

	encodingIndex := indexForEncoding(f.encoding)
	encodedDescription, err := Encoders[encodingIndex].ConvertString(f.description)
	if err != nil {
		return bytes
	}
	encodedText, err := Encoders[encodingIndex].ConvertString(f.text)
	if err != nil {
		return bytes
	}

	bytes[0] = encodingIndex
	copy(bytes[1:4], []byte(f.language))
	index := 4
	copy(bytes[index:index+len(encodedDescription)], []byte(encodedDescription))
	index += len(encodedDescription)
	copy(bytes[index:index+len(encodedText)], []byte(encodedText))

	return bytes
}

type ImageFrame struct {
	FrameHead
	DataFrame
	encoding    string
	mimeType    string
	pictureType byte
	description string
}

func NewImageFrame(head FrameHead, data []byte) Framer {
	f := &ImageFrame{FrameHead: head}

	var err error
	encodingIndex := data[0]

	f.encoding = encodingForIndex(encodingIndex)

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
	return f.encoding
}

func (f *ImageFrame) SetEncoding(encoding string) error {
	if indexForEncoding(encoding) < 0 {
		return errors.New("encoding: invalid encoding")
	}

	f.encoding = encoding
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

func (f ImageFrame) Bytes() []byte {
	bytes := make([]byte, f.Size())

	encodingIndex := indexForEncoding(f.encoding)
	encodedDescription, err := Encoders[encodingIndex].ConvertString(f.description)
	if err != nil {
		return bytes
	}

	bytes[0] = encodingIndex
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
