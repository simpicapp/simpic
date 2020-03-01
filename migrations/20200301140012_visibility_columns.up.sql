ALTER TABLE photos
    ADD COLUMN photo_visibility INT NOT NULL DEFAULT 0;

ALTER TABLE albums
    ADD COLUMN album_visibility INT NOT NULL DEFAULT 0;
