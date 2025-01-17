package tui

import "github.com/charmbracelet/lipgloss"

var (
	docStyle = lipgloss.NewStyle().Margin(1)

	errText = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
)
