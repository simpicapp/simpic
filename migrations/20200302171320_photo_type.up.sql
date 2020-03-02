ALTER TABLE photos
    ADD COLUMN photo_format VARCHAR(10);

UPDATE photos
SET photo_format = 'JPEG'
WHERE photo_type = 1;

UPDATE photos
SET photo_format = 'PNG'
WHERE photo_type = 2;

UPDATE photos
SET photo_format = 'GIF'
WHERE photo_type = 3;

UPDATE photos
SET photo_format = 'TIFF'
WHERE photo_type = 4;

ALTER TABLE photos
    DROP COLUMN photo_type;
