package v1

import (
	"io/ioutil"
	"os"
	"testing"
)

var tripTests = []struct {
	title, artist, album, year, comment string
	genre                               byte
}{
	{"Foo", "Bar", "Baz", "2014", "Blah", 1},
	{"Foo\x00Qux", "Bar", "Baz", "2014", "Blah", 1},
}

func TestParseTag_RoundTrip(t *testing.T) {
	for testNum, tt := range tripTests {
		f, err := ioutil.TempFile("", "id3v1")
		if err != nil {
			t.Errorf("test %d: %s", testNum, err)
			continue
		}
		defer os.Remove(f.Name())
		defer f.Close()

		tag := &Tag{title: tt.title, artist: tt.artist, album: tt.album,
			year: tt.year, comment: tt.comment, genre: tt.genre}
		_, err = f.Write(tag.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		_, err = f.Seek(0, 0)
		if err != nil {
			t.Fatal(err)
		}

		resultTag := ParseTag(f)
		if *tag != *resultTag {
			t.Errorf("expected %q, got %q", *tag, *resultTag)
		}
	}
}
