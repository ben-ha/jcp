package copier

import (
	"fmt"
	"io"
	"os"

	discovery "github.com/ben-ha/jcp/discovery"
)

type BlockCopier struct {
	BlockSize uint64
}

type BlockCopierState struct {
	Size             uint64
	BytesTransferred uint64
}

func MakeBlockCopier(blockSizeBytes uint64) Copier {
	return BlockCopier{BlockSize: blockSizeBytes}
}

func (copier BlockCopier) Copy(source discovery.FileInformation, destination discovery.FileInformation, state CopierState) CopierState {
	return copier.CopyWithProgress(source, destination, state, nil)
}

func (copier BlockCopier) CopyWithProgress(source discovery.FileInformation, destination discovery.FileInformation, state CopierState, progress chan<- CopierProgress) CopierState {
	concreteState, castOK := state.State.(BlockCopierState)
	if !castOK && state.State != nil {
		castErr := fmt.Errorf("converting to BlockCopierState failed")
		reportProgress(progress, source, destination, concreteState, castErr)
		return CopierState{State: state.State, Error: castErr}
	}

	if state.State == nil {
		concreteState = BlockCopierState{}
	}

	inputFile, inputErr := os.Open(source.FullPath)
	if inputErr != nil {
		reportProgress(progress, source, destination, concreteState, inputErr)
		return CopierState{State: state.State, Error: inputErr}
	}

	defer inputFile.Close()

	concreteState.Size = uint64(source.Info.Size())

	outputFile, outputErr := os.OpenFile(destination.FullPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if outputErr != nil {
		reportProgress(progress, source, destination, concreteState, outputErr)
		return CopierState{State: state.State, Error: outputErr}
	}

	defer outputFile.Close()

	_, seekErr := inputFile.Seek(int64(concreteState.BytesTransferred), 0)
	if seekErr != nil {
		reportProgress(progress, source, destination, concreteState, seekErr)
		return CopierState{State: state.State, Error: seekErr}
	}

	_, seekErr = outputFile.Seek(int64(concreteState.BytesTransferred), 0)
	if seekErr != nil {
		reportProgress(progress, source, destination, concreteState, seekErr)
		return CopierState{State: state.State, Error: seekErr}
	}

	reportProgress(progress, source, destination, concreteState, nil)
	var readErr *error = nil
	blockBuffer := make([]byte, copier.BlockSize)
	for readErr == nil {
		read, err := inputFile.Read(blockBuffer)
		if err != nil {
			if err == io.EOF {
				reportProgress(progress, source, destination, concreteState, io.EOF)
				break
			}
			return CopierState{State: concreteState, Error: err}
		}
		written, writeErr := outputFile.Write(blockBuffer[0:read])
		if writeErr != nil {
			reportProgress(progress, source, destination, concreteState, writeErr)
			return CopierState{State: concreteState, Error: writeErr}
		}

		if read != written {
			blockDifferent := fmt.Errorf("write size is different: read=%v, write=%v", read, written)
			reportProgress(progress, source, destination, concreteState, blockDifferent)
			return CopierState{State: concreteState, Error: blockDifferent}
		}
		concreteState.BytesTransferred += uint64(read)
		reportProgress(progress, source, destination, concreteState, nil)
	}

	reportProgress(progress, source, destination, concreteState, io.EOF)
	return CopierState{State: concreteState, Error: io.EOF}
}

func (state BlockCopierState) IsDone() bool {
	return state.Size == state.BytesTransferred
}

func reportProgress(progressChan chan<- CopierProgress, source discovery.FileInformation, dest discovery.FileInformation, currentState BlockCopierState, err error) {
	if progressChan == nil {
		return
	}
	progressChan <- CopierProgress{Source: source.FullPath, Dest: dest.FullPath, Size: currentState.Size, BytesTransferred: currentState.BytesTransferred, Error: err}
}
