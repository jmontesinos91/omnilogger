package log_message

import (
	"github.com/go-playground/validator/v10"
)

// Payload payload example
type Payload struct {
	ID      int    `json:"id" validate:"required"`
	Message string `json:"message" validate:"required"`
}

// Response Holds the response for a created payout
type Response struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

type Filter struct {
	ID *int
}

func (r *Payload) SanitizeAndValidate(validate *validator.Validate) error {
	return validate.Struct(r)
}
