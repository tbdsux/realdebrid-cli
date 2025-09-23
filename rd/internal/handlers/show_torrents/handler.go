package showTorrents

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
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			m.IsQuitting = true
			return m, tea.Quit

		case tea.KeyEnter.String():
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

func formatListDesc(i realdebrid.Torrent) string {
	if i.Status == "downloaded" {
		return fmt.Sprintf("ID: %s\tSize: %s", i.ID, internal.ByteCountSI(i.Bytes))

	}

	return fmt.Sprintf("ID: %s\tSize: %s\t Status: %s \tProgress: %.2f%%", i.ID, internal.ByteCountSI(i.Bytes), i.Status, i.Progress)
}

func ShowTorrentsList(dls []realdebrid.Torrent, page int) (*realdebrid.Torrent, error) {
	var items []list.Item

	for _, i := range dls {
		items = append(items, item{
			title: i.Filename,
			desc:  formatListDesc(i),
			id:    i.ID,
		})
	}

	m := showModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = fmt.Sprintf("RealDebrid Torrents List | Page: %v", page)

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
	var selected realdebrid.Torrent
	for _, v := range dls {
		if v.ID == output.Choice.id {
			selected = v
		}
	}

	return &selected, nil
}
