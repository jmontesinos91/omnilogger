package logs

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmontesinos91/oevents/eventfactory"
	"github.com/jmontesinos91/ologs/logger"
	tracekey "github.com/jmontesinos91/ologs/logger/v2"
	"github.com/jmontesinos91/omnilogger/internal/repositories/logs"
	"github.com/jmontesinos91/omnilogger/internal/utils/export"
	"github.com/jmontesinos91/omnilogger/internal/utils/format"
	"github.com/jmontesinos91/osecurity/sts"
	"github.com/jmontesinos91/terrors"
	lop "github.com/samber/lo/parallel"
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
func (s *DefaultService) GetByID(ctx context.Context, ID *string, filter Filter) (*Response, error) {
	requestID := ctx.Value(middleware.RequestIDKey).(string)

	// Create logic of the controller
	if ID == nil {
		return nil, terrors.New(terrors.ErrBadRequest, "", map[string]string{})
	}

	repoFilter := ToRepoFilter(filter)

	model, err := s.logsRepo.FindByID(ctx, ID, repoFilter)
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

	return ToResponse(model, filter.Lang), nil
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

	return ToResponse(model, payload.Lang), nil
}

// Retrieve logs with filter
func (s *DefaultService) Retrieve(ctx context.Context, filter Filter) (*PaginatedRes, error) {
	requestID := ctx.Value(middleware.RequestIDKey).(string)

	repoFilter := ToRepoFilter(filter)

	res, total, err := s.logsRepo.Retrieve(ctx, repoFilter)
	if err != nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"Retrieve",
			"Error while retrieve logs: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
			},
			err)
		return nil, terrors.New(terrors.ErrInternalService, "Internal error service", map[string]string{})
	}

	items := lop.Map(res, func(p logs.Model, _ int) Response {
		return *ToResponse(&p, filter.Lang)
	})

	currentPage := filter.Page
	if currentPage == 0 {
		currentPage = 1
	}

	return &PaginatedRes{
		Data:  items,
		Size:  filter.Size,
		Total: total,
		Page:  currentPage,
	}, nil
}

// CreateLogFromKafka creates a new log from kafka
func (s *DefaultService) CreateLogFromKafka(ctx context.Context, payload *eventfactory.LogCreatedPayload) error {

	tenantCatJSON, errBind := eventfactory.ToTenantCatJson(payload.TenantCat)
	if errBind != nil {
		s.log.Error(
			logrus.ErrorLevel,
			"CreateLogFromKafka",
			"Error marshalling TenantCat to JSON",
			errBind)
		return terrors.InternalService("tenant_cat_error", "Error marshalling TenantCat to JSON", nil)
	}

	model := &Payload{
		IpAddress:   payload.IpAddress,
		ClientHost:  payload.ClientHost,
		Provider:    payload.Provider,
		Level:       payload.Level,
		Message:     payload.Message,
		Description: payload.Description,
		Resource:    payload.Resource,
		Path:        payload.Path,
		Action:      payload.Action,
		Data:        payload.Data,
		OldData:     payload.OldData,
		UserID:      payload.UserID,
		TenantCat:   tenantCatJSON,
	}

	// Create model for repository
	data, err := ToModel(model)
	if err != nil {
		s.log.Error(
			logrus.ErrorLevel,
			"ToModel",
			"Failed to map payload data to model",
			err)

		return terrors.InternalService("metadata_error", "Failed to map payload data to model", nil)
	}

	err = s.logsRepo.Create(ctx, data)
	if err != nil {
		s.log.Error(
			logrus.ErrorLevel,
			"CreateLogFromKafka",
			"Failed to create log from kafka",
			err)

		return terrors.InternalService("metadata_error", "Error storing model for log", nil)
	}

	return nil
}

func (s *DefaultService) Export(ctx context.Context, filter Filter) ([]byte, error) {
	requestID := ctx.Value(middleware.RequestIDKey).(string)
	claims := ctx.Value(&sts.Claim).(sts.Claims)

	repoFilter := ToRepoFilter(filter)

	res, err := s.logsRepo.Export(ctx, repoFilter)
	if err != nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"Retrieve",
			"Error while retrieve logs: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
			},
			err)
		return nil, terrors.New(terrors.ErrInternalService, "Internal error service", map[string]string{})
	}

	items := lop.Map(res, func(p logs.Model, _ int) Response {
		return *ToResponse(&p, filter.Lang)
	})

	genericMapper := func(item Response) format.ExcelRow {
		return format.ExcelRow{
			Cells: []interface{}{
				item.ID,
				item.IpAddress,
				item.ClientHost,
				item.Provider,
				item.Level,
				item.Message,
				item.LogMessage,
				item.Description,
				item.Path,
				item.Resource,
				item.Action,
				item.Data,
				item.OldData,
				item.TenantCat,
				item.UserID,
				item.CreatedAt,
			},
		}
	}

	excelBytes, err := export.DataToExcel("logs", items, genericMapper)
	if err != nil {
		s.log.WithContext(logrus.ErrorLevel,
			"HandleExport",
			"Failed to generate excel bytes from data",
			logger.Context{
				tracekey.TrackingID: requestID,
				tracekey.UserID:     claims.UserID,
				tracekey.Role:       claims.Role,
			},
			err)
		return []byte{}, err
	}

	return excelBytes, nil
}
