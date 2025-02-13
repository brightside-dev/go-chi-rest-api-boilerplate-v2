-- +goose Up
-- +goose StatementBegin
CREATE TABLE workout_reactions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    workout_id BIGINT UNSIGNED NOT NULL,
    profile_id BIGINT UNSIGNED NOT NULL,
    reaction ENUM('like', 'bicep_flex', 'fire', 'cold', 'star') NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_reactions_workout FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
    CONSTRAINT fk_reactions_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_workout_reactions_workout_id ON workout_reactions (workout_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_workout_reactions_profile_id ON workout_reactions (profile_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workout_reactions;
-- +goose StatementEnd
