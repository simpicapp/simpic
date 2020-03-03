ALTER TABLE photos
    DROP COLUMN photo_width;

ALTER TABLE photos
    DROP COLUMN photo_height;

ALTER TABLE photos
    DROP COLUMN photo_size;

ALTER TABLE photos
    DROP COLUMN photo_format;

CREATE TABLE photo_formats
(
    photo_uuid VARCHAR(36) REFERENCES photos (photo_uuid) ON DELETE CASCADE,
    format_purpose INT,
    format_format VARCHAR(10),
    format_width INT,
    format_height INT,
    format_size INT,
    PRIMARY KEY (photo_uuid, format_purpose, format_format)
);
