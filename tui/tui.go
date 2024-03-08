package tui

import (
	"io"
	"strings"

	"github.com/ben-ha/jcp/copier"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	padding  = 2
	maxWidth = 80
)

type UIModel struct {
	Transfers []UITransfer
}

type UITransferMsg struct {
	Progress copier.CopierProgress
}

type UITransfer struct {
	CopierProgress copier.CopierProgress
	Progress       progress.Model
	shownCompleted bool
}

func (m UIModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m = m.UpdateWidth(msg)
		return m, nil

	case UITransferMsg:
		m = m.UpdateModel(msg)
		return m, nil

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		m, cmd := m.UpdateProgress(msg)
		return m, cmd

	default:
		return m, nil
	}
}

func (model UIModel) UpdateModel(msg UITransferMsg) UIModel {
	// Find the transfer
	transferFound := false
	for idx, transfer := range model.Transfers {
		if transfer.CopierProgress.Source == msg.Progress.Source && transfer.CopierProgress.Dest == msg.Progress.Dest {
			transferFound = true
			transfer.CopierProgress = msg.Progress
			if transfer.CopierProgress.Error == io.EOF {
				transfer.Progress.SetPercent(1)
			} else {
				transfer.Progress.SetPercent(float64(transfer.CopierProgress.BytesTransferred) / float64(transfer.CopierProgress.Size))
			}
		}

		if transfer.Progress.Percent() == 1 && transfer.shownCompleted {
			// Remove transfer if done
			newTransfers := model.Transfers[:idx]
			if len(model.Transfers) > (idx + 1) {
				newTransfers = append(newTransfers, model.Transfers[idx+1:]...)
			}
			model.Transfers = newTransfers
		}
	}

	if !transferFound {
		// Add the transfer
		currentProgress := progress.New(progress.WithDefaultGradient())
		currentProgress.SetPercent(float64(msg.Progress.BytesTransferred) / float64(msg.Progress.Size))
		model.Transfers = append(model.Transfers, UITransfer{CopierProgress: msg.Progress, Progress: currentProgress})
	}

	return model
}

func (model UIModel) UpdateProgress(msg progress.FrameMsg) (UIModel, tea.Cmd) {
	var cmds []tea.Cmd
	for _, transfer := range model.Transfers {
		newModel, cmd := transfer.Progress.Update(msg)
		transfer.Progress = newModel.(progress.Model)
		cmds = append(cmds, cmd)
	}

	return model, tea.Batch(cmds...)
}

func (model UIModel) UpdateWidth(msg tea.WindowSizeMsg) UIModel {
	for _, transfer := range model.Transfers {
		transfer.Progress.Width = msg.Width - padding*2 - 4
		if transfer.Progress.Width > maxWidth {
			transfer.Progress.Width = maxWidth
		}
	}

	return model
}

func (m UIModel) View() string {
	pad := strings.Repeat(" ", padding)
	str := "\n"
	for _, transfer := range m.Transfers {
		line := pad + transfer.CopierProgress.Source + transfer.Progress.View() + transfer.CopierProgress.Dest + "\n\n"
		str += line

		if transfer.Progress.Percent() == 1 {
			transfer.shownCompleted = true
		}
	}
	return str
}
