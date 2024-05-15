package main

import (
	"fmt"
	"github.com/kouhei-github/csv-to-spreadssheet/di"
	"github.com/kouhei-github/csv-to-spreadssheet/pkg/google"
)

func main() {
	sheetClient := google.NewSheetService("./tmp/service-account.json")
	handler := di.NewInjection(sheetClient)

	if err := handler.Run(); err != nil {
		fmt.Println(err)
	}
}
