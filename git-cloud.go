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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cF := parser.ParseLine(scanner.Text())
		fmt.Println(cF)
		file := cF.FetchFile()
		fmt.Println(file)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
