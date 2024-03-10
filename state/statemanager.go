package state

import (
	"encoding/json"
	"os"
	"path"
	"sync"
	"time"

	"github.com/ben-ha/jcp/logic"
)

var stateManagerMutex sync.Mutex

const ValidStateWindowInDays = 1

const JcpStateDirectoryName = "jcp"
const JcpStateFileName = "state.json"

func InitializeState() (*JcpState, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	jcpDir := path.Join(cacheDir, JcpStateDirectoryName)

	err = os.Mkdir(jcpDir, os.ModePerm)
	if !os.IsExist(err) {
		return nil, err
	}

	return loadState(path.Join(jcpDir, JcpStateFileName)), nil
}

func loadState(fileName string) *JcpState {
	data, err := os.ReadFile(fileName)
	loadedState := &JcpState{CopyStates: make(map[CopySourceKey]map[CopyDestinationKey]JcpCopyState)}
	if err == nil {
		json.Unmarshal(data, loadedState)
	}

	loadedState.Clean()
	loadedState.StatePath = fileName

	return loadedState
}

func (copierState *JcpState) SaveState() {
	copierState.Clean()
	copierState.saveState(copierState.StatePath)
}

func (copierState *JcpState) saveState(fileName string) {
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
		copierState.CopyStates[progress.Progress.Source] = map[CopyDestinationKey]JcpCopyState{}
	}

	copierState.CopyStates[progress.Progress.Source][progress.Progress.Dest] = newState
}

func (copierState *JcpState) Clean() {
	stateManagerMutex.Lock()
	defer stateManagerMutex.Unlock()

	newStates := make(map[CopySourceKey]map[CopyDestinationKey]JcpCopyState)

	for src := range copierState.CopyStates {
		if newStates[src] == nil {
			newStates[src] = make(map[CopyDestinationKey]JcpCopyState)
		}

		for dest := range copierState.CopyStates[src] {
			if copierState.CopyStates[src][dest].ShouldKeep() {
				newStates[src][dest] = copierState.CopyStates[src][dest]
			}
		}
	}

	copierState.CopyStates = newStates
}

func (copyState JcpCopyState) ShouldKeep() bool {
	if copyState.LastUpdate.Before(time.Now().AddDate(0, 0, -ValidStateWindowInDays)) {
		return false
	}

	if copyState.Percent == 1 {
		return false
	}

	return true
}
