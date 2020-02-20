ALTER TABLE users
    DROP COLUMN user_session_key;

CREATE TABLE sessions
(
    session_id  SERIAL PRIMARY KEY,
    session_key VARCHAR(128) UNIQUE,
    user_id INT,
    session_created TIMESTAMP,
    session_expires TIMESTAMP,
    session_ip VARCHAR(128),
    session_user_agent VARCHAR(256)
)
