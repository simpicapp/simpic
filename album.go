package simpic

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Album struct {
	Uuid    uuid.UUID  `json:"id" db:"album_uuid"`
	Name    string     `json:"name" db:"album_name"`
	Cover   *uuid.UUID `json:"cover_photo,omitempty" db:"photo_uuid,omitempty"`
	Owner   int        `json:"owner_id" db:"album_owner"`
	Created time.Time  `json:"created" db:"album_created"`
}

type AlbumPhoto struct {
	Photo
	Creator int       `json:"creator_id" db:"content_creator"`
	Order   int       `json:"order" db:"content_order"`
	Added   time.Time `json:"added" db:"content_added"`
}

func NewAlbum(name string, owner int) *Album {
	return &Album{
		Uuid:    uuid.NewV4(),
		Name:    name,
		Owner:   owner,
		Created: time.Now(),
	}
}
