package copier

import (
	discovery "github.com/ben-ha/jcp/discovery"
)

type Copier interface {
	Copy(source discovery.FileInformation, destination discovery.FileInformation, state CopierState) CopierState
	CopyWithProgress(source discovery.FileInformation, destination discovery.FileInformation, state CopierState, progress chan<- CopierProgress) CopierState
}

type CopierState struct {
	Error *error
	State any
}

type CopierProgress struct {
	Size             uint64
	BytesTransferred uint64
}
