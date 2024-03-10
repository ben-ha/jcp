package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/ben-ha/jcp/logic"
	"github.com/ben-ha/jcp/state"
	"github.com/ben-ha/jcp/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func updaterFunc(jcp logic.Jcp, ui *tea.Program, state *state.JcpState) {
	for {
		update, more := <-jcp.ProgressChannel
		if !more {
			break
		}

		if update.JcpError != nil {
			if update.JcpError == io.EOF {
				break
			}

			panic(fmt.Sprintf("An error occurred: %v", update.JcpError.Error()))
		}

		state.Update(update)
		ui.Send(tui.UITransferMsg{Progress: update.Progress})
	}

	ui.Quit()
}

func StartUI(prog *tea.Program, uiComplete *sync.WaitGroup) {
	defer uiComplete.Done()
	if _, err := prog.Run(); err != nil {
		panic(fmt.Sprintf("Error in UI: %v", err))
	}
}

func main() {
	src := os.Args[1]
	dst := os.Args[2]

	state, err := state.InitializeState()
	if err != nil {
		panic(fmt.Sprintf("State: %v", err))
	}

	defer state.SaveState()

	jcp := logic.MakeJcp(10, *state)

	ui := tea.NewProgram(tui.UIModel{})
	uiComplete := sync.WaitGroup{}
	uiComplete.Add(1)

	go StartUI(ui, &uiComplete)

	go updaterFunc(jcp, ui, state)

	err = jcp.StartCopy(src, dst)
	if err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}

	uiComplete.Wait()
}
