package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

type Package struct {
    Path     string `json:"Path"`
    Version  string `json:"Version"`
    Main     bool   `json:"Main"`
    Dir      string `json:"Dir"`
    Indirect bool   `json:"Indirect"`
}

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
            var packages []Package

            for decoder.More() {
                var pkg Package
                if err := decoder.Decode(&pkg); err != nil {
                    return fmt.Errorf("error parsing JSON: %v", err)
                }
                packages = append(packages, pkg)
            }

            fmt.Println("Installed packages:")
            fmt.Println("==================")
            
            for _, pkg := range packages {
                version := pkg.Version
                if version == "" {
                    version = "No version information"
                }
                
                fmt.Printf("Package: %s\n", pkg.Path)
                fmt.Printf("Version: %s\n", version)
                if pkg.Dir != "" {
                    fmt.Printf("Directory: %s\n", pkg.Dir)
                }
                if pkg.Indirect {
                    fmt.Println("Type: Indirect dependency")
                }
                fmt.Println("------------------")
            }
            
            fmt.Printf("\nTotal packages: %d\n", len(packages))
            return nil
        },
    }
    return listCmd
}