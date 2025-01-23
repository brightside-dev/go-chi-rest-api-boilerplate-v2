-- +goose Up
-- +goose StatementBegin
CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    domain VARCHAR(255) NOT NULL,
    level VARCHAR(255) NOT NULL,
    message VARCHAR(255) NOT NULL,
    context TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE logs;
-- +goose StatementEnd
