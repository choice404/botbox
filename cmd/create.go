/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new Bot Box project",
	Long: `Creates a directory containing a new Bot Box project.
  The project directory will contain the following file structure:

  projectName/
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
        |- cogs.py
  `,
	Run: func(cmd *cobra.Command, args []string) {
		var createNewProject bool
		botBoxExistsForm := huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("The current directory is in an existing Bot Box project."),
				huh.NewConfirm().
					Title("Would you like to create a new Bot Box project?").
					Affirmative("yes").
					Negative("no").
					Value(&createNewProject),
			),
		)

		_, err := utils.FindBotConf()
		if err == nil {
			botBoxExistsForm.Run()
			if !createNewProject {
				return
			}
		}
		model := utils.CreateModel(createProjectCallback)
		utils.CupSleeve(model)
	},
}

func createProjectCallback(model *utils.Model) []error {
	var errors []error
	values := model.ModelValues
	rootDir := values.Map["botName"]

	if !filepath.IsAbs(*rootDir) {
		cwd, err := os.Getwd()
		if err != nil {
			errors = append(errors, fmt.Errorf("error getting current directory: %w", err))
			return errors
		}
		*rootDir = filepath.Join(cwd, *rootDir)
	}

	if _, err := os.Stat(*rootDir); err == nil && !os.IsNotExist(err) {
		errors = append(errors, fmt.Errorf("directory already exists: %s", *rootDir))
		return errors
	} else if os.IsNotExist(err) {
		err = os.MkdirAll(*rootDir, os.ModePerm)
		if err != nil {
			errors = append(errors, fmt.Errorf("error creating directory: %w", err))
			return errors
		}
	} else {
		errors = append(errors, fmt.Errorf("error checking directory %s: %w", *rootDir, err))
		return errors
	}

	utils.CreateProject(*rootDir, values)

	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)
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
