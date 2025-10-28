package worker

import (
	"context"
	"github.com/jmontesinos91/oevents/broker"
	"github.com/jmontesinos91/ologs/logger"
	tracekey "github.com/jmontesinos91/ologs/logger/v2"
	"github.com/jmontesinos91/omnilogger/config"

	"github.com/sirupsen/logrus"
)

// ConsumerOptions struct
type ConsumerOptions struct {
	Logger               *logger.ContextLogger
	Broker               broker.MessagingBrokerProvider
	KafkaConfigs         config.KafkaConsumerConfigurations
	EventRoutingStrategy IRoutingStrategy
}

// MainConsumer struct
type MainConsumer struct {
	log                  *logger.ContextLogger
	broker               broker.MessagingBrokerProvider
	kafkaConfigs         config.KafkaConsumerConfigurations
	eventRoutingStrategy IRoutingStrategy
}

// NewMainConsumer generate an instance of MainConsumer
func NewMainConsumer(opts ConsumerOptions) *MainConsumer {
	return &MainConsumer{
		log:                  opts.Logger,
		broker:               opts.Broker,
		kafkaConfigs:         opts.KafkaConfigs,
		eventRoutingStrategy: opts.EventRoutingStrategy,
	}
}

// Start processing of the events
func (m *MainConsumer) Start(ctx context.Context) {
	if !m.kafkaConfigs.Enabled {
		m.log.Log(logrus.WarnLevel, "Start", "Kafka consumer not enabled, ignoring request to start consumer")
		return
	}

	// The workers channel, must be a bounded channel to avoid running out of memory
	workerChannel := make(chan broker.OmniViewMessage, m.kafkaConfigs.MaxRecords)

	go m.eventHandler(workerChannel)

	// Subscribe to the topic
	m.broker.Subscribe(ctx, m.kafkaConfigs.MaxRecords, workerChannel)
}

/*
eventHandler filters workers depending on the event topics, inside it applies an event routing strategy
for example we have the following kafka events:

	Event: A
	Event: B

	eventRoutingStrategy.Apply(Event)
*/
func (m *MainConsumer) eventHandler(workerChannel <-chan broker.OmniViewMessage) {
	for msg := range workerChannel {

		eventType := msg.Event.EventType
		m.log.WithContext(
			logrus.InfoLevel,
			"eventHandler",
			"Received event of type "+eventType+"",
			logger.Context{
				tracekey.EventID: msg.Event.ID,
			},
			nil)

		err := m.eventRoutingStrategy.Apply(msg.Event)
		if err != nil {
			m.log.WithContext(
				logrus.ErrorLevel,
				"eventHandler",
				"Error while trying to handle event:",
				logger.Context{
					tracekey.EventID: msg.Event.ID,
				},
				err)
		}

		// Ack the message
		msg.Ack.Done()
	}
}
