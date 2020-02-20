ALTER TABLE albums
    RENAME COLUMN album_cover TO photo_uuid;

ALTER TABLE albums
    ALTER COLUMN photo_uuid TYPE VARCHAR(36);
