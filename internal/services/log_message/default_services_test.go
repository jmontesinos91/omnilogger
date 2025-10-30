package log_message

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/omnilogger/internal/repositories/log_message"
	"github.com/jmontesinos91/omnilogger/internal/repositories/log_message/logsmessagemock"
	"github.com/jmontesinos91/terrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetByID(t *testing.T) {
	id := 12345
	ctxLogger := logger.NewContextLogger("TestGetByID", "debug", logger.TextFormat)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	type repositoryOpts struct {
		logsRepo     *logsmessagemock.IRepository
		logsRepoFunc func() *logsmessagemock.IRepository
	}
	type args struct {
		ctx       context.Context
		validator *validator.Validate
		ID        *int
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
				logsRepoFunc: func() *logsmessagemock.IRepository {
					repoMock := &logsmessagemock.IRepository{}
					repoMock.On("FindByID", mock.Anything, mock.Anything).
						Return(&log_message.Model{
							ID:      12345,
							Message: "Test description",
						}, nil)
					return repoMock
				},
			},
			args: args{
				ctx:       ctx,
				ID:        &id,
				validator: validator.New(),
			},
			err: false,
			asserts: func(t *testing.T, ap assertsParams) bool {
				return assert.NoError(t, ap.err) &&
					assert.NotNil(t, ap.result) &&
					assert.Equal(t, 12345, ap.result.ID) &&
					ap.logsRepo.AssertCalled(t, "FindByID", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name: "Empty ID",
			repositoryOpts: repositoryOpts{
				logsRepoFunc: func() *logsmessagemock.IRepository {
					return &logsmessagemock.IRepository{}
				},
			},
			args: args{
				ctx:       ctx,
				ID:        nil,
				validator: validator.New(),
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

			service := NewDefaultService(ctxLogger, tc.args.validator, tc.repositoryOpts.logsRepo)
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

func TestCreate(t *testing.T) {
	ctxLogger := logger.NewContextLogger("TestCreate", "debug", logger.TextFormat)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")
	type repositoryOpts struct {
		logsRepo     *logsmessagemock.IRepository
		logsRepoFunc func() *logsmessagemock.IRepository
	}
	type args struct {
		ctx       context.Context
		payload   *Payload
		validator *validator.Validate
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
				logsRepoFunc: func() *logsmessagemock.IRepository {
					repoMock := &logsmessagemock.IRepository{}
					repoMock.On("Create", mock.Anything, mock.Anything).Return(nil)
					return repoMock
				},
			},
			args: args{
				ctx: ctx,
				payload: &Payload{
					Message: "Test description",
				},
				validator: validator.New(),
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
				logsRepoFunc: func() *logsmessagemock.IRepository {
					repoMock := &logsmessagemock.IRepository{}
					repoMock.On("Create", mock.Anything, mock.Anything).
						Return(terrors.New(terrors.ErrInternalService, "DB Error", nil))
					return repoMock
				},
			},
			args: args{
				ctx: ctx,
				payload: &Payload{
					Message: "Test description",
				},
				validator: validator.New(),
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

			service := NewDefaultService(ctxLogger, tc.args.validator, tc.repositoryOpts.logsRepo)
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

func TestUpdate(t *testing.T) {
	ctxLogger := logger.NewContextLogger("TestCreate", "debug", logger.TextFormat)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")
	id := 12345
	type repositoryOpts struct {
		logsRepo     *logsmessagemock.IRepository
		logsRepoFunc func() *logsmessagemock.IRepository
	}
	type args struct {
		ctx       context.Context
		payload   *Payload
		lang      string
		id        *int
		validator *validator.Validate
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
				logsRepoFunc: func() *logsmessagemock.IRepository {
					repoMock := &logsmessagemock.IRepository{}
					repoMock.On("FindByIDAndLang", mock.Anything, &id, "en").
						Return(&log_message.Model{
							ID:      12345,
							Message: "Message test",
							Lang:    "en",
						}, nil)
					repoMock.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
					return repoMock
				},
			},
			args: args{
				ctx:  ctx,
				id:   &id,
				lang: "en",
				payload: &Payload{
					ID:      12345,
					Message: "Test description",
					Lang:    "en",
				},
				validator: validator.New(),
			},
			err: false,
			asserts: func(t *testing.T, ap assertsParams) bool {
				return assert.NoError(t, ap.err) &&
					assert.NotNil(t, ap.result) &&
					ap.logsRepo.AssertCalled(t, "FindByIDAndLang", mock.Anything, mock.Anything, mock.Anything) &&
					ap.logsRepo.AssertCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name: "Error storing data in repository",
			repositoryOpts: repositoryOpts{
				logsRepoFunc: func() *logsmessagemock.IRepository {
					repoMock := &logsmessagemock.IRepository{}
					repoMock.On("FindByIDAndLang", mock.Anything, &id, "en").
						Return(&log_message.Model{
							ID:      12345,
							Message: "Test desc",
							Lang:    "en",
						}, nil)
					repoMock.On("Update", mock.Anything, mock.Anything, mock.Anything).
						Return(terrors.New(terrors.ErrInternalService, "DB Error", nil))
					return repoMock
				},
			},
			args: args{
				ctx:  ctx,
				id:   &id,
				lang: "en",
				payload: &Payload{
					ID:      12345,
					Message: "Test description",
					Lang:    "en",
				},
				validator: validator.New(),
			},
			err: true,
			asserts: func(t *testing.T, ap assertsParams) bool {
				var terr *terrors.Error
				return assert.ErrorAs(t, ap.err, &terr) &&
					assert.Equal(t, "Internal error service", terr.Message) &&
					assert.Nil(t, ap.result)
			},
		},
		{
			name: "Error on validation empty message",
			args: args{
				ctx:  ctx,
				id:   &id,
				lang: "en",
				payload: &Payload{
					ID:   12345,
					Lang: "en",
				},
				validator: validator.New(),
			},
			err: true,
			asserts: func(t *testing.T, ap assertsParams) bool {
				var terr *terrors.Error
				return assert.ErrorAs(t, ap.err, &terr) &&
					assert.Equal(t, "Key: 'Payload.Message' Error:Field validation for 'Message' failed on the 'required' tag", terr.Message) &&
					assert.Nil(t, ap.result)
			},
		},
		{
			name: "Error on validation empty id",
			args: args{
				ctx:  ctx,
				id:   &id,
				lang: "en",
				payload: &Payload{
					Message: "Test description",
					Lang:    "en",
				},
				validator: validator.New(),
			},
			err: true,
			asserts: func(t *testing.T, ap assertsParams) bool {
				var terr *terrors.Error
				return assert.ErrorAs(t, ap.err, &terr) &&
					assert.Equal(t, "Key: 'Payload.ID' Error:Field validation for 'ID' failed on the 'required' tag", terr.Message) &&
					assert.Nil(t, ap.result)
			},
		},
		{
			name: "Error on validation empty lang",
			args: args{
				ctx: ctx,
				id:  &id,
				payload: &Payload{
					Message: "Test description",
					Lang:    "en",
				},
				validator: validator.New(),
			},
			err: true,
			asserts: func(t *testing.T, ap assertsParams) bool {
				var terr *terrors.Error
				return assert.ErrorAs(t, ap.err, &terr) &&
					assert.Equal(t, "Error: Invalid message id or lang code", terr.Message) &&
					assert.Nil(t, ap.result)
			},
		},
		{
			name: "Error on validation empty lang in payload",
			args: args{
				ctx:  ctx,
				id:   &id,
				lang: "en",
				payload: &Payload{
					ID:      id,
					Message: "Test description",
				},
				validator: validator.New(),
			},
			err: true,
			asserts: func(t *testing.T, ap assertsParams) bool {
				var terr *terrors.Error
				return assert.ErrorAs(t, ap.err, &terr) &&
					assert.Equal(t, "Key: 'Payload.Lang' Error:Field validation for 'Lang' failed on the 'required' tag", terr.Message) &&
					assert.Nil(t, ap.result)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.repositoryOpts.logsRepoFunc != nil {
				tc.repositoryOpts.logsRepo = tc.repositoryOpts.logsRepoFunc()
			}

			service := NewDefaultService(ctxLogger, tc.args.validator, tc.repositoryOpts.logsRepo)
			result, err := service.Update(tc.args.ctx, tc.args.id, tc.args.lang, tc.args.payload)

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
