package command

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juancwu/gopack/tui"
	"github.com/spf13/cobra"
)

func get() *cobra.Command {
	var selectResult bool
	getCmd := &cobra.Command{
		Use:     "get",
		Short:   "Search and install first in result",
		Long:    "Search and install first in query result with a confirmation. There is a chance to look all results.",
		Example: "gopack install PKG_NAME",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m := tui.NewInstallModel(args, selectResult)
			p := tea.NewProgram(m)
			_, err := p.Run()
			return err
		},
	}

	getCmd.Flags().BoolVarP(&selectResult, "select", "s", false, "Show list of results and allow manual selection")

	return getCmd
}
