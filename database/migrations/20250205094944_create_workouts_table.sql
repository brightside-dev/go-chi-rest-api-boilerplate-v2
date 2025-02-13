-- +goose Up
-- +goose StatementBegin
CREATE TABLE workouts (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    profile_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    mental_energy_level TINYINT UNSIGNED NOT NULL CHECK (mental_energy_level BETWEEN 1 AND 10),
    physical_energy_level TINYINT UNSIGNED NOT NULL CHECK (physical_energy_level BETWEEN 1 AND 10),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_workouts_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_workouts_profile_id ON workouts (profile_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workouts;
-- +goose StatementEnd
