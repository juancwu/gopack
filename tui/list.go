package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Package struct {
    Path     string `json:"Path"`
    Version  string `json:"Version"`
    Dir      string `json:"Dir"`
}

func (i Package) Title() string       { return i.Path }
func (i Package) Description() string { 
	version := i.Version
	if version == ""{
		version = "Unknown"
	}
	return fmt.Sprintf("Version %s, Directory: %s", version, i.Dir)
}
func (i Package) FilterValue() string { return i.Path }

type listModel struct {
	List list.Model
}

func NewListModel(packages []Package) listModel {
    items := make([]list.Item, len(packages))
    for i, pkg := range packages {
        items[i] = pkg
    }

    l := list.New(items, list.NewDefaultDelegate(), 0, 0)
    l.Title = "Installed Packages"

    return listModel{
        List: l,
    }
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	return docStyle.Render(m.List.View())
}