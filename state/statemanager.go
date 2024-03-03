package state

import (
	"encoding/json"
	"os"
	"sync"
)

var stateManagerMutex sync.RWMutex

func LoadState(fileName string) *CopierState {
	data, err := os.ReadFile(fileName)
	loadedState := &CopierState{}
	if err == nil {
		json.Unmarshal(data, loadedState)
	}

	return loadedState
}

func (copierState *CopierState) Save(fileName string) {
	serialized, err := json.Marshal(copierState)

	if err == nil {
		os.WriteFile(fileName, serialized, os.ModePerm)
	}
}

func (copierState *CopierState) StartCopy(source string, destination string) *CopyState {
	stateManagerMutex.Lock()
	defer stateManagerMutex.Unlock()

	newState := CopyState{}
	if copierState.CopyStates[source] == nil {
		copierState.CopyStates[source] = map[CopyDestinationKey]CopyState{}
	}

	copierState.CopyStates[source][destination] = newState
	return &newState
}

func (copierState *CopierState) EndCopy(copyState *CopyState) {
	stateManagerMutex.Lock()
	defer stateManagerMutex.Unlock()

}
