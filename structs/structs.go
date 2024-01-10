package structs

import (
	"errors"
	"time"
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
type CloudFileInfo struct {
	Name         string
	Size         int64
	DateModified time.Time
}

func ParseService(identifier string) (Service, error) {
	switch identifier {
	case "gdrive":
		return GDrive, nil
	}
	return Undefined, errors.New("Unrecognised cloud service identifier")
}

func (s Service) String() string {
	switch s {
	case GDrive:
		return "gdrive"
	}
	return "undefined"
}
