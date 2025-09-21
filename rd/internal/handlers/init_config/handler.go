package initconfig

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type configModel struct {
	TextInput textinput.Model
	err       error

	Quitting bool
	Success  bool
}

type commandErrMsg error

func initialConfigModel() configModel {
	ti := textinput.New()
	ti.Placeholder = "Enter your API key"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 60

	return configModel{
		TextInput: ti,
		err:       nil,
	}
}

func (m configModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.Success = true
			return m, tea.Quit

		case tea.KeyCtrlC, tea.KeyEsc:
			m.Quitting = true
			return m, tea.Quit
		}

	case commandErrMsg:
		m.err = msg
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m configModel) View() string {
	return fmt.Sprintf("Initial Configuration\nPlease enter your RealDebrid API Key:\n\n%s\n\n%s",
		m.TextInput.View(),
		"Press Enter to submit, Ctrl+C or Esc to exit.",
	) + "\n\n\n"
}

func AskConfigForSetup() (*configModel, error) {
	p := tea.NewProgram(initialConfigModel())
	r, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("Error running program: %v", err)
	}

	output := r.(configModel)
	if output.err != nil {
		return nil, fmt.Errorf("Error: %v", output.err)
	}

	return &output, nil
}
