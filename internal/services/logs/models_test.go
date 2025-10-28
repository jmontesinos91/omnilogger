package logs

import (
	"encoding/json"
	"github.com/jmontesinos91/omnilogger/domains/pagination"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnmarshalPayload(t *testing.T) {

	payloadString := `{
		"ip_address": "192.168.0.1",
		"client_host": "localhost",
		"provider": "TestProvider",
		"level": 1,
		"message": 100,
		"description": "Test Description",
		"path": "/v1/resource",
		"resource": "RESOURCE",
		"action": "CREATE",
		"data": "{\"key\":\"value\"}",
		"old_data": "{\"key\":\"old_value\"}",
		"tenant_cat": "[{\"id\":1,\"name\":\"test\"}]",
		"user_id": "12",
		"target": "LOGS TARGET"
	}`

	var response Payload

	err := json.Unmarshal([]byte(payloadString), &response)

	assert.NoError(t, err)
	assert.Equal(t, "192.168.0.1", response.IpAddress)
	assert.Equal(t, "localhost", response.ClientHost)
	assert.Equal(t, "TestProvider", response.Provider)
	assert.Equal(t, 1, response.Level)
	assert.Equal(t, 100, response.Message)
	assert.Equal(t, "Test Description", response.Description)
	assert.Equal(t, "/v1/resource", response.Path)
	assert.Equal(t, "RESOURCE", response.Resource)
	assert.Equal(t, "CREATE", response.Action)
	assert.Equal(t, "{\"key\":\"value\"}", response.Data)
	assert.Equal(t, "{\"key\":\"old_value\"}", response.OldData)
	assert.Equal(t, "[{\"id\":1,\"name\":\"test\"}]", response.TenantCat)
	assert.Equal(t, "12", response.UserID)
	assert.Equal(t, "LOGS TARGET", response.Target)
}

func TestUnmarshalResponse(t *testing.T) {

	responseString := `{
		"id": "123e4567-e89b-12d3-a456-426614174000",
		"ipAddress": "192.168.0.1",
		"clientHost": "localhost",
		"provider": "TestProvider",
		"level": 1,
		"message": 100,
		"description": "Test Description",
		"path": "/v1/resource",
		"resource": "RESOURCE",
		"action": "CREATE",
		"data": "{\"key\":\"value\"}",
		"oldData": "{\"key\":\"old_value\"}",
		"tenantCat": "[{\"id\":1,\"name\":\"test\"}]",
		"userId": "12",
		"target": "LOGS TARGET",
		"createdAt": "2024-11-15T12:00:00Z"
	}`

	var response Response

	err := json.Unmarshal([]byte(responseString), &response)

	assert.NoError(t, err)
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", response.ID)
	assert.Equal(t, "192.168.0.1", response.IpAddress)
	assert.Equal(t, "localhost", response.ClientHost)
	assert.Equal(t, "TestProvider", response.Provider)
	assert.Equal(t, 1, response.Level)
	assert.Equal(t, 100, response.Message)
	assert.Equal(t, "Test Description", response.Description)
	assert.Equal(t, "/v1/resource", response.Path)
	assert.Equal(t, "RESOURCE", response.Resource)
	assert.Equal(t, "CREATE", response.Action)
	assert.Equal(t, "{\"key\":\"value\"}", response.Data)
	assert.Equal(t, "12", response.UserID)
	assert.NotNil(t, response.CreatedAt)
}

func TestFilter(t *testing.T) {

	startAt, _ := time.Parse(time.RFC3339, "2024-11-15T12:00:00Z")
	endAt, _ := time.Parse(time.RFC3339, "2024-11-16T12:00:00Z")
	filter := Filter{
		Level:    []string{"INFO", "ERROR"},
		Message:  []int{100, 200},
		Provider: []string{"ProviderA", "ProviderB"},
		Action:   []string{"CREATE", "UPDATE"},
		Path:     "/v1/resource",
		Resource: "RESOURCE",
		TenantID: []int{1, 2},
		UserID:   []string{"1", "2"},
		Target:   []string{"LOGS TARGET"},
		StartAt:  startAt,
		EndAt:    endAt,
		Filter: pagination.Filter{
			Page:     1,
			Size:     20,
			Offset:   0,
			SortBy:   "created_at",
			SortDesc: true,
		},
	}

	jsonData, err := json.Marshal(filter)
	assert.NoError(t, err)

	var deserializedFilter Filter
	err = json.Unmarshal(jsonData, &deserializedFilter)

	assert.NoError(t, err)
	assert.Equal(t, filter, deserializedFilter)
}

func TestPaginatedResSerialization(t *testing.T) {

	jsonString := `{
		"data": [
			{
				"id": "123e4567-e89b-12d3-a456-426614174000",
				"message": 1,
				"provider": "ProviderA"
			},
			{
				"id": "223e4567-e89b-12d3-a456-426614174001",
				"message": 2,
				"provider": "ProviderB"
			}
		],
		"max": 10,
		"total": 25,
		"currentPage": 2
	}`

	var paginatedRes PaginatedRes

	err := json.Unmarshal([]byte(jsonString), &paginatedRes)

	assert.NoError(t, err)
	assert.Len(t, paginatedRes.Data, 2)

	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", paginatedRes.Data[0].ID)
	assert.Equal(t, 1, paginatedRes.Data[0].Message)
	assert.Equal(t, "ProviderA", paginatedRes.Data[0].Provider)

	assert.Equal(t, "223e4567-e89b-12d3-a456-426614174001", paginatedRes.Data[1].ID)
	assert.Equal(t, 2, paginatedRes.Data[1].Message)
	assert.Equal(t, "ProviderB", paginatedRes.Data[1].Provider)

	assert.Equal(t, 10, paginatedRes.Size)
	assert.Equal(t, 25, paginatedRes.Total)
	assert.Equal(t, 2, paginatedRes.Page)
}
