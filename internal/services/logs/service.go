package logs

import (
	"context"
	"github.com/jmontesinos91/oevents/eventfactory"
)

// IService Manage log interfaces
type IService interface {
	Create(ctx context.Context, payload *Payload) (*Response, error)
	GetByID(ctx context.Context, id *string, filter Filter) (*Response, error)
	Retrieve(ctx context.Context, filter Filter) (*PaginatedRes, error)
	CreateLogFromKafka(ctx context.Context, logCreated *eventfactory.LogCreatedPayload) error
}
