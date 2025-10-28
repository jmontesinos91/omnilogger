package log_message

import (
	"github.com/jmontesinos91/omnilogger/internal/repositories/log_message"
)

func ToModel(payload *Payload) *log_message.Model {
	return &log_message.Model{
		ID:      payload.ID,
		Message: payload.Message,
	}
}

func ToResponse(model *log_message.Model) *Response {
	return &Response{
		ID:      model.ID,
		Message: model.Message,
	}
}

func ToModelUpdate(model *log_message.Model, payload Payload) *log_message.Model {
	model.Message = payload.Message
	return model
}
