package command

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juancwu/gopack/tui"
	"github.com/juancwu/gopack/util"
	"github.com/spf13/cobra"
)

// install will get all the dependencies in an existing go.mod file
func install() *cobra.Command {
	installCmd := &cobra.Command{
		Use:     "install",
		Short:   "Installs all dependencies defined in go.mod",
		Long:    "Fastest way to get all dependencies in a new machine when a go.mod is present",
		Aliases: []string{"i"},
		RunE: func(cmd *cobra.Command, args []string) error {
			modules, err := util.ParseGoMod()
			if err != nil {
				return err
			}

			m := tui.NewGoModModel(modules)
			p := tea.NewProgram(m)
			_, err = p.Run()
			return err
		},
	}
	return installCmd
}
