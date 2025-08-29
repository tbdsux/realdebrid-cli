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
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type showModel struct {
	list list.Model
}

func (m showModel) Init() tea.Cmd {
	return nil
}

func (m showModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
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

func ShowTorrentsList(dls []realdebrid.Torrent, page int) error {
	var items []list.Item

	for _, i := range dls {
		items = append(items, item{
			title: i.Filename,
			desc:  fmt.Sprintf("ID: %s\tSize: %s\t", i.ID, internal.ByteCountSI(i.Bytes)),
		})
	}

	m := showModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = fmt.Sprintf("RealDebrid Torrents List | Page: %v", page)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
