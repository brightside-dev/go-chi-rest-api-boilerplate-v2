-- +goose Up
-- +goose StatementBegin
CREATE TABLE sets (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    workout_id BIGINT UNSIGNED NOT NULL,
    exercise_id BIGINT UNSIGNED NOT NULL,
    set_number INT NOT NULL,
    duration INT NOT NULL,
    weight_kg DECIMAL(10, 2) NOT NULL,
    weight_lb DECIMAL(10, 2) NOT NULL,
    reps INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_sets_workout FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
    CONSTRAINT fk_sets_exercise FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_sets_workout_id ON sets (workout_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_sets_exercise_id ON sets (exercise_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_sets_workout_exercise ON sets (workout_id, exercise_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sets;
-- +goose StatementEnd
