package spreadsheet

import (
	"fmt"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
)

type spreadSheet struct {
	srv *sheets.Service
}

type SpreadSheet interface {
	GetRangeValues(spreadId, sheetRange string) ([][]interface{}, error)
	CreateSpreadSheet(title string) (string, error)
	BatchWrite(spreadsheetId string, records [][]interface{}) error
}

// NewSpreadsheet creates a new Spreadsheet object with the given JSON key file path.
// It initializes the necessary service for accessing the Google Sheets API.
//
// jsonPath (string): The file path to the JSON key file.
//
// Returns:
//
//	SpreadSheet: A Spreadsheet object that implements the GetRangeValues method.
func NewSpreadsheet(client *http.Client) SpreadSheet {
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Failed to create config %v", err)
	}
	return &spreadSheet{srv: srv}
}

// GetRangeValues retrieves the values from a specific range in a Google Spreadsheet.
//
// spreadId (string): The ID of the Spreadsheet.
// sheetRange (string): The range to retrieve values from (e.g., "Sheet1!A1:B2").
//
// Returns:
//
//	[][]interface{}: The retrieved values, where each row is represented by a slice of interface{}.
//	error: An error if the values cannot be retrieved.
func (s *spreadSheet) GetRangeValues(spreadId, sheetRange string) ([][]interface{}, error) {
	spreadData, err := s.srv.Spreadsheets.Values.Get(
		spreadId,
		sheetRange).Do()
	if err != nil {
		return [][]interface{}{}, err
	}

	return spreadData.Values, nil
}

func (s *spreadSheet) CreateSpreadSheet(title string) (string, error) {
	spreadsheet, err := s.srv.Spreadsheets.Create(&sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: title,
		},
	}).Do()
	if err != nil {
		return "", err
	}
	return spreadsheet.SpreadsheetId, nil
}

func (s *spreadSheet) BatchWrite(spreadsheetId string, records [][]interface{}) error {
	writeRange := fmt.Sprintf("Sheet1!A1:%d%d", len(records[0]), len(records))
	values := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         records,
	}
	_, err := s.srv.Spreadsheets.Values.Update(
		spreadsheetId, writeRange, values).ValueInputOption(
		"USER_ENTERED").Do()
	if err != nil {
		return err
	}
	return nil
}

func (s *spreadSheet) columnName(n int) string {
	name := ""
	for n > 0 {
		n--
		name = string(rune('A'+n%26)) + name
		n /= 26
	}
	return name
}
