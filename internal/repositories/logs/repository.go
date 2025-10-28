package logs

import (
	"context"
)

// IRepository interface
type IRepository interface {
	FindByID(ctx context.Context, ID *string) (*Model, error)
	Create(ctx context.Context, model *Model) error
	Retrieve(ctx context.Context, filter Filter) ([]Model, int, error)
}
