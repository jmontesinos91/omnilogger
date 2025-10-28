package logs

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jmontesinos91/omnilogger/internal/repositories/logs"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func ToModel(payload *Payload) (*logs.Model, error) {
	date := time.Now().UTC()
	return &logs.Model{
		ID:          uuid.NewString(),
		IpAddress:   payload.IpAddress,
		ClientHost:  payload.ClientHost,
		Provider:    payload.Provider,
		Level:       payload.Level,
		Message:     payload.Message,
		Description: payload.Description,
		Path:        payload.Path,
		Resource:    payload.Resource,
		Action:      payload.Action,
		Data:        payload.Data,
		UserID:      payload.UserID,
		CreatedAt:   &date,
	}, nil
}

func ToResponse(model *logs.Model) *Response {
	data := json.RawMessage(model.Data)

	return &Response{
		ID:          model.ID,
		IpAddress:   model.IpAddress,
		ClientHost:  model.ClientHost,
		Provider:    model.Provider,
		Level:       model.Level,
		Message:     model.Message,
		Description: model.Description,
		Path:        model.Path,
		Resource:    model.Resource,
		Action:      model.Action,
		Data:        string(data),
		UserID:      model.UserID,
		CreatedAt:   model.CreatedAt,
	}
}
