package command

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func install() *cobra.Command {
	installCmd := &cobra.Command{
		Use:     "install",
		Short:   "Search and install first in result",
		Long:    "Search and install first in query result with a confirmation. There is a chance to look all results.",
		Example: "gopack install PKG_NAME",
		Aliases: []string{"i"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Infof("Searching package: %s", args[0])
			return nil
		},
	}

	return installCmd
}
