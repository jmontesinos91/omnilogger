package log_message

import (
	"github.com/go-playground/validator/v10"
	"github.com/jmontesinos91/omnilogger/domains/pagination"
)

// Payload payload example
type Payload struct {
	ID      int    `json:"id" validate:"required"`
	Message string `json:"message" validate:"required"`
	Lang    string `json:"lang" validate:"required"`
}

// Response Holds the response for a created payout
type Response struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	Lang    string `json:"lang"`
}

type Filter struct {
	ID   *int
	Lang string
	pagination.Filter
}

type PaginatedRes struct {
	Data  []Response `json:"data"`
	From  int        `json:"from"`
	Size  int        `json:"size"`
	Total int        `json:"total"`
	Page  int        `json:"current_page"`
}

func (r *Payload) SanitizeAndValidate(validate *validator.Validate) error {
	return validate.Struct(r)
}
