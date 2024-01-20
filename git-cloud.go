package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path"

	commands "github.com/bamiesking/git-cloud/commands"
	parser "github.com/bamiesking/git-cloud/parser"
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

	switch os.Args[1] {
	case "fetch":
		fetch.Parse(os.Args[2:])
	case "pull":
		pull.Parse(os.Args[2:])
	case "diff":
		diff.Parse(os.Args[2:])
	}

	// Read in the entries in .gitcloud
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cF, err := parser.ParseLine(scanner.Text())
		if err != nil {
			log.Print(err)
			continue
		}
		commands.Fetch(cF)
	}

	err = scanner.Err()
	utils.Handle(err)
}
