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
		m, cmd := m.UpdateModel(msg)
		return m, cmd

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		m, cmd := m.UpdateProgress(msg)
		return m, cmd

	default:
		return m, nil
	}
}

func (model UIModel) UpdateModel(msg UITransferMsg) (UIModel, tea.Cmd) {
	// Find the transfer
	uiCommands := make([]tea.Cmd, 0)
	transferFound := false
	transfersToKeep := make([]UITransfer, 0)
	for _, transfer := range model.Transfers {
		if transfer.CopierProgress.Source == msg.Progress.Source && transfer.CopierProgress.Dest == msg.Progress.Dest {
			transferFound = true
			transfer.CopierProgress = msg.Progress
			var cmd tea.Cmd
			if transfer.CopierProgress.Error == io.EOF {
				cmd = transfer.Progress.SetPercent(1)
			} else {
				cmd = transfer.Progress.SetPercent(float64(transfer.CopierProgress.BytesTransferred) / float64(transfer.CopierProgress.Size))
			}
			uiCommands = append(uiCommands, cmd)
		}

		if !(transfer.Progress.Percent() == 1) {
			transfersToKeep = append(transfersToKeep, transfer)
		}
	}

	model.Transfers = transfersToKeep

	if !transferFound {
		// Add the transfer
		currentProgress := progress.New(progress.WithDefaultGradient())
		cmd := currentProgress.SetPercent(float64(msg.Progress.BytesTransferred) / float64(msg.Progress.Size))
		uiCommands = append(uiCommands, cmd)
		model.Transfers = append(model.Transfers, UITransfer{CopierProgress: msg.Progress, Progress: currentProgress})
	}

	return model, tea.Batch(uiCommands...)
}

func (model UIModel) UpdateProgress(msg progress.FrameMsg) (UIModel, tea.Cmd) {
	var cmds []tea.Cmd
	for idx := range model.Transfers {
		newModel, cmd := model.Transfers[idx].Progress.Update(msg)
		model.Transfers[idx].Progress = newModel.(progress.Model)
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
		line := pad + transfer.CopierProgress.Source + pad + transfer.Progress.View() + pad + transfer.CopierProgress.Dest + "\n"
		str += line
	}
	return str
}
