package service

import (
	"context"
	"errors"

	"github.com/grachmannico95/mileapp-test-be/internal/config"
	"github.com/grachmannico95/mileapp-test-be/internal/model"
	"github.com/grachmannico95/mileapp-test-be/internal/repository"
	"github.com/grachmannico95/mileapp-test-be/internal/util"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*model.User, string, string, error)
	Register(ctx context.Context, email, password string) (*model.User, error)
}

type authServiceImpl struct {
	userRepo repository.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo repository.UserRepository, config *config.Config) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *authServiceImpl) Login(ctx context.Context, email, password string) (*model.User, string, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}

	if user == nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	if !util.CheckPasswordHash(password, user.Password) {
		return nil, "", "", errors.New("invalid email or password")
	}

	jwtToken, err := util.GenerateJWT(user.ID, user.Email, s.config.JWT.Secret, s.config.JWT.Expiry)
	if err != nil {
		return nil, "", "", err
	}

	csrfToken := util.GenerateCSRFToken(s.config.CSRF.Secret)

	return user, jwtToken, csrfToken, nil
}

func (s *authServiceImpl) Register(ctx context.Context, email, password string) (*model.User, error) {
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := model.NewUser(email, hashedPassword)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
