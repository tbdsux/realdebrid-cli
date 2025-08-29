package showDownloads

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/imroc/req/v3"
	"github.com/tbdsux/realdebrid-cli/realdebrid"
)

const (
	padding  = 2
	maxWidth = 80
)

var (
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
	messageStyle = lipgloss.NewStyle().MarginLeft(1)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

func finalPause() tea.Cmd {
	return tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg {
		return nil
	})
}

type downloadProgressModel struct {
	progress progress.Model
	err      error
	Quitting bool
	Fail     bool
}

type downloadErrMsg error

type downloadMsg struct {
	err          error
	progressSize float64
}

func (m downloadProgressModel) Init() tea.Cmd {
	return nil
}

func (m downloadProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case downloadErrMsg:
		m.err = msg
		m.Fail = true
		m.Quitting = true
		return m, tea.Quit

	case downloadMsg:
		var cmds []tea.Cmd

		if msg.progressSize >= 1.0 {
			cmds = append(cmds, tea.Sequence(finalPause(), tea.Quit))
		}

		cmds = append(cmds, m.progress.SetPercent(float64(msg.progressSize)))
		return m, tea.Batch(cmds...)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func ShowSuccessDLMessage(msg string) string {
	checkmark := successStyle.Render("✓")
	message := messageStyle.Render(msg)
	return "\n" + fmt.Sprintf("\n  %s %s\n", checkmark, message) + "\n"
}

func ShowFailDLMessage(msg string) string {
	cross := errorStyle.Render("✗")
	message := messageStyle.Render(msg)
	return "\n" + fmt.Sprintf("\n  %s %s\n", cross, message) + "\n"
}

func (m downloadProgressModel) View() string {
	if m.err != nil {
		cross := errorStyle.Render("✗")
		message := messageStyle.Render(m.err.Error())
		return "\n" + fmt.Sprintf("\n  %s %s\n", cross, message) + "\n"
	}

	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press ESC or Ctrl+C to quit (will stop download)")
}

type downloader struct {
	dlLink   string
	fileName string
	p        *tea.Program
}

func (d *downloader) Start() {
	client := req.C()

	callback := func(info req.DownloadInfo) {
		if info.Response.StatusCode >= 300 {
			d.p.Send(downloadErrMsg(fmt.Errorf("Status :: %v", info.Response.Status)))
		}

		if info.Response.Response != nil {
			d.p.Send(downloadMsg{
				progressSize: float64(info.DownloadedSize) / float64(info.Response.ContentLength),
			})
		}
	}

	client.R().
		SetOutputFile(d.fileName).
		SetDownloadCallback(callback).
		Get(d.dlLink)
}

func DoDownloadFile(f realdebrid.Download) (*downloadProgressModel, error) {
	model := downloadProgressModel{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	p := tea.NewProgram(model)

	dler := &downloader{
		fileName: f.Filename,
		dlLink:   f.Download,
		p:        p,
	}

	go dler.Start()

	r, err := p.Run()
	if err != nil {
		return nil, err
	}

	output := r.(downloadProgressModel)

	return &output, nil
}
