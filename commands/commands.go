package commands

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/bamiesking/git-cloud/services/gdrive"
	"github.com/bamiesking/git-cloud/structs"
	"github.com/bamiesking/git-cloud/utils"
)

func Fetch(cF structs.CloudFile) structs.CloudFileInfo {
	switch cF.Service {
	case structs.GDrive:
		return gdrive.FetchGDrive(cF.Handle)
	}
	return structs.CloudFileInfo{}
}

func Pull(cF structs.CloudFile) {
	var chunkSize int64 = 4096
	cacheBuffer := make([]byte, chunkSize)
	treeBuffer := make([]byte, chunkSize)
	Fetch(cF)
	gitPath, err := utils.GitRepoPath()
	if err != nil {
		log.Fatal(err)
	}
	treePath := path.Join(gitPath, cF.Path)
	cachePath := path.Join(gitPath, ".git/cloud/cache", cF.Handle)
	err = os.MkdirAll(path.Dir(treePath), os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}
	treeFile, err := os.OpenFile(treePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	cacheFile, err := os.OpenFile(cachePath, os.O_RDONLY, os.ModeExclusive)
	if err != nil {
		log.Fatal(err)
	}
	treeReadWriteSeeker := io.ReadWriteSeeker(treeFile)
	cacheReadSeeker := io.ReadSeeker(cacheFile)
	pos, err := compareFiles(cacheReadSeeker, cacheBuffer, treeReadWriteSeeker, treeBuffer)
	fmt.Printf("Pos: %d\n", pos)
	if err != nil {
		log.Fatal(err)
	}
	if pos >= 0 {
		pos, err = cacheReadSeeker.Seek(-1*chunkSize, io.SeekCurrent)
		if err != nil {
			log.Fatal(err)
		}
		_, err := treeReadWriteSeeker.Seek(pos, io.SeekStart)
		if err != nil {
			log.Fatal(err)
		}
		size, err := io.Copy(treeReadWriteSeeker, cacheReadSeeker)
		if err != nil {
			log.Fatal(err)
		}
		fileStat, err := cacheFile.Stat()
		if err != nil {
			log.Fatal(err)
		}
		if pos+size != fileStat.Size() {
			log.Fatal("Tree file and cache file are not the same size.")
		}
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
