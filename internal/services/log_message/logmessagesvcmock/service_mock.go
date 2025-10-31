package logmessagesvcmock

import (
	"context"
	"github.com/jmontesinos91/omnilogger/internal/services/log_message"
)

type IService struct {
	GetErr        error
	CreateErr     error
	UpdateErr     error
	RetrieveErr   error
	DeleteLangErr error
}

func (m *IService) GetByID(ctx context.Context, id *int) (*log_message.Response, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return &log_message.Response{ID: *id, Message: "ok"}, nil
}

func (m *IService) Create(ctx context.Context, payload *log_message.Payload) (*log_message.Response, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	return &log_message.Response{ID: 1, Message: "created"}, nil
}

func (m *IService) Update(ctx context.Context, id *int, lang string, payload *log_message.Payload) (*log_message.Response, error) {
	if m.UpdateErr != nil {
		return nil, m.UpdateErr
	}
	return &log_message.Response{ID: *id, Message: "updated"}, nil
}

func (m *IService) Retrieve(ctx context.Context, filter log_message.Filter) (*log_message.PaginatedRes, error) {
	if m.RetrieveErr != nil {
		return nil, m.RetrieveErr
	}
	return &log_message.PaginatedRes{Data: []log_message.Response{{ID: 1, Message: "updated"}}}, nil
}

func (m *IService) DeleteLang(ctx context.Context, id *int, lang string) error {
	if m.RetrieveErr != nil {
		return m.RetrieveErr
	}
	return nil
}
