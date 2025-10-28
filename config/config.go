package config

import (
	"strings"

	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/osecurity/services/omnibackend"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/sirupsen/logrus"
)

// ServerConfigurations Server configurations
type ServerConfigurations struct {
	Port          int    `koanf:"port"`
	BaseDirectory string `koanf:"base-directory"`
	Host          string `koanf:"host"`
}

// KeysConfigurations asymmetric keys
type KeysConfigurations struct {
	Public string `koanf:"public"`
}

// Service configurations
type Service struct {
	Name string `koanf:"name"`
}

// DatabaseConfigurations database configurations
type DatabaseConfigurations struct {
	Dsn        string `koanf:"dsn"`
	Pool       int    `koanf:"pool"`
	Database   string `koanf:"database"`
	Collection string `koanf:"collection"`
}

// KafkaConfigurations struct
type KafkaConfigurations struct {
	SecuredMode bool                        `koanf:"secured-mode"`
	Servers     string                      `koanf:"servers"`
	User        string                      `koanf:"user"`
	Password    string                      `koanf:"pass"`
	ClientName  string                      `koanf:"client-name"`
	Consumer    KafkaConsumerConfigurations `koanf:"consumer"`
}

// KafkaConsumerConfigurations Kafka consumer configurations
type KafkaConsumerConfigurations struct {
	Enabled    bool     `koanf:"enabled"`
	Group      string   `koanf:"group"`
	Topics     []string `koanf:"topics"`
	MaxRecords int      `koanf:"max-records"`
}

// Configurations Application wide configurations
type Configurations struct {
	Server   ServerConfigurations               `koanf:"server"`
	Keys     KeysConfigurations                 `koanf:"keys"`
	Service  Service                            `koanf:"service"`
	Database DatabaseConfigurations             `koanf:"database"`
	OmniView omnibackend.OmniViewConfigurations `koanf:"omniview"`
	Kafka    KafkaConfigurations                `koanf:"kafka"`
}

// LoadConfig Loads configurations depending upon the environment
func LoadConfig(logger *logger.ContextLogger) *Configurations {
	k := koanf.New(".")
	err := k.Load(file.Provider("resources/config.yml"), yaml.Parser())
	if err != nil {
		logger.Error(logrus.FatalLevel, "LoadConfig", "Failed to locate configurations.", err)
	}

	// Searches for env variables and will transform them into koanf format
	// e.g. SERVER_PORT variable will be server.port: value
	err = k.Load(env.Provider("", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(s), "_", ".")
	}), nil)
	if err != nil {
		logger.Error(logrus.FatalLevel, "LoadConfig", "Failed to replace environment variables. ", err)
	}

	var configuration Configurations

	err = k.Unmarshal("", &configuration)
	if err != nil {
		logger.Error(logrus.FatalLevel, "LoadConfig", "Failed to load configurations. ", err)
	}

	return &configuration
}
