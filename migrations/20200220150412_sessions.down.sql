ALTER TABLE users
    ADD COLUMN user_session_key BYTEA;

DROP TABLE sessions;
