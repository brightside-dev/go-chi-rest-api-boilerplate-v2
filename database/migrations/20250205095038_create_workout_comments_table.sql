-- +goose Up
-- +goose StatementBegin
CREATE TABLE workout_comments (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    workout_id BIGINT UNSIGNED NOT NULL,
    profile_id BIGINT UNSIGNED NOT NULL,
    comment TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_comments_workout FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_profile FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_workout_comments_workout_id ON workout_comments(workout_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_workout_comments_profile_id ON workout_comments (profile_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workout_comments;
-- +goose StatementEnd
