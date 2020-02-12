CREATE TABLE photos
(
    photo_id       SERIAL PRIMARY KEY,
    photo_uuid     VARCHAR(36) UNIQUE,
    photo_filename VARCHAR(256),
    photo_width    INT,
    photo_height   INT,
    photo_size     INT,
    photo_uploaded TIMESTAMP,
    photo_type     INT
);
