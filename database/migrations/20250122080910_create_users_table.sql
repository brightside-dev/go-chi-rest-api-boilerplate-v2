-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    birthday DATE NOT NULL,
    country VARCHAR(100) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX idx_users_email ON users (email);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_users_name ON users (last_name, first_name);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_users_country ON users (country);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_users_created_at ON users (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
