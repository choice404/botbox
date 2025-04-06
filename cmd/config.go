/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Create files for the project",
	Long: `Create different types of files for the discord bot.
  Main files
  Cogs
  Config`,
	Run: func(cmd *cobra.Command, args []string) {
		rootDir, err := FindBotConf()
		if err != nil {
			fmt.Println("Current directory is not in a botbox project.")
			return
		}
		displayConfig(rootDir)
	},
}

func displayConfig(rootDir string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}
	fmt.Println("Bot Box project configuration:")
	fmt.Println("Root Directory:", rootDir)
	fmt.Println("Bot Name:", config.BotInfo.Name)
	fmt.Println("Command Prefix:", config.BotInfo.CommandPrefix)
	fmt.Println("Author:", config.BotInfo.Author)
	fmt.Println("Description:", config.BotInfo.Description)
	fmt.Println("Cogs:")
	for _, cog := range config.Cogs {
		fmt.Printf(" - %s(%s)\n", cog.File, cog.Name)
		for _, command := range cog.Commands {
			fmt.Printf("   - Command: %s\n", command)
		}
	}
}

func init() {
	rootCmd.AddCommand(configCmd)
}

/*
Copyright © 2025 2025 Austin "Choice404" Choi

https://github.com/choice404/botbox

Bot Box
A discord bot template generator to help create discord bots quickly and easily

This code is licensed under the MIT License.
Please see the LICENSE file in the root directory of this project for the full license details.
*/
