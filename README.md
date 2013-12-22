# id3

ID3 library for Go. Work in progress.

# Install

The platform ($GOROOT/bin) "go get" tool is the best method to install.

    go get github.com/mikkyang/id3-go

This downloads and installs the package into your $GOPATH. If you only want to
recompile, use "go install".

    go install github.com/mikkyang/id3-go

# Usage

An import allows access to the package.

    import (
        id3 "github.com/mikkyang/id3-go"
    )

# Opening a File

To access the tag of a file, first open the file using the package's `Open`
function.

    mp3File, err := id3.Open("All-In.mp3")

It's also a good idea to ensure that the file is closed using `defer`.

    defer mp3File.Close()

# Accessing Frames

Some commonly used frames have fields in the file for easier access. These
frames are `Title`, `Artist`, `Album`, and `Genre`.

    mp3File.Artist = "Okasian"
    fmt.Println(mp3File.Artist)

## Other Frames

Other frames can be accessed directly by using the `Frames` field of the
file, which is a map, keyed by ID3 frame identifiers, of slices containing
`Framer` interfaces. This interfaces allow read access to general details of
the file.

    lyricsFrame := mp3File.Frames["USLT"][0]
    lyrics := lyricsFrame.String()

If more specific information is needed, or frame-specific write access is
needed, then the interface must be cast into the appropriate underlying type.
The example provided does not check for errors, but it is recommended to do
so.

    lyricsFrame := mp3File.Frames["USLT"][0].(*id3.UnsynchTextFrame)
