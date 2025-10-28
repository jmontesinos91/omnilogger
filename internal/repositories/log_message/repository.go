package log_message

import (
	"context"
)

// IRepository interface
type IRepository interface {
	FindByID(ctx context.Context, ID *int) (*Model, error)
	Create(ctx context.Context, model *Model) error
	Update(ctx context.Context, ID *int, model *Model) error
}
