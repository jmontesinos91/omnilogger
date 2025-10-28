package logs

import (
	"github.com/jmontesinos91/omnilogger/internal/repositories/logs"
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
