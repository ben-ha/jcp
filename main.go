package main

import (
	"fmt"

	state "github.com/ben-ha/jcp/state"
)

func main() {
	st := &state.CopierState{CopyStates: make(map[string]map[string]state.CopyState)}
	st.StartCopy("bla", "bbb")
	st.Save("/workspaces/try.txt")
	fmt.Println("Hello world!")
}
