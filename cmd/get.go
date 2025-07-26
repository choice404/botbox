/*
Copyright © 2025 Austin "Choice404" Choi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	getGlobal bool
	getLocal  bool
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Long: `Get a configuration value using dot notation for nested keys.

Examples:
  botbox config get bot.name                    # local (default)
  botbox config get -l bot.author               # local (explicit)
  botbox config get -g user.default_user       # global`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigGet,
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	if getGlobal {
		return handleGlobalConfigGet(key)
	} else {
		return handleLocalConfigGet(key)
	}
}

func handleGlobalConfigGet(key string) error {
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

	value := utils.GetGlobalConfigValue(key)
	if value == nil {
		fmt.Printf("global %s = <not set>\n", key)
	} else {
		fmt.Printf("global %s = %v\n", key, value)
	}

	return nil
}

func handleLocalConfigGet(key string) error {
	_, err := utils.FindBotConf()
	if err != nil {
		return fmt.Errorf("not in a botbox project: %w", err)
	}

	if !isValidLocalConfigKey(key) {
		return fmt.Errorf("invalid local config key: %s", key)
	}

	value, err := utils.GetLocalConfigValue(key)
	if err != nil {
		return fmt.Errorf("failed to get local config value: %w", err)
	}

	fmt.Printf("local %s = %v\n", key, value)
	return nil
}

func init() {
	configCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVarP(&getGlobal, "global", "g", false, "Get global configuration")
	getCmd.Flags().BoolVarP(&getLocal, "local", "l", false, "Get local configuration")
	getCmd.MarkFlagsMutuallyExclusive("global", "local")
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
