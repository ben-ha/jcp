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
	inputFile, inputErr := os.Open(source)
	if inputErr != nil {
		return CopierState{State: state.State, Error: &inputErr}
	}
	outputFile, outputErr := os.OpenFile(destination, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if outputErr != nil {
		return CopierState{State: state.State, Error: &outputErr}
	}

	concreteState, castOK := state.State.(BlockCopierState)
	if !castOK {
		castErr := fmt.Errorf("converting to BlockCopierState failed")
		return CopierState{State: state.State, Error: &castErr}
	}

	_, seekErr := inputFile.Seek(int64(concreteState.BytesTransferred), 0)
	if seekErr != nil {
		return CopierState{State: state.State, Error: &seekErr}
	}

	_, seekErr = outputFile.Seek(int64(concreteState.BytesTransferred), 0)
	if seekErr != nil {
		return CopierState{State: state.State, Error: &seekErr}
	}

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
	}

	return CopierState{State: concreteState, Error: nil}
}
