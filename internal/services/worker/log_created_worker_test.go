package worker

import (
	"context"
	"testing"

	"github.com/jmontesinos91/oevents"
	"github.com/jmontesinos91/oevents/broker/brokermock"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/omnilogger/internal/repositories/logs/logsmock"
	"github.com/jmontesinos91/omnilogger/internal/services/logs"
	"github.com/jmontesinos91/terrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogCreatedWorker_Handle(t *testing.T) {
	ctxLogger := logger.NewContextLogger("OMNILOGGER", "test", logger.TextFormat)
	ctx := context.Background()

	type fields struct {
		logsRepo     *logsmock.IRepository
		streamClient *brokermock.MessagingBrokerProvider
	}
	type args struct {
		event oevents.OmniViewEvent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		asserts func(*testing.T, error, *logsmock.IRepository)
	}{
		{
			name: "Happy Path Log Created Worker - Success",
			fields: fields{
				logsRepo: func() *logsmock.IRepository {
					repoMock := new(logsmock.IRepository)
					repoMock.On("Create", mock.Anything, mock.Anything).Return(nil)
					return repoMock
				}(),
			},
			args: args{
				event: oevents.OmniViewEvent{
					ID: "12345",
					Data: map[string]any{
						"IpAddress":   "192.168.1.1",
						"ClientHost":  "localhost",
						"Provider":    "example",
						"Level":       1,
						"Message":     1,
						"Description": "This is a test",
						"Resource":    "test-resource",
						"Path":        "/test",
						"Action":      "CREATE",
						"Data":        `{"key": "value"}`,
						"OldData":     `{"old_key": "old_value"}`,
						"UserID":      "12345",
						"TenantCat": []map[string]interface{}{
							{"id": 1, "name": "Tenant A"},
						},
					},
				},
			},
			wantErr: false,
			asserts: func(t *testing.T, err error, repoMock *logsmock.IRepository) {
				assert.NoError(t, err)
				repoMock.AssertCalled(t, "Create", mock.Anything, mock.Anything)
			},
		},
		{
			name: "Log created without TenantCat",
			fields: fields{
				logsRepo: func() *logsmock.IRepository {
					repoMock := new(logsmock.IRepository)
					repoMock.On("Create", mock.Anything, mock.Anything).Return(nil)
					return repoMock
				}(),
			},
			args: args{
				event: oevents.OmniViewEvent{
					ID: "12345",
					Data: map[string]any{
						"IpAddress":  "192.168.1.1",
						"ClientHost": "localhost",
						"Provider":   "example",
						"Level":      1,
						"Message":    1,
						"UserID":     "12345",
					},
				},
			},
			wantErr: false,
			asserts: func(t *testing.T, err error, repoMock *logsmock.IRepository) {
				assert.NoError(t, err)
				repoMock.AssertCalled(t, "Create", mock.Anything, mock.Anything)
			},
		},
		{
			name: "Log created with multiple TenantCat",
			fields: fields{
				logsRepo: func() *logsmock.IRepository {
					repoMock := new(logsmock.IRepository)
					repoMock.On("Create", mock.Anything, mock.Anything).Return(nil)
					return repoMock
				}(),
			},
			args: args{
				event: oevents.OmniViewEvent{
					ID: "12345",
					Data: map[string]any{
						"IpAddress":  "192.168.1.1",
						"ClientHost": "localhost",
						"Provider":   "example",
						"Level":      1,
						"Message":    1,
						"UserID":     "12345",
						"TenantCat": []map[string]interface{}{
							{"id": 1, "name": "Tenant A"},
							{"id": 2, "name": "Tenant B"},
							{"id": 3, "name": "Tenant C"},
							{"id": 4, "name": "Tenant D"},
						},
					},
				},
			},
			wantErr: false,
			asserts: func(t *testing.T, err error, repoMock *logsmock.IRepository) {
				assert.NoError(t, err)
				repoMock.AssertCalled(t, "Create", mock.Anything, mock.Anything)
			},
		},
		{
			name: "Error Serializing Event Data",
			fields: fields{
				logsRepo: func() *logsmock.IRepository {
					repoMock := new(logsmock.IRepository)
					repoMock.On("Create", mock.Anything, mock.Anything).Return(terrors.InternalService("tenant_cat_error", "Error marshalling TenantCat to JSON", nil))
					return repoMock
				}(),
			},
			args: args{
				event: oevents.OmniViewEvent{
					ID: "12345",
					Data: map[string]any{
						"IpAddress":   "192.168.1.1",
						"ClientHost":  "localhost",
						"Provider":    "example",
						"Level":       1,
						"Message":     1,
						"Description": "This is a test",
						"Resource":    "test-resource",
						"Path":        "/test",
						"Action":      "CREATE",
						"Data":        `{"key": "value"}`,
						"OldData":     `{"old_key": "old_value"}`,
						"UserID":      "12345",
						"TenantCat": []map[string]interface{}{
							{"id": 1, "name": "Tenant A"},
						},
					},
				},
			},
			wantErr: true,
			asserts: func(t *testing.T, err error, repoMock *logsmock.IRepository) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "internal_service.metadata_error: Error storing model for log")
				repoMock.AssertCalled(t, "Create", mock.Anything, mock.Anything)
			},
		},
		{
			name: "Error in ToLogCreatedPayload",
			fields: fields{
				logsRepo: func() *logsmock.IRepository {
					repoMock := new(logsmock.IRepository)
					repoMock.On("Create", mock.Anything, mock.Anything).Return(nil)
					return repoMock
				}(),
			},
			args: args{
				event: oevents.OmniViewEvent{
					ID: "12345",
					Data: map[string]any{
						"InvalidField": "invalid",
					},
				},
			},
			wantErr: false,
			asserts: func(t *testing.T, err error, repoMock *logsmock.IRepository) {
				assert.NoError(t, err)
				repoMock.AssertCalled(t, "Create", mock.Anything, mock.Anything)
			},
		},
		{
			name: "Error creating log in repository",
			fields: fields{
				logsRepo: func() *logsmock.IRepository {
					repoMock := new(logsmock.IRepository)
					repoMock.On("Create", mock.Anything, mock.Anything).
						Return(terrors.InternalService("metadata_error", "Failed to store log", nil))
					return repoMock
				}(),
			},
			args: args{
				event: oevents.OmniViewEvent{
					ID: "12345",
					Data: map[string]any{
						"IpAddress":  "192.168.1.1",
						"ClientHost": "localhost",
						"Provider":   "example",
						"Level":      1,
						"Message":    1,
						"UserID":     "12345",
						"TenantCat": []map[string]interface{}{
							{"id": 1, "name": "Tenant A"},
						},
					},
				},
			},
			wantErr: true,
			asserts: func(t *testing.T, err error, repoMock *logsmock.IRepository) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "internal_service.metadata_error: Error storing model for log")
				repoMock.AssertCalled(t, "Create", mock.Anything, mock.Anything)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logSvc := logs.NewDefaultService(ctxLogger, tt.fields.logsRepo)
			worker := NewLogCreatedWorker(ctxLogger, logSvc, tt.fields.streamClient)

			err := worker.Handle(ctx, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogCreatedWorker.Handle() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.asserts != nil {
				tt.asserts(t, err, tt.fields.logsRepo)
			}
		})
	}
}
