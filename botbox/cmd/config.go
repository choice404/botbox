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
		fmt.Println("botbox config")
	},
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
