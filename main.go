package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/ben-ha/jcp/logic"
	"github.com/ben-ha/jcp/sleepless"
	"github.com/ben-ha/jcp/state"
	"github.com/ben-ha/jcp/tui"
	tea "github.com/charmbracelet/bubbletea"
)

var Version string = "dev"

func updaterFunc(jcp logic.Jcp, ui *tea.Program, state *state.JcpState) {
	defer ui.Quit()

	for {
		update, more := <-jcp.ProgressChannel
		if !more {
			break
		}

		if update.JcpError != nil {
			if update.JcpError == io.EOF {
				ui.Send(tui.UIErrorMessage("Done"))
			} else {
				ui.Send(tui.UIErrorMessage(update.JcpError.Error()))
			}
			break
		}

		state.Update(update)
		ui.Send(tui.UITransferMsg{Progress: update.Progress})
	}
}

func StartUI(prog *tea.Program, uiComplete *sync.WaitGroup) {
	defer uiComplete.Done()
	if _, err := prog.Run(); err != nil {
		panic(fmt.Sprintf("Error in UI: %v", err))
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Jcp %s\n", Version)
		fmt.Printf("Usage: %v <src> <dest>\n", os.Args[0])
		return
	}

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

	sleepEnabler, _ := sleepless.PreventSleep("jcp", "copy in progress")
	defer sleepEnabler()

	go StartUI(ui, &uiComplete)

	go updaterFunc(jcp, ui, state)

	jcp.StartCopy(src, dst)

	uiComplete.Wait()
}
