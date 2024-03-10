package state

import (
	"time"

	"github.com/ben-ha/jcp/logic"
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

func MakeNewCopyState(progress logic.JcpProgress) JcpCopyState {
	percent := float64(progress.Progress.BytesTransferred) / float64(progress.Progress.Size)
	return JcpCopyState{OpaqueState: progress.OpaqueState, CopierType: BlockCopier, LastUpdate: time.Now(), Percent: percent}
}
