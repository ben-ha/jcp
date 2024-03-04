package copier

type Copier interface {
	Copy(source string, destination string, state CopierState) CopierState
}

type CopierState struct {
	Error *error
	State any
}
