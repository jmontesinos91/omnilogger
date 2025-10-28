package logs

import (
	"github.com/jmontesinos91/omnilogger/domains/pagination"
	"github.com/jmontesinos91/omnilogger/internal/repositories/logs"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToModel(t *testing.T) {

	type args struct { //nolint:wsl
		payload *Payload
	}
	type expected struct {
		model *logs.Model
	}
	cases := []struct { //nolint:wsl
		name     string
		args     args
		empty    bool
		expected expected
	}{
		{
			name: "Happy Path",
			args: args{
				payload: &Payload{
					IpAddress:  "192.168.0.1",
					ClientHost: "localhost",
					Provider:   "TestProvider",
					Level:      1,
					Message:    100,
					Path:       "/v1/resource",
					Resource:   "RESOURCE",
					Action:     "CREATE",
					Data:       "{}",
					UserID:     "user@example.com",
				},
			},
			empty: false,
			expected: expected{model: &logs.Model{
				ID:         "123e4567-e89b-12d3-a456-426614174000",
				IpAddress:  "192.168.0.1",
				ClientHost: "localhost",
				Provider:   "TestProvider",
				Level:      1,
				Message:    100,
				Path:       "/v1/resource",
				Resource:   "RESOURCE",
				Action:     "CREATE",
				Data:       "{}",
				UserID:     "user@example.com",
			}},
		},
		{
			name: "WithOut TenantCat Information",
			args: args{
				payload: &Payload{
					IpAddress:  "192.168.0.1",
					ClientHost: "localhost",
					Provider:   "TestProvider",
					Level:      1,
					Message:    100,
					Path:       "/v1/resource",
					Resource:   "RESOURCE",
					Action:     "CREATE",
					Data:       "{}",
					UserID:     "12",
				},
			},
			empty: false,
			expected: expected{model: &logs.Model{
				ID:         "123e4567-e89b-12d3-a456-426614174000",
				IpAddress:  "192.168.0.1",
				ClientHost: "localhost",
				Provider:   "TestProvider",
				Level:      1,
				Message:    100,
				Path:       "/v1/resource",
				Resource:   "RESOURCE",
				Action:     "CREATE",
				Data:       "{}",
				UserID:     "12",
			}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			result, err := ToModel(tc.args.payload)

			if !tc.empty {
				assert.Equal(t, tc.expected.model.IpAddress, result.IpAddress)
				assert.Equal(t, tc.expected.model.ClientHost, result.ClientHost)
				assert.Equal(t, tc.expected.model.Provider, result.Provider)
				assert.Equal(t, tc.expected.model.Level, result.Level)
				assert.Equal(t, tc.expected.model.Message, result.Message)
				assert.Equal(t, tc.expected.model.Path, result.Path)
				assert.Equal(t, tc.expected.model.Resource, result.Resource)
				assert.Equal(t, tc.expected.model.Action, result.Action)
				assert.Equal(t, tc.expected.model.Data, result.Data)
				assert.Equal(t, tc.expected.model.UserID, result.UserID)
				assert.NoError(t, err)
			} else {
				assert.NotNil(t, result)
			}

		})
	}
}

func TestToResponse(t *testing.T) {
	type args struct {
		model *logs.Model
	}
	type expected struct {
		response *Response
	}
	cases := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "Happy Path with english logMessage",
			args: args{
				model: &logs.Model{
					ID:          "123e4567-e89b-12d3-a456-426614174000",
					IpAddress:   "192.168.0.1",
					ClientHost:  "localhost",
					Provider:    "TestProvider",
					Level:       1,
					Message:     100,
					Description: "Log message test",
					Path:        "/v1/resource",
					Resource:    "RESOURCE",
					Action:      "CREATE",
					Data:        "{\"key\":\"value\"}",
					UserID:      "user@example.com",
					CreatedAt:   &time.Time{},
				},
			},
			expected: expected{
				response: &Response{
					ID:          "123e4567-e89b-12d3-a456-426614174000",
					IpAddress:   "192.168.0.1",
					ClientHost:  "localhost",
					Provider:    "TestProvider",
					Level:       1,
					Message:     100,
					Description: "Log message test",
					Path:        "/v1/resource",
					Resource:    "RESOURCE",
					Action:      "CREATE",
					Data:        "{\"key\":\"value\"}",
					UserID:      "user@example.com",
					CreatedAt:   &time.Time{},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			result := ToResponse(tc.args.model)

			assert.Equal(t, tc.expected.response.ID, result.ID)
			assert.Equal(t, tc.expected.response.IpAddress, result.IpAddress)
			assert.Equal(t, tc.expected.response.ClientHost, result.ClientHost)
			assert.Equal(t, tc.expected.response.Provider, result.Provider)
			assert.Equal(t, tc.expected.response.Level, result.Level)
			assert.Equal(t, tc.expected.response.Message, result.Message)
			assert.Equal(t, tc.expected.response.Description, result.Description)
			assert.Equal(t, tc.expected.response.Path, result.Path)
			assert.Equal(t, tc.expected.response.Resource, result.Resource)
			assert.Equal(t, tc.expected.response.Action, result.Action)
			assert.Equal(t, tc.expected.response.Data, result.Data)
			assert.Equal(t, tc.expected.response.UserID, result.UserID)
			assert.Equal(t, tc.expected.response.CreatedAt, result.CreatedAt)
		})
	}
}

