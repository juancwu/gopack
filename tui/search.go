package tui

import (
	"fmt"
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
	ti    textinput.Model
	state searchModelState
	keys  searchModelKeyMap
	help  help.Model
	im    tea.Model // installModel, not defined as installModel type because Go doesn't accept it
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
		m.im, cmd = m.im.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Enter):
			m.im = NewInstallModel([]string{m.ti.Value()}, false)
			m.state = searchingState
			cmd = m.im.Init()
			return m, cmd
		}
	}

	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

func (m searchModel) View() string {
	switch m.state {
	case inputState:
		return strings.Join([]string{m.ti.View(), m.help.View(m.keys)}, "\n\n")
	case searchingState:
		return m.im.View()
	}

	return fmt.Sprintf("Invalid Search Model State: %s", m.state)
}
