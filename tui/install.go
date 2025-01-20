package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/juancwu/gopack/util"
)

type installModel struct {
	spinner             spinner.Model
	list                list.Model
	isSearching         bool
	isInstalling        bool
	selectFirst         bool
	queries             []string
	results             []list.Item
	searchingTerm       string
	installingTerm      string
	installationHistory []installResult
	current_query_idx   int
	err                 error
	isDone              bool
	// name is the model name
	name string
	// asComponent represnets if the model is being used as part of a component to a parent model
	asComponent bool
}

type installResult struct {
	title   string
	success bool
}

func (s installResult) Title() string { return s.title }

// searchResult represents a single search result
type searchResult string

// Implementation of list.DefaultItem and list.Item interfaces for SearchResult
func (s searchResult) Title() string       { return string(s) }
func (s searchResult) Description() string { return "" }
func (s searchResult) FilterValue() string { return "" }

func NewInstallModel(queries []string, selectFirst bool) installModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return installModel{
		spinner:             s,
		list:                list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		isSearching:         true,
		isInstalling:        false,
		searchingTerm:       queries[0],
		installingTerm:      "",
		installationHistory: []installResult{},
		selectFirst:         selectFirst,
		queries:             queries,
		current_query_idx:   0,
		err:                 nil,
		isDone:              false,
		name:                "Install Model",
	}
}

func (m installModel) Init() tea.Cmd {
	// start first search here
	return tea.Batch(
		m.spinner.Tick,
		searchCmd(m.queries[m.current_query_idx]),
	)
}

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !m.isSearching {
				return m.install()
			}
		case "ctrl+c": // force quit
			if m.asComponent {
				return m, quitCmd(msg.String(), m.name)
			}
			return m, tea.Quit
		case "q", "esc": // normal quit
			if !m.isSearching || !m.isInstalling {
				if m.asComponent {
					return m, quitCmd(msg.String(), m.name)
				}
				return m, tea.Quit
			}
		default:
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	case afterSearchMsg:
		m.results = msg.results
		if m.selectFirst && len(msg.results) > 0 {
			return m.install()
		} else {
			return m.showSearchResults()
		}
	case afterInstallMsg:
		m = m.recordHistory(msg.Err)
		// search the next query
		return m.search()
	default:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m installModel) View() string {
	var builder strings.Builder

	// render the installation history
	builder.WriteString(m.renderHistory())

	if builder.Len() > 0 {
		builder.WriteString("\n")
	}

	if m.isDone {
		builder.WriteString("Done!\n")
		return wrapper.Render(builder.String())
	}

	if m.isSearching {
		builder.WriteString(m.spinner.View() + fmt.Sprintf(" Searching '%s'\n", m.searchingTerm))
	}
	if m.isInstalling {
		builder.WriteString(m.spinner.View() + fmt.Sprintf(" Installing '%s'\n", m.installingTerm))
	}
	if !m.isSearching && !m.isInstalling {
		// show search result
		builder.WriteString(m.list.View())
	}

	return wrapper.Render(builder.String())
}

func (m installModel) History() string {
	return m.renderHistory()
}

func (m *installModel) SetAsComponent(enabled bool) {
	m.asComponent = enabled
}

func (m installModel) install() (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var item list.Item

	if m.selectFirst {
		item = m.results[0]
	} else {
		item = m.list.SelectedItem()
	}

	if s, ok := item.(searchResult); ok {
		title := string(s)
		cmd = m.installCmd(title)
		m.searchingTerm = ""
		m.installingTerm = title
		m.isSearching = false
		m.isInstalling = true
		m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	}

	return m, cmd
}

func (m installModel) search() (tea.Model, tea.Cmd) {
	if m.current_query_idx < len(m.queries)-1 {
		m.isSearching = true
		m.isInstalling = false
		m.current_query_idx += 1
		m.searchingTerm = m.queries[m.current_query_idx]
		return m, searchCmd(m.searchingTerm)
	}
	return m.end(nil)
}

func (m installModel) showSearchResults() (tea.Model, tea.Cmd) {
	m.isSearching = false
	m.isInstalling = false

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false

	keys := list.DefaultKeyMap()

	// don't want to show these in help section, they are disabled
	keys.Filter.SetEnabled(false)
	keys.CancelWhileFiltering.SetEnabled(false)
	keys.ForceQuit.SetEnabled(false)
	keys.AcceptWhileFiltering.SetEnabled(false)
	keys.ClearFilter.SetEnabled(false)

	m.list = list.New(m.results, delegate, 50, 20)
	m.list.SetFilteringEnabled(false)
	m.list.Title = "Search Result: " + m.searchingTerm
	m.list.KeyMap = keys

	m.searchingTerm = ""
	m.installingTerm = ""

	return m, nil
}

func (m installModel) recordHistory(err error) installModel {
	var s string
	if err != nil {
		s = fmt.Sprintf("Error installing '%s': %s", m.installingTerm, err.Error())
	} else {
		s = fmt.Sprintf("Successfully installed '%s'", m.installingTerm)
	}
	m.installationHistory = append(m.installationHistory, installResult{title: s, success: err == nil})
	return m
}

func (m installModel) renderHistory() string {
	var builder strings.Builder
	for _, record := range m.installationHistory {
		if record.success {
			builder.WriteString(okText.Render(record.title) + "\n")
		} else {
			builder.WriteString(errText.Render(record.title) + "\n")
		}
	}
	return builder.String()
}

func (m installModel) end(err error) (tea.Model, tea.Cmd) {
	m.isSearching = false
	m.isInstalling = false
	m.searchingTerm = ""
	m.installingTerm = ""
	m.isDone = true
	m.err = err
	if m.asComponent {
		return m, quitCmd("", m.name)
	}
	return m, tea.Quit
}

type afterSearchMsg struct {
	results []list.Item
}

// searchCmd searches the go packages and returns a tea.Msg so that the searchCmd model
// can update the TUI.
func searchCmd(term string) tea.Cmd {
	return func() tea.Msg {
		results := util.Search(term)
		msg := afterSearchMsg{
			results: make([]list.Item, len(results)),
		}
		for i, res := range results {
			msg.results[i] = searchResult(res)
		}
		return msg
	}
}

type afterInstallMsg struct {
	Err error
}

func (m installModel) installCmd(term string) tea.Cmd {
	return func() tea.Msg {
		pkg := util.GetPkgUrl(term)
		err := util.RunGoGet(pkg)
		return afterInstallMsg{Err: err}
	}
}

type quitMsg struct {
	// Model represents the model name that sent the msg
	Model string
	// Key is the key combination input
	Key string
}

func quitCmd(key string, model string) tea.Cmd {
	return func() tea.Msg {
		return quitMsg{Model: model, Key: key}
	}
}
