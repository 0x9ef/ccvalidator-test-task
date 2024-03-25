package http

import (
	"errors"

	"github.com/0x9ef/card-validator/internal/service"
	"github.com/gin-gonic/gin"
)

type credicCardValidatorRoutes struct {
	RouterContext
}

func newCreditCardValidatorRoutes(options *RouterOptions) {
	r := &credicCardValidatorRoutes{
		RouterContext{
			services: options.Services,
			logger:   options.Logger.Named("creditCardValidatorRoutes"),
		},
	}

	p := options.Handler.Group("/validate")
	{
		p.POST("", errorHandlerMiddleware(options, r.validate))
	}
}

const (
	ErrCodeValidateInvalidExpMonth ErrorCode = "C01"
	ErrCodeValidateInvalidExpYear  ErrorCode = "C02"
	ErrCodeValidateInvalidNumber   ErrorCode = "C03"
	ErrCodeValidateCardExpired     ErrorCode = "C04"
	ErrCodeValidateUndefined       ErrorCode = "X00"
)

type validateRequestBody struct {
	CardNumber string `json:"cn"`
	ExpMonth   int    `json:"expm"`
	ExpYear    int    `json:"expy"`
}

type validateRequestResponse struct {
	IsValid bool       `json:"valid"`
	Error   *httpError `json:"error,omitempty"`
}

func (r *credicCardValidatorRoutes) validate(c *gin.Context) (interface{}, *httpError) {
	logger := r.logger.
		Named("validate").
		WithContext(c, "RequestID")

	var body validateRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Info("cannot parse JSON body", "err", err)
		return nil, &httpError{Type: ErrorTypeClient, Message: "invalid request body"}
	}
	logger.Debug("request body parsed successfully")

	var httpErr *httpError
	err := r.services.CardValidatorService.Validate(c, body.CardNumber, body.ExpYear, body.ExpMonth)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrValidateCardCardExpired):
			httpErr = &httpError{Type: ErrorTypeClient, Code: ErrCodeValidateCardExpired, Message: err.Error()}
			logger.Info("card has expired")
		case errors.Is(err, service.ErrValidateCardInvalidExpMonth):
			httpErr = &httpError{Type: ErrorTypeClient, Code: ErrCodeValidateInvalidExpMonth, Message: err.Error()}
			logger.Info("invalid expiration month")
		case errors.Is(err, service.ErrValidateCardInvalidExpYear):
			httpErr = &httpError{Type: ErrorTypeClient, Code: ErrCodeValidateInvalidExpYear, Message: err.Error()}
			logger.Info("invalid expiration year")
		default:
			httpErr = &httpError{Type: ErrorTypeServer, Code: ErrCodeValidateUndefined, Message: err.Error()}
			// log unexpected errors with ERROR severity, otherwise - INFO
			logger.Error("failed to validate credit card", "err", err)
		}
	}

	logger.Info("successfully validated credit card information")
	return &validateRequestResponse{
		IsValid: err == nil,
		Error:   httpErr,
	}, nil
}
