package parser

import (
	serv "github.com/bamiesking/git-cloud/service"
	"strings"
)

func ParseLine(scannerText string) serv.CloudFile {
	// TODO: make this more robust
	args := strings.Fields(scannerText)
	handleArgs := strings.Split(args[1], ":")
	service := serv.ParseService(handleArgs[0])
	cf := serv.CloudFile{Path: args[0], Service: service, Handle: handleArgs[1]}
	return cf
}
