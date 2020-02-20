ALTER TABLE albums
    RENAME COLUMN photo_uuid TO album_cover;

ALTER TABLE albums
    ALTER COLUMN album_cover TYPE INT;
