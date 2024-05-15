package usecase

import (
	"fmt"
	"github.com/kouhei-github/csv-to-spreadssheet/pkg/excels"
	"github.com/kouhei-github/csv-to-spreadssheet/plugins/google_drive"
	"github.com/kouhei-github/csv-to-spreadssheet/plugins/spreadsheet"
	"sync"
)

type SpreadUseCase interface {
	Run() error
}

type spreadUseCase struct {
	spread spreadsheet.SpreadSheet
	drive  google_drive.GoogleDrive
	excel  excels.Excel
}

// NewSpreadUseCase creates a new instance of the SpreadUseCase interface.
// It initializes a spreadUseCase object with the provided SpreadSheet and GoogleDrive implementations.
//
// spread (SpreadSheet): The SpreadSheet implementation to be used by the spreadUseCase.
// drive (GoogleDrive): The GoogleDrive implementation to be used by the spreadUseCase.
//
// Returns:
//
//	SpreadUseCase: A spreadUseCase object that implements the Run method.
func NewSpreadUseCase(
	spread spreadsheet.SpreadSheet,
	drive google_drive.GoogleDrive,
	excel excels.Excel,
) SpreadUseCase {
	return &spreadUseCase{
		spread: spread,
		drive:  drive,
		excel:  excel,
	}
}

type stats struct {
	num int
	err error
}

// Run executes the business logic of the spreadUseCase object.
// It gets the values from a specific range in the spreadsheet using the GetRangeValues method of the SpreadSheet interface.
// It prints the first value from the obtained range values.
//
// Returns:
//
//	error: An error if any occurred during the execution of the method, otherwise nil.
func (s *spreadUseCase) Run() error {
	records, columns, err := s.excel.GetDataFromCsv()
	if err != nil {
		return err
	}

	// 並列処理
	lastCols := s.excel.GetLastColumn(records)
	points, max := s.excel.RemoveDuplicates(lastCols)

	data := make([]int, 0, max)
	results := make(chan stats, max)
	var wg sync.WaitGroup

	// Convert each page to an image in a separate goroutine
	for _, point := range points {
		wg.Add(1)
		go func(point string, wg *sync.WaitGroup, results chan<- stats) {
			defer wg.Done()
			var insert [][]interface{}
			insert = append(insert, columns)
			for _, record := range records {
				if record[len(record)-1] == point {
					slice := make([]interface{}, len(record))
					for j, e := range record {
						slice[j] = e
					}
					insert = append(insert, slice)
				}
			}

			title := fmt.Sprintf("案件名 (ID: %s)", point)

			spreadId, err := s.spread.CreateSpreadSheet(title)
			if err != nil {
				fmt.Println(fmt.Sprintf("%s: [%s]", spreadId, err.Error()))
				results <- stats{0, err}
			}

			fmt.Println(spreadId)

			if err = s.drive.AttachFullAccess(spreadId); err != nil {
				fmt.Println(fmt.Sprintf("%s: [%s]", spreadId, err.Error()))

				results <- stats{0, err}
			}
			err = s.spread.BatchWrite(spreadId, insert)
			if err != nil {
				fmt.Println(fmt.Sprintf("%s: [%s]", spreadId, err.Error()))

				results <- stats{0, err}
			}

			results <- stats{len(data), nil}
		}(point, &wg, results)
	}

	// Close the results channel after all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results as they complete
	for stats := range results {
		if stats.err != nil {
			fmt.Println(stats.err)
			continue
		}
		data = append(data, stats.num)
	}

	fmt.Println(data)

	return nil
}
