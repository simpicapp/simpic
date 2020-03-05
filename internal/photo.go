package internal

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Photo struct {
	Id         uuid.UUID  `json:"id" db:"photo_uuid"`
	FileName   string     `json:"file_name" db:"photo_filename"`
	UploadTime *time.Time `json:"uploaded" db:"photo_uploaded"`
	TakenTime  *time.Time `json:"taken" db:"photo_taken"`
	Uploader   int        `json:"user_id" db:"photo_uploader"`
	Visibility Visibility `json:"visibility" db:"photo_visibility"`
	Processed  int        `json:"-" db:"photo_processed"`
}

type FormatPurpose int

const (
	PurposePreview FormatPurpose = iota + 1
	PurposeScreen
	PurposeOriginal
)

type Format struct {
	Photo   uuid.UUID     `json:"-" db:"photo_uuid"`
	Purpose FormatPurpose `json:"purpose" db:"format_purpose"`
	Format  string        `json:"format" db:"format_format"`
	Width   uint          `json:"width" db:"format_width"`
	Height  uint          `json:"height" db:"format_height"`
	Size    int64         `json:"size" db:"format_size"`
}

type PhotoWithFormats struct {
	*Photo
	Formats []Format `json:"formats"`
}

type ExifTag struct {
	Photo uuid.UUID `db:"photo_uuid"`
	Field string    `db:"exif_field"`
	Value string    `db:"exif_value"`
}

func NewPhoto(fileName string) *Photo {
	t := time.Now()
	return &Photo{
		Id:         uuid.NewV4(),
		FileName:   fileName,
		UploadTime: &t,
	}
}
