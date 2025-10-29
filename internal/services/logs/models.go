package logs

import (
	"time"

	"github.com/jmontesinos91/omnilogger/domains/pagination"
)

// Payload payload example
type Payload struct {
	IpAddress   string `json:"ip_address" validate:"required"`
	ClientHost  string `json:"client_host" validate:"required"`
	Provider    string `json:"provider" validate:"required"`
	Level       int    `json:"level" validate:"required"`
	Message     int    `json:"message" validate:"required"`
	Description string `json:"description"`
	Path        string `json:"path" validate:"required"`
	Resource    string `json:"resource" validate:"required"`
	Action      string `json:"action" validate:"required"`
	Data        string `json:"data" validate:"required"`
	OldData     string `json:"old_data" validate:"required"`
	TenantCat   string `json:"tenant_cat"`
	UserID      string `json:"user_id" validate:"required"`
	Target      string `json:"target" validate:"required"`
}

// Response Holds the response for a created payout
type Response struct {
	ID          string      `json:"id"`
	IpAddress   string      `json:"ipAddress"`
	ClientHost  string      `json:"clientHost"`
	Provider    string      `json:"provider"`
	Level       int         `json:"level"`
	Message     int         `json:"message"`
	Description string      `json:"description"`
	Path        string      `json:"path"`
	Resource    string      `json:"resource"`
	Action      string      `json:"action"`
	Data        string      `json:"data"`
	OldData     string      `json:"oldData"`
	TenantCat   string      `json:"tenantCat"`
	UserID      string      `json:"userId"`
	Target      string      `json:"target"`
	CreatedAt   *time.Time  `json:"createdAt,omitempty"`
	LogMessage  interface{} `json:"logMessage"`
}

type Filter struct {
	Level    []string
	Message  []int
	Provider []string
	Action   []string
	Path     string
	Resource string
	TenantID []int
	UserID   []string
	Target   []string
	StartAt  time.Time
	EndAt    time.Time
	pagination.Filter
}

type PaginatedRes struct {
	Data  []Response `json:"data"`
	Size  int        `json:"max"`
	Total int        `json:"total"`
	Page  int        `json:"currentPage"`
}
