package state

type CopierState struct {
	CopyStates map[CopySourceKey](map[CopyDestinationKey]CopyState)
}

type CopySourceKey = string
type CopyDestinationKey = string
type FileStateSourceKey = string
type FileStateDestinationKey = string

type CopyState struct {
	FileStates     map[FileStateSourceKey]map[FileStateDestinationKey]FileCopyState
	DiscoveryState []string
}

type FileCopyState struct {
	Size             uint64
	BytesTransferred uint64
}
