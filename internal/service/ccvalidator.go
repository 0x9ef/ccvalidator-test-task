package service

import (
	"context"

	"github.com/0x9ef/card-validator/config"

	"github.com/0x9ef/card-validator/pkg/ccvalidator"
)

type ccValidatorService struct {
	cfg *config.Config
}

func NewCreditCardValidator(cfg *config.Config) *ccValidatorService {
	return &ccValidatorService{cfg: cfg}
}

func (s *ccValidatorService) Validate(_ context.Context, number string, year int, month int) error {
	// We don't use logger here as far as we dont log any sensitive data
	// ...

	// We don't require user to provide CVV, so we are going to use mocked value for all cards
	const mockCVV = 111
	return ccvalidator.
		NewLuhnValidator(s.cfg.Validate.AllowTestCards).
		Validate(number, mockCVV, month, year)
}
