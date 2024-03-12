package copier

import (
	discovery "github.com/ben-ha/jcp/discovery"
)

type Copier interface {
	Copy(source discovery.FileInformation, destination discovery.FileInformation, state CopierState) CopierState
	CopyWithProgress(source discovery.FileInformation, destination discovery.FileInformation, state CopierState, progress chan<- CopierProgress) CopierState
}

type CopierState struct {
	Error error
	State any
}

type CopierProgress struct {
	Source           string
	Dest             string
	Size             uint64
	BytesTransferred uint64
	Error            error
	OpaqueState      any
}

func IsCopyRequired(source discovery.FileInformation, destination discovery.FileInformation) bool {
	if source.Info == nil || destination.Info == nil {
		return true
	}

	return source.Info.Size() != destination.Info.Size()
}
