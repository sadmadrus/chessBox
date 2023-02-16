-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions(
    uid INT NOT NULL PRIMARY KEY, 
    start_date INTEGER,
    end_date INTEGER,
    w_player INTEGER,
    b_player INTEGER,
    description INTEGER,
    session_state INTEGER,
    position TEXT,

    FOREIGN KEY(w_player) REFERENCES users(uid),
    FOREIGN KEY(b_player) REFERENCES users(uid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
