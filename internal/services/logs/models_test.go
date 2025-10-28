package logs

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
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
		"user_id": "12"
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
	assert.Equal(t, "12", response.UserID)
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
		"userId": "12",
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
