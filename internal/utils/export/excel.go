package export

import (
	"bytes"
	"fmt"
	"reflect"
	"time"

	"github.com/jmontesinos91/omnilogger/internal/utils/format"
	"github.com/xuri/excelize/v2"
)

// DataToExcel returns ExcelRow information represented in byte array to be sent via acted-stream
func DataToExcel[T any](sheetName string, data []T, mapperOf func(T) format.ExcelRow) ([]byte, error) {

	var buffer bytes.Buffer

	// WriteRow is a recursive function to print data into the spreadsheet in different levels
	var writeRow func(row format.ExcelRow, level int) error

	f := excelize.NewFile()

	// Wrap f.Close() in a closure
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing the file: %v\n", err)
		}
	}()

	sheet, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	// Date style format
	dateStyle, err := f.NewStyle(&excelize.Style{NumFmt: 22})
	if err != nil {
		return nil, err
	}

	// Header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  14,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#1F4E78"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	if err := f.SetColWidth(sheetName, "A", "Z", 40); err != nil {
		return nil, err
	}

	// Validate empty data
	if len(data) == 0 {
		if err := f.Write(&buffer); err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	}

	// Starter row index
	rowIdx := 1

	// Write struct field names as headers (based on reflection of the first element's type)
	if len(data) > 0 {

		// Get first element of data
		typ := reflect.TypeOf(data[0])
		if typ.Kind() == reflect.Ptr {
			// validate if the first element it's a pointer and get its value
			typ = typ.Elem()
		}
		headers := ExtractHeadersRecursively(typ)

		// Save initial headers
		headerRow := format.ExcelRow{Cells: make([]interface{}, len(headers))}
		for i, header := range headers {
			headerRow.Cells[i] = header
		}
		colIdx := 0
		for _, cell := range headerRow.Cells {
			cellPosition := fmt.Sprintf("%s%d", string(rune('A'+colIdx)), rowIdx)
			if cell != nil {
				if err := f.SetCellValue(sheetName, cellPosition, cell); err != nil {
					return nil, err
				}
				if err := f.SetCellStyle(sheetName, cellPosition, cellPosition, headerStyle); err != nil {
					return nil, err
				}
			}
			colIdx++
		}
		rowIdx++
	}

	writeRow = func(row format.ExcelRow, level int) error {
		// Write the row cells
		colIdx := 0
		for _, cell := range row.Cells {

			// Calculate the cell position based on column and level
			cellPosition := fmt.Sprintf("%s%d", string(rune('A'+colIdx+level)), rowIdx)

			if cell != nil && !(reflect.ValueOf(cell).Kind() == reflect.Ptr && reflect.ValueOf(cell).IsNil()) { //nolint:staticcheck

				v := reflect.ValueOf(cell)
				if v.Kind() == reflect.Ptr && !v.IsNil() {
					// get the pointer value and sets it to cell
					cell = v.Elem().Interface()
				}

				switch v := cell.(type) {
				case time.Time:
					if err := f.SetCellValue(sheetName, cellPosition, v); err != nil {
						return err
					}

					// Set date format
					if err := f.SetCellStyle(sheetName, cellPosition, cellPosition, dateStyle); err != nil {
						return err
					}
				case *time.Time:
					if v != nil {
						if err := f.SetCellValue(sheetName, cellPosition, *v); err != nil {
							return err
						}

						// Set date format
						if err := f.SetCellStyle(sheetName, cellPosition, cellPosition, dateStyle); err != nil {
							return err
						}
					}
				default:
					if cell != "<nil>" {
						if err := f.SetCellValue(sheetName, cellPosition, cell); err != nil {
							return err
						}
					}

				}

			}
			colIdx++
		}
		rowIdx++

		// Count all columns to maintain proper alignment
		columnsUsed := len(row.Cells)

		// Write subgroups
		for _, group := range row.Groups {
			if err := writeRow(group, level+columnsUsed); err != nil {
				return err
			}
		}
		return nil
	}

	// Process all generic data
	for _, item := range data {
		row := mapperOf(item)
		if err := writeRow(row, 0); err != nil {
			return nil, err
		}
	}

	f.SetActiveSheet(sheet)

	// f.Write writes all spreadsheet bytes into buffer
	if err := f.Write(&buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ExtractHeadersRecursively(typ reflect.Type) []string {
	var headers []string

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Check if field is exported
		if !field.IsExported() {
			continue
		}

		// If field is time.Time type just get field name
		if field.Type.PkgPath() == "time" && field.Type.Name() == "Time" {
			headers = append(headers, field.Name)
			continue
		}

		// If field is slice and the first element is a struct get all its fields []customStruct
		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
			subHeaders := ExtractHeadersRecursively(field.Type.Elem())
			headers = append(headers, subHeaders...)

			// If field is a pointer to a slice of structs, get all its fields *[]customStruct
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Slice && field.Type.Elem().Elem().Kind() == reflect.Struct {
			subHeaders := ExtractHeadersRecursively(field.Type.Elem().Elem())
			headers = append(headers, subHeaders...)
			continue
			// If field is a slice of pointers to structs, handle it recursively []*customStruct
		} else if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Ptr && field.Type.Elem().Elem().Kind() == reflect.Struct {
			subHeaders := ExtractHeadersRecursively(field.Type.Elem().Elem())
			headers = append(headers, subHeaders...)
			continue

			// If field is struct get all field names
		} else if field.Type.Kind() == reflect.Struct {
			subHeaders := ExtractHeadersRecursively(field.Type)
			headers = append(headers, subHeaders...)

		} else {
			// Add field name
			headers = append(headers, field.Name)
		}
	}

	return headers
}
