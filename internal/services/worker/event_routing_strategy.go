package worker

import (
	"context"
	"fmt"

	"github.com/jmontesinos91/oevents"
	"github.com/jmontesinos91/oevents/eventfactory"
	"github.com/jmontesinos91/ologs/logger"
	tracekey "github.com/jmontesinos91/ologs/logger/v2"
	"github.com/sirupsen/logrus"
)

// EventRoutingStrategy struct
type EventRoutingStrategy struct {
	log              *logger.ContextLogger
	defaultWorker    IWorker
	logCreatedWorker IWorker
}

// EventRoutingStrategyOpts configuration object to initialize the Routing strategies
type EventRoutingStrategyOpts struct {
	Logger           *logger.ContextLogger
	DefaultWorker    IWorker
	LogCreatedWorker IWorker
}

// NewEventRoutingStrategy generates an instance of EventRoutingStrategy
func NewEventRoutingStrategy(opts EventRoutingStrategyOpts) *EventRoutingStrategy {
	return &EventRoutingStrategy{
		log:              opts.Logger,
		defaultWorker:    opts.DefaultWorker,
		logCreatedWorker: opts.LogCreatedWorker,
	}
}

// Apply method apply the correct strategy depending on event type
func (s *EventRoutingStrategy) Apply(event oevents.OmniViewEvent) error {

	eventType := event.EventType
	s.log.WithContext(
		logrus.InfoLevel,
		"Apply",
		"Routing event of type "+eventType+"",
		logger.Context{
			tracekey.EventID: event.ID,
		},
		nil)

	var eventWorker IWorker
	switch eventType {

	case eventfactory.LogCreatedEvent:
		eventWorker = s.logCreatedWorker

	default:
		eventWorker = s.defaultWorker
	}

	err := eventWorker.Handle(context.Background(), event)
	if err != nil {
		s.log.WithContext(
			logrus.ErrorLevel,
			"Apply",
			fmt.Sprintf("Error while executing worker for event [%s]:", eventType),
			logger.Context{
				tracekey.EventID: event.ID,
			},
			err)
	}

	return nil
}
