package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageFilter_SanitizePageFilter(t *testing.T) {
	type fields struct {
		From     int
		Size     int
		SortDesc bool
		SortBy   string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		asserts func(*testing.T, error, *Filter) bool
	}{
		{
			name: "Check if From has MinimumFromValue",
			fields: fields{
				From: MinimumFromValue - 1,
			},
			wantErr: false,
			asserts: func(t *testing.T, err error, pf *Filter) bool {
				return assert.Equal(t, MinimumFromValue, pf.Offset)
			},
		},
		{
			name: "Check if Size has DefaultSizeValue",
			fields: fields{
				Size: MinimumSizeValue - 1,
			},
			wantErr: false,
			asserts: func(t *testing.T, err error, pf *Filter) bool {
				return assert.Equal(t, DefaultSizeValue, pf.Size)
			},
		},
		{
			name: "Check if Size has MaximumSizeValue",
			fields: fields{
				Size: MaximumSizeValue + 1,
			},
			wantErr: false,
			asserts: func(t *testing.T, err error, pf *Filter) bool {
				return assert.Equal(t, MaximumSizeValue, pf.Size)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				Offset:   tt.fields.From,
				Size:     tt.fields.Size,
				SortDesc: tt.fields.SortDesc,
				SortBy:   tt.fields.SortBy,
			}
			err := f.SanitizePageFilter()
			if (err != nil) != tt.wantErr {
				t.Errorf("PageFilter.SanitizePageFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.asserts(t, err, f) {
				t.Errorf("Assert error on test = '%v'", tt.name)
			}
		})
	}
}
