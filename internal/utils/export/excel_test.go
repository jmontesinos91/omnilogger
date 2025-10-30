package export_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/jmontesinos91/omnilogger/internal/utils/export"
	"github.com/jmontesinos91/omnilogger/internal/utils/format"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

// Payment test struct
type Payment struct {
	Amount  float64
	PayDate time.Time
}

// Car test struct
type Car struct {
	Name       string
	Model      string
	FinalPrice float64
	Payments   []*Payment
}

// Owner test struct
type Owner struct {
	Crashes        *int
	CriminalRecord *string
	FirstName      string
	LastName       string
	Acquisition    *time.Time
	Cars           *[]Car
}

func TestDataToExcel(t *testing.T) {

	AcquisitionPtr := time.Date(2020, 1, 15, 10, 0, 0, 0, time.UTC)

	criminalRecordNilText := new(string)
	*criminalRecordNilText = "<nil>"

	criminalRecordPtr := new(string)
	*criminalRecordPtr = "Crashed under alcohol effects"

	crashesPtr := new(int)
	*crashesPtr = 1

	tests := []struct {
		name          string
		excelFileName string
		data          []Owner
		mapper        func(Owner) format.ExcelRow
		expected      [][]interface{}
		expectErr     bool
	}{
		{
			name:          "Happy path - valid data with 3 levels",
			excelFileName: "Excel test",
			data: []Owner{
				{
					Crashes:        nil,
					CriminalRecord: nil,
					FirstName:      "John",
					LastName:       "Doe",
					Acquisition:    &AcquisitionPtr,
					Cars: &[]Car{
						{
							Name:       "Toyota Corolla",
							Model:      "2020",
							FinalPrice: 20000,
							Payments: []*Payment{
								{Amount: 5000, PayDate: time.Date(2020, 2, 10, 12, 0, 0, 0, time.UTC)},
								{Amount: 15000, PayDate: time.Date(2021, 2, 10, 12, 0, 0, 0, time.UTC)},
							},
						},
					},
				},
			},
			mapper: func(owner Owner) format.ExcelRow {
				carRows := []format.ExcelRow{}
				for _, car := range *owner.Cars {
					paymentRows := []format.ExcelRow{}
					for _, payment := range car.Payments {
						paymentRows = append(paymentRows, format.ExcelRow{
							Cells: []interface{}{payment.Amount, payment.PayDate},
						})
					}
					carRows = append(carRows, format.ExcelRow{
						Cells:  []interface{}{car.Name, car.Model, car.FinalPrice},
						Groups: paymentRows,
					})
				}
				return format.ExcelRow{
					Cells:  []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
					Groups: carRows,
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition", "Name", "Model", "FinalPrice", "Amount", "PayDate"},
				{nil, nil, "John", "Doe", "1/15/20 10:00", nil, nil, nil, nil, nil},
				{nil, nil, nil, nil, nil, "Toyota Corolla", "2020", 20000.0, nil, nil},
				{nil, nil, nil, nil, nil, nil, nil, nil, 5000.0, "2/10/20 12:00"},
				{nil, nil, nil, nil, nil, nil, nil, nil, 15000.0, "2/10/21 12:00"},
			},
			expectErr: false,
		},
		{
			name:          "Happy path - valid data with 1 level",
			excelFileName: "Excel test",
			data: []Owner{
				{
					Crashes:        nil,
					CriminalRecord: nil,
					FirstName:      "Jane",
					LastName:       "Smith",
					Acquisition:    &AcquisitionPtr,
				},
			},
			mapper: func(owner Owner) format.ExcelRow {
				return format.ExcelRow{
					Cells: []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition"},
				{nil, nil, "Jane", "Smith", "1/15/20 10:00"},
			},
			expectErr: false,
		},
		{
			name:          "Happy path - valid data with pointers",
			excelFileName: "Excel test",
			data: []Owner{
				{
					Crashes:        crashesPtr,
					CriminalRecord: criminalRecordPtr,
					FirstName:      "Jane",
					LastName:       "Smith",
					Acquisition:    &AcquisitionPtr,
					Cars: &[]Car{
						{
							Name:       "Toyota Corolla",
							Model:      "2020",
							FinalPrice: 20000,
							Payments: []*Payment{
								{Amount: 5000, PayDate: time.Date(2020, 2, 10, 12, 0, 0, 0, time.UTC)},
								{Amount: 15000, PayDate: time.Date(2021, 2, 10, 12, 0, 0, 0, time.UTC)},
							},
						},
					},
				},
			},
			mapper: func(owner Owner) format.ExcelRow {
				carRows := []format.ExcelRow{}
				for _, car := range *owner.Cars {
					paymentRows := []format.ExcelRow{}
					for _, payment := range car.Payments {
						paymentRows = append(paymentRows, format.ExcelRow{
							Cells: []interface{}{payment.Amount, payment.PayDate},
						})
					}
					carRows = append(carRows, format.ExcelRow{
						Cells:  []interface{}{car.Name, car.Model, car.FinalPrice},
						Groups: paymentRows,
					})
				}
				return format.ExcelRow{
					Cells:  []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
					Groups: carRows,
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition", "Name", "Model", "FinalPrice", "Amount", "PayDate"},
				{nil, nil, "Jane", "Smith", "1/15/20 10:00", nil, nil, nil, nil, nil},
				{nil, nil, nil, nil, nil, "Toyota Corolla", "2020", 20000.0, nil, nil},
				{nil, nil, nil, nil, nil, nil, nil, nil, 5000.0, "2/10/20 12:00"},
				{nil, nil, nil, nil, nil, nil, nil, nil, 15000.0, "2/10/21 12:00"},
			},
			expectErr: false,
		},
		{
			name:          "Happy path - correct sheet name",
			excelFileName: "Excel test",
			data: []Owner{
				{
					Crashes:        crashesPtr,
					CriminalRecord: criminalRecordNilText,
					FirstName:      "Jane",
					LastName:       "Smith",
					Acquisition:    nil,
				},
			},
			mapper: func(owner Owner) format.ExcelRow {
				return format.ExcelRow{
					Cells: []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition"},
				{nil, nil, "Jane", "Smith", nil},
			},
			expectErr: false,
		},
		{
			name:          "Error - long sheet name",
			excelFileName: "Excel test filenametoolargetobeprocessedbygoexelizelib",
			data: []Owner{
				{
					Crashes:        crashesPtr,
					CriminalRecord: criminalRecordNilText,
					FirstName:      "Jane",
					LastName:       "Smith",
					Acquisition:    nil,
				},
			},
			mapper: func(owner Owner) format.ExcelRow {
				return format.ExcelRow{
					Cells: []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition"},
				{nil, nil, "Jane", "Smith", nil},
			},
			expectErr: true,
		},
		{
			name:          "Error - special characters",
			excelFileName: "excel[file]{}test\\/?*",
			data: []Owner{
				{
					Crashes:        crashesPtr,
					CriminalRecord: criminalRecordNilText,
					FirstName:      "Jane",
					LastName:       "Smith",
					Acquisition:    nil,
				},
			},
			mapper: func(owner Owner) format.ExcelRow {
				return format.ExcelRow{
					Cells: []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition"},
				{nil, nil, "Jane", "Smith", nil},
			},
			expectErr: true,
		},
		{
			name:          "Error - empty sheet name",
			excelFileName: "",
			data: []Owner{
				{
					Crashes:        crashesPtr,
					CriminalRecord: criminalRecordNilText,
					FirstName:      "Jane",
					LastName:       "Smith",
					Acquisition:    nil,
				},
			},
			mapper: func(owner Owner) format.ExcelRow {
				return format.ExcelRow{
					Cells: []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition"},
				{nil, nil, "Jane", "Smith", nil},
			},
			expectErr: true,
		},
		{
			name:          "Empty data",
			excelFileName: "excel test",
			data:          nil,
			mapper:        func(owner Owner) format.ExcelRow { return format.ExcelRow{} },
			expected:      [][]interface{}{},
			expectErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the DataToExcel function
			excelBytes, err := export.DataToExcel(tt.excelFileName, tt.data, tt.mapper)

			// Check for errors
			if tt.expectErr {
				assert.Error(t, err, "Expected an error but got none")
				assert.Nil(t, excelBytes, "Expected nil bytes on error")
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, excelBytes, "Expected valid Excel bytes")

			// Open the Excel file in memory
			f, err := excelize.OpenReader(bytes.NewReader(excelBytes))
			assert.NoError(t, err)

			// Read rows from the Excel sheet
			rows, err := f.GetRows(tt.excelFileName)
			assert.NoError(t, err)

			// Validate the rows match expected
			assert.Equal(t, len(tt.expected), len(rows), "Row count mismatch")
			for i, expectedRow := range tt.expected {
				actualRow := rows[i]
				for j, expectedCell := range expectedRow {
					if expectedCell != nil {
						assert.Equal(t, fmt.Sprintf("%v", expectedCell), actualRow[j], "Cell mismatch")
					}
				}
			}
		})
	}
}

func TestDataToExcel_WithPointerRoot(t *testing.T) {
	AcquisitionPtr := time.Date(2020, 1, 15, 10, 0, 0, 0, time.UTC)

	criminalRecordNilText := new(string)
	*criminalRecordNilText = "<nil>"

	criminalRecordPtr := new(string)
	*criminalRecordPtr = "Crashed under alcohol effects"

	crashesPtr := new(int)
	*crashesPtr = 1

	tests := []struct {
		name          string
		excelFileName string
		data          []*Owner
		mapper        func(*Owner) format.ExcelRow
		expected      [][]interface{}
		expectErr     bool
	}{
		{
			name:          "Happy path - valid data with 3 levels",
			excelFileName: "Excel test",
			data: []*Owner{
				{
					Crashes:        nil,
					CriminalRecord: nil,
					FirstName:      "John",
					LastName:       "Doe",
					Acquisition:    &AcquisitionPtr,
					Cars: &[]Car{
						{
							Name:       "Toyota Corolla",
							Model:      "2020",
							FinalPrice: 20000,
							Payments: []*Payment{
								{Amount: 5000, PayDate: time.Date(2020, 2, 10, 12, 0, 0, 0, time.UTC)},
								{Amount: 15000, PayDate: time.Date(2021, 2, 10, 12, 0, 0, 0, time.UTC)},
							},
						},
					},
				},
			},
			mapper: func(owner *Owner) format.ExcelRow {
				carRows := []format.ExcelRow{}
				for _, car := range *owner.Cars {
					paymentRows := []format.ExcelRow{}
					for _, payment := range car.Payments {
						paymentRows = append(paymentRows, format.ExcelRow{
							Cells: []interface{}{payment.Amount, payment.PayDate},
						})
					}
					carRows = append(carRows, format.ExcelRow{
						Cells:  []interface{}{car.Name, car.Model, car.FinalPrice},
						Groups: paymentRows,
					})
				}
				return format.ExcelRow{
					Cells:  []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
					Groups: carRows,
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition", "Name", "Model", "FinalPrice", "Amount", "PayDate"},
				{nil, nil, "John", "Doe", "1/15/20 10:00", nil, nil, nil, nil, nil},
				{nil, nil, nil, nil, nil, "Toyota Corolla", "2020", 20000.0, nil, nil},
				{nil, nil, nil, nil, nil, nil, nil, nil, 5000.0, "2/10/20 12:00"},
				{nil, nil, nil, nil, nil, nil, nil, nil, 15000.0, "2/10/21 12:00"},
			},
			expectErr: false,
		},
		{
			name:          "Happy path - valid data with 1 level",
			excelFileName: "Excel test",
			data: []*Owner{
				{
					Crashes:        nil,
					CriminalRecord: nil,
					FirstName:      "Jane",
					LastName:       "Smith",
					Acquisition:    &AcquisitionPtr,
				},
			},
			mapper: func(owner *Owner) format.ExcelRow {
				return format.ExcelRow{
					Cells: []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition"},
				{nil, nil, "Jane", "Smith", "1/15/20 10:00"},
			},
			expectErr: false,
		},
		{
			name:          "Happy path - valid data with pointers",
			excelFileName: "Excel test",
			data: []*Owner{
				{
					Crashes:        crashesPtr,
					CriminalRecord: criminalRecordPtr,
					FirstName:      "Jane",
					LastName:       "Smith",
					Acquisition:    &AcquisitionPtr,
					Cars: &[]Car{
						{
							Name:       "Toyota Corolla",
							Model:      "2020",
							FinalPrice: 20000,
							Payments: []*Payment{
								{Amount: 5000, PayDate: time.Date(2020, 2, 10, 12, 0, 0, 0, time.UTC)},
								{Amount: 15000, PayDate: time.Date(2021, 2, 10, 12, 0, 0, 0, time.UTC)},
							},
						},
					},
				},
			},
			mapper: func(owner *Owner) format.ExcelRow {
				carRows := []format.ExcelRow{}
				for _, car := range *owner.Cars {
					paymentRows := []format.ExcelRow{}
					for _, payment := range car.Payments {
						paymentRows = append(paymentRows, format.ExcelRow{
							Cells: []interface{}{payment.Amount, payment.PayDate},
						})
					}
					carRows = append(carRows, format.ExcelRow{
						Cells:  []interface{}{car.Name, car.Model, car.FinalPrice},
						Groups: paymentRows,
					})
				}
				return format.ExcelRow{
					Cells:  []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
					Groups: carRows,
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition", "Name", "Model", "FinalPrice", "Amount", "PayDate"},
				{nil, nil, "Jane", "Smith", "1/15/20 10:00", nil, nil, nil, nil, nil},
				{nil, nil, nil, nil, nil, "Toyota Corolla", "2020", 20000.0, nil, nil},
				{nil, nil, nil, nil, nil, nil, nil, nil, 5000.0, "2/10/20 12:00"},
				{nil, nil, nil, nil, nil, nil, nil, nil, 15000.0, "2/10/21 12:00"},
			},
			expectErr: false,
		},
		{
			name:          "Happy path - correct sheet name",
			excelFileName: "Excel test",
			data: []*Owner{
				{
					Crashes:        crashesPtr,
					CriminalRecord: criminalRecordNilText,
					FirstName:      "Jane",
					LastName:       "Smith",
					Acquisition:    nil,
				},
			},
			mapper: func(owner *Owner) format.ExcelRow {
				return format.ExcelRow{
					Cells: []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition"},
				{nil, nil, "Jane", "Smith", nil},
			},
			expectErr: false,
		},
		{
			name:          "Error - long sheet name",
			excelFileName: "Excel test filenametoolargetobeprocessedbygoexelizelib",
			data: []*Owner{
				{
					Crashes:        crashesPtr,
					CriminalRecord: criminalRecordNilText,
					FirstName:      "Jane",
					LastName:       "Smith",
					Acquisition:    nil,
				},
			},
			mapper: func(owner *Owner) format.ExcelRow {
				return format.ExcelRow{
					Cells: []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition"},
				{nil, nil, "Jane", "Smith", nil},
			},
			expectErr: true,
		},
		{
			name:          "Error - special characters",
			excelFileName: "excel[file]{}test\\/?*",
			data: []*Owner{
				{
					Crashes:        crashesPtr,
					CriminalRecord: criminalRecordNilText,
					FirstName:      "Jane",
					LastName:       "Smith",
					Acquisition:    nil,
				},
			},
			mapper: func(owner *Owner) format.ExcelRow {
				return format.ExcelRow{
					Cells: []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition"},
				{nil, nil, "Jane", "Smith", nil},
			},
			expectErr: true,
		},
		{
			name:          "Error - empty sheet name",
			excelFileName: "",
			data: []*Owner{
				{
					Crashes:        crashesPtr,
					CriminalRecord: criminalRecordNilText,
					FirstName:      "Jane",
					LastName:       "Smith",
					Acquisition:    nil,
				},
			},
			mapper: func(owner *Owner) format.ExcelRow {
				return format.ExcelRow{
					Cells: []interface{}{owner.Crashes, owner.CriminalRecord, owner.FirstName, owner.LastName, owner.Acquisition},
				}
			},
			expected: [][]interface{}{
				{"Crashes", "CriminalRecord", "FirstName", "LastName", "Acquisition"},
				{nil, nil, "Jane", "Smith", nil},
			},
			expectErr: true,
		},
		{
			name:          "Empty data",
			excelFileName: "excel test",
			data:          nil,
			mapper:        func(owner *Owner) format.ExcelRow { return format.ExcelRow{} },
			expected:      [][]interface{}{},
			expectErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the DataToExcel function
			excelBytes, err := export.DataToExcel(tt.excelFileName, tt.data, tt.mapper)

			// Check for errors
			if tt.expectErr {
				assert.Error(t, err, "Expected an error but got none")
				assert.Nil(t, excelBytes, "Expected nil bytes on error")
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, excelBytes, "Expected valid Excel bytes")

			// Open the Excel file in memory
			f, err := excelize.OpenReader(bytes.NewReader(excelBytes))
			assert.NoError(t, err)

			// Read rows from the Excel sheet
			rows, err := f.GetRows(tt.excelFileName)
			assert.NoError(t, err)

			// Validate the rows match expected
			assert.Equal(t, len(tt.expected), len(rows), "Row count mismatch")
			for i, expectedRow := range tt.expected {
				actualRow := rows[i]
				for j, expectedCell := range expectedRow {
					if expectedCell != nil {
						assert.Equal(t, fmt.Sprintf("%v", expectedCell), actualRow[j], "Cell mismatch")
					}
				}
			}
		})
	}
}
