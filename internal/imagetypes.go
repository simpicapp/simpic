package internal

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type PhotoType int

const (
	Unknown PhotoType = iota
	Jpeg
	Png
	Gif
)
