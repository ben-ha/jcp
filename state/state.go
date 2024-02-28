package state

type CopierState struct {
	CopyStates map[CopyStateKey]*CopyState
}

type CopyStateKey struct {
	Source      string
	Destination string
}

type CopyState struct {
	FileStates     map[FileStateKey]*FileCopyState
	DiscoveryState []string
}

type FileStateKey struct {
	Source      string
	Destination string
}

type FileCopyState struct {
	Size             uint64
	BytesTransferred uint64
}
