package logs

import (
	"context"
	"github.com/jmontesinos91/oevents/eventfactory"
	"github.com/jmontesinos91/omnilogger/domains/pagination"
	"github.com/jmontesinos91/omnilogger/internal/repositories/log_message"
	"github.com/jmontesinos91/omnilogger/internal/repositories/logs"
	"testing"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/omnilogger/internal/repositories/logs/logsmock"
	"github.com/jmontesinos91/terrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	ctxLogger := logger.NewContextLogger("TestCreate", "debug", logger.TextFormat)

	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	type repositoryOpts struct {
		logsRepo     *logsmock.IRepository
		logsRepoFunc func() *logsmock.IRepository
	}

	type args struct {
		ctx     context.Context
		payload *Payload
	}

	type assertsParams struct {
		repositoryOpts
		args
		result *Response
		err    error
	}

	cases := []struct {
		name           string
		repositoryOpts repositoryOpts
		args           args
		err            bool
		asserts        func(*testing.T, assertsParams) bool
	}{
		{
			name: "Happy path",
			repositoryOpts: repositoryOpts{
				logsRepoFunc: func() *logsmock.IRepository {
					repoMock := &logsmock.IRepository{}
					repoMock.On("Create", mock.Anything, mock.Anything).Return(nil)
					return repoMock
				},
			},
			args: args{
				ctx: ctx,
				payload: &Payload{
					IpAddress:   "192.168.1.1",
					ClientHost:  "localhost",
					Provider:    "ExampleProvider",
					Level:       1,
					Message:     2,
					Description: "Test description",
					Path:        "/example",
					Resource:    "resource-path",
					Action:      "CREATE",
					Data:        `{"key": "value"}`,
					UserID:      "1234",
				},
			},
			err: false,
			asserts: func(t *testing.T, ap assertsParams) bool {
				return assert.NoError(t, ap.err) &&
					assert.NotNil(t, ap.result) &&
					ap.logsRepo.AssertCalled(t, "Create", mock.Anything, mock.Anything)
			},
		},
		{
			name: "Error storing data in repository",
			repositoryOpts: repositoryOpts{
				logsRepoFunc: func() *logsmock.IRepository {
					repoMock := &logsmock.IRepository{}
					repoMock.On("Create", mock.Anything, mock.Anything).
						Return(terrors.New(terrors.ErrInternalService, "DB Error", nil))
					return repoMock
				},
			},
			args: args{
				ctx: ctx,
				payload: &Payload{
					IpAddress:   "192.168.1.1",
					ClientHost:  "localhost",
					Provider:    "ExampleProvider",
					Level:       1,
					Message:     2,
					Description: "Test description",
					Path:        "/example",
					Resource:    "resource-path",
					Action:      "CREATE",
					Data:        `{"key": "value"}`,
					UserID:      "1234",
				},
			},
			err: true,
			asserts: func(t *testing.T, ap assertsParams) bool {
				var terr *terrors.Error
				return assert.ErrorAs(t, ap.err, &terr) &&
					assert.Equal(t, "Internal error service", terr.Message) &&
					assert.Nil(t, ap.result)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.repositoryOpts.logsRepoFunc != nil {
				tc.repositoryOpts.logsRepo = tc.repositoryOpts.logsRepoFunc()
			}

			service := NewDefaultService(ctxLogger, tc.repositoryOpts.logsRepo)
			result, err := service.Create(tc.args.ctx, tc.args.payload)

			assertsParams := assertsParams{
				repositoryOpts: tc.repositoryOpts,
				args:           tc.args,
				result:         result,
				err:            err,
			}

			if !tc.asserts(t, assertsParams) {
				t.Errorf("Assert error on test case: %s", tc.name)
			}
		})
	}
}

