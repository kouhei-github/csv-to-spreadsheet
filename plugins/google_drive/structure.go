package google_drive

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"log"
	"net/http"
)

type GoogleDrive interface {
	AttachFullAccess(spreadId string) error
}

type googleDrive struct {
	srv *drive.Service
}

func NewGoogleDrive(client *http.Client) GoogleDrive {
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Failed to create config %v", err)
	}
	return &googleDrive{
		srv: srv,
	}
}

func (d *googleDrive) AttachFullAccess(spreadId string) error {
	permission := &drive.Permission{
		Type: "anyone",
		Role: "writer",
	}

	_, err := d.srv.Permissions.Create(spreadId, permission).Do()
	if err != nil {
		fmt.Println("here")
		return err
	}
	return nil
}
