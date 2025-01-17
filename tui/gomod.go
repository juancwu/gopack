package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/juancwu/gopack/util"
)

type gomodModel struct {
	spinner  spinner.Model
	progress progress.Model
	history  []installResult
	modules  []string
	idx      int
	isDone   bool
}

func NewGoModModel(modules []string) gomodModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	p := progress.New(progress.WithDefaultGradient())
	return gomodModel{
		spinner:  s,
		progress: p,
		history:  []installResult{},
		modules:  modules,
		idx:      0,
		isDone:   false,
	}
}

func (m gomodModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.installCmd(m.modules[m.idx]),
	)
}
func (m gomodModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case afterInstallMsg:
		// record installation history
		m = m.recordHistory(msg.Err)
		// install the next module if any
		return m.install()
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)

		if m.isDone {
			return m, tea.Quit
		}

		return m, cmd
	default:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m gomodModel) View() string {
	var builder strings.Builder

	// render the installation history
	for _, record := range m.history {
		if record.success {
			builder.WriteString(okText.Render(record.title) + "\n")
		} else {
			builder.WriteString(errText.Render(record.title) + "\n")
		}
	}

	if builder.Len() > 0 {
		builder.WriteString("\n")
	}

	if m.isDone {
		builder.WriteString("Done!\n")
		return wrapper.Render(builder.String())
	}

	builder.WriteString(m.spinner.View() + fmt.Sprintf(" %s Installing '%s'\n", m.progress.View(), m.modules[m.idx]))

	return wrapper.Render(builder.String())
}

func (m gomodModel) installCmd(module string) tea.Cmd {
	return func() tea.Msg {
		err := util.RunGoGet(module)
		return afterInstallMsg{Err: err}
	}
}

func (m gomodModel) install() (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.idx < len(m.modules)-1 {
		m.idx += 1
		n := float64(len(m.modules))
		cmds = append(cmds, m.progress.SetPercent(float64(m.idx)/n), m.installCmd(m.modules[m.idx]))
	} else {
		m.isDone = true
		cmds = append(cmds, m.progress.SetPercent(1.0))
	}

	return m, tea.Batch(cmds...)
}

func (m gomodModel) recordHistory(err error) gomodModel {
	var s string
	if err != nil {
		s = fmt.Sprintf("Error installing '%s': %s", m.modules[m.idx], err.Error())
	} else {
		s = fmt.Sprintf("Successfully installed '%s'", m.modules[m.idx])
	}
	m.history = append(m.history, installResult{title: s, success: err == nil})
	return m
}
