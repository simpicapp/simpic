CREATE TABLE albums
(
    album_id      SERIAL PRIMARY KEY,
    album_uuid    VARCHAR(36) UNIQUE NOT NULL,
    album_name    VARCHAR(128),
    album_cover   INT,
    album_owner   INT,
    album_created TIMESTAMP
);

CREATE TABLE album_contents
(
    album_id        INT NOT NULL,
    photo_id        INT NOT NULL,
    content_creator INT,
    content_order   INT,
    content_added   TIMESTAMP,
    PRIMARY KEY (album_id, photo_id)
);
