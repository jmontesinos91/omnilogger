package logs

import (
	"encoding/json"
	"fmt"
	"github.com/jmontesinos91/omnilogger/domains/pagination"
	"net/http"
	"strconv"
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

	tenantIds, err := payloadToTenantIds(payload)
	if err != nil {
		return nil, err
	}

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
		OldData:     payload.OldData,
		TenantCat:   payload.TenantCat,
		TenantID:    string(tenantIds),
		UserID:      payload.UserID,
		Target:      payload.Target,
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

func ToRepoFilter(filter Filter) logs.Filter {

	from := ((filter.Page * filter.Size) - filter.Size) + 1

	return logs.Filter{
		Level:    filter.Level,
		Message:  filter.Message,
		Provider: filter.Provider,
		Action:   filter.Action,
		Path:     filter.Path,
		Resource: filter.Resource,
		TenantID: filter.TenantID,
		UserID:   filter.UserID,
		Target:   filter.Target,
		StartAt:  filter.StartAt,
		EndAt:    filter.EndAt,
		From:     from,
		Size:     filter.Size,
	}
}

func ToParseFilterRequest(r *http.Request) (Filter, error) {
	var startAt time.Time
	var endAt time.Time
	var format = "2006-01-02T15:04:05"

	query := r.URL.Query()
	provider := query["provider[]"]
	level := query["level[]"]
	action := query["action[]"]
	resource := query.Get("resource")
	path := query.Get("path")
	message := query["message[]"]
	tenantId := query["tenant_id[]"]
	userId := query["user_id[]"]
	target := query["target[]"]
	startAtString := query.Get("start_at")
	endAtString := query.Get("end_at")

	tenantIds, err := strArrToIntArr(tenantId)
	if err != nil {
		return Filter{}, err
	}

	messageIds, err := strArrToIntArr(message)
	if err != nil {
		return Filter{}, err
	}

	if startAtString != "" {
		parsedDate, err := time.Parse(format, startAtString)
		if err != nil {
			return Filter{}, err
		}
		startAt = parsedDate
	}

	if endAtString != "" {
		parsedDate, err := time.Parse(format, endAtString)
		if err != nil {
			return Filter{}, err
		}
		endAt = parsedDate
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
		Provider: provider,
		Message:  messageIds,
		Level:    level,
		Action:   action,
		Path:     path,
		Resource: resource,
		TenantID: tenantIds,
		UserID:   userId,
		Target:   target,
		StartAt:  startAt,
		EndAt:    endAt,
		Filter:   page,
	}, nil
}

func strArrToIntArr(strArr []string) ([]int, error) {
	var intArray []int
	for _, str := range strArr {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("error converting string to int: %v", err)
		}
		intArray = append(intArray, val)
	}
	return intArray, nil
}

func payloadToTenantIds(payload *Payload) ([]byte, error) {
	if payload.TenantCat != "" {
		var items []Item
		err := json.Unmarshal([]byte(payload.TenantCat), &items)
		if err != nil {
			return nil, err
		}

		var ids []int
		for _, item := range items {
			ids = append(ids, item.ID)
		}
		return json.Marshal(ids)
	} else {
		return []byte{}, nil
	}
}
