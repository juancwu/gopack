package command

import (
	"github.com/juancwu/gopack/tui"
	"github.com/spf13/cobra"
)

func create() *cobra.Command {
    createCmd := &cobra.Command{
        Use:     "create",
        Short:   "Create a new Go project",
        Long:    "Create a new Go project with default directory structure and go.mod file",
        Example: "gopack create",
        RunE: func(cmd *cobra.Command, args []string) error {
            model := tui.NewCreateModel()
            _, err := model.Run()
            if err != nil {
                return err
            }
            
            return nil
        },
    }
    return createCmd
}