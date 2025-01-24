package command

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/juancwu/gopack/tui"
	"github.com/juancwu/gopack/util"
	"github.com/spf13/cobra"
)

func list() *cobra.Command {
    listCmd := &cobra.Command{
        Use:     "list",
        Short:   "List all the packages that was installed and used",
        Long:    "List all the packages that was installed and used. And also going to show the path they they are installed and the version.",
        Example: "gopack list",
        RunE: func(cmd *cobra.Command, args []string) error {
            packages, err := util.GetDependencyList()
            if err != nil {
                fmt.Println("Error getting dependency list: ", err)
            }
            
            m := tui.NewListModel(packages)
            m.List.Title = "Installed Packages"

            p := tea.NewProgram(m, tea.WithAltScreen())
            if _, err := p.Run(); err != nil {
                fmt.Println("Error running program:", err)
            }
            return nil
        },
    }
    return listCmd
}