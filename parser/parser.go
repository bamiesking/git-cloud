package parser

import (
	"errors"
	"strings"

	"github.com/bamiesking/git-cloud/structs"
)

func ParseLine(scannerText string) (structs.CloudFile, error) {
	args := strings.Fields(scannerText)
	if len(args) != 2 {
		return structs.CloudFile{}, errors.New("Unable to parse line")
	}
	handleArgs := strings.Split(args[1], ":")
	if len(handleArgs) != 2 {
		return structs.CloudFile{}, errors.New("Unable to parse cloud handle")
	}
	service, err := structs.ParseService(handleArgs[0])
	if err != nil {
		return structs.CloudFile{}, err

	}
	cf := structs.CloudFile{Path: args[0], Service: service, Handle: handleArgs[1]}
	return cf, nil
}
