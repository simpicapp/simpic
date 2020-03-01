package internal

type Visibility int

const (
	VisPublic Visibility = iota
	VisUnlisted
	VisPrivate
)
