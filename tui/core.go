package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/juancwu/gopack/util"
)

type model struct {
	textInput  textinput.Model
	spinner    spinner.Model
	choices    []string
	cursor     int
	startIdx   int
	err        error
	step       int
	selection  string
	installing bool
}

var (
	installing bool = false
	finished   bool = false
	results    []string
)

func NewModel() model {
	ti := textinput.New()
	ti.Placeholder = "Package name"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		choices:    []string{},
		textInput:  ti,
		spinner:    s,
		err:        nil,
		step:       0,
		selection:  "",
		installing: false,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) FirstStepUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			pkgName := m.textInput.Value()
			results = util.Search(pkgName)
			m.choices = results
			m.step++
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) SecondStepUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.choices) == 0 {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			selection := util.GetPkgUrl(m.choices[m.cursor])
			if selection == "" {
				return m, tea.Quit
			}
			m.selection = selection
			m.step++
		}
	}

	const pageSize = 5
	m.startIdx = max(0, min(len(m.choices)-pageSize, m.cursor-(pageSize/2)))

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) ThirdStepUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !installing {
		m.installing = true
		go m.install()
	} else if finished {
		return m, tea.Quit
	} else {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) install() {
	installing = true
	util.RunGoInstall(m.selection)
	finished = true
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.step == 0 {
		return m.FirstStepUpdate(msg)
	} else if m.step == 1 {
		return m.SecondStepUpdate(msg)
	} else {
		return m.ThirdStepUpdate(msg)
	}
}

func (m model) View() string {
	// The header
	var s string
	if m.step == 0 {
		s = fmt.Sprintf("Enter package name:\n\n%s\n\n%s\n", m.textInput.View(), "(esc to quit)")
	} else if m.step == 1 {
		s = "Choose package to install:\n\n"

		// Iterate over our choices
		for i, choice := range results[m.startIdx:min(m.startIdx+5, len(results))] {

			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.cursor == m.startIdx+i {
				cursor = ">" // cursor!
			}

			// Render the row
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		// Number of results
		s += fmt.Sprintf("\nNumber of results: %d\n", len(m.choices))

		// The footer
		s += "\nPress q to quit.\n"
	} else {
		if m.installing {
			s = fmt.Sprintf("\n\n %s Installing package %s\n\n", m.spinner.View(), m.selection)
		}
	}

	// Send the UI for rendering
	return s
}
