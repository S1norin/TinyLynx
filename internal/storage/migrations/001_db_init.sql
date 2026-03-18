-- +goose Up
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    original_link TEXT NOT NULL,
    short_code TEXT NOT NULL
);

-- +goose Down
DROP TABLE links;