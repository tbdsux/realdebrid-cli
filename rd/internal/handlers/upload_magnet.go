package handlers

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

type uploadMagnetMsg struct {
	success bool
	result  *realdebrid.AddTorrent
	err     error
}

type uploadMagnetModel struct {
	spinner spinner.Model
	err     error
	message string

	Loading  bool
	Quitting bool
	TaskDone bool

	Result *realdebrid.AddTorrent

	MagnetLink string

	rd *realdebrid.RealDebridClient
}

func (m uploadMagnetModel) doUploadMagnet() tea.Msg {
	res, err := m.rd.AddMagnet(m.MagnetLink)
	if err != nil {
		return uploadMagnetMsg{
			success: false,
			result:  nil,
			err:     err,
		}
	}

	return uploadMagnetMsg{
		success: true,
		result:  res,
		err:     nil,
	}
}

func initialUploadMagnetModel(magnetLink string, rd *realdebrid.RealDebridClient) uploadMagnetModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return uploadMagnetModel{
		spinner:    s,
		MagnetLink: magnetLink,
		Loading:    true,
		rd:         rd,
		message:    "Uploading magnet link...",
	}
}

func (m uploadMagnetModel) Init() tea.Cmd {
	return tea.Batch(m.doUploadMagnet,
		m.spinner.Tick,
	)
}

func (m uploadMagnetModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.Quitting = true
		return m, nil

	case uploadMagnetMsg:
		m.Loading = false
		m.TaskDone = true
		m.err = msg.err

		if msg.success {
			m.message = fmt.Sprintf("Magnet link uploaded successfully: %s", msg.result.ID)
			m.Result = msg.result
		} else {
			m.message = fmt.Sprintf("Error uploading magnet link: %v", msg.err)

		}

		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	}

}

func (m uploadMagnetModel) View() string {
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

func HandleUploadMagnetLink(magnetLink string, rd *realdebrid.RealDebridClient) (*uploadMagnetModel, error) {
	p := tea.NewProgram(initialUploadMagnetModel(magnetLink, rd))
	r, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("Error running program: %v", err)
	}

	output := r.(uploadMagnetModel)
	if output.err != nil {
		return nil, fmt.Errorf("Error: %v", output.err)
	}

	return &output, nil
}
