package showDownloads

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tbdsux/realdebrid-cli/rd/internal"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc, id string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type showModel struct {
	list list.Model

	IsQuitting bool
	Choice     item
}

func (m showModel) Init() tea.Cmd {
	return nil
}

func (m showModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.IsQuitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.Choice = i
			}
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m showModel) View() string {
	return docStyle.Render(m.list.View())
}

func ShowDownloadsList(dls []realdebrid.Download) (*realdebrid.Download, error) {
	var items []list.Item

	for _, i := range dls {
		items = append(items, item{
			title: i.Filename,
			desc:  fmt.Sprintf("ID: %s\tSize: %s\t", i.ID, internal.ByteCountSI(i.FileSize)),
			id:    i.ID,
		})
	}

	m := showModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "RealDebrid Downloads List"

	p := tea.NewProgram(m, tea.WithAltScreen())

	model, err := p.Run()
	if err != nil {
		return nil, err
	}

	output := model.(showModel)

	if output.IsQuitting {
		return nil, nil
	}

	// Filter result selected
	var selected realdebrid.Download
	for _, v := range dls {
		if v.ID == output.Choice.id {
			selected = v
		}
	}

	return &selected, nil
}
