CREATE TABLE users
(
    id         VARCHAR(36) PRIMARY KEY,
    email      TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name  TEXT NOT NULL
);

CREATE TABLE refresh_tokens
(
    id         VARCHAR(36) PRIMARY KEY,
    user_id    VARCHAR(36)              NOT NULL REFERENCES users (id),
    token      TEXT                     NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);