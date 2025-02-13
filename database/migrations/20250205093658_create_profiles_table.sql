-- +goose Up
-- +goose StatementBegin
CREATE TABLE profiles (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,                
    user_id BIGINT UNSIGNED NOT NULL,             
    display_name VARCHAR(255) NOT NULL,
    privacy ENUM('public', 'private') NOT NULL DEFAULT 'public',
    avatar_version INT NOT NULL DEFAULT 1,
    is_notifications_enabled TINYINT(1) NOT NULL DEFAULT 1,
    fitness_experience VARCHAR(255) NOT NULL,
    experience_points INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, 
    CONSTRAINT fk_profiles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_profiles_user_id ON profiles (user_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_profiles_display_name ON profiles (display_name);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_profiles_experience_points ON profiles (experience_points DESC);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_profiles_created_at ON profiles (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE profiles;
-- +goose StatementEnd
