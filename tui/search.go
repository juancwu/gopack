package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/juancwu/gopack/util"
)

type SearchModel struct {
	spinner           spinner.Model
	list              list.Model
	searching         bool
	installing        bool
	pickFirst         bool
	queries           []string
	current_query_idx int
	err               error
}

// SearchResult represents a single search result
type SearchResult struct {
	title       string
	description string
}

// Implementation of list.DefaultItem and list.Item interfaces for SearchResult
func (s SearchResult) Title() string       { return s.title }
func (s SearchResult) Description() string { return s.description }
func (s SearchResult) FilterValue() string { return s.title }

func NewSearchModel(queries []string, pickFirst bool) SearchModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return SearchModel{
		spinner:           s,
		list:              list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		searching:         true,
		installing:        false,
		pickFirst:         pickFirst,
		queries:           queries,
		current_query_idx: 0,
		err:               nil,
	}
}

func (m SearchModel) Init() tea.Cmd {
	// start first search here
	return tea.Batch(
		m.spinner.Tick,
		search(m.queries[m.current_query_idx]),
	)
}

type searchMsg struct {
	results []list.Item
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !m.searching {
				if s, ok := m.list.SelectedItem().(SearchResult); ok {
					cmds = append(cmds, install(s.Title()))
					m.searching = false
					m.installing = true
					m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
				}
			}
		case "ctrl+c":
			return m, tea.Quit
		default:
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}
	case searchMsg:
		m.searching = false
		if m.pickFirst && len(msg.results) > 0 {
			m.installing = true
			cmds = append(cmds, install(msg.results[0].FilterValue()))
		} else {
			m.list = list.New(msg.results, list.NewDefaultDelegate(), 50, 20)
			m.list.SetFilteringEnabled(false)
			m.list.Title = "Search Result: " + m.queries[m.current_query_idx]
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}
	case installMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, tea.Quit
		}
		if m.current_query_idx < len(m.queries)-1 {
			m.searching = true
			m.installing = false
			m.current_query_idx += 1
			cmds = append(cmds, search(m.queries[m.current_query_idx]))
		} else {
			m.searching = false
			m.installing = false
			return m, tea.Quit
		}
	default:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m SearchModel) View() string {
	s := "Bye"
	if m.err != nil {
		s = "Error: " + m.err.Error()
		s = errText.Render(s)
	} else if m.searching {
		s = m.spinner.View() + " Searching..."
	} else if m.installing {
		s = m.spinner.View() + " Installing..."
	} else if len(m.list.Items()) > 0 {
		s = m.list.View()
	}
	return docStyle.Render(s)
}

// search searches the go packages and returns a tea.Msg so that the search model
// can update the TUI.
func search(term string) tea.Cmd {
	return func() tea.Msg {
		results := util.Search(term)
		msg := searchMsg{
			results: make([]list.Item, len(results)),
		}
		for i, res := range results {
			msg.results[i] = SearchResult{
				title:       res,
				description: "",
			}
		}
		return msg
	}
}

type installMsg struct {
	Err error
}

func install(term string) tea.Cmd {
	return func() tea.Msg {
		pkg := util.GetPkgUrl(term)
		err := util.RunGoGet(pkg)
		return installMsg{Err: err}
	}
}
