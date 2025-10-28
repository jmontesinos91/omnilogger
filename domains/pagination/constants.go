package pagination

// MinimumFromValue Minimum value for 'from' parameter that could have
const MinimumFromValue = 0

// Pagination size limits and defaults values
const (

	// MinimumSizeValue Minimum value for 'size' parameter that could have
	MinimumSizeValue = 1
	// MaximumSizeValue Maximum value for `size` parameter that could have
	MaximumSizeValue = 100
	// DefaultSizeValue Default value given to `size` parameter if the given one is not valid
	DefaultSizeValue = 10
)
