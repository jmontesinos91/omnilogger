package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/omnilogger/config"
	"github.com/jmontesinos91/omnilogger/internal/adapters/api"
	"github.com/jmontesinos91/omnilogger/internal/adapters/db"
	lmrepository "github.com/jmontesinos91/omnilogger/internal/repositories/log_message"
	repository "github.com/jmontesinos91/omnilogger/internal/repositories/logs"
	"github.com/jmontesinos91/omnilogger/internal/services/log_message"
	"github.com/jmontesinos91/omnilogger/internal/services/logs"
	"github.com/jmontesinos91/osecurity/services/omnibackend"
	"github.com/jmontesinos91/osecurity/sts"
)

func main() {
	// Logger
	contextLogger := logger.NewContextLogger("OMNILOGGER", "debug", logger.TextFormat)

	// Configs
	configs := config.LoadConfig(contextLogger)

	// STS Client
	omniService := omnibackend.NewOmniViewService(contextLogger, configs.OmniView)

	stsClient := sts.NewDefaultISTSClient(contextLogger, omniService)

	// Http Router
	httpServer := api.NewHTTPServer(contextLogger, configs.Server, configs.Service, stsClient)

	// Validator
	validate := validator.New()

	// DB Connection
	conn := db.NewDatabaseConnection(contextLogger, configs.Database)

	// -- Start dependency injection section --

	// - Initialize repository -
	omniLoggerRepo := repository.NewDatabaseRepository(contextLogger, conn)
	logMessageRepo := lmrepository.NewDatabaseRepository(contextLogger, conn)

	// - Initialize service -
	omniLoggerSvc := logs.NewDefaultService(contextLogger, omniLoggerRepo)
	logMessageSvc := log_message.NewDefaultService(contextLogger, validate, logMessageRepo)

	api.NewHealthController(httpServer)
	api.NewOmniLoggerController(httpServer, validate, omniLoggerSvc, stsClient)
	api.NewLogMessageController(httpServer, validate, logMessageSvc, stsClient)
	// -- End dependency injection section --

	// Let the party started!
	httpServer.Start()
}
