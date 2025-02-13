package model

import "time"

// CREATE TABLE verification_codes (
//     id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
//     user_id BIGINT UNSIGNED NOT NULL,
//     email VARCHAR(255) NOT NULL,
//     code VARCHAR(255) NOT NULL,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     expires_at TIMESTAMP NOT NULL,
//     CONSTRAINT fk_verification_codes_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
// );

type VerificationCode struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
