package logs

import (
	"context"

	"github.com/jmontesinos91/ologs/logger"
	tracekey "github.com/jmontesinos91/ologs/logger/v2"

	"github.com/jmontesinos91/omnilogger/internal/repositories/logs"
	"github.com/jmontesinos91/terrors"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

// DefaultService struct
type DefaultService struct {
	log      *logger.ContextLogger
	logsRepo logs.IRepository
}

// NewDefaultService creates a new instance of DefaultService log
func NewDefaultService(l *logger.ContextLogger, s logs.IRepository) *DefaultService {
	return &DefaultService{
		log:      l,
		logsRepo: s,
	}
}

// GetByID gets a records by ID
func (s *DefaultService) GetByID(ctx context.Context, ID *string) (*Response, error) {
	requestID := ctx.Value(middleware.RequestIDKey).(string)

	// Create logic of the controller
	if ID == nil {
		return nil, terrors.New(terrors.ErrBadRequest, "", map[string]string{})
	}

	model, err := s.logsRepo.FindByID(ctx, ID)
	if err != nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"GetByID",
			"Error while retrieve log: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
			},
			err)
		return nil, terrors.New(terrors.ErrNotFound, "Log not found", map[string]string{})
	}

	return ToResponse(model), nil
}

// Create model
func (s *DefaultService) Create(ctx context.Context, payload *Payload) (*Response, error) {
	requestID := ctx.Value(middleware.RequestIDKey).(string)

	// Create model for repository
	model, err := ToModel(payload)
	if err != nil {
		s.log.Error(
			logrus.ErrorLevel,
			"ToModel",
			"Failed to map payload data to model",
			err)

		return nil, terrors.InternalService("metadata_error", "Failed to map payload data to model", nil)
	}

	// Store in DB
	err = s.logsRepo.Create(ctx, model)
	if err != nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"Create",
			"Invalid log request payload: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
			},
			err)
		return nil, terrors.New(terrors.ErrInternalService, "Internal error service", map[string]string{})
	}

	return ToResponse(model), nil
}
