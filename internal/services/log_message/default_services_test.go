package log_message

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/omnilogger/domains/pagination"
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

func TestRetrieve(t *testing.T) {
	ctxLogger := logger.NewContextLogger("TestRetrieve", "debug", logger.TextFormat)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	cases := []struct {
		name     string
		repoFunc func() *logsmessagemock.IRepository
		filter   Filter
		assertFn func(t *testing.T, repo *logsmessagemock.IRepository, res *PaginatedRes, err error)
	}{
		{
			name: "Happy path - one item",
			repoFunc: func() *logsmessagemock.IRepository {
				repoMock := &logsmessagemock.IRepository{}
				repoMock.On("Retrieve", mock.Anything, mock.Anything).
					Return([]log_message.Model{
						{ID: 12345, Message: "msg", Lang: "en"},
					}, 1, nil)
				return repoMock
			},
			filter: Filter{Filter: pagination.Filter{Page: 1, Size: 10}},
			assertFn: func(t *testing.T, repo *logsmessagemock.IRepository, res *PaginatedRes, err error) {
				if !assert.NoError(t, err) {
					return
				}
				if !assert.NotNil(t, res) {
					return
				}
				assert.Equal(t, 1, len(res.Data))
				assert.Equal(t, 12345, res.Data[0].ID)
				assert.Equal(t, 1, res.Total)
				repo.AssertCalled(t, "Retrieve", mock.Anything, mock.Anything)
			},
		},
		{
			name: "Not found - empty result",
			repoFunc: func() *logsmessagemock.IRepository {
				repoMock := &logsmessagemock.IRepository{}
				repoMock.On("Retrieve", mock.Anything, mock.Anything).
					Return([]log_message.Model{}, 0, nil)
				return repoMock
			},
			filter: Filter{Filter: pagination.Filter{Page: 1, Size: 10}},
			assertFn: func(t *testing.T, repo *logsmessagemock.IRepository, res *PaginatedRes, err error) {
				if !assert.NoError(t, err) {
					return
				}
				if !assert.NotNil(t, res) {
					return
				}
				assert.Equal(t, 0, len(res.Data))
				assert.Equal(t, 0, res.Total)
				repo.AssertCalled(t, "Retrieve", mock.Anything, mock.Anything)
			},
		},
		{
			name: "DB error - repository fails",
			repoFunc: func() *logsmessagemock.IRepository {
				repoMock := &logsmessagemock.IRepository{}
				repoMock.On("Retrieve", mock.Anything, mock.Anything).
					Return(nil, 0, terrors.New(terrors.ErrInternalService, "DB Error", nil))
				return repoMock
			},
			filter: Filter{Filter: pagination.Filter{Page: 1, Size: 10}},
			assertFn: func(t *testing.T, repo *logsmessagemock.IRepository, res *PaginatedRes, err error) {
				if !assert.Error(t, err) {
					return
				}
				assert.Nil(t, res)
				var terr *terrors.Error
				assert.ErrorAs(t, err, &terr)
				assert.Equal(t, "DB Error", terr.Message)
				repo.AssertCalled(t, "Retrieve", mock.Anything, mock.Anything)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.repoFunc()
			service := NewDefaultService(ctxLogger, validator.New(), repo)
			res, err := service.Retrieve(ctx, tc.filter)
			tc.assertFn(t, repo, res, err)
		})
	}
}

func TestDeleteLang(t *testing.T) {
	ctxLogger := logger.NewContextLogger("TestDeleteLang", "debug", logger.TextFormat)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")
	id := 12345

	cases := []struct {
		name     string
		repoFunc func() *logsmessagemock.IRepository
		assertFn func(t *testing.T, repo *logsmessagemock.IRepository, err error)
	}{
		{
			name: "Happy path - delete lang succeeds",
			repoFunc: func() *logsmessagemock.IRepository {
				repoMock := &logsmessagemock.IRepository{}
				repoMock.On("DeleteLang", mock.Anything, &id, "es").Return(nil)
				return repoMock
			},
			assertFn: func(t *testing.T, repo *logsmessagemock.IRepository, err error) {
				if !assert.NoError(t, err) {
					return
				}
				repo.AssertCalled(t, "DeleteLang", mock.Anything, &id, "es")
			},
		},
		{
			name: "DB error - repository fails",
			repoFunc: func() *logsmessagemock.IRepository {
				repoMock := &logsmessagemock.IRepository{}
				repoMock.On("DeleteLang", mock.Anything, &id, "es").
					Return(terrors.New(terrors.ErrInternalService, "DB Error", nil))
				return repoMock
			},
			assertFn: func(t *testing.T, repo *logsmessagemock.IRepository, err error) {
				if !assert.Error(t, err) {
					return
				}
				var terr *terrors.Error
				if !assert.ErrorAs(t, err, &terr) {
					return
				}
				assert.Equal(t, "Internal error service", terr.Message)
				repo.AssertCalled(t, "DeleteLang", mock.Anything, &id, "es")
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.repoFunc()
			service := NewDefaultService(ctxLogger, validator.New(), repo)
			err := service.DeleteLang(ctx, &id, "es")
			tc.assertFn(t, repo, err)
		})
	}
}
