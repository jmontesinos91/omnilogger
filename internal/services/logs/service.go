package logs

import (
	"context"
)

// IService Manage log interfaces
type IService interface {
	Create(ctx context.Context, payload *Payload) (*Response, error)
	GetByID(ctx context.Context, id *string) (*Response, error)
}