func TestGetByID(t *testing.T) {

	ctxLogger := logger.NewContextLogger("TestGetByID", "debug", logger.TextFormat)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	type repositoryOpts struct {
		logsRepo     *logsmock.IRepository
		logsRepoFunc func() *logsmock.IRepository
	}

	type args struct {
		ctx      context.Context
		ID       *string
		filter   Filter
		expected *Response
	}

	type assertsParams struct {
		repositoryOpts
		args
		result *Response
		err    error
	}

	cases := []struct {
		name           string
		repositoryOpts repositoryOpts
		args           args
		err            bool
		asserts        func(*testing.T, assertsParams) bool
	}{
		{
			name: "Happy path",
			repositoryOpts: repositoryOpts{
				logsRepoFunc: func() *logsmock.IRepository {
					repoMock := &logsmock.IRepository{}
					repoMock.On("FindByID", mock.Anything, mock.Anything, mock.Anything).
						Return(&logs.Model{
							ID:          "12345",
							IpAddress:   "192.168.1.1",
							ClientHost:  "localhost",
							Provider:    "ExampleProvider",
							Level:       1,
							Message:     2,
							Description: "Test description",
							Path:        "/example",
							Resource:    "resource-path",
							Action:      "GET",
							LogMessage: []*log_message.Model{
								{
									ID:      100,
									Message: "Example log message",
									Lang:    "en",
								},
								{
									ID:      101,
									Message: "Ejemplo de mensaje de registro",
									Lang:    "es",
								},
							},
						}, nil)
					return repoMock
				},
			},
			args: args{
				ctx: ctx,
				ID:  stringPtr("12345"),
				filter: Filter{
					Lang: "en",
				},
				expected: &Response{
					ID:          "12345",
					IpAddress:   "192.168.1.1",
					ClientHost:  "localhost",
					Provider:    "ExampleProvider",
					Level:       1,
					Message:     2,
					Description: "Test description",
					Path:        "/example",
					Resource:    "resource-path",
					Action:      "GET",
					LogMessage: &log_message.Model{
						ID:      100,
						Message: "Example log message",
						Lang:    "en",
					},
				},
			},
			err: false,
			asserts: func(t *testing.T, ap assertsParams) bool {
				return assert.NoError(t, ap.err) &&
					assert.NotNil(t, ap.result) &&
					assert.Equal(t, ap.result, ap.args.expected) &&
					assert.Equal(t, "12345", ap.result.ID) &&
					ap.logsRepo.AssertCalled(t, "FindByID", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name: "Empty ID",
			repositoryOpts: repositoryOpts{
				logsRepoFunc: func() *logsmock.IRepository {
					return &logsmock.IRepository{}
				},
			},
			args: args{
				ctx: ctx,
				ID:  nil,
			},
			err: true,
			asserts: func(t *testing.T, ap assertsParams) bool {
				return assert.Error(t, ap.err) &&
					assert.Nil(t, ap.result) &&
					assert.Equal(t, ap.result, ap.args.expected) &&
					assert.Contains(t, ap.err.Error(), terrors.ErrBadRequest) &&
					ap.logsRepo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name: "Happy path with different message filter",
			repositoryOpts: repositoryOpts{
				logsRepoFunc: func() *logsmock.IRepository {
					repoMock := &logsmock.IRepository{}
					repoMock.On("FindByID", mock.Anything, mock.Anything, mock.Anything).
						Return(&logs.Model{
							ID:          "12345",
							IpAddress:   "192.168.1.1",
							ClientHost:  "localhost",
							Provider:    "ExampleProvider",
							Level:       1,
							Message:     2,
							Description: "",
							Path:        "/example",
							Resource:    "resource-path",
							Action:      "GET",
							LogMessage: []*log_message.Model{
								{
									ID:      100,
									Message: "Example log message",
								},
								{
									ID:      101,
									Message: "Ejemplo de mensaje de registro",
								},
							},
						}, nil)
					return repoMock
				},
			},
			args: args{
				ctx: ctx,
				ID:  stringPtr("12345"),
				filter: Filter{
					Lang: "en",
				},
				expected: &Response{
					ID:          "12345",
					IpAddress:   "192.168.1.1",
					ClientHost:  "localhost",
					Provider:    "ExampleProvider",
					Level:       1,
					Message:     2,
					Description: "",
					Path:        "/example",
					Resource:    "resource-path",
					Action:      "GET",
					LogMessage:  (*log_message.Model)(nil),
				},
			},
			err: false,
			asserts: func(t *testing.T, ap assertsParams) bool {
				return assert.NoError(t, ap.err) &&
					assert.NotNil(t, ap.result) &&
					assert.Equal(t, ap.result, ap.args.expected) &&
					assert.Equal(t, "12345", ap.result.ID) &&
					ap.logsRepo.AssertCalled(t, "FindByID", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name: "Happy path with empty log message",
			repositoryOpts: repositoryOpts{
				logsRepoFunc: func() *logsmock.IRepository {
					repoMock := &logsmock.IRepository{}
					repoMock.On("FindByID", mock.Anything, mock.Anything, mock.Anything).
						Return(&logs.Model{
							ID:          "12345",
							IpAddress:   "192.168.1.1",
							ClientHost:  "localhost",
							Provider:    "ExampleProvider",
							Level:       1,
							Message:     2,
							Description: "",
							Path:        "/example",
							Resource:    "resource-path",
							Action:      "GET",
							LogMessage:  []*log_message.Model{},
						}, nil)
					return repoMock
				},
			},
			args: args{
				ctx: ctx,
				ID:  stringPtr("12345"),
				filter: Filter{
					Message: []int{2, 4},
				},
				expected: &Response{
					ID:          "12345",
					IpAddress:   "192.168.1.1",
					ClientHost:  "localhost",
					Provider:    "ExampleProvider",
					Level:       1,
					Message:     2,
					Description: "",
					Path:        "/example",
					Resource:    "resource-path",
					Action:      "GET",
					LogMessage:  (*log_message.Model)(nil),
				},
			},
			err: false,
			asserts: func(t *testing.T, ap assertsParams) bool {
				return assert.NoError(t, ap.err) &&
					assert.NotNil(t, ap.result) &&
					assert.Equal(t, ap.result, ap.args.expected) &&
					assert.Equal(t, "12345", ap.result.ID) &&
					ap.logsRepo.AssertCalled(t, "FindByID", mock.Anything, mock.Anything, mock.Anything)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.repositoryOpts.logsRepoFunc != nil {
				tc.repositoryOpts.logsRepo = tc.repositoryOpts.logsRepoFunc()
			}

			service := NewDefaultService(ctxLogger, tc.repositoryOpts.logsRepo)
			result, err := service.GetByID(tc.args.ctx, tc.args.ID, tc.args.filter)

			assertsParams := assertsParams{
				repositoryOpts: tc.repositoryOpts,
				args:           tc.args,
				result:         result,
				err:            err,
			}

			if !tc.asserts(t, assertsParams) {
				t.Errorf("Assert error on test case: %s", tc.name)
			}
		})
	}
}

func TestRetrieve(t *testing.T) {

	ctxLogger := logger.NewContextLogger("TestRetrieve", "debug", logger.TextFormat)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	type repositoryOpts struct {
		logsRepo     *logsmock.IRepository
		logsRepoFunc func() *logsmock.IRepository
	}

	type args struct {
		ctx    context.Context
		filter Filter
	}

	type assertsParams struct {
		repositoryOpts
		args
		result *PaginatedRes
		err    error
	}

	cases := []struct {
		name           string
		repositoryOpts repositoryOpts
		args           args
		err            bool
		asserts        func(*testing.T, assertsParams) bool
	}{
		{
			name: "Happy path",
			repositoryOpts: repositoryOpts{
				logsRepoFunc: func() *logsmock.IRepository {
					repoMock := &logsmock.IRepository{}
					repoMock.On("Retrieve", mock.Anything, mock.Anything).
						Return([]logs.Model{
							{
								ID:          "12345",
								IpAddress:   "192.168.1.1",
								ClientHost:  "localhost",
								Provider:    "ExampleProvider",
								Level:       1,
								Message:     2,
								Description: "Test description",
								Path:        "/example",
								Resource:    "resource-path",
								Action:      "RETRIEVE",
							},
						}, 1, nil)
					return repoMock
				},
			},
			args: args{
				ctx: ctx,
				filter: Filter{
					Level:    []string{"INFO", "ERROR"},
					Message:  []int{2, 4},
					Provider: []string{"ExampleProvider"},
					Action:   []string{"CREATE", "RETRIEVE"},
					Path:     "/example",
					Resource: "resource-path",
					TenantID: []int{1, 2},
					UserID:   []string{"1234", "5678"},
					Lang:     "en",
					StartAt:  time.Now().Add(-24 * time.Hour),
					EndAt:    time.Now(),
					Filter: pagination.Filter{
						Page:     1,
						Size:     10,
						Offset:   0,
						SortBy:   "created_at",
						SortDesc: true,
					},
				},
			},
			err: false,
			asserts: func(t *testing.T, ap assertsParams) bool {
				return assert.NoError(t, ap.err) &&
					assert.NotNil(t, ap.result) &&
					assert.Len(t, ap.result.Data, 1) &&
					assert.Equal(t, 1, ap.result.Total) &&
					ap.logsRepo.AssertCalled(t, "Retrieve", mock.Anything, mock.Anything)
			},
		},
		{
			name: "Error retrieving logs",
			repositoryOpts: repositoryOpts{
				logsRepoFunc: func() *logsmock.IRepository {
					repoMock := &logsmock.IRepository{}
					repoMock.On("Retrieve", mock.Anything, mock.Anything).
						Return(nil, 0, terrors.New(terrors.ErrInternalService, "Error retrieving logs", nil))
					return repoMock
				},
			},
			args: args{
				ctx: ctx,
				filter: Filter{
					Level:    []string{"ERROR"},
					Provider: []string{"ExampleProvider"},
					Filter: pagination.Filter{
						Page: 1,
						Size: 10,
					},
				},
			},
			err: true,
			asserts: func(t *testing.T, ap assertsParams) bool {
				var terr *terrors.Error
				return assert.ErrorAs(t, ap.err, &terr) &&
					assert.Nil(t, ap.result) &&
					assert.Contains(t, terr.Message, "Internal error service") &&
					ap.logsRepo.AssertCalled(t, "Retrieve", mock.Anything, mock.Anything)
			},
		},
		{
			name: "Empty filter",
			repositoryOpts: repositoryOpts{
				logsRepoFunc: func() *logsmock.IRepository {
					repoMock := &logsmock.IRepository{}
					repoMock.On("Retrieve", mock.Anything, mock.Anything).
						Return([]logs.Model{}, 0, nil)
					return repoMock
				},
			},
			args: args{
				ctx: ctx,
				filter: Filter{
					Filter: pagination.Filter{},
				},
			},
			err: false,
			asserts: func(t *testing.T, ap assertsParams) bool {
				return assert.NoError(t, ap.err) &&
					assert.NotNil(t, ap.result) &&
					assert.Len(t, ap.result.Data, 0) &&
					assert.Equal(t, 0, ap.result.Total) &&
					assert.Equal(t, 1, ap.result.Page) &&
					ap.logsRepo.AssertCalled(t, "Retrieve", mock.Anything, mock.Anything)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.repositoryOpts.logsRepoFunc != nil {
				tc.repositoryOpts.logsRepo = tc.repositoryOpts.logsRepoFunc()
			}

			service := NewDefaultService(ctxLogger, tc.repositoryOpts.logsRepo)
			result, err := service.Retrieve(tc.args.ctx, tc.args.filter)

			assertsParams := assertsParams{
				repositoryOpts: tc.repositoryOpts,
				args:           tc.args,
				result:         result,
				err:            err,
			}

			if !tc.asserts(t, assertsParams) {
				t.Errorf("Assert error on test case: %s", tc.name)
			}
		})
	}
}

func TestCreateLogFromKafka(t *testing.T) {

	ctx := context.Background()
	ctxLogger := logger.NewContextLogger("TestCreateLogFromKafka", "debug", logger.TextFormat)

	type repositoryOpts struct {
		logsRepo     *logsmock.IRepository
		logsRepoFunc func() *logsmock.IRepository
	}

	type args struct {
		ctx     context.Context
		payload *eventfactory.LogCreatedPayload
	}

	type assertsParams struct {
		repositoryOpts
		args
		err error
	}

	cases := []struct {
		name           string
		repositoryOpts repositoryOpts
		args           args
		err            bool
		asserts        func(*testing.T, error, assertsParams) bool
	}{
		{
			name: "Happy path",
			repositoryOpts: repositoryOpts{
				logsRepoFunc: func() *logsmock.IRepository {
					repoMock := &logsmock.IRepository{}
					repoMock.On("Create", mock.Anything, mock.Anything).Return(nil)
					return repoMock
				},
			},
			args: args{
				ctx: ctx,
				payload: &eventfactory.LogCreatedPayload{
					IpAddress:   "192.168.1.1",
					ClientHost:  "localhost",
					Provider:    "ExampleProvider",
					Level:       1,
					Message:     2,
					Description: "Test description",
					Resource:    "resource-path",
					Path:        "/example",
					Action:      "CREATE",
					Data:        `{"key": "value"}`,
					OldData:     `{"old_key": "old_value"}`,
					UserID:      "1234",
					TenantCat: []eventfactory.TenantItem{
						{ID: 1, Name: "Tenant A"},
						{ID: 2, Name: "Tenant B"},
					},
				},
			},
			err: false,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				return assert.NoError(t, err) &&
					ap.logsRepo.AssertCalled(t, "Create", mock.Anything, mock.Anything)
			},
		},
		{
			name: "Error on repository Create",
			repositoryOpts: repositoryOpts{
				logsRepoFunc: func() *logsmock.IRepository {
					repoMock := &logsmock.IRepository{}
					repoMock.On("Create", mock.Anything, mock.Anything).
						Return(terrors.InternalService("metadata_error", "Error storing model for log", nil))
					return repoMock
				},
			},
			args: args{
				ctx: ctx,
				payload: &eventfactory.LogCreatedPayload{
					IpAddress:   "192.168.1.1",
					ClientHost:  "localhost",
					Provider:    "ExampleProvider",
					Level:       1,
					Message:     2,
					Description: "Test description",
					Resource:    "resource-path",
					Path:        "/example",
					Action:      "CREATE",
					Data:        `{"key": "value"}`,
					OldData:     `{"old_key": "old_value"}`,
					UserID:      "1234",
					TenantCat: []eventfactory.TenantItem{
						{ID: 1, Name: "Tenant A"},
					},
				},
			},
			err: true,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				var terr *terrors.Error
				return assert.ErrorAs(t, err, &terr) &&
					assert.Equal(t, "Error storing model for log", terr.Message)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			if tc.repositoryOpts.logsRepoFunc != nil {
				tc.repositoryOpts.logsRepo = tc.repositoryOpts.logsRepoFunc()
			}

			service := NewDefaultService(ctxLogger, tc.repositoryOpts.logsRepo)
			err := service.CreateLogFromKafka(tc.args.ctx, tc.args.payload)

			assertsParams := assertsParams{
				repositoryOpts: tc.repositoryOpts,
				args:           tc.args,
				err:            err,
			}

			if !tc.asserts(t, err, assertsParams) {
				t.Errorf("Assert error on test case: %s", tc.name)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
