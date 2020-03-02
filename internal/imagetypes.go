package internal

import (
	"bufio"
	_ "golang.org/x/image/tiff"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

type PhotoType int

const (
	TypeUnknown PhotoType = iota
	TypeJpeg
	TypePng
	TypeGif
	TypeTiff

	gifMagic    = "GIF8?a"
	jpegMagic   = "\xff\xd8"
	pngMagic    = "png"
	tiffLeMagic = "II\x2A\x00"
	tiffBeMagic = "MM\x00\x2A"
)

type format struct {
	t PhotoType
	m string
}

var formats = []format{
	{TypeJpeg, jpegMagic},
	{TypeGif, gifMagic},
	{TypePng, pngMagic},
	{TypeTiff, tiffLeMagic},
	{TypeTiff, tiffBeMagic},
}

type peekReader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

func asPeekReader(r io.Reader) peekReader {
	if pr, ok := r.(peekReader); ok {
		return pr
	}
	return bufio.NewReader(r)
}

func match(magic string, b []byte) bool {
	for i, c := range b {
		if magic[i] != c && magic[i] != '?' {
			return false
		}
	}
	return true
}

func SniffPhotoType(r io.Reader) PhotoType {
	pr := asPeekReader(r)
	for _, f := range formats {
		b, err := pr.Peek(len(f.m))
		if err == nil && match(f.m, b) {
			return f.t
		}
	}
	return TypeUnknown
}
