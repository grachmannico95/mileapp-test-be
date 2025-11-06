package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/grachmannico95/mileapp-test-be/internal/config"
	"github.com/grachmannico95/mileapp-test-be/internal/model"
	"github.com/grachmannico95/mileapp-test-be/internal/util"
	"github.com/grachmannico95/mileapp-test-be/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAuthService_Login_Success(t *testing.T) {
	// Setup
	mockUserRepo := mocks.NewMockUserRepository(t)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expiry: 15 * time.Minute,
		},
		CSRF: config.CSRFConfig{
			Secret: "csrf-secret",
		},
	}

	authService := NewAuthService(mockUserRepo, cfg)

	// Test data
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := util.HashPassword(password)
	userID := primitive.NewObjectID()

	user := &model.User{
		ID:        userID,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	mockUserRepo.EXPECT().
		FindByEmail(mock.Anything, email).
		Return(user, nil).
		Once()

	// Execute
	resultUser, jwtToken, csrfToken, err := authService.Login(context.Background(), email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.Equal(t, email, resultUser.Email)
	assert.NotEmpty(t, jwtToken)
	assert.NotEmpty(t, csrfToken)

	// Verify JWT token is valid
	claims, err := util.ValidateJWT(jwtToken, cfg.JWT.Secret)
	assert.NoError(t, err)
	assert.Equal(t, userID.Hex(), claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Setup
	mockUserRepo := mocks.NewMockUserRepository(t)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expiry: 15 * time.Minute,
		},
		CSRF: config.CSRFConfig{
			Secret: "csrf-secret",
		},
	}

	authService := NewAuthService(mockUserRepo, cfg)

	// Test data
	email := "nonexistent@example.com"
	password := "password123"

	// Mock expectations
	mockUserRepo.EXPECT().
		FindByEmail(mock.Anything, email).
		Return(nil, nil).
		Once()

	// Execute
	resultUser, jwtToken, csrfToken, err := authService.Login(context.Background(), email, password)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultUser)
	assert.Empty(t, jwtToken)
	assert.Empty(t, csrfToken)
	assert.Equal(t, "invalid email or password", err.Error())
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	// Setup
	mockUserRepo := mocks.NewMockUserRepository(t)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expiry: 15 * time.Minute,
		},
		CSRF: config.CSRFConfig{
			Secret: "csrf-secret",
		},
	}

	authService := NewAuthService(mockUserRepo, cfg)

	// Test data
	email := "test@example.com"
	correctPassword := "password123"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := util.HashPassword(correctPassword)

	user := &model.User{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	mockUserRepo.EXPECT().
		FindByEmail(mock.Anything, email).
		Return(user, nil).
		Once()

	// Execute
	resultUser, jwtToken, csrfToken, err := authService.Login(context.Background(), email, wrongPassword)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultUser)
	assert.Empty(t, jwtToken)
	assert.Empty(t, csrfToken)
	assert.Equal(t, "invalid email or password", err.Error())
}

func TestAuthService_Login_RepositoryError(t *testing.T) {
	// Setup
	mockUserRepo := mocks.NewMockUserRepository(t)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expiry: 15 * time.Minute,
		},
		CSRF: config.CSRFConfig{
			Secret: "csrf-secret",
		},
	}

	authService := NewAuthService(mockUserRepo, cfg)

	// Test data
	email := "test@example.com"
	password := "password123"

	// Mock expectations
	mockUserRepo.EXPECT().
		FindByEmail(mock.Anything, email).
		Return(nil, errors.New("database error")).
		Once()

	// Execute
	resultUser, jwtToken, csrfToken, err := authService.Login(context.Background(), email, password)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultUser)
	assert.Empty(t, jwtToken)
	assert.Empty(t, csrfToken)
	assert.Equal(t, "database error", err.Error())
}

func TestAuthService_Register_Success(t *testing.T) {
	// Setup
	mockUserRepo := mocks.NewMockUserRepository(t)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expiry: 15 * time.Minute,
		},
		CSRF: config.CSRFConfig{
			Secret: "csrf-secret",
		},
	}

	authService := NewAuthService(mockUserRepo, cfg)

	// Test data
	email := "newuser@example.com"
	password := "password123"

	// Mock expectations
	mockUserRepo.EXPECT().
		FindByEmail(mock.Anything, email).
		Return(nil, nil).
		Once()

	mockUserRepo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(user *model.User) bool {
			return user.Email == email && user.Password != password
		})).
		Return(nil).
		Once()

	// Execute
	resultUser, err := authService.Register(context.Background(), email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.Equal(t, email, resultUser.Email)
	assert.NotEqual(t, password, resultUser.Password) // Password should be hashed
}

func TestAuthService_Register_EmailExists(t *testing.T) {
	// Setup
	mockUserRepo := mocks.NewMockUserRepository(t)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expiry: 15 * time.Minute,
		},
		CSRF: config.CSRFConfig{
			Secret: "csrf-secret",
		},
	}

	authService := NewAuthService(mockUserRepo, cfg)

	// Test data
	email := "existing@example.com"
	password := "password123"

	existingUser := &model.User{
		ID:    primitive.NewObjectID(),
		Email: email,
	}

	// Mock expectations
	mockUserRepo.EXPECT().
		FindByEmail(mock.Anything, email).
		Return(existingUser, nil).
		Once()

	// Execute
	resultUser, err := authService.Register(context.Background(), email, password)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultUser)
	assert.Equal(t, "email already exists", err.Error())
}
