package log_message

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jmontesinos91/ologs/logger"
	tracekey "github.com/jmontesinos91/ologs/logger/v2"
	"github.com/jmontesinos91/omnilogger/internal/repositories/log_message"
	"github.com/jmontesinos91/terrors"
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

func (s *DefaultService) Update(ctx context.Context, id *int, payload *Payload) (*Response, error) {
	requestID := ctx.Value(middleware.RequestIDKey).(string)

	if err := payload.SanitizeAndValidate(s.validate); err != nil {
		return nil, terrors.New(terrors.ErrBadRequest, err.Error(), map[string]string{})
	}

	model, err := s.logMessagesRepo.FindByID(ctx, id)
	if err != nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"Update",
			"Record not found",
			logger.Context{
				tracekey.TrackingID: requestID,
				"ID":                id,
			}, err)

		if model == nil {
			return nil, terrors.New(terrors.ErrNotFound, "Record not found", map[string]string{})
		}

		return nil, terrors.New(terrors.ErrInternalService, "Internal error service", map[string]string{})
	}

	model = ToModelUpdate(model, *payload)
	errU := s.logMessagesRepo.Update(ctx, id, model)
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
