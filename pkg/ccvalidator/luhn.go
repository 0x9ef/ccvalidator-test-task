package ccvalidator

import (
	"errors"
	"strconv"
	"strings"

	creditcard "github.com/durango/go-credit-card"
)

type luhnValidator struct {
	allowTestCards bool
}

var _ Validator = (*luhnValidator)(nil)

func NewLuhnValidator(allowTestCards bool) *luhnValidator {
	return &luhnValidator{allowTestCards: allowTestCards}
}

var (
	ErrInvalidCardNumber = errors.New("Invalid card number")
	ErrInvalidCVV        = errors.New("Invalid CVV")
	ErrInvalidExpYear    = errors.New("Invalid expiration year.")
	ErrInvalidExpMonth   = errors.New("Invalid expiration month.")
	ErrCardExpired       = errors.New("Card is expired.")
)

func (v *luhnValidator) Validate(number string, cvv int, month int, year int) error {
	card := creditcard.Card{
		// Handle possible cases when spaces
		Number: strings.Replace(number, " ", "", -1),
		Cvv:    strconv.Itoa(cvv),
		Month:  strconv.Itoa(month),
		Year:   strconv.Itoa(year),
	}

	// As far as go-credit-card doesn't have predefine error variables, we have to
	// manually switch on the error message. I fully understand that it is the bad practice,
	// because the error message can be changed furthure. The ideal solution would be to fork the repo and
	// provide own error variables like ErrInvalidCardNumber, ErrInvalidCvv, etc...
	err := card.Validate(v.allowTestCards)
	if err != nil {
		switch strings.ToLower(err.Error()) {
		case "invalid year":
			return ErrInvalidExpYear
		case "invalid month":
			return ErrInvalidExpMonth
		case "invalid cvv":
			return ErrInvalidCVV
		case "credit card has expired":
			return ErrCardExpired
		case "invalid credit card number":
			return ErrInvalidCardNumber
		default:
			return err
		}
	}
	return nil
}
