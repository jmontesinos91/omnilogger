package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/omnilogger/internal/services/logs"
	"github.com/jmontesinos91/omnilogger/internal/services/logs/logssvcmock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestOmniLoggerController_TableDriven(t *testing.T) {
	ctxLogger := logger.NewContextLogger("TestControllerOmniLogger", "debug", logger.TextFormat)
	type tc struct {
		name            string
		handler         string // "create", "get"
		method          string
		path            string
		body            string
		reqID           string
		chiParams       map[string]string
		mockSvc         *logssvcmock.IService
		expectedCode    int
		expectedNot     []int
		expectedCounter float64
		expectedRespID  int // empty -> don't check
	}

	tests := []tc{
		{
			name:            "HandleCreate_Success",
			handler:         "create",
			method:          http.MethodPost,
			path:            "/v1/logs",
			body:            `{"message":1}`,
			reqID:           "rid-2",
			mockSvc:         &logssvcmock.IService{},
			expectedCode:    http.StatusCreated,
			expectedCounter: 1,
			expectedRespID:  1,
		},
		{
			name:            "HandleGet_Success_IncrementsCounter",
			handler:         "get",
			method:          http.MethodGet,
			path:            "/v1/logs/abc",
			chiParams:       map[string]string{"id": "3"},
			mockSvc:         &logssvcmock.IService{},
			expectedCode:    http.StatusOK,
			expectedCounter: 1,
			expectedRespID:  3,
		},
		{
			name:            "HandleCreate_BadJSON_ReturnsBadRequest",
			handler:         "create",
			method:          http.MethodPost,
			path:            "/v1/logs",
			body:            "{{invalid-json",
			reqID:           "rid-2",
			mockSvc:         &logssvcmock.IService{},
			expectedNot:     []int{http.StatusCreated},
			expectedCounter: 1,
		},
		{
			name:            "HandleCreate_ServiceError_PropagatesError",
			handler:         "create",
			method:          http.MethodPost,
			path:            "/v1/logs",
			body:            `{"message":"hello"}`,
			reqID:           "rid-3",
			mockSvc:         &logssvcmock.IService{CreateErr: errors.New("svc fail")},
			expectedNot:     []int{http.StatusCreated},
			expectedCounter: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSvc == nil {
				tt.mockSvc = &logssvcmock.IService{}
			}

			counter := prometheus.NewCounter(prometheus.CounterOpts{
				Name: "test_omnilogger_counter_" + strings.ReplaceAll(tt.name, " ", "_"),
				Help: "test counter",
			})

			// Inicializar ContextLogger con un *logrus.Logger para evitar nil pointer deref.
			sc := &OmniLoggerController{
				log:           ctxLogger,
				validate:      validator.New(),
				logsSvc:       tt.mockSvc,
				stsClient:     nil,
				counterMetric: counter,
			}

			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			if tt.reqID != "" {
				req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, tt.reqID))
			}
			if len(tt.chiParams) > 0 {
				req = withChiRouteParams(req, tt.chiParams)
			}

			rr := httptest.NewRecorder()

			switch tt.handler {
			case "create":
				sc.handleCreate(rr, req)
			case "get":
				sc.handleGetLog(rr, req)
			default:
				t.Fatalf("unknown handler %s", tt.handler)
			}

			if tt.expectedCode != 0 {
				if rr.Code != tt.expectedCode {
					t.Fatalf("esperado status %d, obtenido %d, body: %s", tt.expectedCode, rr.Code, rr.Body.String())
				}
			}

			if len(tt.expectedNot) > 0 {
				for _, code := range tt.expectedNot {
					if rr.Code == code {
						t.Fatalf("no se esperaba status %d para el caso %s, body: %s", rr.Code, tt.name, rr.Body.String())
					}
				}
			}

			if got := testutil.ToFloat64(counter); got != tt.expectedCounter {
				t.Fatalf("counter esperado %v, obtenido %v", tt.expectedCounter, got)
			}

			if tt.expectedRespID > 0 {
				var resp logs.Response
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("respuesta JSON no válida: %v", err)
				}
				responseID, err := strconv.Atoi(resp.ID)
				if err != nil {
					t.Fatalf("wrong id: %v", err)
				}
				if responseID != tt.expectedRespID {
					t.Fatalf("esperado ID %d, obtenido %d", tt.expectedRespID, responseID)
				}
			}

			if tt.handler == "create" && tt.body == "{{invalid-json" {
				if tt.mockSvc.CreateCalled {
					t.Fatalf("no se esperaba llamada a Create con payload malformado")
				}
			}
		})
	}
}
