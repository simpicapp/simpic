ALTER TABLE photos
    DROP CONSTRAINT photos_pkey;

ALTER TABLE photos
    DROP COLUMN photo_id;

ALTER TABLE photos
    ADD PRIMARY KEY (photo_uuid);


ALTER TABLE albums
    DROP CONSTRAINT albums_pkey;

ALTER TABLE albums
    DROP COLUMN album_id;

ALTER TABLE albums
    ADD PRIMARY KEY (album_uuid);


ALTER TABLE album_contents
    ALTER COLUMN album_id TYPE VARCHAR(36);

ALTER TABLE album_contents
    ALTER COLUMN photo_id TYPE VARCHAR(36);

ALTER TABLE album_contents
    RENAME COLUMN photo_id TO photo_uuid;

ALTER TABLE album_contents
    RENAME COLUMN album_id TO album_uuid;
