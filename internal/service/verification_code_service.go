package service

import (
	"math/rand"

	"github.com/brightside-dev/ronin-fitness-be/database/client"
)

type VerificationCodeService interface {
	GenerateVerificationCode() string
}

type verificationCodeService struct {
	DB client.DatabaseService
}

func NewVerificationCodeService(db client.DatabaseService) VerificationCodeService {
	return &verificationCodeService{DB: db}
}

func (s *verificationCodeService) GenerateVerificationCode() string {
	// Define characters for letters (A-Z) and digits (0-9)
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"

	// Create a slice for the result
	code := make([]rune, 5)

	// Randomly place 3 letters and 2 digits
	for i := 0; i < 3; i++ {
		code[i] = rune(letters[rand.Intn(len(letters))])
	}
	for i := 3; i < 5; i++ {
		code[i] = rune(digits[rand.Intn(len(digits))])
	}

	// Shuffle the slice to randomize the positions of letters and digits
	rand.Shuffle(5, func(i, j int) {
		code[i], code[j] = code[j], code[i]
	})

	return string(code)
}
