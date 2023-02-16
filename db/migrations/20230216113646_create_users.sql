-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    uid INT NOT NULL PRIMARY KEY, 
    email TEXT(100) NOT NULL unique,
    password TEXT NOT NULL,
    name TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
