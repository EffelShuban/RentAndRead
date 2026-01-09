package service

import (
	"context"
	"errors"
	"rent-and-read/internal/dto"
	"rent-and-read/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Register(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		req := dto.RegisterRequest{
			Email:    "test@example.com",
			Name:     "Test User",
			Password: "password123",
		}
		mockRepo.On("Create", ctx, mock.MatchedBy(func(u models.User) bool {
			return u.Email == req.Email && u.Name == req.Name && u.Password != ""
		})).Return(nil).Once()

		err := service.Register(ctx, req)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ValidationError", func(t *testing.T) {
		req := dto.RegisterRequest{Email: "invalid-email"}
		err := service.Register(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
	})
}

func TestUserService_Login(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	validUser := models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	t.Run("Success", func(t *testing.T) {
		req := dto.LoginRequest{Email: "test@example.com", Password: "password123"}
		mockRepo.On("FindByEmail", ctx, req.Email).Return(validUser, nil).Once()

		token, err := service.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		req := dto.LoginRequest{Email: "unknown@example.com", Password: "password123"}
		mockRepo.On("FindByEmail", ctx, req.Email).Return(models.User{}, models.ErrUserNotFound).Once()

		token, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, models.ErrUnauthorized, err)
		assert.Empty(t, token)
	})

	t.Run("WrongPassword", func(t *testing.T) {
		req := dto.LoginRequest{Email: "test@example.com", Password: "wrongpassword"}
		mockRepo.On("FindByEmail", ctx, req.Email).Return(validUser, nil).Once()

		token, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, models.ErrUnauthorized, err)
		assert.Empty(t, token)
	})
}

func TestUserService_GetBalance(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := NewUserService(mockRepo)
	ctx := context.Background()
	userID := uint(1)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, userID).Return(models.User{ID: userID, Balance: 50000}, nil).Once()

		balance, err := service.GetBalance(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, 50000, balance)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, userID).Return(models.User{}, errors.New("db error")).Once()

		balance, err := service.GetBalance(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, 0, balance)
	})
}

func TestUserService_ProcessPayment(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := NewUserService(mockRepo)
	ctx := context.Background()
	extID := "TOPUP-1-12345"

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("HandlePaymentSuccess", ctx, extID).Return(nil).Once()
		err := service.ProcessPayment(ctx, extID)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.On("HandlePaymentSuccess", ctx, extID).Return(errors.New("db error")).Once()
		err := service.ProcessPayment(ctx, extID)
		assert.Error(t, err)
	})
}