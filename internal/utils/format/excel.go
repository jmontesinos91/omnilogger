package format

// ExcelRow represents a generic struct to map information in excel spreadsheet
type ExcelRow struct {
	Cells  []interface{}
	Groups []ExcelRow
}
