package worker

import (
	"context"

	"github.com/jmontesinos91/oevents"
	"github.com/jmontesinos91/oevents/broker"
	"github.com/jmontesinos91/oevents/eventfactory"
	"github.com/jmontesinos91/ologs/logger"
	tracekey "github.com/jmontesinos91/ologs/logger/v2"
	"github.com/jmontesinos91/omnilogger/internal/services/logs"
	"github.com/sirupsen/logrus"
)

// LogCreatedWorker Processes a record to be saved from somewhere in the system.
type LogCreatedWorker struct {
	log          *logger.ContextLogger
	logSvc       logs.IService
	streamClient broker.MessagingBrokerProvider
}

// NewLogCreatedWorker Generates a LogCreatedWorker instance
func NewLogCreatedWorker(l *logger.ContextLogger, ls logs.IService, sc broker.MessagingBrokerProvider) *LogCreatedWorker {
	return &LogCreatedWorker{
		log:          l,
		logSvc:       ls,
		streamClient: sc,
	}
}

// Handle handles incoming logs to be created
func (w *LogCreatedWorker) Handle(ctx context.Context, event oevents.OmniViewEvent) error {
	eventPayload, err := eventfactory.ToLogCreatedPayload(event.Data)
	if err != nil {
		w.log.WithContext(
			logrus.ErrorLevel,
			"Handle",
			"Error parsing event to LogCreatedPayload:",
			logger.Context{
				tracekey.EventID: event.ID,
			},
			err)
		return nil
	}

	w.log.WithContext(
		logrus.InfoLevel,
		"Handle",
		"Processing event of type log_created_by_logs",
		logger.Context{
			tracekey.EventID: event.ID,
		},
		err)

	errCFK := w.logSvc.CreateLogFromKafka(ctx, eventPayload)
	if errCFK != nil {
		w.log.WithContext(
			logrus.ErrorLevel,
			"Handle",
			"Error creating the log",
			logger.Context{
				tracekey.EventID: event.ID,
			},
			errCFK)

		return errCFK
	}

	w.log.Log(logrus.InfoLevel, "Handle", "log created with event ID: "+event.ID)

	return nil
}
