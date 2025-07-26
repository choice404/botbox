/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package cmd

import (
	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

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
        |- helloWorld.py (demo cog)
        |- cogs.py
`,
	Run: func(cmd *cobra.Command, args []string) {
		model := utils.CreateModel(CreateProjectInitCallback)
		utils.CupSleeve(model)
	},
}

func CreateProjectInitCallback(model *utils.Model) []error {
	values := model.ModelValues
	// TODO: Update so the CreateProject function does the overwrite form using the custom Tea/Huh manager.
	// Or maybe not after some testing, still gonna figure it out a bit
	utils.CreateProject("./", values)
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}

/*
Copyright © 2025 2025 Austin "Choice404" Choi

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
