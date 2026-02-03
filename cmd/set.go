/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package cmd

import (
	"fmt"
	"strings"

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	setGlobal bool
	setLocal  bool
)

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value using dot notation for nested keys.

Local configuration (default):
  - Modify bot settings like name, description, prefix, and author
  - Changes are saved to the project's botbox.conf file

Global configuration (use -g flag):
  - Update CLI preferences and default values
  - Configure update behavior and development tools
  - Set user information for new projects

Boolean values accept: true/false, t/f, yes/no, y/n, 1/0

Examples:
  botbox config set bot.name "My Awesome Bot"           # Set local bot name
  botbox config set -g cli.auto_update true            # Enable auto-updates
  botbox config set -l bot.command_prefix "!"          # Set local prefix`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSet,
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	valueStr := args[1]

	if setGlobal {
		return handleGlobalConfigSet(key, valueStr)
	} else {
		return handleLocalConfigSet(key, valueStr)
	}
}

func handleGlobalConfigSet(key, valueStr string) error {
	exists, err := utils.GlobalConfigExists()
	if err != nil {
		return fmt.Errorf("failed to check global config: %w", err)
	}
	if !exists {
		return fmt.Errorf("global config does not exist. Run 'botbox config init' first")
	}

	if !isValidGlobalConfigKey(key) {
		return fmt.Errorf("invalid global config key: %s", key)
	}

	value, err := parseValue(valueStr, key)
	if err != nil {
		return fmt.Errorf("invalid value for key %s: %w", key, err)
	}

	if err := utils.SetGlobalConfigValue(key, value); err != nil {
		return fmt.Errorf("failed to set global config value: %w", err)
	}

	fmt.Printf("✅ Set global %s = %v\n", key, value)
	return nil
}

func handleLocalConfigSet(key, valueStr string) error {
	_, err := utils.FindBotConf()
	if err != nil {
		return fmt.Errorf("not in a botbox project: %w", err)
	}

	if !isValidLocalConfigKey(key) {
		return fmt.Errorf("invalid local config key: %s", key)
	}

	if err := utils.SetLocalConfigValue(key, valueStr); err != nil {
		return fmt.Errorf("failed to set local config value: %w", err)
	}

	fmt.Printf("✅ Set local %s = %s\n", key, valueStr)
	return nil
}

func isValidGlobalConfigKey(key string) bool {
	validKeys := map[string]bool{
		"cli.check_updates": true,
		"cli.auto_update":   true,

		"user.default_user":    true,
		"user.github_username": true,

		"display.scroll_enabled": true,
		"display.color_scheme":   true,

		"defaults.command_prefix": true,
		"defaults.auto_git_init":  true,

		"dev.editor": true,
	}

	return validKeys[key]
}

func isValidLocalConfigKey(key string) bool {
	validKeys := map[string]bool{
		"bot.name":           true,
		"bot.description":    true,
		"bot.command_prefix": true,
		"bot.author":         true,
	}

	return validKeys[key]
}

func parseValue(valueStr, key string) (any, error) {
	boolKeys := map[string]bool{
		"cli.check_updates":      true,
		"cli.auto_update":        true,
		"display.scroll_enabled": true,
		"defaults.auto_git_init": true,
	}

	if boolKeys[key] {
		switch strings.ToLower(valueStr) {
		case "true", "t", "yes", "y", "1":
			return true, nil
		case "false", "f", "no", "n", "0":
			return false, nil
		default:
			return nil, fmt.Errorf("expected boolean value (true/false), got: %s", valueStr)
		}
	}

	return valueStr, nil
}

func init() {
	configCmd.AddCommand(setCmd)

	setCmd.Flags().BoolVarP(&setGlobal, "global", "g", false, "Set global configuration")
	setCmd.Flags().BoolVarP(&setLocal, "local", "l", false, "Set local configuration")
	setCmd.MarkFlagsMutuallyExclusive("global", "local")
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
