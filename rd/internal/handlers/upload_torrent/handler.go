package uploadtorrent

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

type commandErrMsg error

type uploadTorrentMsg struct {
	success bool
	result  *realdebrid.AddTorrent
	err     error
}

type uploadTorrentModel struct {
	spinner spinner.Model
	err     error
	message string

	Loading  bool
	Quitting bool
	TaskDone bool

	Result *realdebrid.AddTorrent

	TorrentFile string

	rd *realdebrid.RealDebridClient
}

func (m uploadTorrentModel) doUploadTorrent() tea.Msg {
	res, err := m.rd.AddTorrent(m.TorrentFile)
	if err != nil {
		return uploadTorrentMsg{
			success: false,
			result:  nil,
			err:     err,
		}
	}

	return uploadTorrentMsg{
		success: true,
		result:  res,
		err:     nil,
	}
}

func initialUploadTorrentModel(torrentFile string, rd *realdebrid.RealDebridClient) uploadTorrentModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return uploadTorrentModel{
		spinner:     s,
		TorrentFile: torrentFile,
		Loading:     true,
		rd:          rd,
		message:     "Uploading torrent...",
	}
}

func (m uploadTorrentModel) Init() tea.Cmd {
	return tea.Batch(
		m.doUploadTorrent,
		m.spinner.Tick,
	)
}

func (m uploadTorrentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case uploadTorrentMsg:
		m.Loading = false
		m.TaskDone = true
		m.err = msg.err

		if msg.success {
			m.message = fmt.Sprintf("Torrent uploaded successfully: %s", msg.result.ID)
			m.Result = msg.result
		} else {
			m.message = fmt.Sprintf("Error uploading torrent: %v", msg.err)
		}

		return m, tea.Quit

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m uploadTorrentModel) View() string {
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

func HandleUploadTorrent(torrentFile string, rd *realdebrid.RealDebridClient) (*uploadTorrentModel, error) {
	p := tea.NewProgram(initialUploadTorrentModel(torrentFile, rd))
	r, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("Error running program: %v", err)
	}

	output := r.(uploadTorrentModel)
	if output.err != nil {
		return nil, fmt.Errorf("Error: %v", output.err)
	}

	return &output, nil
}
