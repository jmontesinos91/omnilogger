package log_message

import (
	"github.com/jmontesinos91/omnilogger/internal/repositories/log_message"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToModel(t *testing.T) {

	type args struct { //nolint:wsl
		payload *Payload
	}
	type expected struct {
		model *log_message.Model
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
					ID:      123456,
					Message: "test message payload",
				},
			},
			empty: false,
			expected: expected{model: &log_message.Model{
				ID:      123456,
				Message: "test message payload",
			}},
		},
		{
			name: "WithOut Message Information",
			args: args{
				payload: &Payload{
					ID:      123456,
					Message: "",
				},
			},
			empty: false,
			expected: expected{model: &log_message.Model{
				ID:      123456,
				Message: "",
			}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			result := ToModel(tc.args.payload)

			if !tc.empty {
				assert.Equal(t, tc.expected.model.ID, result.ID)
				assert.Equal(t, tc.expected.model.Message, result.Message)
			} else {
				assert.NotNil(t, result)
			}

		})
	}
}

func TestToResponse(t *testing.T) {
	type args struct {
		model *log_message.Model
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
			name: "Happy Path",
			args: args{
				model: &log_message.Model{
					ID:      123456,
					Message: "Testing message",
				},
			},
			expected: expected{
				response: &Response{
					ID:      123456,
					Message: "Testing message",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			result := ToResponse(tc.args.model)

			assert.Equal(t, tc.expected.response.ID, result.ID)
			assert.Equal(t, tc.expected.response.Message, result.Message)
		})
	}
}

func TestToModelUpdate(t *testing.T) {
	type args struct {
		model   *log_message.Model
		payload Payload
	}
	type expected struct {
		model *log_message.Model
	}
	cases := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "Happy Path",
			args: args{
				model: &log_message.Model{
					ID:      1234,
					Message: "test1",
				},
				payload: Payload{
					ID:      12345,
					Message: "Test 2",
				},
			},
			expected: expected{
				model: &log_message.Model{
					ID:      1234,
					Message: "Test 2",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToModelUpdate(tc.args.model, tc.args.payload)

			assert.Equal(t, tc.expected.model.ID, result.ID)
		})
	}
}
