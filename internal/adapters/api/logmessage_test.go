package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmontesinos91/omnilogger/internal/services/log_message"
	"github.com/jmontesinos91/omnilogger/internal/services/log_message/logmessagesvcmock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// helper: add chi route context with params to request
func withChiRouteParams(req *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestLogMessageController_TableDriven(t *testing.T) {
	type tc struct {
		name            string
		handler         string // "create", "get", "update", "retrieve"
		method          string
		path            string
		body            string
		reqID           string
		chiParams       map[string]string
		svc             *logmessagesvcmock.IService
		expectedCode    int
		expectedNot     []int
		expectedCounter float64
		expectedRespID  int // 0 means don't check
	}

	tests := []tc{
		{
			name:            "HandleCreate_Success_IncrementsCounter",
			handler:         "create",
			method:          http.MethodPost,
			path:            "/v1/log_messages",
			body:            `{"message":"hello"}`,
			reqID:           "rid-1",
			svc:             &logmessagesvcmock.IService{},
			expectedCode:    http.StatusCreated,
			expectedCounter: 1,
			expectedRespID:  1,
		},
		{
			name:            "HandleGet_Success",
			handler:         "get",
			method:          http.MethodGet,
			path:            "/v1/log_messages/42",
			chiParams:       map[string]string{"id": "42"},
			svc:             &logmessagesvcmock.IService{},
			expectedCode:    http.StatusOK,
			expectedCounter: 0, // handleGet does not increment metric
			expectedRespID:  42,
		},
		{
			name:            "HandleUpdate_UsesLangURLParam",
			handler:         "update",
			method:          http.MethodPost,
			path:            "/v1/log_messages/5/es",
			body:            `{"message":"updated"}`,
			chiParams:       map[string]string{"id": "5", "lang": "es"},
			svc:             &logmessagesvcmock.IService{},
			expectedCode:    http.StatusAccepted,
			expectedCounter: 1,
			expectedRespID:  5,
		},
		{
			name:            "HandleUpdate_BadJSON_IncrementsCounterAndReturnsBadRequest",
			handler:         "update",
			method:          http.MethodPost,
			path:            "/v1/log_messages/5",
			body:            "{{invalid-json",
			chiParams:       map[string]string{"id": "5"},
			svc:             &logmessagesvcmock.IService{},
			expectedNot:     []int{http.StatusAccepted, http.StatusOK},
			expectedCounter: 1,
		},
		{
			name:            "HandleCreate_ServiceError_PropagatesError",
			handler:         "create",
			method:          http.MethodPost,
			path:            "/v1/log_messages",
			body:            `{"message":"hello"}`,
			reqID:           "rid-2",
			svc:             &logmessagesvcmock.IService{CreateErr: errors.New("svc fail")},
			expectedNot:     []int{http.StatusCreated},
			expectedCounter: 1,
		},
		{
			name:            "HandleRetrieve_Success",
			handler:         "retrieve",
			method:          http.MethodGet,
			path:            "/v1/log_messages?lang=en&max=0&limit=10&page=1",
			body:            "",
			svc:             &logmessagesvcmock.IService{},
			expectedCode:    http.StatusOK,
			expectedCounter: 1,
		},
		{
			name:            "HandleRetrieve_ServiceError_PropagatesError",
			handler:         "retrieve",
			method:          http.MethodGet,
			path:            "/v1/log_messages?lang=en&max=0&limit=10&page=1",
			body:            "",
			svc:             &logmessagesvcmock.IService{RetrieveErr: errors.New("svc fail")},
			expectedNot:     []int{http.StatusOK},
			expectedCounter: 1,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			counter := prometheus.NewCounter(prometheus.CounterOpts{
				Name: "test_counter_" + strings.ReplaceAll(tt.name, " ", "_"),
				Help: "test counter",
			})

			sc := &LogMessageController{
				log:           nil,
				validate:      nil,
				logMessageSvc: tt.svc,
				stsClient:     nil,
				counterMetric: counter,
			}

			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			// set RequestID when provided (used by handleCreate)
			if tt.reqID != "" {
				req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, tt.reqID))
			}

			// add chi params if any
			if len(tt.chiParams) > 0 {
				req = withChiRouteParams(req, tt.chiParams)
			}

			rr := httptest.NewRecorder()

			switch tt.handler {
			case "create":
				sc.handleCreate(rr, req)
			case "get":
				sc.handleGet(rr, req)
			case "update":
				sc.handleUpdate(rr, req)
			case "retrieve":
				sc.handleRetrieve(rr, req)
			default:
				t.Fatalf("unknown handler %s", tt.handler)
			}

			// check expected exact code when provided
			if tt.expectedCode != 0 {
				if rr.Code != tt.expectedCode {
					t.Fatalf("esperado status %d, obtenido %d, body: %s", tt.expectedCode, rr.Code, rr.Body.String())
				}
			}

			// check expected not-in codes
			if len(tt.expectedNot) > 0 {
				for _, code := range tt.expectedNot {
					if rr.Code == code {
						t.Fatalf("no se esperaba status %d para el caso %s", rr.Code, tt.name)
					}
				}
			}

			// check counter value
			if got := testutil.ToFloat64(counter); got != tt.expectedCounter {
				t.Fatalf("counter esperado %v, obtenido %v", tt.expectedCounter, got)
			}

			// optional: validate response ID when required
			if tt.expectedRespID != 0 {
				var resp log_message.Payload
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("respuesta JSON no válida: %v", err)
				}
				if resp.ID != tt.expectedRespID {
					t.Fatalf("esperado ID %d, obtenido %d", tt.expectedRespID, resp.ID)
				}
			}
		})
	}
}
