package logs

import (
	"context"
	"github.com/jmontesinos91/omnilogger/internal/repositories/logs"
	"testing"

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
		ctx context.Context
		ID  *string
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
						}, nil)
					return repoMock
				},
			},
			args: args{
				ctx: ctx,
				ID:  stringPtr("12345"),
			},
			err: false,
			asserts: func(t *testing.T, ap assertsParams) bool {
				return assert.NoError(t, ap.err) &&
					assert.NotNil(t, ap.result) &&
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
					assert.Contains(t, ap.err.Error(), terrors.ErrBadRequest) &&
					ap.logsRepo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything, mock.Anything)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.repositoryOpts.logsRepoFunc != nil {
				tc.repositoryOpts.logsRepo = tc.repositoryOpts.logsRepoFunc()
			}

			service := NewDefaultService(ctxLogger, tc.repositoryOpts.logsRepo)
			result, err := service.GetByID(tc.args.ctx, tc.args.ID)

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

func stringPtr(s string) *string {
	return &s
}
