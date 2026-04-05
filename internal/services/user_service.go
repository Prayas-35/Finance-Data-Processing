package services

import (
	"context"
	"fmt"

	"github.com/Prayas-35/Finance-Data-Processing/internal/models"
	"github.com/Prayas-35/Finance-Data-Processing/internal/repositories"
	"github.com/Prayas-35/Finance-Data-Processing/internal/utils"
)

type UserService struct {
	users *repositories.UserRepository
}

func NewUserService(users *repositories.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) Create(ctx context.Context, name, email, password, role string) (string, error) {
	if name == "" || email == "" || password == "" {
		return "", fmt.Errorf("name, email and password are required")
	}

	switch models.Role(role) {
	case models.RoleViewer, models.RoleAnalyst, models.RoleAdmin:
	default:
		return "", fmt.Errorf("invalid role")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}
	return s.users.Create(ctx, name, email, hash, models.Role(role))
}

func (s *UserService) List(ctx context.Context, limit, offset int) ([]models.User, error) {
	return s.users.List(ctx, limit, offset)
}

func (s *UserService) Update(ctx context.Context, id, name, role string) error {
	if id == "" || name == "" {
		return fmt.Errorf("id and name are required")
	}
	return s.users.Update(ctx, id, name, models.Role(role))
}

func (s *UserService) SetActive(ctx context.Context, id string, active bool) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	return s.users.UpdateActive(ctx, id, active)
}

func (s *UserService) EnsureSeedAdmin(ctx context.Context, email, name, password string) error {
	if email == "" || name == "" || password == "" {
		return nil
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	return s.users.SeedAdmin(ctx, email, name, hash)
}
