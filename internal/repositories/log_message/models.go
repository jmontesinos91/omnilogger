package log_message

import (
	"github.com/uptrace/bun"
)

// Model Database model for log messages
type Model struct {
	bun.BaseModel `bun:"table:log_messages"`

	ID      int    `bun:"id,pk" json:"id"`
	Message string `bun:"message" json:"message"`
	Lang    string `bun:"lang" json:"lang"`
}

type Filter struct {
	ID   *int
	Lang string
	From int
	Size int
}
