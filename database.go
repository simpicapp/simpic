package simpic

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type Database struct {
	db *sqlx.DB
}

func OpenDatabase(dsn, migrationPath string) (*Database, error) {
	db, err := sqlx.Connect("postgres", dsn)
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
	driver, err := postgres.WithInstance(d.db.DB, &postgres.Config{})
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

func (d *Database) GetPhoto(id uuid.UUID) (*Photo, error) {
	var photo Photo
	return &photo,
		d.db.Get(&photo,
			`SELECT
				photo_uuid,
				photo_filename,
				photo_width,
				photo_height,
				photo_uploaded,
				photo_type
			FROM photos
			WHERE photo_uuid = $1
			LIMIT 1
		`, id.String())
}

func (d *Database) GetPhotosByTime(offset, count int) ([]Photo, error) {
	var photos []Photo
	return photos, d.db.Select(&photos,
		`SELECT
				photo_uuid,
				photo_filename,
				photo_width,
				photo_height,
				photo_uploaded,
				photo_type
			FROM photos
			ORDER BY photo_uploaded DESC
			LIMIT $1
			OFFSET $2
	`, count, offset)
}

func (d *Database) StorePhoto(photo *Photo) error {
	_, err := d.db.Exec(
		`INSERT INTO photos (
			photo_uuid, photo_filename,
			photo_width, photo_height,
			photo_uploaded, photo_type
		) VALUES ($1, $2, $3, $4, $5, $6)`,
		photo.Id, photo.FileName,
		photo.Width, photo.Height,
		photo.Timestamp, photo.Type)
	return err
}
