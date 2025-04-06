/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package cmd

import (
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a Bot Box project",
	Long: `Initializes a Bot Box project in the current directory and prompts the user for information about the bot as well as setup other default configurations in a botbox.conf file.
  Will also create the initial project strucutre:

  /
  |- README.md
  |- botbox.conf
  |- run.sh
  |- LICENSE (optional)
  |- doppler.yaml (optional)
  |- src/
     |- main.py
     |- cogs/
        |- __init__.py
        |- helloWorld.py
`,
	Run: func(cmd *cobra.Command, args []string) {
		BotBoxCreate(CreateProjectInitWrapper)
	},
}

func CreateProjectInitWrapper() {
	CreateProject("./")
}

func init() {
	rootCmd.AddCommand(initCmd)
}

/*
Copyright © 2025 2025 Austin "Choice404" Choi

Bot Box
A discord bot template generator to help create discord bots quickly and easily

This code is licensed under the MIT License.
Please see the LICENSE file in the root directory of this project for the full license details.
*/
