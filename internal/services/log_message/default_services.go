package log_message

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jmontesinos91/ologs/logger"
	tracekey "github.com/jmontesinos91/ologs/logger/v2"
	"github.com/jmontesinos91/omnilogger/internal/repositories/log_message"
	"github.com/jmontesinos91/terrors"
	lop "github.com/samber/lo/parallel"
	"github.com/sirupsen/logrus"
)

// DefaultService struct
type DefaultService struct {
	log             *logger.ContextLogger
	validate        *validator.Validate
	logMessagesRepo log_message.IRepository
}

// NewDefaultService creates a new instance of DefaultService log message
func NewDefaultService(l *logger.ContextLogger, v *validator.Validate, s log_message.IRepository) *DefaultService {
	return &DefaultService{
		log:             l,
		validate:        v,
		logMessagesRepo: s,
	}
}

// GetByID service to get byID
func (s *DefaultService) GetByID(ctx context.Context, ID *int) (*Response, error) {
	requestID := ctx.Value(middleware.RequestIDKey).(string)

	// Create logic of the controller
	if ID == nil {
		return nil, terrors.New(terrors.ErrBadRequest, "", map[string]string{})
	}

	model, err := s.logMessagesRepo.FindByID(ctx, ID)
	if err != nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"GetByID",
			"Invalid payout request payload: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
			},
			err)
		return nil, terrors.New(terrors.ErrInternalService, "Internal error service", map[string]string{})
	}

	return ToResponse(model), nil
}

// Create model
func (s *DefaultService) Create(ctx context.Context, payload *Payload) (*Response, error) {
	requestID := ctx.Value(middleware.RequestIDKey).(string)

	// Create model for repository
	dbModel := ToModel(payload)

	// Store in DB
	err := s.logMessagesRepo.Create(ctx, dbModel)
	if err != nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"Create",
			"Invalid log message request payload: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
			},
			err)

		return nil, terrors.New(terrors.ErrInternalService, "Internal error service", map[string]string{})
	}

	return ToResponse(dbModel), nil
}

// Update model
func (s *DefaultService) Update(ctx context.Context, id *int, lang string, payload *Payload) (*Response, error) {
	requestID := ctx.Value(middleware.RequestIDKey).(string)

	if id == nil || lang == "" {
		s.log.WithContext(
			logrus.ErrorLevel,
			"Update",
			"Invalid message id or lang code",
			logger.Context{
				tracekey.TrackingID: requestID,
				"ID":                id,
				"Lang":              lang,
			}, nil)
		return nil, terrors.New(terrors.ErrBadRequest, "Error: Invalid message id or lang code", map[string]string{})
	}

	if err := payload.SanitizeAndValidate(s.validate); err != nil {
		return nil, terrors.New(terrors.ErrBadRequest, err.Error(), map[string]string{})
	}

	model, err := s.logMessagesRepo.FindByIDAndLang(ctx, id, lang)
	if err != nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"Update",
			"Record not found",
			logger.Context{
				tracekey.TrackingID: requestID,
				"ID":                id,
				"Lang":              lang,
			}, err)

		if model == nil {
			return nil, terrors.New(terrors.ErrNotFound, "Record not found", map[string]string{})
		}

		return nil, terrors.New(terrors.ErrInternalService, "Internal error service", map[string]string{})
	}

	model = ToModelUpdate(model, *payload)
	errU := s.logMessagesRepo.Update(ctx, id, lang, model)
	if errU != nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"Update",
			"Error while persisting log message: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
				"ID":                id,
			}, errU)

		return nil, terrors.New(terrors.ErrInternalService, "Internal error service", map[string]string{})
	}

	return ToResponse(model), nil
}

// Retrieve service to get log messages with filters
func (s *DefaultService) Retrieve(ctx context.Context, filter Filter) (*PaginatedRes, error) {
	requestID := ctx.Value(middleware.RequestIDKey).(string)

	repoFilter := ToRepoFilter(filter)

	res, total, err := s.logMessagesRepo.Retrieve(ctx, repoFilter)
	if err != nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"Retrieve",
			"Error retrieving log messages: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
			}, err)
		return nil, err
	}

	items := lop.Map(res, func(p log_message.Model, _ int) Response {
		return *ToResponse(&p)
	})

	return &PaginatedRes{
		Data:  items,
		Size:  filter.Size,
		Total: total,
		Page:  filter.Page,
	}, nil
}

// DeleteLang service to delete a language from a log message
func (s *DefaultService) DeleteLang(ctx context.Context, id *int, lang string) error {
	requestID := ctx.Value(middleware.RequestIDKey).(string)

	if id == nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"DeleteLang",
			"Missing id param: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
				"ID":                id,
				"Lang":              lang,
			}, nil)

		return terrors.New(terrors.ErrBadRequest, "Missing id param", map[string]string{})
	}

	if lang == "" {
		s.log.WithContext(
			logrus.ErrorLevel,
			"DeleteLang",
			"Missing lang param: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
				"ID":                id,
				"Lang":              lang,
			}, nil)

		return terrors.New(terrors.ErrBadRequest, "Missing lang param", map[string]string{})
	}

	err := s.logMessagesRepo.DeleteLang(ctx, id, lang)
	if err != nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"DeleteLang",
			"Error deleting lang: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
				"ID":                id,
				"Lang":              lang,
			}, err)

		return terrors.New(terrors.ErrInternalService, "Internal error service", map[string]string{})
	}

	return nil
}

// DeleteMessage service to delete a log message
func (s *DefaultService) DeleteMessage(ctx context.Context, id *int) error {
	requestID := ctx.Value(middleware.RequestIDKey).(string)

	if id == nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"DeleteMessage",
			"Missing id param: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
				"ID":                id,
			}, nil)
		return terrors.New(terrors.ErrBadRequest, "", map[string]string{})
	}

	err := s.logMessagesRepo.DeleteMessage(ctx, id)
	if err != nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"DeleteMessage",
			"Error deleting message: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
				"ID":                id,
			}, err)

		return terrors.New(terrors.ErrInternalService, "Internal error service", map[string]string{})
	}

	return nil
}
