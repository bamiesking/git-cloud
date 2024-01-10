package commands

import (
	"github.com/bamiesking/git-cloud/services/gdrive"
	"github.com/bamiesking/git-cloud/structs"
)

func Fetch(cF structs.CloudFile) structs.CloudFileInfo {
	switch cF.Service {
	case structs.GDrive:
		return gdrive.FetchGDrive(cF.Handle)
	}
	return structs.CloudFileInfo{}
}
