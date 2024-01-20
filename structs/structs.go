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
	Handle  string
	Service Service
}

type CloudFileInfo struct {
	DateModified time.Time
	Name         string
	Size         int64
}

func ParseService(identifier string) (Service, error) {
	switch identifier {
	case "gdrive":
		return GDrive, nil
	}
	return Undefined, errors.New("unrecognised cloud service identifier")
}

func (s Service) String() string {
	switch s {
	case GDrive:
		return "gdrive"
	}
	return "undefined"
}
