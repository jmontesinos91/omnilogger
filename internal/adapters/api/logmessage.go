package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jmontesinos91/ologs/logger"
	tracekey "github.com/jmontesinos91/ologs/logger/v2"
	"github.com/jmontesinos91/omnilogger/internal/services/log_message"
	"github.com/jmontesinos91/osecurity/sts"
	"github.com/jmontesinos91/terrors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

// LogMessageController OmniLogger controller
type LogMessageController struct {
	log           *logger.ContextLogger
	validate      *validator.Validate
	logMessageSvc log_message.IService
	stsClient     sts.ISTSClient
	counterMetric prometheus.Counter
}

// NewLogMessageController Constructor
func NewLogMessageController(server *HTTPServer, validator *validator.Validate, ss log_message.IService, sts sts.ISTSClient) *LogMessageController {
	sc := &LogMessageController{
		log:           server.Logger,
		validate:      validator,
		logMessageSvc: ss,
		stsClient:     sts,
		counterMetric: promauto.NewCounter(prometheus.CounterOpts{
			Name: "log_messages_reqs_total",
			Help: "The total number of requests to omni logger endpoints",
		}),
	}

	// Endpoints if we add JWTVerifyMiddleWare, we add the secure
	server.Router.Group(func(r chi.Router) {
		// r.Use(JwtVerifyMiddleware(server.Logger, sts))
		r.Get("/v1/log_messages", sc.handleRetrieve)
		r.Get("/v1/log_messages/{id}", sc.handleGet)
		r.Post("/v1/log_messages", sc.handleCreate)
		r.Post("/v1/log_messages/{id}", sc.handleUpdate)
		r.Delete("/v1/log_messages/{id}/{lang}", sc.handleDeleteLang)
		r.Delete("/v1/log_messages/{id}", sc.handleDeleteMessage)
	})

	return sc
}

func (sc *LogMessageController) handleGet(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	idRes, err := sc.logMessageSvc.GetByID(r.Context(), &id)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}

	RenderJSON(r.Context(), w, http.StatusOK, idRes)
}

func (sc *LogMessageController) handleCreate(w http.ResponseWriter, r *http.Request) {
	// Increment metric
	sc.counterMetric.Inc()
	requestID := r.Context().Value(middleware.RequestIDKey).(string)

	var payload log_message.Payload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		sc.log.WithContext(
			logrus.ErrorLevel,
			"handleCreate",
			"Error while parsing request payload: %v",
			logger.Context{
				tracekey.TrackingID: requestID,
			},
			err)
		RenderError(r.Context(), w, err)
		return
	}

	// Call the service
	res, err := sc.logMessageSvc.Create(r.Context(), &payload)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}

	RenderJSON(r.Context(), w, http.StatusCreated, res)
}

func (sc *LogMessageController) handleUpdate(w http.ResponseWriter, r *http.Request) {
	// Increment metric
	sc.counterMetric.Inc()
	var payload log_message.Payload

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	lang := chi.URLParam(r, "lang")
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		terr := terrors.BadRequest(terrors.ErrBadRequest, "Malformed body", map[string]string{})
		RenderError(r.Context(), w, terr)
		return
	}

	res, err := sc.logMessageSvc.Update(r.Context(), &id, lang, &payload)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}

	RenderJSON(r.Context(), w, http.StatusAccepted, res)
}

func (sc *LogMessageController) handleRetrieve(w http.ResponseWriter, r *http.Request) {
	// Increment metric
	sc.counterMetric.Inc()

	filter, err := log_message.ToParseFilterRequest(r)
	if err != nil {
		sc.log.Error(logrus.ErrorLevel, "handleRetrieve", "Invalid request parameters", err)
		RenderError(r.Context(), w, err)
		return
	}

	res, err := sc.logMessageSvc.Retrieve(r.Context(), filter)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}

	RenderJSON(r.Context(), w, http.StatusOK, res)
}

func (sc *LogMessageController) handleDeleteLang(w http.ResponseWriter, r *http.Request) {
	sc.counterMetric.Inc()

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	lang := chi.URLParam(r, "lang")

	err := sc.logMessageSvc.DeleteLang(r.Context(), &id, lang)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}

	RenderJSON(r.Context(), w, http.StatusNoContent, nil)
}

func (sc *LogMessageController) handleDeleteMessage(w http.ResponseWriter, r *http.Request) {
	sc.counterMetric.Inc()

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := sc.logMessageSvc.DeleteMessage(r.Context(), &id)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}

	RenderJSON(r.Context(), w, http.StatusNoContent, nil)
}
