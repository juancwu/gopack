package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juancwu/gopack/util"
)

type packageItem struct {
	pkg util.Package
}

type listModel struct {
	List list.Model
}

func (i packageItem) Title() string { return i.pkg.Path }
func (i packageItem) Description() string {
	version := i.pkg.Version
	directory := i.pkg.Dir
	if version == "" {
		version = "Unknown"
	}
	if directory == "" {
		directory = "Unknown"
	}
	return fmt.Sprintf("Version %s, Directory: %s", version, directory)
}
func (i packageItem) FilterValue() string { return i.pkg.Path }

func NewListModel(packages []util.Package) listModel {
	items := make([]list.Item, len(packages))
	for i, pkg := range packages {
		items[i] = packageItem{pkg: pkg}
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
