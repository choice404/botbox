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

var rootCmd = &cobra.Command{
	Use:   "botbox",
	Short: "A powerful CLI tool for Discord bot development",
	Long: `BotBox is a comprehensive CLI tool that helps you scaffold, configure, and manage 
Discord bot projects quickly and efficiently. It provides automated cog generation, 
project initialization, configuration management, and built-in utilities to streamline 
your Discord bot development workflow.

Built with a cog-based architecture for modularity and featuring automatic updates, 
global configuration management, and seamless project upgrades.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "update" {
			checkForUpdatesIfEnabled()
		}
	},
}

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
