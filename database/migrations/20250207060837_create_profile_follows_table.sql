-- +goose Up
-- +goose StatementBegin
CREATE TABLE profile_follows (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    profile_id BIGINT UNSIGNED NOT NULL,
    follower_profile_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_profile_followed FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
    CONSTRAINT fk_profile_follower FOREIGN KEY (follower_profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_profile_follows_profile_id ON profile_follows (profile_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_profile_follows_follower_profile_id ON profile_follows (follower_profile_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
