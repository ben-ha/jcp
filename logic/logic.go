package logic

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/ben-ha/jcp/copier"
	"github.com/ben-ha/jcp/discovery"
)

const BLOCKSIZE = 1024 * 1024 // 1MB

func StartCopy(src string, dest string, concurrencyLimit uint) (<-chan copier.CopierProgress, error) {
	progressChannel := make(chan copier.CopierProgress)
	isDir, isDirErr := IsDirectory(src)
	if isDirErr != nil {
		return nil, isDirErr
	}

	if isDir {
		err := startDirectoryCopy(src, dest, concurrencyLimit, progressChannel)
		if err != nil {
			return nil, err
		}
	} else {
		err := startFileCopy(src, dest, copier.CopierState{}, progressChannel)
		if err != nil {
			return nil, err
		}
	}

	return progressChannel, nil
}

func startDirectoryCopy(src string, dest string, concurrencyLimit uint, progressChannel chan copier.CopierProgress) error {
	concurrencyLimiter := make(chan int, concurrencyLimit)
	discoverer, err := discovery.MakeBfsDiscoverer(src)
	if err != nil {
		return err
	}

	srcBasePath := path.Dir(src)
	destBasePath := path.Dir(dest)
	go func() {
		var currentErr error = nil
		var currentFile discovery.FileInformation
		var transferCompletion sync.WaitGroup
		for currentErr == nil {
			currentFile, currentErr = discoverer.Next()
			if currentErr != nil {
				if currentErr != io.EOF {
					panic(fmt.Sprintf("Failed: %v", currentErr))
				}

				break
			}
			concurrencyLimiter <- 1

			if currentFile.Info.IsDir() {
				os.MkdirAll(path.Join(destBasePath, RemoveBaseDirectory(srcBasePath, currentFile.FullPath)), fs.ModePerm)
				<-concurrencyLimiter
				continue
			}

			transferCompletion.Add(1)
			destFile, _ := discovery.MakeFileInformation(path.Join(destBasePath, RemoveBaseDirectory(srcBasePath, currentFile.FullPath)))
			go func() {
				defer transferCompletion.Done()
				startFileCopyByInfo(currentFile, destFile, copier.CopierState{}, progressChannel)
				<-concurrencyLimiter
			}()
		}

		transferCompletion.Wait()
		close(progressChannel)
	}()

	return nil
}

func startFileCopy(source string, destination string, state copier.CopierState, progress chan<- copier.CopierProgress) error {
	srcInfo, err := discovery.MakeFileInformation(source)
	if err != nil {
		return err
	}

	dstInfo, _ := discovery.MakeFileInformation(destination)

	go startFileCopyByInfo(srcInfo, dstInfo, state, progress)
	return nil
}

func startFileCopyByInfo(source discovery.FileInformation, destination discovery.FileInformation, state copier.CopierState, progress chan<- copier.CopierProgress) {
	cp := copier.MakeBlockCopier(BLOCKSIZE)
	cp.CopyWithProgress(source, destination, state, progress)
}

func RemoveBaseDirectory(base string, input string) string {
	noBase := strings.Replace(input, base, "", 1)
	return noBase
}
