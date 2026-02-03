/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package cmd

import (
	"fmt"
	"sort"

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	listGlobal bool
	listLocal  bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration keys and their values",
	Long: `Display all available configuration keys and their current values.

Local configuration (default):
  Shows all bot settings, project information, and cog configurations
  for the current Bot Box project.

Global configuration (use -g flag):
  Shows all CLI settings, user preferences, default values, and 
  development tool configurations.

This command provides a comprehensive overview of all configurable 
options and their current state, useful for debugging configuration 
issues or understanding available settings.`,
	RunE: runConfigList,
}

func runConfigList(cmd *cobra.Command, args []string) error {
	if listGlobal {
		return handleGlobalConfigList()
	} else {
		return handleLocalConfigList()
	}
}

func handleGlobalConfigList() error {
	exists, err := utils.GlobalConfigExists()
	if err != nil {
		return fmt.Errorf("failed to check global config: %w", err)
	}
	if !exists {
		return fmt.Errorf("global config does not exist. Run 'botbox config init' first")
	}

	keys := []string{
		"cli.check_updates", "cli.auto_update",
		"user.default_user", "user.github_username",
		"display.scroll_enabled", "display.color_scheme",
		"defaults.command_prefix", "defaults.auto_git_init",
		"dev.editor",
	}

	sort.Strings(keys)

	fmt.Println("Global Configuration:")
	for _, key := range keys {
		value := utils.GetGlobalConfigValue(key)
		if value == nil {
			fmt.Printf("  %s = <not set>\n", key)
		} else {
			fmt.Printf("  %s = %v\n", key, value)
		}
	}

	return nil
}

func handleLocalConfigList() error {
	_, err := utils.FindBotConf()
	if err != nil {
		return fmt.Errorf("not in a botbox project: %w", err)
	}

	keys := []string{
		"bot.name", "bot.description", "bot.command_prefix", "bot.author",
	}

	fmt.Println("Local Configuration:")
	for _, key := range keys {
		value, err := utils.GetLocalConfigValue(key)
		if err != nil {
			fmt.Printf("  %s = <error: %v>\n", key, err)
		} else {
			fmt.Printf("  %s = %v\n", key, value)
		}
	}

	return nil
}

func init() {
	configCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listGlobal, "global", "g", false, "List global configuration")
	listCmd.Flags().BoolVarP(&listLocal, "local", "l", false, "List local configuration")
	listCmd.MarkFlagsMutuallyExclusive("global", "local")
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
