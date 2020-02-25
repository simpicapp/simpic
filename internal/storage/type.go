package storage

type StoreKind int

const (
	KindPhoto StoreKind = iota
	KindThumbnail
)

var (
	StoreKinds = []StoreKind{KindPhoto, KindThumbnail}
)
