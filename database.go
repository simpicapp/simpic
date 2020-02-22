package simpic

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	uuid "github.com/satori/go.uuid"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

type Database struct {
	db sqlbuilder.Database
}

func OpenDatabase(dsn, migrationPath string) (*Database, error) {
	url, err := postgresql.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	conn, err := postgresql.Open(url)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	database := &Database{db: conn}
	if err := database.migrate(migrationPath); err != nil {
		return nil, err
	}

	return database, nil
}

func (d *Database) migrate(migrationPath string) error {
	driver, err := postgres.WithInstance(d.db.Driver().(*sql.DB), &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationPath), "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

//region AddedPhotos

func (d *Database) Add(photo *Photo) (err error) {
	_, err = d.db.Collection("photos").Insert(photo)
	return
}

func (d *Database) GetPhoto(id uuid.UUID) (photo *Photo, err error) {
	err = d.db.Collection("photos").Find("photo_uuid", id).One(&photo)
	return
}

func (d *Database) GetPhotosByTime(offset, count int) (photos []Photo, err error) {
	err = d.db.Collection("photos").Find().OrderBy("-photo_uploaded").Offset(offset).Limit(count).All(&photos)
	return
}

func (d *Database) DeletePhoto(photo *Photo) error {
	return d.db.Collection("photos").Find("photo_uuid", photo.Id).Delete()
}

//endregion

//region Albums

func (d *Database) AddAlbum(album *Album) (err error) {
	_, err = d.db.Collection("albums").Insert(album)
	return
}

func (d *Database) UpdateAlbum(album *Album) error {
	return d.db.Collection("albums").Find("album_uuid", album.Id).Update(album)
}

func (d *Database) GetAlbum(id uuid.UUID) (album *Album, err error) {
	err = d.db.Collection("albums").Find("album_uuid", id).One(&album)
	return
}

func (d *Database) GetAlbums(offset, count int) (albums []Album, err error) {
	err = d.db.Collection("albums").Find().OrderBy("album_name").Offset(offset).Limit(count).All(&albums)
	return
}

func (d *Database) GetAlbumOrderMax(album uuid.UUID) (int, error) {
	var result struct {
		Max *int `db:"max"`
	}

	err := d.db.Select(db.Raw(`max(content_order) "max"`)).From("album_contents").
		Where("album_uuid = ?", album).One(&result)

	if result.Max == nil {
		return 0, err
	} else {
		return *result.Max, err
	}
}

func (d *Database) GetAlbumPhotos(album uuid.UUID, offset, count int) (photos []AlbumPhoto, err error) {
	err = d.db.SelectFrom("album_contents").
		Join("photos").Using("photo_uuid").
		Where("album_uuid = ?", album).OrderBy("-content_order").
		Offset(offset).Limit(count).All(&photos)
	return
}

func (d *Database) AddAlbumPhotos(photos []AlbumEntry) error {
	batch := d.db.InsertInto("album_contents").Amend(onConflictDoNothing).Batch(20)
	for _, photo := range photos {
		batch.Values(photo)
	}
	batch.Done()
	return batch.Wait()
}

func (d *Database) RemoveAlbumPhotos(album uuid.UUID, photos []uuid.UUID) error {
	return d.db.Collection("album_contents").
		Find(db.And(db.Cond{"album_uuid": album}, db.Cond{"photo_uuid": photos})).
		Delete()
}

//endregion

//region Users

func (d *Database) AddUser(user *User) (err error) {
	_, err = d.db.Collection("users").Insert(user)
	return
}

func (d *Database) GetUser(username string) (user *User, err error) {
	err = d.db.Collection("users").Find("user_name", username).One(&user)
	return
}

//endregion

//region Sessions

func (d *Database) AddSession(session *Session) (err error) {
	_, err = d.db.Collection("sessions").Insert(session)
	return
}

func (d *Database) GetSession(sessionKey string) (session *SessionUser, err error) {
	err = d.db.SelectFrom("sessions").Join("users").Using("user_id").
		Where("session_key = ? AND session_expires > NOW()", sessionKey).
		One(&session)
	return
}

//endregion

func onConflictDoNothing(queryIn string) (queryOut string) {
	return fmt.Sprintf("%s ON CONFLICT DO NOTHING", queryIn)
}
