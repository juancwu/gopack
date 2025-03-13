package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// version is set during build using ldflags
var version = "dev"

const (
	repoOwner      = "juancwu"
	repoName       = "gopack"
	githubAPI      = "https://api.github.com/repos/%s/%s/releases/latest"
)

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func update() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update gopack to the latest version",
		Long:  "Check GitHub for a new version of gopack and update if available",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Current version:", version)
			fmt.Println("Checking for updates...")

			// Get latest release info from GitHub
			apiURL := fmt.Sprintf(githubAPI, repoOwner, repoName)
			resp, err := http.Get(apiURL)
			if err != nil {
				return fmt.Errorf("error checking for updates: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("error checking for updates: HTTP %d", resp.StatusCode)
			}

			var release Release
			if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
				return fmt.Errorf("error parsing GitHub response: %v", err)
			}

			latestVersion := release.TagName
			fmt.Println("Latest version:", latestVersion)

			// Compare versions
			if latestVersion == version {
				fmt.Println("You already have the latest version!")
				return nil
			}

			// Find the appropriate asset for the current OS/arch
			osName := runtime.GOOS
			archName := runtime.GOARCH
			
			var downloadURL string
			assetName := fmt.Sprintf("gopack_%s_%s", osName, archName)
			
			for _, asset := range release.Assets {
				if strings.Contains(asset.Name, assetName) {
					downloadURL = asset.BrowserDownloadURL
					break
				}
			}

			if downloadURL == "" {
				return fmt.Errorf("no suitable binary found for %s/%s", osName, archName)
			}

			// Confirm update
			fmt.Printf("Do you want to update from %s to %s? [y/N] ", version, latestVersion)
			var confirmation string
			fmt.Scanln(&confirmation)
			
			if strings.ToLower(confirmation) != "y" {
				fmt.Println("Update cancelled.")
				return nil
			}

			// Download the new binary
			fmt.Println("Downloading update...")
			
			// Get executable path
			execPath, err := os.Executable()
			if err != nil {
				return fmt.Errorf("error finding executable path: %v", err)
			}

			// Download to a temporary file
			tmpFile := execPath + ".new"
			out, err := os.Create(tmpFile)
			if err != nil {
				return fmt.Errorf("error creating temporary file: %v", err)
			}
			defer out.Close()

			resp, err = http.Get(downloadURL)
			if err != nil {
				return fmt.Errorf("error downloading update: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("error downloading update: HTTP %d", resp.StatusCode)
			}

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				return fmt.Errorf("error saving update: %v", err)
			}
			out.Close()

			// Make the new binary executable
			err = os.Chmod(tmpFile, 0755)
			if err != nil {
				return fmt.Errorf("error setting permissions: %v", err)
			}

			// Replace the old binary
			if runtime.GOOS == "windows" {
				// On Windows, we can't replace a running executable
				// So we create a batch file to do it after we exit
				batchFile := execPath + ".bat"
				batch := fmt.Sprintf(`@echo off
ping -n 3 127.0.0.1 > nul
move /y "%s" "%s"
del "%s"
`, tmpFile, execPath, batchFile)
				
				err = os.WriteFile(batchFile, []byte(batch), 0755)
				if err != nil {
					return fmt.Errorf("error creating update script: %v", err)
				}
				
				cmd := exec.Command("cmd", "/c", "start", "/b", batchFile)
				err = cmd.Start()
				if err != nil {
					return fmt.Errorf("error launching update script: %v", err)
				}
				
				fmt.Println("Update downloaded. It will be installed when you close this program.")
			} else {
				// On Unix-like systems, we can replace the binary directly
				err = os.Rename(tmpFile, execPath)
				if err != nil {
					return fmt.Errorf("error installing update: %v", err)
				}
				
				fmt.Println("Update installed successfully!")
			}
			
			return nil
		},
	}

	return updateCmd
}