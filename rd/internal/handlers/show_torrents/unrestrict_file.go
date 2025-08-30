package showTorrents

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

type commandErrMsg error

type unrestrictFileMsg struct {
	success bool
	result  *realdebrid.UnrestrictedLink
	err     error
}

type unrestrictFileModel struct {
	spinner spinner.Model
	err     error
	message string

	Loading  bool
	Quitting bool
	TaskDone bool

	Result *realdebrid.UnrestrictedLink

	Link string

	rd *realdebrid.RealDebridClient
}

func (m unrestrictFileModel) doUnrestrictFile() tea.Msg {
	res, err := m.rd.UnrestricLink(&realdebrid.UnrestrictProps{
		Link: m.Link,
	})
	if err != nil {
		return unrestrictFileMsg{
			success: false,
			result:  nil,
			err:     err,
		}
	}

	return unrestrictFileMsg{
		success: true,
		result:  res,
		err:     nil,
	}
}

func initialUnrestrictFileModel(link string, rd *realdebrid.RealDebridClient) unrestrictFileModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return unrestrictFileModel{
		spinner: s,
		Link:    link,
		Loading: true,
		rd:      rd,
		message: "Initializing file download...",
	}
}

func (m unrestrictFileModel) Init() tea.Cmd {
	return tea.Batch(m.doUnrestrictFile,
		m.spinner.Tick,
	)
}

func (m unrestrictFileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case unrestrictFileMsg:
		m.Loading = false
		m.TaskDone = true
		m.err = msg.err

		if msg.success {
			m.message = fmt.Sprintf("Download initialize complete: %s", msg.result.ID)
			m.Result = msg.result
		} else {
			m.message = fmt.Sprintf("Error download init: %v", msg.err)

		}

		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	}

}

func (m unrestrictFileModel) View() string {
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

func HandleUnrestrictFileLink(link string, rd *realdebrid.RealDebridClient) (*unrestrictFileModel, error) {
	p := tea.NewProgram(initialUnrestrictFileModel(link, rd))
	r, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("Error running program: %v", err)
	}

	output := r.(unrestrictFileModel)
	if output.err != nil {
		return nil, fmt.Errorf("Error: %v", output.err)
	}

	return &output, nil
}
