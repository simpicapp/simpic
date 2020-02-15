CREATE TABLE users
(
    user_id            SERIAL PRIMARY KEY,
    user_name          VARCHAR(128) UNIQUE,
    user_password_salt BYTEA,
    user_password_hash BYTEA,
    user_session_key   BYTEA,
    user_admin         BOOLEAN
)
