package logs

import (
	"time"

	"github.com/jmontesinos91/omnilogger/internal/repositories/log_message"
	"github.com/uptrace/bun"
)

// Model Database model for logs
type Model struct {
	bun.BaseModel `bun:"table:logs"`

	ID          string               `bun:"id,pk"`
	IpAddress   string               `bun:"ip_address"`
	ClientHost  string               `bun:"client_host"`
	Provider    string               `bun:"provider"`
	Level       int                  `bun:"level"`
	Message     int                  `bun:"message"`
	Description string               `bun:"description"`
	Path        string               `bun:"path"`
	Resource    string               `bun:"resource"`
	Action      string               `bun:"action"`
	Data        string               `bun:"data"`
	OldData     string               `bun:"old_data"`
	TenantCat   string               `bun:"tenant_cat"`
	TenantID    string               `bun:"tenant_id"`
	UserID      string               `bun:"user_id"`
	Target      string               `bun:"target"`
	CreatedAt   *time.Time           `bun:"created_at"`
	LogMessage  []*log_message.Model `bun:"rel:has-many,join:message=id"`
}

type Filter struct {
	Message  []int
	Level    []string
	Provider []string
	Action   []string
	Path     string
	Resource string
	TenantID []int
	UserID   []string
	Target   []string
	StartAt  time.Time
	EndAt    time.Time
	From     int
	Size     int
}
