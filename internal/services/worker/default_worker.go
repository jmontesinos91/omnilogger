package worker

import (
	"context"
	"fmt"

	"github.com/jmontesinos91/oevents"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/sirupsen/logrus"
)

// DefaultWorker struct
type DefaultWorker struct {
	log *logger.ContextLogger
}

// NewDefaultWorker create an instance of worker for unexpected events
func NewDefaultWorker(l *logger.ContextLogger) *DefaultWorker {
	return &DefaultWorker{
		log: l,
	}
}

// Handle handles the processing for unknown event types
func (w *DefaultWorker) Handle(ctx context.Context, event oevents.OmniViewEvent) error {
	w.log.WithContext(
		logrus.InfoLevel,
		"Handle",
		fmt.Sprintf("Event of type [%s] is not supported in logs", event.EventType),
		logger.Context{},
		nil)
	return nil
}
