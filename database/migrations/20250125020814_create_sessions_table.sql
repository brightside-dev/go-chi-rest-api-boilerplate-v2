-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    token VARCHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX sessions_expiry_idx ON `sessions` (expiry)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
