package models

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExist = errors.New("user already exist")
	ErrUnauthorized = errors.New("invalid email/password")
	ErrBadUserRequest = errors.New("bad user request")
	ErrBookOutOfStock = errors.New("book out of stock")
	ErrInsufficientBalance = errors.New("insuficient wallet balance")
)