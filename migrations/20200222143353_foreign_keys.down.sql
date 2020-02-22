ALTER TABLE album_contents
    DROP CONSTRAINT album_contents_album_uuid_fkey;

ALTER TABLE album_contents
    DROP CONSTRAINT album_contents_photo_uuid_fkey;

ALTER TABLE albums
    DROP CONSTRAINT albums_photo_uuid_fkey;
