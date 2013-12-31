// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"io"
)

const (
	V2FrameHeaderSize = 6
)

var (
	// V2FrameTypeMap specifies the frame IDs and constructors allowed in ID3v2.2
	V2FrameTypeMap = map[string]FrameType{
		"BUF": FrameType{id: "BUF", description: "Recommended buffer size", constructor: NewDataFrame},
		"CNT": FrameType{id: "CNT", description: "Play counter", constructor: NewDataFrame},
		"COM": FrameType{id: "COM", description: "Comments", constructor: NewDataFrame},
		"CRA": FrameType{id: "CRA", description: "Audio encryption", constructor: NewDataFrame},
		"CRM": FrameType{id: "CRM", description: "Encrypted meta frame", constructor: NewDataFrame},
		"ETC": FrameType{id: "ETC", description: "Event timing codes", constructor: NewDataFrame},
		"EQU": FrameType{id: "EQU", description: "Equalization", constructor: NewDataFrame},
		"GEO": FrameType{id: "GEO", description: "General encapsulated object", constructor: NewDataFrame},
		"IPL": FrameType{id: "IPL", description: "Involved people list", constructor: NewDataFrame},
		"LNK": FrameType{id: "LNK", description: "Linked information", constructor: NewDataFrame},
		"MCI": FrameType{id: "MCI", description: "Music CD Identifier", constructor: NewDataFrame},
		"MLL": FrameType{id: "MLL", description: "MPEG location lookup table", constructor: NewDataFrame},
		"PIC": FrameType{id: "PIC", description: "Attached picture", constructor: NewDataFrame},
		"POP": FrameType{id: "POP", description: "Popularimeter", constructor: NewDataFrame},
		"REV": FrameType{id: "REV", description: "Reverb", constructor: NewDataFrame},
		"RVA": FrameType{id: "RVA", description: "Relative volume adjustment", constructor: NewDataFrame},
		"SLT": FrameType{id: "SLT", description: "Synchronized lyric/text", constructor: NewDataFrame},
		"STC": FrameType{id: "STC", description: "Synced tempo codes", constructor: NewDataFrame},
		"TAL": FrameType{id: "TAL", description: "Album/Movie/Show title", constructor: NewTextFrame},
		"TBP": FrameType{id: "TBP", description: "BPM (Beats Per Minute)", constructor: NewTextFrame},
		"TCM": FrameType{id: "TCM", description: "Composer", constructor: NewTextFrame},
		"TCO": FrameType{id: "TCO", description: "Content type", constructor: NewTextFrame},
		"TCR": FrameType{id: "TCR", description: "Copyright message", constructor: NewTextFrame},
		"TDA": FrameType{id: "TDA", description: "Date", constructor: NewTextFrame},
		"TDY": FrameType{id: "TDY", description: "Playlist delay", constructor: NewTextFrame},
		"TEN": FrameType{id: "TEN", description: "Encoded by", constructor: NewTextFrame},
		"TFT": FrameType{id: "TFT", description: "File type", constructor: NewTextFrame},
		"TIM": FrameType{id: "TIM", description: "Time", constructor: NewTextFrame},
		"TKE": FrameType{id: "TKE", description: "Initial key", constructor: NewTextFrame},
		"TLA": FrameType{id: "TLA", description: "Language(s)", constructor: NewTextFrame},
		"TLE": FrameType{id: "TLE", description: "Length", constructor: NewTextFrame},
		"TMT": FrameType{id: "TMT", description: "Media type", constructor: NewTextFrame},
		"TOA": FrameType{id: "TOA", description: "Original artist(s)/performer(s)", constructor: NewTextFrame},
		"TOF": FrameType{id: "TOF", description: "Original filename", constructor: NewTextFrame},
		"TOL": FrameType{id: "TOL", description: "Original Lyricist(s)/text writer(s)", constructor: NewTextFrame},
		"TOR": FrameType{id: "TOR", description: "Original release year", constructor: NewTextFrame},
		"TOT": FrameType{id: "TOT", description: "Original album/Movie/Show title", constructor: NewTextFrame},
		"TP1": FrameType{id: "TP1", description: "Lead artist(s)/Lead performer(s)/Soloist(s)/Performing group", constructor: NewTextFrame},
		"TP2": FrameType{id: "TP2", description: "Band/Orchestra/Accompaniment", constructor: NewTextFrame},
		"TP3": FrameType{id: "TP3", description: "Conductor/Performer refinement", constructor: NewTextFrame},
		"TP4": FrameType{id: "TP4", description: "Interpreted, remixed, or otherwise modified by", constructor: NewTextFrame},
		"TPA": FrameType{id: "TPA", description: "Part of a set", constructor: NewTextFrame},
		"TPB": FrameType{id: "TPB", description: "Publisher", constructor: NewTextFrame},
		"TRC": FrameType{id: "TRC", description: "ISRC (International Standard Recording Code)", constructor: NewTextFrame},
		"TRD": FrameType{id: "TRD", description: "Recording dates", constructor: NewTextFrame},
		"TRK": FrameType{id: "TRK", description: "Track number/Position in set", constructor: NewTextFrame},
		"TSI": FrameType{id: "TSI", description: "Size", constructor: NewTextFrame},
		"TSS": FrameType{id: "TSS", description: "Software/hardware and settings used for encoding", constructor: NewTextFrame},
		"TT1": FrameType{id: "TT1", description: "Content group description", constructor: NewTextFrame},
		"TT2": FrameType{id: "TT2", description: "Title/Songname/Content description", constructor: NewTextFrame},
		"TT3": FrameType{id: "TT3", description: "Subtitle/Description refinement", constructor: NewTextFrame},
		"TXT": FrameType{id: "TXT", description: "Lyricist/text writer", constructor: NewTextFrame},
		"TXX": FrameType{id: "TXX", description: "User defined text information frame", constructor: NewDescTextFrame},
		"TYE": FrameType{id: "TYE", description: "Year", constructor: NewTextFrame},
		"UFI": FrameType{id: "UFI", description: "Unique file identifier", constructor: NewDataFrame},
		"ULT": FrameType{id: "ULT", description: "Unsychronized lyric/text transcription", constructor: NewDataFrame},
		"WAF": FrameType{id: "WAF", description: "Official audio file webpage", constructor: NewDataFrame},
		"WAR": FrameType{id: "WAR", description: "Official artist/performer webpage", constructor: NewDataFrame},
		"WAS": FrameType{id: "WAS", description: "Official audio source webpage", constructor: NewDataFrame},
		"WCM": FrameType{id: "WCM", description: "Commercial information", constructor: NewDataFrame},
		"WCP": FrameType{id: "WCP", description: "Copyright/Legal information", constructor: NewDataFrame},
		"WPB": FrameType{id: "WPB", description: "Publishers official webpage", constructor: NewDataFrame},
		"WXX": FrameType{id: "WXX", description: "User defined URL link frame", constructor: NewDataFrame},
	}
)

func NewV2Frame(reader io.Reader) Framer {
	data := make([]byte, V2FrameHeaderSize)
	if n, err := io.ReadFull(reader, data); n < V2FrameHeaderSize || err != nil {
		return nil
	}

	id := string(data[:3])
	t, ok := V2FrameTypeMap[id]
	if !ok {
		return nil
	}

	size, err := normint(data[3:6])
	if err != nil {
		return nil
	}

	h := FrameHead{
		FrameType: t,
		size:      size,
	}

	frameData := make([]byte, size)
	if n, err := io.ReadFull(reader, frameData); n < int(size) || err != nil {
		return nil
	}

	return t.constructor(h, frameData)
}

func V2Bytes(f Framer) []byte {
	headBytes := make([]byte, V2FrameHeaderSize)

	copy(headBytes[:3], []byte(f.Id()))
	copy(headBytes[3:6], normbytes(int32(f.Size()))[1:])

	return append(headBytes, f.Bytes()...)
}
