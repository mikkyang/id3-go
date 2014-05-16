package v2

import (
	"testing"
)

func TestUnsynchTextFrame_SetEncoding(t *testing.T) {
	f := NewUnsychTextFrame(V23CommonFrame["Comments"], "Foo", "Bar")
	size := f.Size()

	err := f.SetEncoding("UTF-16")
	if err != nil {
		t.Fatal(err)
	}
	newSize := f.Size()
	if newSize-size != 1 {
		t.Errorf("expected size to increase to %d, but it was %d", size+1, newSize)
	}

	size = newSize
	err := f.SetEncoding("UTF-16")
	if err != nil {
		t.Fatal(err)
	}
	newSize := f.Size()
	if newSize-size != -1 {
		t.Errorf("expected size to decrease to %d, but it was %d", size-1, newSize)
	}
}
