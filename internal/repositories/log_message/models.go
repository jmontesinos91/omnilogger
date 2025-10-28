package log_message

import (
	"github.com/uptrace/bun"
)

// Model Database model for log messages
type Model struct {
	bun.BaseModel `bun:"table:log_messages"`

	ID      int    `bun:"id,pk" json:"id"`
	Message string `bun:"message" json:"message"`
}
