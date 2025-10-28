package pagination

import "time"

// Filter is the base filter model
type Filter struct {
	QParam       string    `json:"q"`
	Page         int       `json:"page"`
	Size         int       `json:"max"`
	Offset       int       `json:"offset"`
	SortBy       string    `json:"sortBy"`
	SortDesc     bool      `json:"sortDesc"`
	CreatedAtMin time.Time `json:"createdAtMin"`
	CreatedAtMax time.Time `json:"createdAtMax"`
	UpdatedAtMin time.Time `json:"updatedAtMin"`
	UpdatedAtMax time.Time `json:"updatedAtMax"`
}

// SanitizePageFilter Handles the sanitization for the values in query parameters
func (f *Filter) SanitizePageFilter() error {
	// Set default offset
	if f.Offset < MinimumFromValue {
		f.Offset = MinimumFromValue
	}
	// Set default size
	if f.Size < MinimumSizeValue {
		f.Size = DefaultSizeValue
	} else if f.Size > MaximumSizeValue {
		f.Size = MaximumSizeValue
	}

	return nil
}

// PaginatedRes generic model to paginate all API responses
type PaginatedRes struct {
	Data        interface{} `json:"data"`
	CurrentPage int         `json:"currentPage"`
	Pages       int         `json:"pages"`
	Total       int         `json:"total"`
}
