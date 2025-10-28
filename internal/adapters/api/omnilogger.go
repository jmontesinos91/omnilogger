package api

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jmontesinos91/ologs/logger"
	tracekey "github.com/jmontesinos91/ologs/logger/v2"
	"github.com/jmontesinos91/omnilogger/internal/services/logs"
	"github.com/jmontesinos91/osecurity/sts"
	"github.com/jmontesinos91/terrors"
	"github.com/sirupsen/logrus"
)

// OmniLoggerController OmniLogger controller
type OmniLoggerController struct {
	log           *logger.ContextLogger
	validate      *validator.Validate
	logsSvc       logs.IService
	stsClient     sts.ISTSClient
	counterMetric prometheus.Counter
}

// NewOmniLoggerController Constructor
func NewOmniLoggerController(server *HTTPServer, validator *validator.Validate, ss logs.IService, sts sts.ISTSClient) *OmniLoggerController {
	sc := &OmniLoggerController{
		log:       server.Logger,
		validate:  validator,
		logsSvc:   ss,
		stsClient: sts,
		counterMetric: promauto.NewCounter(prometheus.CounterOpts{
			Name: "omni_logger_reqs_total",
			Help: "The total number of requests to omni logger endpoints",
		}),
	}

	// Endpoints if we add JWTVerifyMiddleWare, we add the secure
	server.Router.Group(func(r chi.Router) {
		r.Use(JwtVerifyMiddleware(server.Logger, sts))
		r.Get("/v1/logs/{id}", sc.handleGetLog)
		r.Post("/v1/logs", sc.handleCreate)
	})

	return sc
}

func (sc *OmniLoggerController) handleGetLog(w http.ResponseWriter, r *http.Request) {
	// Increment metric
	sc.counterMetric.Inc()

	id := chi.URLParam(r, "id")
	// Increment metric
	sc.counterMetric.Inc()
	idRes, err := sc.logsSvc.GetByID(r.Context(), &id)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}

	RenderJSON(r.Context(), w, http.StatusOK, idRes)
}

func (sc *OmniLoggerController) handleCreate(w http.ResponseWriter, r *http.Request) {
	// Increment metric
	sc.counterMetric.Inc()

	sc.log.Log(logrus.InfoLevel, "handleCreate", "start endpoint")

	var payload logs.Payload
	requestID := r.Context().Value(middleware.RequestIDKey).(string)

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		sc.log.WithContext(
			logrus.ErrorLevel,
			"handleCreate",
			"Error while parsing request payload: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
			},
			err)
		terr := terrors.BadRequest(terrors.ErrBadRequest, "Malformed body", map[string]string{})
		RenderError(r.Context(), w, terr)
		return
	}
	// Increment metric
	sc.counterMetric.Inc()
	// Call the service
	res, err := sc.logsSvc.Create(r.Context(), &payload)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}

	RenderJSON(r.Context(), w, http.StatusCreated, res)
}
