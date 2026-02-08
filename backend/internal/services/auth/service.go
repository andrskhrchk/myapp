package auth

import (
	"context"
	"fmt"
	"time"

	jwt "github.com/andrskhrchk/myapp/pkg/jwt"

	"github.com/andrskhrchk/myapp/internal/domain"
	"github.com/andrskhrchk/myapp/internal/repository/postgres"
	"github.com/andrskhrchk/myapp/internal/transport/dto"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo postgres.UserRepository
	jwtMgr   *jwt.TokenManager
}

func NewAuthService(userRepo postgres.UserRepository, jwtMgr *jwt.TokenManager) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtMgr:   jwtMgr,
	}
}

func (s *AuthService) Register(ctx context.Context, regData *dto.RegisterDTO) (*domain.User, string, error) {
	if _, err := s.userRepo.GetUserByEmail(ctx, regData.Email); err != nil {
		return nil, "", fmt.Errorf("user already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(regData.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	user := &domain.User{
		Email:        regData.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    regData.FirstName,
		LastName:     regData.LastName,
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := s.jwtMgr.CreateToken(int64(user.ID), 24*time.Hour)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
