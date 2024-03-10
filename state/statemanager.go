package state

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/ben-ha/jcp/logic"
)

var stateManagerMutex sync.Mutex

const ValidStateWindowInDays = 1

func LoadState(fileName string) *JcpState {
	data, err := os.ReadFile(fileName)
	loadedState := &JcpState{}
	if err == nil {
		json.Unmarshal(data, loadedState)
	}

	loadedState.Clean()

	return loadedState
}

func (copierState *JcpState) Save(fileName string) {
	serialized, err := json.Marshal(copierState)

	if err == nil {
		os.WriteFile(fileName, serialized, os.ModePerm)
	}
}

func (copierState *JcpState) Update(progress logic.JcpProgress) {
	stateManagerMutex.Lock()
	defer stateManagerMutex.Unlock()

	newState := MakeNewCopyState(progress)
	if copierState.CopyStates[progress.Progress.Source] == nil {
		copierState.CopyStates[progress.Progress.Dest] = map[CopyDestinationKey]JcpCopyState{}
	}

	copierState.CopyStates[progress.Progress.Source][progress.Progress.Source] = newState
}

func (copierState *JcpState) Clean() {
	stateManagerMutex.Lock()
	defer stateManagerMutex.Unlock()

	newStates := make(map[CopySourceKey]map[CopyDestinationKey]JcpCopyState)

	for src := range copierState.CopyStates {
		for dest := range copierState.CopyStates[src] {
			if copierState.CopyStates[src][dest].LastUpdate.After(time.Now().AddDate(0, 0, -ValidStateWindowInDays)) {
				newStates[src][dest] = copierState.CopyStates[src][dest]
			}
		}
	}

	copierState.CopyStates = newStates
}
