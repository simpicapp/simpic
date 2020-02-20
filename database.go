package simpic

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	uuid "github.com/satori/go.uuid"
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

	db, err := postgresql.Open(url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	database := &Database{db: db}
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

// Photos

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

// Albums

func (d *Database) AddAlbum(album *Album) (err error) {
	_, err = d.db.Collection("albums").Insert(album)
	return
}

func (d *Database) GetAlbums(offset, count int) (albums []Album, err error) {
	err = d.db.Collection("albums").Find().OrderBy("album_name").Offset(offset).Limit(count).All(&albums)
	return
}

// Users

func (d *Database) AddUser(user *User) (err error) {
	_, err = d.db.Collection("users").Insert(user)
	return
}

func (d *Database) GetUser(username string) (user *User, err error) {
	err = d.db.Collection("users").Find("user_name", username).One(&user)
	return
}

// Sessions

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
