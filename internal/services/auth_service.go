package services

import (
	"context"
	"fmt"

	"github.com/Prayas-35/Finance-Data-Processing/internal/auth"
	"github.com/Prayas-35/Finance-Data-Processing/internal/repositories"
	"github.com/Prayas-35/Finance-Data-Processing/internal/utils"
)

type AuthService struct {
	users      *repositories.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthService(users *repositories.UserRepository, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{users: users, jwtManager: jwtManager}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	if !user.IsActive {
		return "", fmt.Errorf("user is inactive")
	}
	if !utils.VerifyPassword(password, user.PasswordHash) {
		return "", fmt.Errorf("invalid credentials")
	}
	return s.jwtManager.Generate(user.ID, string(user.Role))
}
