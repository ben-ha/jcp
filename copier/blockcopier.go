package copier

import (
	"fmt"
	"io"
	"os"
)

type BlockCopier struct {
	BlockSize uint64
}

type BlockCopierState struct {
	Size             uint64
	BytesTransferred uint64
}

func (copier BlockCopier) Copy(source string, destination string, state CopierState) CopierState {
	return copier.CopyWithProgress(source, destination, state, nil)
}

func (copier BlockCopier) CopyWithProgress(source string, destination string, state CopierState, progress chan<- CopierProgress) CopierState {
	if progress != nil {
		defer close(progress)
	}

	concreteState, castOK := state.State.(BlockCopierState)
	if !castOK {
		castErr := fmt.Errorf("converting to BlockCopierState failed")
		return CopierState{State: state.State, Error: &castErr}
	}

	inputFile, inputErr := os.Open(source)
	if inputErr != nil {
		return CopierState{State: state.State, Error: &inputErr}
	}

	defer inputFile.Close()

	if concreteState.Size == 0 {
		stat, statErr := inputFile.Stat()
		if statErr != nil {
			return CopierState{State: state.State, Error: &statErr}
		}

		concreteState.Size = uint64(stat.Size())
	}

	outputFile, outputErr := os.OpenFile(destination, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if outputErr != nil {
		return CopierState{State: state.State, Error: &outputErr}
	}

	defer outputFile.Close()

	_, seekErr := inputFile.Seek(int64(concreteState.BytesTransferred), 0)
	if seekErr != nil {
		return CopierState{State: state.State, Error: &seekErr}
	}

	_, seekErr = outputFile.Seek(int64(concreteState.BytesTransferred), 0)
	if seekErr != nil {
		return CopierState{State: state.State, Error: &seekErr}
	}

	reportProgress(progress, concreteState)
	var readErr *error = nil
	blockBuffer := make([]byte, copier.BlockSize)
	for readErr == nil {
		read, err := inputFile.Read(blockBuffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return CopierState{State: concreteState, Error: &err}
		}
		written, writeErr := outputFile.Write(blockBuffer[0:read])
		if writeErr != nil {
			return CopierState{State: concreteState, Error: &writeErr}
		}

		if read != written {
			blockDifferent := fmt.Errorf("write size is different: read=%v, write=%v", read, written)
			return CopierState{State: concreteState, Error: &blockDifferent}
		}
		concreteState.BytesTransferred += uint64(read)
		reportProgress(progress, concreteState)
	}

	return CopierState{State: concreteState, Error: nil}
}

func (state BlockCopierState) IsDone() bool {
	return state.Size == state.BytesTransferred
}

func reportProgress(progressChan chan<- CopierProgress, currentState BlockCopierState) {
	if progressChan == nil {
		return
	}
	progressChan <- CopierProgress{Size: currentState.Size, BytesTransferred: currentState.BytesTransferred}
}
