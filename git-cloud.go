package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	commands "github.com/bamiesking/git-cloud/commands"
	parser "github.com/bamiesking/git-cloud/parser"
	"github.com/bamiesking/git-cloud/structs"
	utils "github.com/bamiesking/git-cloud/utils"
)

func main() {
	// Verify that we are in a git repo
	gitPath, err := utils.GitRepoPath()
	utils.Handle(err)

	// Open/create .gitcloud file
	file, err := os.Open(path.Join(gitPath, ".gitcloud"))
	utils.Handle(err)
	defer file.Close()

	// Make caching directory
	err = os.MkdirAll(path.Join(gitPath, ".git/cloud/cache"), os.ModePerm)
	utils.Handle(err)

	fetch := flag.NewFlagSet("fetch", flag.ExitOnError)
	pull := flag.NewFlagSet("pull", flag.ExitOnError)
	diff := flag.NewFlagSet("diff", flag.ExitOnError)

	if len(os.Args) < 2 {
		os.Exit(1)
	}

	var command func(structs.CloudFile)
	switch os.Args[1] {
	case "fetch":
		fetch.Parse(os.Args[2:])
		command = commands.Fetch
	case "pull":
		pull.Parse(os.Args[2:])
		command = commands.Pull
	case "diff":
		diff.Parse(os.Args[2:])
		command = commands.Diff
	default:
		fmt.Printf("git-cloud: '%s' is not a valid subcommand. See 'git cloud --help'\n", os.Args[1])
		os.Exit(0)
	}

	// Read in the entries in .gitcloud
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cF, err := parser.ParseLine(scanner.Text())
		if err != nil {
			log.Print(err)
			continue
		}
		command(cF)
	}

	err = scanner.Err()
	utils.Handle(err)
}
