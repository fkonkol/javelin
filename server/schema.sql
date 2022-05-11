-- Represents user's identity
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Status of friendship between users
CREATE TYPE friends_status AS ENUM (
    'PENDING_SENT',
    'PENDING_RECEIVED',
    'ACCEPTED',
    'BLOCKED'
);
-- Represents friendship between two users
CREATE TABLE friends (
    uuid VARCHAR(32) NOT NULL,
    status friends_status NOT NULL,
    user_id INT NOT NULL REFERENCES users(id),
    friend_id INT NOT NULL REFERENCES users(id),
    PRIMARY KEY(user_id, friend_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
