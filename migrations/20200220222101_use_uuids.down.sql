ALTER TABLE photos
    DROP CONSTRAINT photos_pkey;

ALTER TABLE photos
    ADD COLUMN photo_id SERIAL;

ALTER TABLE photos
    ADD PRIMARY KEY (photo_id);


ALTER TABLE albums
    DROP CONSTRAINT albums_pkey;

ALTER TABLE albums
    ADD COLUMN album_id SERIAL;

ALTER TABLE albums
    ADD PRIMARY KEY (album_id);


ALTER TABLE album_contents
    ALTER COLUMN album_uuid TYPE INT;

ALTER TABLE album_contents
    ALTER COLUMN photo_uuid TYPE INT;

ALTER TABLE album_contents
    RENAME COLUMN photo_uuid TO photo_id;

ALTER TABLE album_contents
    RENAME COLUMN album_uuid TO album_id;
