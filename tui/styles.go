package tui

import "github.com/charmbracelet/lipgloss"

var (
	wrapper = lipgloss.NewStyle().Margin(1)

	okText   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	errText  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)
