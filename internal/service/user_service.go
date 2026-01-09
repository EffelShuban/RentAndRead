package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"rent-and-read/config"
	"rent-and-read/internal/dto"
	"rent-and-read/internal/models"
	"rent-and-read/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface{
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (string, error)
	ProcessPayment(ctx context.Context, externalID string) error
	GetBalance(ctx context.Context, userID uint) (int, error)
}

type userService struct{
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService{
	return &userService{
		repo: repo,
	}
}

func (s *userService) ProcessPayment(ctx context.Context, externalID string) error {
    return s.repo.HandlePaymentSuccess(ctx, externalID)
}

func (s *userService) GetBalance(ctx context.Context, userID uint) (int, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return user.Balance, nil
}

func (s *userService) Register(ctx context.Context, req dto.RegisterRequest) error{
	err := validateRegisterRequest(req)
	if err != nil{
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil{
		return err
	}

	u := models.User{
		Email: req.Email,
		Name: req.Name,
		Password: string(hashed),
	}

	return s.repo.Create(ctx, u)
}

func (s *userService) Login(ctx context.Context, req dto.LoginRequest) (string, error){
	err := validateLoginRequest(req)
	if err != nil{
		return "", models.ErrBadUserRequest
	}

	u, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil{
		if errors.Is(err, models.ErrUserNotFound) {
			return "", models.ErrUnauthorized
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil{
		return "", models.ErrUnauthorized
	}

	token, err := config.GenerateToken(u.ID)
	if err != nil{
		return "", err
	}

	return token, nil
}

func validateRegisterRequest(req dto.RegisterRequest)error{
	if req.Email== ""{
		return fmt.Errorf("%v: email is required", models.ErrBadUserRequest)
	}

	if _, err := mail.ParseAddress(req.Email); err != nil{
		return fmt.Errorf("%v: invalid email format", models.ErrBadUserRequest)
	}

	if req.Password == ""{
		return fmt.Errorf("%v: password is required", models.ErrBadUserRequest)
	}

	if req.Name == ""{
		return fmt.Errorf("%v: name is required", models.ErrBadUserRequest)
	}

	return nil
}

func validateLoginRequest(req dto.LoginRequest) error{
	if req.Email== ""{
		return fmt.Errorf("%v: email is required", models.ErrBadUserRequest)
	}

	if _, err := mail.ParseAddress(req.Email); err != nil{
		return fmt.Errorf("%v: invalid email format", models.ErrBadUserRequest)
	}

	if req.Password == ""{
		return fmt.Errorf("%v: password is required", models.ErrBadUserRequest)
	}

	return nil
}