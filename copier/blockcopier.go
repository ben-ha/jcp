package copier

import (
	"errors"
	"os"
)

type BlockCopier struct
{
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
		castErr := errors.New("Converting to BlockCopierState failed")
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
}