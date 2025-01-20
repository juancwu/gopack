package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type searchModel struct {
	ti textinput.Model
}

func NewSearchModel() searchModel {
	ti := textinput.New()
	ti.Placeholder = "Search Packages"

	return searchModel{
		ti: ti,
	}
}

func (m searchModel) Init() tea.Cmd {
	return nil
}

func (m searchModel) Update(mst tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m searchModel) View() string {
	return ""
}
