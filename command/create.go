package command

import (
	"fmt"

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
            fm, err := model.Run()
            if err != nil {
                return err
            }
            _, ok := fm.(*tui.CreateModel)
            if !ok {
                return fmt.Errorf("unexpected model type returned")
            }
            
            return nil
        },
    }
    return createCmd
}