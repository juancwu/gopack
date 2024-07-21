package command

import (
	"github.com/charmbracelet/log"
	"github.com/juancwu/gopack/util"
	"github.com/spf13/cobra"
)

func install() *cobra.Command {
	installCmd := &cobra.Command{
		Use:     "install",
		Short:   "Search and install first in result",
		Long:    "Search and install first in query result with a confirmation. There is a chance to look all results.",
		Example: "gopack install PKG_NAME",
		Aliases: []string{"i"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, pkgName := range args {
				log.Warnf("Installing first match of '%s' ignoring namespace.", pkgName)
				matches := util.Search(pkgName)
				if len(matches) < 1 {
					log.Errorf("Failed to search package: %s", pkgName)
					continue
				}
				log.Infof("Found '%s', installing...", matches[0])
				url := util.GetPkgUrl(matches[0])
				if err := util.RunGoGet(url); err != nil {
					log.Errorf("Failed to run 'go get': %v", err)
					return err
				}
			}
			return nil
		},
	}

	return installCmd
}
