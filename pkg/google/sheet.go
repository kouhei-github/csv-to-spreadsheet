package google

import (
	"context"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
)

type sheetService struct {
	jsonPath string
}

type Client interface {
	Service() *http.Client
}

// NewSheetService creates a new instance of SheetService with the provided JSON key file path.
// It initializes the necessary service for accessing the Google Sheets API.
//
// jsonPath (string): The file path to the JSON key file.
//
// Returns:
//
//	SheetService: A SheetService object.
func NewSheetService(jsonPath string) Client {
	return &sheetService{
		jsonPath: jsonPath,
	}
}

// Service initializes and returns a *google.Service object for accessing the Google Sheets API.
//
// It reads the client secret file from `s.jsonPath` and uses it to create a JWT config.
// The config is used to create a client with the necessary scopes. Finally, the client is used
// to create the *google.Service object.
//
// Returns:
//
//	*google.Service: A *google.Service object for accessing the Google Sheets API.
func (s *sheetService) Service() *http.Client {
	ctx := context.Background()
	b, err := os.ReadFile(s.jsonPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.JWTConfigFromJSON(
		b,
		"https://www.googleapis.com/auth/spreadsheets",
		"https://www.googleapis.com/auth/drive")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := config.Client(ctx)

	return client
}
