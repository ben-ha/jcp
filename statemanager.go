package jcp

import (
	"sync"
)

var stateManagerMutex sync.RWMutex

func LoadState() *CopierState {
	return &CopierState{}
}

func (*CopierState) Save(fileName string) {

}

func (copierState *CopierState) StartCopy(source string, destination string) *CopyState {
	stateManagerMutex.Lock()
	defer stateManagerMutex.Unlock()

	newCopyKey := CopyStateKey{
		Source:      source,
		Destination: destination,
	}

	newState := &CopyState{}
	copierState.CopyStates[newCopyKey] = newState
	return newState
}

func (copierState *CopierState) EndCopy(copyState *CopyState) {
	stateManagerMutex.Lock()
	defer stateManagerMutex.Unlock()

}
