package storage

type StoreKind int

const (
	KindRaw StoreKind = iota
	KindThumbnailJpeg
	KindScreenJpeg
)

var (
	StoreKinds = []StoreKind{KindRaw, KindThumbnailJpeg, KindScreenJpeg}
)
