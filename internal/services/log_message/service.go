package log_message

import (
	"context"
)

// IService Manage log message interfaces
type IService interface {
	GetByID(ctx context.Context, id *int) (*Response, error)
	Create(ctx context.Context, payload *Payload) (*Response, error)
	Update(ctx context.Context, id *int, payload *Payload) (*Response, error)
}
