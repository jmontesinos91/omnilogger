package logssvcmock

import (
	"context"
	"github.com/jmontesinos91/omnilogger/internal/services/logs"
)

type IService struct {
	GetErr         error
	CreateErr      error
	RetrieveErr    error
	GetCalled      bool
	CreateCalled   bool
	RetrieveCalled bool
	RetrieveRes    *logs.PaginatedRes
}

func (m *IService) GetByID(ctx context.Context, id *string) (*logs.Response, error) {
	m.GetCalled = true
	if m.GetErr != nil {
		return nil, m.GetErr
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
