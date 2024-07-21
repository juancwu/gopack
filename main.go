package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"

	"github.com/juancwu/gopack/command"
	"github.com/juancwu/gopack/tui"
)

func main() {
	log.SetReportCaller(false)
	log.SetReportTimestamp(false)
	if len(os.Args) == 1 {
		m := tui.NewModel()
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	} else {
		err := command.Execute()
		if err != nil {
			log.Fatal(err)
		}
	}
}
