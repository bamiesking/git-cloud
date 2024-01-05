package service

import (
	"context"
	"log"
	"os"
	"time"

	gdrive "github.com/bamiesking/git-cloud/service/gdrive"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type Service int

const (
	Undefined Service = iota
	GDrive
)

type CloudFile struct {
	Path    string
	Service Service
	Handle  string
}

func (cF CloudFile) FetchFile() CloudFileInfo {
	switch cF.Service {
	case GDrive:
		return fetchGDriveFile(cF.Handle)
	}
	return CloudFileInfo{}
}

type CloudFileInfo struct {
	name         string
	size         int64
	dateModified time.Time
}

func ParseService(identifier string) Service {
	switch identifier {
	case "gdrive":
		return GDrive
	}
	log.Printf("Unrecognised Service identifier: %s", identifier)
	return Undefined
}

func (s Service) String() string {
	switch s {
	case GDrive:
		return "gdrive"
	}
	return "undefined"
}

// Google Drive
func fetchGDriveFile(handle string) CloudFileInfo {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := gdrive.GetClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	r, err := srv.Files.Get(handle).Fields(googleapi.Field("name,size,modifiedTime")).Do()
	info := CloudFileInfo{}
	if r != nil {
		info.name = r.Name
		info.size = r.Size
		t, err := time.Parse(time.RFC3339, r.ModifiedTime)
		if err != nil {
			log.Printf("Failed to parse ModifiedTime for file %s (%s)", r.Name, handle)
			return info
		}
		info.dateModified = t
	}
	return info
}
