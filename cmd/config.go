/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	configName        bool
	configDir         bool
	configAuthor      bool
	configDescription bool
	configPrefix      bool
	configCogs        bool

	configCmd = &cobra.Command{
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
			if configName || configDir || configAuthor || configDescription || configPrefix || configCogs {
				displayPartialConfig(rootDir)
			} else {
				displayConfig(rootDir)
			}
		},
	}
)

func displayConfig(rootDir string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	fmt.Println("Bot Box project configuration:")
	fmt.Println("Bot Name:", config.BotInfo.Name)
	fmt.Println("Root Directory:", rootDir)
	fmt.Println("Author:", config.BotInfo.Author)
	fmt.Println("Description:", config.BotInfo.Description)
	fmt.Println("Command Prefix:", config.BotInfo.CommandPrefix)
	fmt.Println("Cogs:")
	for _, cog := range config.Cogs {
		fmt.Printf(" - %s(%s)\n", cog.File, cog.Name)
		for _, command := range cog.Commands {
			fmt.Printf("   - Command: %s\n", command)
		}
	}
}

func displayPartialConfig(rootDir string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	fmt.Println("Bot Box project configuration:")
	if configName {
		fmt.Println("Bot Name:", config.BotInfo.Name)
	}
	if configDir {
		fmt.Println("Root Directory:", rootDir)
	}
	if configAuthor {
		fmt.Println("Author:", config.BotInfo.Author)
	}
	if configDescription {
		fmt.Println("Description:", config.BotInfo.Description)
	}
	if configPrefix {
		fmt.Println("Command Prefix:", config.BotInfo.CommandPrefix)
	}
	if configCogs {
		fmt.Println("Cogs:")
		for _, cog := range config.Cogs {
			fmt.Printf(" - %s(%s)\n", cog.File, cog.Name)
			for _, command := range cog.Commands {
				fmt.Printf("   - Command: %s\n", command)
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVarP(&configName, "name", "n", false, "Display the bot name")
	configCmd.Flags().BoolVarP(&configDir, "dir", "d", false, "Display the bot directory")
	configCmd.Flags().BoolVarP(&configAuthor, "author", "a", false, "Display the bot author")
	configCmd.Flags().BoolVarP(&configDescription, "description", "D", false, "Display the bot description")
	configCmd.Flags().BoolVarP(&configPrefix, "prefix", "p", false, "Display the command prefix")
	configCmd.Flags().BoolVarP(&configCogs, "cogs", "c", false, "Display the cogs")
	configCmd.Flags().BoolP("help", "h", false, "Help for config")
}

/*
Copyright © 2025 2025 Austin "Choice404" Choi

https://github.com/choice404/botbox

Bot Box
A discord bot template generator to help create discord bots quickly and easily

This code is licensed under the MIT License.
Please see the LICENSE file in the root directory of this project for the full license details.
*/
