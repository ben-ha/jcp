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

type JcpProgress struct {
	JcpError error
	Progress copier.CopierProgress
}

type Jcp struct {
	ProgressChannel  chan JcpProgress
	ConcurrencyLimit uint
}

func MakeJcp(concurrencyLimit uint) Jcp {
	progressChannel := make(chan JcpProgress)
	return Jcp{ProgressChannel: progressChannel, ConcurrencyLimit: concurrencyLimit}
}

func (jcp Jcp) StartCopy(src string, dest string) error {
	copierProgressChannel := make(chan copier.CopierProgress)
	jcp.startProcessingProgress(copierProgressChannel)

	isDir, isDirErr := IsDirectory(src)
	if isDirErr != nil {
		return isDirErr
	}

	if isDir {
		err := jcp.startDirectoryCopy(src, dest, jcp.ConcurrencyLimit, copierProgressChannel)
		if err != nil {
			return err
		}
	} else {
		err := jcp.startFileCopy(src, dest, copier.CopierState{}, copierProgressChannel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (jcp Jcp) startDirectoryCopy(src string, dest string, concurrencyLimit uint, progressChannel chan copier.CopierProgress) error {
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
				jcp.startFileCopyByInfo(currentFile, destFile, copier.CopierState{}, progressChannel)
				<-concurrencyLimiter
			}()
		}

		transferCompletion.Wait()

		close(progressChannel)
	}()

	return nil
}

func (jcp Jcp) startFileCopy(source string, destination string, state copier.CopierState, progress chan<- copier.CopierProgress) error {
	srcInfo, err := discovery.MakeFileInformation(source)
	if err != nil {
		return err
	}

	dstInfo, _ := discovery.MakeFileInformation(destination)
	go func() {
		newState := jcp.startFileCopyByInfo(srcInfo, dstInfo, state, progress)
		err := newState.Error
		if err == nil {
			err = io.EOF
		}
		jcp.reportError(err)
	}()

	return nil
}

func (jcp Jcp) startFileCopyByInfo(source discovery.FileInformation, destination discovery.FileInformation, state copier.CopierState, progress chan<- copier.CopierProgress) copier.CopierState {
	cp := copier.MakeBlockCopier(BLOCKSIZE)
	newState := cp.CopyWithProgress(source, destination, state, progress)
	return newState
}

func (jcp Jcp) reportError(err error) {
	jcp.ProgressChannel <- JcpProgress{JcpError: err, Progress: copier.CopierProgress{}}
}

func RemoveBaseDirectory(base string, input string) string {
	noBase := strings.Replace(input, base, "", 1)
	return noBase
}

func (jcp Jcp) startProcessingProgress(copierProgress chan copier.CopierProgress) {
	go func() {
		for {
			msg := <-copierProgress
			jcp.ProgressChannel <- JcpProgress{JcpError: nil, Progress: msg}
		}
	}()
}
