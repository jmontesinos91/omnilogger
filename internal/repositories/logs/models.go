package logs

import (
	"time"

	"github.com/uptrace/bun"
)

// Model Database model for logs
type Model struct {
	bun.BaseModel `bun:"table:logs"`

	ID          string     `bun:"id,pk"`
	IpAddress   string     `bun:"ip_address"`
	ClientHost  string     `bun:"client_host"`
	Provider    string     `bun:"provider"`
	Level       int        `bun:"level"`
	Message     int        `bun:"message"`
	Description string     `bun:"description"`
	Path        string     `bun:"path"`
	Resource    string     `bun:"resource"`
	Action      string     `bun:"action"`
	Data        string     `bun:"data"`
	UserID      string     `bun:"user_id"`
	CreatedAt   *time.Time `bun:"created_at"`
}
