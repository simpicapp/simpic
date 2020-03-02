package internal

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Photo struct {
	Id         uuid.UUID  `json:"id" db:"photo_uuid"`
	FileName   string     `json:"file_name" db:"photo_filename"`
	Width      uint       `json:"width" db:"photo_width"`
	Height     uint       `json:"height" db:"photo_height"`
	Timestamp  time.Time  `json:"timestamp" db:"photo_uploaded"`
	Format     string     `json:"type" db:"photo_format"`
	Uploader   int        `json:"user_id" db:"photo_uploader"`
	Visibility Visibility `json:"visibility" db:"photo_visibility"`
	Processed  int        `json:"-" db:"photo_processed"`
}

type ExifTag struct {
	Photo uuid.UUID `db:"photo_uuid"`
	Field string    `db:"exif_field"`
	Value string    `db:"exif_value"`
}

func NewPhoto(fileName string) *Photo {
	return &Photo{
		Id:        uuid.NewV4(),
		FileName:  fileName,
		Timestamp: time.Now(),
	}
}
