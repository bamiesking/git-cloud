package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/bamiesking/git-cloud/services/gdrive"
	"github.com/bamiesking/git-cloud/structs"
	"github.com/bamiesking/git-cloud/utils"
)

func Fetch(cF structs.CloudFile) {
	switch cF.Service {
	case structs.GDrive:
		gdrive.FetchGDrive(cF.Handle)
	}
}

func Pull(cF structs.CloudFile) {
	var chunkSize int64 = 4096
	cacheBuffer := make([]byte, chunkSize)
	treeBuffer := make([]byte, chunkSize)
	Fetch(cF)
	gitPath, err := utils.GitRepoPath()
	utils.Handle(err)
	treePath := path.Join(gitPath, cF.Path)
	cachePath := path.Join(gitPath, ".git/cloud/cache", cF.Handle)
	err = os.MkdirAll(path.Dir(treePath), os.ModeDir)
	utils.Handle(err)
	treeFile, err := os.OpenFile(treePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	utils.Handle(err)
	cacheFile, err := os.OpenFile(cachePath, os.O_RDONLY, os.ModeExclusive)
	utils.Handle(err)
	treeReadWriteSeeker := io.ReadWriteSeeker(treeFile)
	cacheReadSeeker := io.ReadSeeker(cacheFile)
	pos, err := compareFiles(cacheReadSeeker, cacheBuffer, treeReadWriteSeeker, treeBuffer)
	utils.Handle(err)
	if pos >= 0 {
		pos, err = cacheReadSeeker.Seek(-1*chunkSize, io.SeekCurrent)
		utils.Handle(err)
		_, err := treeReadWriteSeeker.Seek(pos, io.SeekStart)
		utils.Handle(err)
		size, err := io.Copy(treeReadWriteSeeker, cacheReadSeeker)
		utils.Handle(err)
		fileStat, err := cacheFile.Stat()
		utils.Handle(err)
		if pos+size != fileStat.Size() {
			log.Fatal("Tree file and cache file are not the same size.")
		}
	}
}

func Diff(cF structs.CloudFile) {
	gitPath, err := utils.GitRepoPath()
	utils.Handle(err)
	treePath := path.Join(gitPath, cF.Path)
	cachePath := path.Join(gitPath, ".git/cloud/cache", cF.Handle)
	treeStat, err := os.Stat(treePath)
	if errors.Is(err, os.ErrNotExist) {
		// Deleted file
		fmt.Printf("- %s\n", cF.Path)
		return
	} else {
		utils.Handle(err)
	}
	cacheStat, err := os.Stat(cachePath)
	if errors.Is(err, os.ErrNotExist) {
		// Added file
		fmt.Printf("+ %s\n", cF.Path)
		return
	} else {
		utils.Handle(err)
	}
	if treeStat.Size() != cacheStat.Size() {
		// Modified file
		fmt.Printf("M %s\n", cF.Path)
		return
	}
	var chunkSize int64 = 4096
	cacheBuffer := make([]byte, chunkSize)
	treeBuffer := make([]byte, chunkSize)
	Fetch(cF)
	err = os.MkdirAll(path.Dir(treePath), os.ModeDir)
	utils.Handle(err)
	treeFile, err := os.OpenFile(treePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	utils.Handle(err)
	cacheFile, err := os.OpenFile(cachePath, os.O_RDONLY, os.ModeExclusive)
	utils.Handle(err)
	treeReader := io.Reader(treeFile)
	cacheReader := io.ReadSeeker(cacheFile)
	diff, err := compareFiles(cacheReader, cacheBuffer, treeReader, treeBuffer)
	utils.Handle(err)
	if diff >= 0 {
		// Modified file
		fmt.Printf("M %s\n", cF.Path)
		return
	}
}

func compareFiles(r1 io.ReadSeeker, b1 []byte, r2 io.Reader, b2 []byte) (int64, error) {
	for {
		_, err1 := r1.Read(b1)
		_, err2 := r2.Read(b2)

		if err1 == nil && err2 == nil {
			if !bytes.Equal(b1, b2) {
				return r1.Seek(0, io.SeekCurrent)
			} else {
				continue
			}
		} else if err1 == io.EOF && err2 == io.EOF {
			return -1, nil
		} else if err1 == io.EOF || err2 == io.EOF {
			return r1.Seek(0, io.SeekCurrent)
		} else {
			log.Fatal(err1, err2)
		}
	}
}
