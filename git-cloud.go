package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"

	parser "github.com/bamiesking/git-cloud/parser"
	utils "github.com/bamiesking/git-cloud/utils"
)

func main() {
	gitPath, err := utils.GitRepoPath()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	file, err := os.Open(path.Join(gitPath, ".gitcloud"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer file.Close()

	err = os.MkdirAll(path.Join(gitPath, ".git/cloud/cache"), os.ModePerm)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cF, err := parser.ParseLine(scanner.Text())
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Println(cF)
		file := cF.FetchFile()
		fmt.Println(file)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
