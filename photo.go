package simpic

import (
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

type Type int

const (
	Unknown Type = iota
	Jpeg
	Png
)

type Photo struct {
	Id        uuid.UUID `db:"photo_uuid"`
	FileName  string    `json:"file_name" db:"photo_filename"`
	Width     int       `json:"width" db:"photo_width"`
	Height    int       `json:"height" db:"photo_height"`
	Timestamp time.Time `json:"timestamp" db:"photo_uploaded"`
	Type      Type      `json:"type" db:"photo_type"`
}

func NewPhoto(fileName string) *Photo {
	return &Photo{
		Id:        uuid.NewV4(),
		FileName:  fileName,
		Width:     0,
		Height:    0,
		Timestamp: time.Now(),
		Type:      typeFromFilename(fileName),
	}
}

func typeFromFilename(fileName string) Type {
	lower := strings.ToLower(fileName)
	if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") {
		return Jpeg
	} else if strings.HasSuffix(lower, ".png") {
		return Png
	} else {
		return Unknown
	}
}
