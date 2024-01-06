package parser

import (
	"errors"
	serv "github.com/bamiesking/git-cloud/service"
	"strings"
)

func ParseLine(scannerText string) (serv.CloudFile, error) {
	args := strings.Fields(scannerText)
	if len(args) != 2 {
		return serv.CloudFile{}, errors.New("Unable to parse line")
	}
	handleArgs := strings.Split(args[1], ":")
	if len(handleArgs) != 2 {
		return serv.CloudFile{}, errors.New("Unable to parse cloud handle")
	}
	service, err := serv.ParseService(handleArgs[0])
	if err != nil {
		return serv.CloudFile{}, err

	}
	cf := serv.CloudFile{Path: args[0], Service: service, Handle: handleArgs[1]}
	return cf, nil
}
