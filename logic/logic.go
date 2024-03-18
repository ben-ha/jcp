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
	"github.com/ben-ha/jcp/state"
)

const BLOCKSIZE = 1024 * 1024 // 1MB

type Jcp struct {
	ProgressChannel  chan state.JcpProgress
	ConcurrencyLimit uint
	JcpState         state.JcpState
}

func MakeJcp(concurrencyLimit uint, jcpState state.JcpState) Jcp {
	progressChannel := make(chan state.JcpProgress)
	return Jcp{ProgressChannel: progressChannel, ConcurrencyLimit: concurrencyLimit, JcpState: jcpState}
}

func (jcp Jcp) StartCopy(src string, dest string) error {
	copierProgressChannel := make(chan copier.CopierProgress)
	jcp.startProcessingProgress(copierProgressChannel)

	isDir, isDirErr := IsDirectory(src)
	if isDirErr != nil {
		return isDirErr
	}

	isDirDest, _ := IsDirectory(dest)

	if isDir {
		err := jcp.startDirectoryCopy(src, dest, copierProgressChannel)
		if err != nil {
			return err
		}
	} else {
		if isDirDest {
			dest = path.Join(dest, path.Base(src))
		}
		err := jcp.startFileCopy(src, dest, copierProgressChannel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (jcp Jcp) startDirectoryCopy(src string, dest string, progressChannel chan copier.CopierProgress) error {
	concurrencyLimiter := make(chan int, jcp.ConcurrencyLimit)
	discoverer, err := discovery.MakeBfsDiscoverer(src)
	if err != nil {
		return err
	}

	srcBasePath := path.Dir(src)
	destBasePath := dest
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
			go func(currentFile discovery.FileInformation, destFile discovery.FileInformation) {
				defer transferCompletion.Done()
				jcp.startFileCopyByInfo(currentFile, destFile, progressChannel)
				<-concurrencyLimiter
			}(currentFile, destFile)
		}

		transferCompletion.Wait()
		jcp.reportError(io.EOF)
		close(progressChannel)
	}()

	return nil
}

func (jcp Jcp) startFileCopy(source string, destination string, progress chan<- copier.CopierProgress) error {
	srcInfo, err := discovery.MakeFileInformation(source)
	if err != nil {
		return err
	}

	dstInfo, _ := discovery.MakeFileInformation(destination)
	go func() {
		newState := jcp.startFileCopyByInfo(srcInfo, dstInfo, progress)
		err := newState.Error
		if err == nil {
			err = io.EOF
		}
		jcp.reportError(err)
	}()

	return nil
}

func (jcp Jcp) startFileCopyByInfo(source discovery.FileInformation, destination discovery.FileInformation, progress chan<- copier.CopierProgress) copier.CopierState {
	cp := copier.MakeBlockCopier(BLOCKSIZE)
	copierState := jcp.JcpState.GetStateForTransfer(source.FullPath, destination.FullPath)
	newState := cp.CopyWithProgress(source, destination, copierState, progress)
	return newState
}

func (jcp Jcp) reportError(err error) {
	jcp.ProgressChannel <- state.JcpProgress{JcpError: err, Progress: copier.CopierProgress{}}
}

func RemoveBaseDirectory(base string, input string) string {
	noBase := strings.Replace(input, base, "", 1)
	return noBase
}

func (jcp Jcp) startProcessingProgress(copierProgress chan copier.CopierProgress) {
	go func() {
		for {
			msg, more := <-copierProgress
			if !more {
				break
			}
			jcp.ProgressChannel <- state.JcpProgress{JcpError: nil, Progress: msg}
		}
	}()
}
