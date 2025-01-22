package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/juancwu/gopack/tui"
	"github.com/spf13/cobra"
)

func list() *cobra.Command {
    listCmd := &cobra.Command{
        Use:     "list",
        Short:   "List all the packages that was installed and used",
        Long:    "List all the packages that was installed and used. And also going to show the path they they are installed and the version.",
        Example: "gopack list",
        RunE: func(cmd *cobra.Command, args []string) error {
            output, err := exec.Command("go", "list", "-m", "-json", "all").Output()
            if err != nil {
                return fmt.Errorf("error executing command: %v", err)
            }

            decoder := json.NewDecoder(bytes.NewReader(output))
            var packages []tui.Package

            for decoder.More() {
                var pkg tui.Package
                if err := decoder.Decode(&pkg); err != nil {
                    return fmt.Errorf("error parsing JSON: %v", err)
                }
                packages = append(packages, pkg)
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