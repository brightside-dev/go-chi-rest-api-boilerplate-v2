-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,  -- Use BIGINT for better indexing
    user_id BIGINT UNSIGNED NOT NULL,   -- Foreign key linking to the users table
    token VARCHAR(255) NOT NULL UNIQUE, -- The actual refresh token (must be unique)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    expires_at TIMESTAMP NOT NULL,      -- Expiration timestamp
    revoked TINYINT(1) NOT NULL DEFAULT 0, -- BOOLEAN fix (MySQL treats BOOLEAN as TINYINT(1))
    ip_address VARCHAR(45),             -- (Optional) IP address of the user
    user_agent TEXT,                     -- (Optional) User agent string (e.g., browser/device info)
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens (expires_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_refresh_tokens_revoked ON refresh_tokens (revoked);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_refresh_tokens;
-- +goose StatementEnd
