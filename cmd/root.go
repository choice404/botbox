/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	GlobalConfig *utils.GlobalConfig
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "botbox",
	Short: "A CLI tool to help create and manage a cog based discord.py bot",
	Long: ` A discord bot template generator to help create discord
bots quickly and easily. Forget about the boilerplate
and focus on what really matters, what your bot will do.

Bot Box is built using Golang, Cobra, Bubble Tea, and Huh,
offering an intuitive cli tool to quickly build
Discord bot projects. It includes a cog-based architecture, ` +
		"`.env`" + ` management, and built-in utilities for automating
bot configuration and extension development.
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "update" {
			checkForUpdatesIfEnabled()
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	utils.SetVersion(rootCmd.Version)
	rootCmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)

	if err := utils.SyncGlobalConfigVersion(); err != nil {
		fmt.Printf("⚠️  Warning: failed to sync version: %v\n", err)
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func checkForUpdatesIfEnabled() {
	if !utils.ShouldCheckForUpdates() {
		return
	}

	latest, hasUpdate, err := utils.CheckForUpdates()
	if err != nil {
		return
	}

	if hasUpdate {
		fmt.Printf("📦 Update available: %s → %s\n", utils.Version, latest.TagName)
		fmt.Printf("Run 'botbox update' to update, or 'botbox config set -g cli.check_updates false' to disable notifications.\n\n")

		if utils.ShouldAutoUpdate() {
			fmt.Println("🔄 Auto-update is enabled, updating now...")
			if err := utils.UpdateBotBox(latest.TagName); err != nil {
				fmt.Printf("❌ Auto-update failed: %v\n", err)
				fmt.Println("Please run 'botbox update' to update.")
				fmt.Println("If you want to disable auto-updates, run 'botbox config set -g cli.auto_update false'.")
			} else {
				fmt.Println("✅ Auto-update completed! Please restart your terminal.")
			}
		}
	}
}

func init() {
	exists, err := utils.GlobalConfigExists()
	if err != nil {
		fmt.Printf("Error checking config: %v\n", err)
		os.Exit(1)
	}

	if !exists {
		fmt.Println("Config file does not exist. Creating...")
		if err := utils.CreateGlobalConfig(); err != nil {
			fmt.Printf("Error creating config: %v\n", err)
			os.Exit(1)
		}
	}

	GlobalConfig, err = utils.LoadGlobalConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
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
