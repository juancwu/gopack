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
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "gop",
		Short: "A simple go package installer",
	}

	rootCmd.AddCommand(get())
	rootCmd.AddCommand(list())

	return rootCmd.ExecuteContext(context.Background())
}