func TestToRepoFilter(t *testing.T) {

	type args struct {
		filter Filter
	}
	type expected struct {
		repoFilter logs.Filter
	}
	cases := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "Happy Path",
			args: args{
				filter: Filter{
					Message:  []int{100, 101},
					Level:    []string{"INFO", "DEBUG"},
					Provider: []string{"TestProvider1", "TestProvider2"},
					Action:   []string{"CREATE", "UPDATE"},
					Path:     "/v1/resource",
					Resource: "RESOURCE",
					TenantID: []int{1, 2, 3},
					UserID:   []string{"1", "12"},
					Target:   []string{"logs target"},
					StartAt:  time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC),
					EndAt:    time.Date(2024, 11, 16, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: expected{
				repoFilter: logs.Filter{
					Message:  []int{100, 101},
					Level:    []string{"INFO", "DEBUG"},
					Provider: []string{"TestProvider1", "TestProvider2"},
					Action:   []string{"CREATE", "UPDATE"},
					Path:     "/v1/resource",
					Resource: "RESOURCE",
					TenantID: []int{1, 2, 3},
					UserID:   []string{"1", "12"},
					Target:   []string{"logs target"},
					StartAt:  time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC),
					EndAt:    time.Date(2024, 11, 16, 0, 0, 0, 0, time.UTC),
					From:     11,
					Size:     10,
				},
			},
		},
		{
			name: "Edge Case - Empty Fields",
			args: args{
				filter: Filter{
					Message:  []int{},
					Level:    []string{},
					Provider: []string{},
					Action:   []string{},
					Path:     "",
					Resource: "",
					TenantID: []int{},
					UserID:   []string{},
					Target:   []string{},
					StartAt:  time.Time{},
					EndAt:    time.Time{},
				},
			},
			expected: expected{
				repoFilter: logs.Filter{
					Message:  []int{},
					Level:    []string{},
					Provider: []string{},
					Action:   []string{},
					Path:     "",
					Resource: "",
					TenantID: []int{},
					UserID:   []string{},
					Target:   []string{},
					StartAt:  time.Time{},
					EndAt:    time.Time{},
					From:     1,
					Size:     20,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToRepoFilter(tc.args.filter)

			assert.Equal(t, tc.expected.repoFilter.Message, result.Message)
			assert.Equal(t, tc.expected.repoFilter.Level, result.Level)
			assert.Equal(t, tc.expected.repoFilter.Provider, result.Provider)
			assert.Equal(t, tc.expected.repoFilter.Action, result.Action)
			assert.Equal(t, tc.expected.repoFilter.Path, result.Path)
			assert.Equal(t, tc.expected.repoFilter.Resource, result.Resource)
			assert.Equal(t, tc.expected.repoFilter.TenantID, result.TenantID)
			assert.Equal(t, tc.expected.repoFilter.UserID, result.UserID)
			assert.Equal(t, tc.expected.repoFilter.StartAt, result.StartAt)
			assert.Equal(t, tc.expected.repoFilter.EndAt, result.EndAt)
		})
	}
}

func TestToParseFilterRequest(t *testing.T) {
	tests := []struct {
		name        string
		queryParams map[string]string
		expectError bool
		errorMsg    string
		expected    Filter
	}{
		{
			name: "Valid Parameters",
			queryParams: map[string]string{
				"provider[]":  "aws",
				"level[]":     "info",
				"action[]":    "create",
				"resource":    "USER",
				"path":        "/v1/user/auth",
				"message[]":   "1001",
				"tenant_id[]": "1",
				"user_id[]":   "user1",
				"target[]":    "logs target",
				"start_at":    "2024-01-01T00:00:00",
				"end_at":      "2024-01-02T00:00:00",
				"lang":        "en",
				"max":         "10",
				"page":        "1",
			},
			expectError: false,
			expected: Filter{
				Provider: []string{"aws"},
				Level:    []string{"info"},
				Action:   []string{"create"},
				Resource: "USER",
				Path:     "/v1/user/auth",
				Message:  []int{1001},
				TenantID: []int{1},
				UserID:   []string{"user1"},
				Target:   []string{"logs target"},
				StartAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				EndAt:    time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				Filter: pagination.Filter{
					Size: 10,
					Page: 1,
				},
			},
		},
		{
			name: "Missing required parameters",
			queryParams: map[string]string{
				"max":  "10",
				"page": "1",
			},
			expectError: false,
			expected: Filter{
				Filter: pagination.Filter{
					Size: 10,
					Page: 1,
				},
			},
		},
		{
			name: "Invalid tenant ID",
			queryParams: map[string]string{
				"tenant_id[]": "invalid",
			},
			expectError: true,
			errorMsg:    "invalid syntax",
		},
		{
			name: "Invalid start date",
			queryParams: map[string]string{
				"start_at": "invalid-date",
			},
			expectError: true,
			errorMsg:    "parsing time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := url.Values{}
			for key, value := range tt.queryParams {
				query.Set(key, value)
			}
			req := &http.Request{
				URL: &url.URL{RawQuery: query.Encode()},
			}

			fr, err := ToParseFilterRequest(req)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, fr)
			}
		})
	}
}
