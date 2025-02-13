-- +goose Up
-- +goose StatementBegin
CREATE TABLE exercises (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    profile_id BIGINT UNSIGNED DEFAULT NULL, -- Nullable for default exercises
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    icon_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_exercises_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_exercises_name ON exercises (name);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_profile_id ON exercises (profile_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE exercises;
-- +goose StatementEnd
