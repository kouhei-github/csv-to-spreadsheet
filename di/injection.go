package di

import (
	"github.com/kouhei-github/csv-to-spreadssheet/interval/usecase"
	"github.com/kouhei-github/csv-to-spreadssheet/pkg/excels"
	"github.com/kouhei-github/csv-to-spreadssheet/pkg/google"
	"github.com/kouhei-github/csv-to-spreadssheet/plugins/google_drive"
	"github.com/kouhei-github/csv-to-spreadssheet/plugins/spreadsheet"
)

// NewInjection returns a new instance of usecase.SpreadUseCase using the provided google.SheetService.
// It creates a sheet using the sheetClient, then creates a spreadsheet using the sheet,
// and finally creates a useCase using the spread. The useCase implements the SpreadUseCase interface.
//
// sheetClient (google.SheetService): The SheetService implementation to be used for creating the spreadsheet and useCase.
//
// Returns:
//
//	usecase.SpreadUseCase: A SpreadUseCase object.
func NewInjection(googleClient google.Client) usecase.SpreadUseCase {
	client := googleClient.Service()
	spread := spreadsheet.NewSpreadsheet(client)
	drive := google_drive.NewGoogleDrive(client)
	excel := excels.NewExcel("./medridge-jobs.csv")
	useCase := usecase.NewSpreadUseCase(spread, drive, excel)

	return useCase
}
