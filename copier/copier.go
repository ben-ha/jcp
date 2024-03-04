package copier

type Copier interface {
	Copy(source string, destination string, state CopierState) CopierState
	CopyWithProgress(source string, destination string, state CopierState, progress chan<- CopierProgress) CopierState
}

type CopierState struct {
	Error *error
	State any
}

type CopierProgress struct {
	Size uint64
	BytesTransferred uint64
}
