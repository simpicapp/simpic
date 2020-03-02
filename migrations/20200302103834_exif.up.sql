ALTER TABLE photos
    ADD COLUMN photo_processed INT NOT NULL DEFAULT 1;

CREATE TABLE photo_exif
(
    photo_uuid VARCHAR(36) REFERENCES photos (photo_uuid) ON DELETE CASCADE,
    exif_field VARCHAR(50) NOT NULL,
    exif_value VARCHAR(128) NOT NULL,
    PRIMARY KEY (photo_uuid, exif_field)
);
