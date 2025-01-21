package tui

import "github.com/charmbracelet/lipgloss"

var (
	wrapper = lipgloss.NewStyle().Margin(1)

	okText  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	errText = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
<<<<<<< HEAD
	docStyle = lipgloss.NewStyle().Margin(1, 2)
=======

	titleStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FF75B7")).
        Bold(true).
        MarginLeft(2)

	inputStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FF75B7")).
        MarginLeft(2)
>>>>>>> 286c42d (Create command added)
)
