package state

import (
	"time"

	"github.com/ben-ha/jcp/copier"
)

type CopySourceKey = string
type CopyDestinationKey = string

type JcpState struct {
	CopyStates map[CopySourceKey](map[CopyDestinationKey]JcpCopyState)
	StatePath  string
}

type CopierType int

const (
	BlockCopier CopierType = 1
)

type JcpCopyState struct {
	OpaqueState any
	CopierType  CopierType
	LastUpdate  time.Time
	Percent     float64
}

type JcpProgress struct {
	JcpError error
	Progress copier.CopierProgress
}

func (jcpState JcpState) GetStateForTransfer(src string, dest string) copier.CopierState {
	if jcpState.CopyStates[src] == nil {
		return copier.CopierState{}
	}

	copierState, err := anyToType[copier.CopierState](jcpState.CopyStates[src][dest].OpaqueState)
	if err != nil {
		panic("Corrupted state. Please delete state.json cache and try again")
	}

	if copierState.State != nil {
		switch jcpState.CopyStates[src][dest].CopierType {
		case BlockCopier:
			copierState.State, err = anyToType[copier.BlockCopierState](copierState.State)
			if err != nil {
				panic("Corrupted state. Please delete state.json cache and try again")
			}
		default:
			panic("Corrupted state. Please delete state.json cache and try again")
		}
	}

	return copierState
}
