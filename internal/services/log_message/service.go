package log_message

import (
	"context"
)

// IService Manage log message interfaces
type IService interface {
	GetByID(ctx context.Context, id *int) (*Response, error)
	Create(ctx context.Context, payload *Payload) (*Response, error)
	Update(ctx context.Context, id *int, lang string, payload *Payload) (*Response, error)
	Retrieve(ctx context.Context, filter Filter) (*PaginatedRes, error)
}
