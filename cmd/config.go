/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package cmd

import (
	"fmt"

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	configName        bool
	configDir         bool
	configAuthor      bool
	configDescription bool
	configPrefix      bool
	configCogs        bool
	configIsOptional  bool

	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Get config for the project",
		Long: `Create different types of files for the discord bot.
Main files
Cogs
Config`,
		Run: func(cmd *cobra.Command, args []string) {

			if configName || configDir || configAuthor || configDescription || configPrefix || configCogs {
				configIsOptional = true
			} else {
				configIsOptional = false
			}

			_, err := utils.FindBotConf()
			if err != nil {
				fmt.Println("Current directory is not in a botbox project.")
				return
			}

			model := utils.ConfigModel(configCallback, configInitCallback)

			utils.CupSleeve(&model)
		},
	}
)

func configInitCallback(modelValues map[string]*string, allFormsModels []map[string]*string) {

	config, err := utils.LoadConfig()
	rootDir, err := utils.FindBotConf()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	if configName || !configIsOptional {
		modelValues["botName"] = &config.BotInfo.Name
	}
	if configDir || !configIsOptional {
		modelValues["rootDir"] = &rootDir
	}
	if configAuthor || !configIsOptional {
		modelValues["botAuthor"] = &config.BotInfo.Author
	}
	if configDescription || !configIsOptional {
		modelValues["botDescription"] = &config.BotInfo.Description
	}
	if configPrefix || !configIsOptional {
		modelValues["botPrefix"] = &config.BotInfo.CommandPrefix
	}
	if configCogs || !configIsOptional {
		var cogConfigs string
		if len(config.Cogs) > 0 {
			cogConfigs, err = utils.CogConfigSliceToJSON(config.Cogs)
			if err != nil {
				fmt.Println("Error converting cogs to JSON:", err)
				return
			}
		} else {
			cogConfigs = "[]" // Empty JSON array
		}
		modelValues["cogs"] = &cogConfigs
	}
}

func configCallback(values map[string]string) {
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
