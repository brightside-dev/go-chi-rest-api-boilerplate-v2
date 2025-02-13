package service

import (
	"context"
	"fmt"

	"github.com/brightside-dev/ronin-fitness-be/internal/model"
	"github.com/brightside-dev/ronin-fitness-be/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AdminUserService interface {
	Login(email, password string) (*model.AdminUser, error)
	Logout()
	Create(firstName, lastName, email, password string) (string, error)
}

type adminUserService struct {
	AdminUserRepository repository.AdminUserRepository
}

func NewAdminUserService(adminUserRepository repository.AdminUserRepository) AdminUserService {
	return &adminUserService{
		AdminUserRepository: adminUserRepository,
	}
}

func (s *adminUserService) Login(email, password string) (*model.AdminUser, error) {
	ctx := context.Background()
	adminUser, err := s.AdminUserRepository.GetByEmail(ctx, email)
	if err != nil {
		return &model.AdminUser{}, fmt.Errorf("invalid password and/or email")
	}

	if adminUser == nil {
		return &model.AdminUser{}, fmt.Errorf("invalid password and/or email")

	}

	// Compare the password
	err = bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(password))
	if err != nil {
		return &model.AdminUser{}, fmt.Errorf("invalid password and/or email")
	}

	return adminUser, nil
}

func (s *adminUserService) Logout() {
}

func (s *adminUserService) Create(firstName, lastName, email, password string) (string, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password")

	}

	// Create a new adminUser
	adminUser := model.AdminUser{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hashedPassword),
	}

	// Save adminUser to database

	ctx := context.Background()
	_, err = s.AdminUserRepository.Create(ctx, &adminUser)
	if err != nil {
		return "", fmt.Errorf("failed to save adminUser to database: %v", err)
	}

	return "Admin user created successfully", nil
}
