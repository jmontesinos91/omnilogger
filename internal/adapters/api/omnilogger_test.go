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
		name                 string
		handler              string // "create", "get", "retrieve"
		method               string
		path                 string
		query                string // include leading "?" when non-empty (used for retrieve)
		body                 string
		reqID                string
		chiParams            map[string]string
		mockSvc              *logssvcmock.IService
		expectedCode         int
		expectedNot          []int
		expectedCounter      float64
		expectedRespID       int    // for single resource responses
		expectRetrieveCalled bool   // for retrieve handler
		expectDataID         string // for retrieve success, expected first Data[0].ID
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
			query:           "?page=1&per_page=10&max=10",
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
		{
			name:    "Retrieve_Success",
			handler: "retrieve",
			method:  http.MethodGet,
			path:    "/v1/logs",
			query:   "?page=1&per_page=10&max=10",
			mockSvc: &logssvcmock.IService{
				RetrieveRes: &logs.PaginatedRes{
					Data: []logs.Response{
						{ID: "42", Message: 1},
					},
					Total: 1,
				},
			},
			expectRetrieveCalled: true,
			expectedCounter:      1,
			expectDataID:         "42",
		},
		{
			name:    "Retrieve_ServiceError",
			handler: "retrieve",
			method:  http.MethodGet,
			path:    "/v1/logs",
			query:   "?page=1&max=10",
			mockSvc: &logssvcmock.IService{
				RetrieveErr: errors.New("svc fail"),
			},
			expectRetrieveCalled: true,
			expectedCounter:      1,
		},
		{
			name:                 "Retrieve_InvalidTenantIDParam",
			handler:              "retrieve",
			method:               http.MethodGet,
			path:                 "/v1/logs",
			query:                "?page=1&max=10&tenant_id[]=abc",
			mockSvc:              &logssvcmock.IService{},
			expectRetrieveCalled: false,
			expectedCounter:      1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSvc == nil {
				tt.mockSvc = &logssvcmock.IService{}
			}

			counter := prometheus.NewCounter(prometheus.CounterOpts{
				Name: "test_omnilogger_" + strings.ReplaceAll(tt.name, " ", "_"),
				Help: "test counter",
			})

			sc := &OmniLoggerController{
				log:           ctxLogger,
				validate:      validator.New(),
				logsSvc:       tt.mockSvc,
				counterMetric: counter,
			}

			var req *http.Request
			if tt.handler == "retrieve" || tt.handler == "get" {
				req = httptest.NewRequest(tt.method, tt.path+tt.query, nil)
			} else {
				if tt.body != "" {
					req = httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
				} else {
					req = httptest.NewRequest(tt.method, tt.path, nil)
				}
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
			case "retrieve":
				sc.handleRetrieve(rr, req)
			default:
				t.Fatalf("unknown handler %s", tt.handler)
			}

			// expected status handling
			if tt.expectedCode != 0 {
				if rr.Code != tt.expectedCode {
					t.Fatalf("expected status %d, got %d, body: %s", tt.expectedCode, rr.Code, rr.Body.String())
				}
			}
			if len(tt.expectedNot) > 0 {
				for _, code := range tt.expectedNot {
					if rr.Code == code {
						t.Fatalf("did not expect status %d for case %s, body: %s", rr.Code, tt.name, rr.Body.String())
					}
				}
			}

			// retrieve call expectations (only meaningful for retrieve cases)
			if tt.handler == "retrieve" {
				if tt.expectRetrieveCalled {
					if !tt.mockSvc.RetrieveCalled {
						t.Fatalf("expected Retrieve to be called on the mock service")
					}
				} else {
					if tt.mockSvc.RetrieveCalled {
						t.Fatalf("did not expect Retrieve to be called on the mock service")
					}
				}
			}

			// invalid JSON create should not call Create
			if tt.handler == "create" && tt.body == "{{invalid-json" {
				if tt.mockSvc.CreateCalled {
					t.Fatalf("did not expect Create to be called with malformed payload")
				}
			}

			// counter assertion
			if got := testutil.ToFloat64(counter); got != tt.expectedCounter {
				t.Fatalf("expected counter %v, got %v", tt.expectedCounter, got)
			}

			// response body assertions
			if tt.expectedRespID > 0 && rr.Code == http.StatusOK {
				var resp logs.Response
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("invalid JSON response: %v", err)
				}
				responseID, err := strconv.Atoi(resp.ID)
				if err != nil {
					t.Fatalf("wrong id: %v", err)
				}
				if responseID != tt.expectedRespID {
					t.Fatalf("expected ID %d, got %d", tt.expectedRespID, responseID)
				}
			}

			if tt.expectDataID != "" && rr.Code == http.StatusOK {
				var resp logs.PaginatedRes
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("invalid JSON response: %v", err)
				}
				if len(resp.Data) == 0 {
					t.Fatalf("expected at least one item in paginated response")
				}
				if resp.Data[0].ID != tt.expectDataID {
					t.Fatalf("expected ID %s, got %s", tt.expectDataID, resp.Data[0].ID)
				}
			}
		})
	}
}
