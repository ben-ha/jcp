package state

import (
	"time"

	"github.com/ben-ha/jcp/logic"
)

type CopySourceKey = string
type CopyDestinationKey = string

type JcpState struct {
	CopyStates map[CopySourceKey](map[CopyDestinationKey]JcpCopyState)
}

type CopierType int

const (
	BlockCopier CopierType = 1
)

type JcpCopyState struct {
	OpaqueState any
	CopierType  CopierType
	LastUpdate  time.Time
}

func MakeNewCopyState(progress logic.JcpProgress) JcpCopyState {
	return JcpCopyState{OpaqueState: progress.OpaqueState, CopierType: BlockCopier, LastUpdate: time.Now()}
}
