package addMagnet

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type commandErrMsg error

type addMagnetModel struct {
	Textarea textarea.Model
	err      error

	Quitting bool
	Success  bool
}

func initialAddMagnetModel() addMagnetModel {
	ti := textarea.New()
	ti.Placeholder = "Add magnet..."
	ti.Focus()
	ti.SetWidth(60)

	return addMagnetModel{
		Textarea: ti,
		err:      nil,
	}
}

func (m addMagnetModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m addMagnetModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.Textarea.Focused() {
				m.Textarea.Blur()
			}
		case tea.KeyCtrlC:
			m.Quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			m.Success = true
			return m, tea.Quit
		default:
			if !m.Textarea.Focused() {
				cmd = m.Textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	case commandErrMsg:
		m.err = msg
		m.Quitting = true
		return m, nil
	}

	m.Textarea, cmd = m.Textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m addMagnetModel) View() string {
	return fmt.Sprintf("Add a Torrent Magnet Link\n\n%s\n\n%s",
		m.Textarea.View(),
		"(Press Enter to submit, Ctrl+C to exit)",
	) + "\n\n\n"
}

func HandleAskMagnetLink() (*addMagnetModel, error) {
	p := tea.NewProgram(initialAddMagnetModel())
	r, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("Error running program: %v", err)
	}

	output := r.(addMagnetModel)
	if output.err != nil {
		return nil, fmt.Errorf("Error: %v", output.err)
	}

	return &output, err
}
