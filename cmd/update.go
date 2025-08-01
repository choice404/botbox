/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package cmd

import (
	"fmt"

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Bot Box to the latest version",
	Long: `Update Bot Box CLI to the latest version from GitHub releases.

This command will:
  1. Check GitHub for the latest Bot Box release
  2. Compare with your current version
  3. Download and install the update using 'go install'
  4. Clean the module cache for a fresh installation
  5. Display release notes and changelog information

The update process is automatic and handles all necessary cleanup. 
After updating, you may need to restart your terminal or run 'hash -r' 
to refresh the command cache.

Bot Box can also check for updates automatically based on your global 
configuration settings (cli.check_updates and cli.auto_update).`,
	RunE: runUpdate,
}

func runUpdate(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Checking for updates...")

	latest, hasUpdate, err := utils.CheckForUpdates()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !hasUpdate {
		fmt.Printf("✅ You're already running the latest version (%s)\n", utils.Version)
		return nil
	}

	fmt.Printf("📦 Update available: %s → %s\n", utils.Version, latest.TagName)
	fmt.Printf("🔗 Release notes: %s\n", latest.HTMLURL)

	if err := utils.UpdateBotBox(latest.TagName); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Println("🎉 Update completed! Please restart your terminal or run 'hash -r' to refresh the command cache.")

	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

/*
Copyright © 2025 Austin "Choice404" Choi

https://github.com/choice404/botbox

Bot Box
A discord bot template generator to help create discord bots quickly and easily

This code is licensed under the MIT License.

MIT License

Copyright (c) 2025 Austin Choi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
