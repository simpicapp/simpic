ALTER TABLE photos
    ADD COLUMN photo_width INT;

ALTER TABLE photos
    ADD COLUMN photo_height INT;

ALTER TABLE photos
    ADD COLUMN photo_size INT;

ALTER TABLE photos
    ADD COLUMN photo_format INT;

DROP TABLE photo_formats;
