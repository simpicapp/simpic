ALTER TABLE album_contents
    ADD FOREIGN KEY (album_uuid) REFERENCES albums (album_uuid) ON DELETE CASCADE;

ALTER TABLE album_contents
    ADD FOREIGN KEY (photo_uuid) REFERENCES photos (photo_uuid) ON DELETE CASCADE;

ALTER TABLE albums
    ADD FOREIGN KEY (photo_uuid) REFERENCES photos (photo_uuid) ON DELETE SET NULL;
