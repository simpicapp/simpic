package internal

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Photo struct {
	Id        uuid.UUID `json:"id" db:"photo_uuid"`
	FileName  string    `json:"file_name" db:"photo_filename"`
	Width     int       `json:"width" db:"photo_width"`
	Height    int       `json:"height" db:"photo_height"`
	Timestamp time.Time `json:"timestamp" db:"photo_uploaded"`
	Type      PhotoType `json:"type" db:"photo_type"`
	Uploader  int       `json:"user_id" db:"photo_uploader"`
}

func NewPhoto(fileName string) *Photo {
	return &Photo{
		Id:        uuid.NewV4(),
		FileName:  fileName,
		Timestamp: time.Now(),
	}
}
