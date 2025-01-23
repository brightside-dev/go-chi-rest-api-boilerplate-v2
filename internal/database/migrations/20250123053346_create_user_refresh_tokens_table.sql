-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_refresh_tokens (
    id SERIAL PRIMARY KEY,                -- Unique identifier for the token entry
    user_id BIGINT UNSIGNED NOT NULL,             -- Foreign key linking to the users table
    token VARCHAR(255) NOT NULL UNIQUE,   -- The actual refresh token (must be unique)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), -- When the token was created
    expires_at TIMESTAMP NOT NULL,        -- When the token will expire
    revoked BOOLEAN NOT NULL DEFAULT FALSE, -- Indicates if the token has been revoked
    ip_address VARCHAR(45),              -- (Optional) IP address of the user
    user_agent TEXT,                      -- (Optional) User agent string (e.g., browser/device info)
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_refresh_tokens;
-- +goose StatementEnd
