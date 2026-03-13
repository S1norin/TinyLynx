-- +goose Up
CREATE TABLE link (
    id SERIAL PRIMARY KEY,
    long TEXT NOT NULL,
    short TEXT NOT NULL,
);

-- +goose Down
DROP TABLE link;