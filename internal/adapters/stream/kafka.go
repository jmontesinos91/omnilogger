package stream

import (
	"github.com/jmontesinos91/oevents/broker"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/omnilogger/config"
	"github.com/sirupsen/logrus"
)

// NewKafkaConnection generate a new MessagingBrokerProvider
func NewKafkaConnection(log *logger.ContextLogger, c config.KafkaConfigurations) (broker.MessagingBrokerProvider, func()) {
	streamConfig := broker.OBrokerConfig{
		Servers:           c.Servers,
		User:              c.User,
		Password:          c.Password,
		ClientName:        c.ClientName,
		ConsumerEnabled:   c.Consumer.Enabled,
		ConsumerGroupName: c.Consumer.Group,
		ConsumeFromTopics: c.Consumer.Topics,
	}

	var stream broker.MessagingBrokerProvider
	var err error

	if c.SecuredMode {
		log.Log(logrus.InfoLevel, "NewKafkaConnection", "Connecting to kafka in secure mode.")
		stream, err = broker.ConnectKafka(streamConfig, log)
	} else {
		log.Log(logrus.InfoLevel, "NewKafkaConnection", "Connecting to kafka in insecure mode.")
		stream, err = broker.ConnectInsecure(streamConfig, log)
	}

	if err != nil {
		log.Error(logrus.FatalLevel, "NewKafkaConnection", "not able to connect to the broker:", err)
	}
	log.Log(logrus.InfoLevel, "NewKafkaConnection", "Kafka Client started!")

	return stream, func() {
		stream.Close() //nolint:errcheck
	}
}
