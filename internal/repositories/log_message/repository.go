package log_message

import (
	"context"
)

// IRepository interface
type IRepository interface {
	FindByID(ctx context.Context, ID *int) (*Model, error)
	FindByIDAndLang(ctx context.Context, ID *int, lang string) (*Model, error)
	Create(ctx context.Context, model *Model) error
	Update(ctx context.Context, ID *int, lang string, model *Model) error
	Retrieve(ctx context.Context, filter Filter) ([]Model, int, error)
}
