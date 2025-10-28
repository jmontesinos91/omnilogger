package api

import (
	"github.com/jmontesinos91/omnilogger/config"
	"net/http"
	"strconv"
	"time"

	"github.com/jmontesinos91/osecurity/sts"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmchiv5/v2"
)

// HTTPServer http server
type HTTPServer struct {
	Logger    *logger.ContextLogger
	sc        config.ServerConfigurations
	Router    *chi.Mux
	stsClient sts.ISTSClient
}

// NewHTTPServer Initializes a new http server
func NewHTTPServer(logger *logger.ContextLogger, serverConf config.ServerConfigurations,
	serviceConf config.Service, client sts.ISTSClient) *HTTPServer {
	router := chi.NewRouter()

	// Enable APM chiv5 Middleware
	router.Use(apmchiv5.Middleware())

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.AllowContentType("application/json"))

	// Set a timeout value on the request models (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Use(client.CorsMiddleware)

	return &HTTPServer{
		Logger:    logger,
		sc:        serverConf,
		Router:    router,
		stsClient: client,
	}
}

// Start Fires the http server
func (r *HTTPServer) Start() {
	listeningAddr := ":" + strconv.Itoa(r.sc.Port)
	r.Logger.Log(logrus.InfoLevel, "Start", "Server listening on port "+listeningAddr+"")

	err := http.ListenAndServe(listeningAddr, r.Router)
	if err != nil {
		r.Logger.Error(logrus.FatalLevel, "Start", "Failed to start http server. ", err)
	}
}
