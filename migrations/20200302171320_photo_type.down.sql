ALTER TABLE photos
    ADD COLUMN photo_type INT;

UPDATE photos
SET photo_type = 1
WHERE photo_format = 'JPEG';


UPDATE photos
SET photo_type = 2
WHERE photo_format = 'PNG';

UPDATE photos
SET photo_type = 3
WHERE photo_format = 'GIF';

UPDATE photos
SET photo_type = 4
WHERE photo_format = 'TIFF';

ALTER TABLE photos
    DROP COLUMN photo_format;
