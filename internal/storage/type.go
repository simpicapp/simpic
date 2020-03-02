package storage

type StoreKind int

const (
	KindRaw StoreKind = iota
	KindThumbnail
	KindScreenJpeg
)

var (
	StoreKinds = []StoreKind{KindRaw, KindThumbnail, KindScreenJpeg}
)
