package handlers

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

type autoSelectTorrentModel struct {
	spinner spinner.Model
	err     error
	message string

	Loading  bool
	Quitting bool
	TaskDone bool

	Result bool

	TorrentId string
	rd        *realdebrid.RealDebridClient
}

type autoSelectRes struct {
	success bool
	err     error
}

func (m autoSelectTorrentModel) doAutoSelectTorrent() tea.Msg {
	err := m.rd.SelectTorrentFiles(m.TorrentId, []string{})
	if err != nil {
		return autoSelectRes{
			success: false,
			err:     err,
		}
	}

	return autoSelectRes{
		success: true,
		err:     nil,
	}

}

func initialAutoSelectTorrent(torrentId string, rd *realdebrid.RealDebridClient) autoSelectTorrentModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = s.Style.Foreground(lipgloss.Color("205"))

	return autoSelectTorrentModel{
		spinner:   s,
		Loading:   true,
		TorrentId: torrentId,
		rd:        rd,
		message:   "Auto selecting torrent files...",
	}
}

func (m autoSelectTorrentModel) Init() tea.Cmd {
	return tea.Batch(
		m.doAutoSelectTorrent,
		m.spinner.Tick,
	)
}

func (m autoSelectTorrentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case commandErrMsg:
		m.err = msg
		return m, nil

	case autoSelectRes:
		m.Loading = false
		m.TaskDone = true
		m.err = msg.err

		if msg.success {
			m.message = "Auto selected files successfully"
			m.Result = true
		} else {
			m.message = fmt.Sprintf("Error auto selecting files: %v", msg.err)
		}

		return m, tea.Quit

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m autoSelectTorrentModel) View() string {
	if m.Quitting {
		return ""
	}

	spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	messageStyle := lipgloss.NewStyle().MarginLeft(1)
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	var output string
	if m.Loading {
		message := messageStyle.Render(m.message)
		output = fmt.Sprintf("\n  %s %s\n", spinnerStyle.Render(m.spinner.View()), message)
	} else if m.TaskDone {
		checkmark := successStyle.Render("✓")
		message := messageStyle.Render(m.message)
		output = fmt.Sprintf("\n  %s %s\n", checkmark, message)
	} else {
		cross := errorStyle.Render("✗")
		message := messageStyle.Render(m.message)
		output = fmt.Sprintf("\n  %s %s\n", cross, message)
	}

	return output
}

func AutoSelectFiles(torrentId string, rd *realdebrid.RealDebridClient) error {
	p := tea.NewProgram(initialAutoSelectTorrent(torrentId, rd))
	r, err := p.Run()
	if err != nil {
		return fmt.Errorf("Error running program: %v", err)
	}

	output := r.(autoSelectTorrentModel)
	if output.err != nil {
		return fmt.Errorf("Error: %v", output.err)
	}

	return nil
}
