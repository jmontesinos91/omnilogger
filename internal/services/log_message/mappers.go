package log_message

import (
	"github.com/jmontesinos91/omnilogger/domains/pagination"
	"github.com/jmontesinos91/omnilogger/internal/repositories/log_message"
	"net/http"
	"strconv"
)

func ToModel(payload *Payload) *log_message.Model {
	return &log_message.Model{
		ID:      payload.ID,
		Message: payload.Message,
		Lang:    payload.Lang,
	}
}

func ToResponse(model *log_message.Model) *Response {
	return &Response{
		ID:      model.ID,
		Message: model.Message,
		Lang:    model.Lang,
	}
}

func ToModelUpdate(model *log_message.Model, payload Payload) *log_message.Model {
	model.Message = payload.Message
	model.Lang = payload.Lang
	return model
}

func ToParseFilterRequest(r *http.Request) (Filter, error) {
	query := r.URL.Query()
	var id *int

	if query.Get("id") != "" {
		idInt, err := strconv.Atoi(query.Get("id"))
		if err == nil {
			return Filter{}, err
		} else {
			id = &idInt
		}
	}

	language := ""

	if language = query.Get("lang"); language == "" {
		language = "en"
	}

	size, err := strconv.Atoi(query.Get("max"))
	if err != nil {
		return Filter{}, err
	}

	pageNumber, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		return Filter{}, err
	}

	page := pagination.Filter{
		Size: size,
		Page: pageNumber,
	}

	return Filter{
		ID:     id,
		Lang:   language,
		Filter: page,
	}, nil
}

func ToRepoFilter(filter Filter) log_message.Filter {
	return log_message.Filter{
		ID:   filter.ID,
		Lang: filter.Lang,
		Size: filter.Size,
	}
}
