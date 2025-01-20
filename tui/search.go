package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type searchModelState string

const (
	inputState     searchModelState = "input"
	searchingState searchModelState = "searching"
)

// searchModelKeyMap implements the help.KeyMap interface
type searchModelKeyMap struct {
	Enter key.Binding
	Quit  key.Binding
}

func (k searchModelKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Quit}
}

func (k searchModelKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter},
		{k.Quit},
	}
}

type searchModel struct {
	ti      textinput.Model
	state   searchModelState
	keys    searchModelKeyMap
	help    help.Model
	history string
	im      tea.Model // installModel, not defined as installModel type because Go doesn't accept it
}

func NewSearchModel() searchModel {
	ti := textinput.New()
	ti.Placeholder = "Search Packages"
	ti.Focus()

	keyMap := searchModelKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Search"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "esc"),
			key.WithHelp("esc", "Quit program"),
		),
	}

	return searchModel{
		ti:    ti,
		state: inputState, // initial state in input mode
		keys:  keyMap,
		help:  help.New(),
	}
}

func (m searchModel) Init() tea.Cmd {
	return nil
}

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// skip key matching when using the installModel
	if m.state == searchingState {
		switch msg := msg.(type) {
		case afterInstallMsg:
			// let the installModel record the installation history
			m.im, _ = m.im.Update(msg)
			// add the latest recorded history
			m.history += m.im.(installModel).History()
			// reset the state for the next input
			m.state = inputState
		default:
			m.im, cmd = m.im.Update(msg)
			return m, cmd
		}

	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Enter):
			m.im = NewInstallModel([]string{m.ti.Value()}, false)
			m.state = searchingState
			m.ti.Reset()
			cmd = m.im.Init()
			return m, cmd
		}
	}

	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

func (m searchModel) View() string {
	// let the installModel handle the rendering
	if m.state == searchingState {
		return m.im.View()
	}

	// render the input and history
	var builder strings.Builder
	builder.WriteString(m.history)
	builder.WriteString(strings.Join([]string{m.ti.View(), m.help.View(m.keys)}, "\n\n"))
	return wrapper.Render(builder.String())
}
