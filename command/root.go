package command

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	cfgFile    = ".gonttrc"
	cfgType    = "yaml"
	timeFormat = "20060102150405"
	timezone   = "UTC"

	dirFlagUsage = "This is the destination the migrations are located. This should be relative to the CWD or an absolute path"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "gopack",
		Short: "A simple go package installer",
	}

	rootCmd.AddCommand(install())

	return rootCmd.ExecuteContext(context.Background())
}
