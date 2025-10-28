package logssvcmock

import (
	"context"
	"github.com/jmontesinos91/oevents/eventfactory"
	"github.com/jmontesinos91/omnilogger/internal/services/logs"
)

type IService struct {
	// Create
	CreateErr    error
	CreateRes    *logs.Response
	CreateCalled bool

	// GetByID
	GetByIDErr    error
	GetByIDRes    *logs.Response
	GetByIDCalled bool

	// Retrieve
	RetrieveErr    error
	RetrieveRes    *logs.PaginatedRes
	RetrieveCalled bool

	// CreateLogFromKafka (new)
	CreateLogFromKafkaErr     error
	CreateLogFromKafkaCalled  bool
	CreateLogFromKafkaPayload *eventfactory.LogCreatedPayload
}

func (m *IService) GetByID(ctx context.Context, id *string) (*logs.Response, error) {
	m.GetByIDCalled = true
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	return &logs.Response{ID: *id, Message: 1}, nil
}

func (m *IService) Create(ctx context.Context, payload *logs.Payload) (*logs.Response, error) {
	m.CreateCalled = true
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	// Return with an ID set to emulate persistence
	return &logs.Response{ID: "1", Message: payload.Message}, nil
}

func (m *IService) Retrieve(ctx context.Context, filter logs.Filter) (*logs.PaginatedRes, error) {
	m.RetrieveCalled = true
	if m.RetrieveErr != nil {
		return nil, m.RetrieveErr
	}
	if m.RetrieveRes != nil {
		return m.RetrieveRes, nil
	}
	// Default empty paginated response (valor cero) para pruebas.
	return &logs.PaginatedRes{}, nil
}

func (m *IService) CreateLogFromKafka(ctx context.Context, logCreated *eventfactory.LogCreatedPayload) error {
	m.CreateLogFromKafkaCalled = true
	m.CreateLogFromKafkaPayload = logCreated
	return m.CreateLogFromKafkaErr
}
