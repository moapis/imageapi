-- +migrate Up
CREATE TABLE images (
    id serial PRIMARY KEY,
    link text NOT NULL,
    created_at timestamptz
);
-- +migrate Down
DROP TABLE images;