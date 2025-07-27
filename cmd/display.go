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

// displayCmd represents the display command
var displayCmd = &cobra.Command{
	Use:   "display [sections...]",
	Short: "Display configuration",
	Long:  `Display local project configuration or global BotBox CLI configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		var configModel utils.Model
		configArgs = args
		if len(args) == 0 || args[0] == "all" {
			allConfigs = true
		} else {
			allConfigs = false
		}

		globalFlag, _ := cmd.Flags().GetBool("global")
		if globalFlag {
			configModel = utils.GlobalConfigModel(configCallback, globalConfigInitCallback)
		} else {
			_, err := utils.FindBotConf()
			if err != nil {
				fmt.Println("Current directory is not in a botbox project.")
				return
			}
			configModel = utils.LocalConfigModel(configCallback, localConfigInitCallback)
		}
		utils.CupSleeve(configModel)
	},
}

func init() {
	configCmd.AddCommand(displayCmd)

	displayCmd.Flags().BoolP("global", "g", false, "Show global CLI configuration")
	displayCmd.Flags().Bool("local", false, "Show local project configuration (default)")

	displayCmd.Flags().StringP("format", "f", "default", "Output format (default, json, yaml)")
	displayCmd.Flags().BoolP("keys-only", "k", false, "Show only configuration keys")

	displayCmd.MarkFlagsMutuallyExclusive("global", "local")

	displayCmd.Flags().BoolP("help", "h", false, "Help for config")
}
