package service

import (
	"context"
	"errors"
)

type Services struct {
	CardValidatorService
}

type CardValidatorService interface {
	Validate(ctx context.Context, number string, year int, month int) error
}

var (
	ErrValidateCardInvalidCardNumber = errors.New("Invalid card number")
	ErrValidateCardInvalidCVV        = errors.New("Invalid CVV")
	ErrValidateCardInvalidExpYear    = errors.New("Invalid expiration year.")
	ErrValidateCardInvalidExpMonth   = errors.New("Invalid expiration month.")
	ErrValidateCardCardExpired       = errors.New("Card is expired.")
)
