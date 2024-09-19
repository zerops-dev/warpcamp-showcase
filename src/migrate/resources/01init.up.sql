-- +migrate Up
CREATE TABLE IF NOT EXISTS messages
(
    id  serial PRIMARY KEY,
    text TEXT NOT NULL
);