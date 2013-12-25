// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"bytes"
	"testing"
)

func TestSynch(t *testing.T) {
	synch := []byte{0x44, 0x7a, 0x70, 0x04}
	const synchResult = 144619524

	if result, err := synchint(synch); result != synchResult {
		t.Errorf("synchint(%v) = %d with error %v, want %d", synch, result, err, synchResult)
	}
	if result := synchbytes(synchResult); !bytes.Equal(result, synch) {
		t.Errorf("synchbytes(%d) = %v, want %v", synchResult, result, synch)
	}
}

func TestNorm(t *testing.T) {
	norm := []byte{0x0b, 0x95, 0xae, 0xb4}
	const normResult = 194358964

	if result, err := normint(norm); result != normResult {
		t.Errorf("normint(%v) = %d with error %v, want %d", norm, result, err, normResult)
	}
	if result := normbytes(normResult); !bytes.Equal(result, norm) {
		t.Errorf("normbytes(%d) = %v, want %v", normResult, result, norm)
	}
}
