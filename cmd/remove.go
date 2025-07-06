/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	cogList       []string
	cogRemoveName string
	cogRemove     CogConfig
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a cog from the project",
	Long: ` This command allows you to remove a specific cog from your bot project.
You can specify the cog to remove by providing its name as an argument or select it from a list of available cogs.
  `,
	Run: func(cmd *cobra.Command, args []string) {

		rootDir, err := FindBotConf()
		if err != nil {
			return
		}

		configPath := filepath.Join(rootDir, "botbox.conf")

		config, err := LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}
		for _, cog := range config.Cogs {
			cogList = append(cogList, cog.Name)
		}

		Banner()
		filename := ""
		if len(args) > 0 {
			filename = args[0]
		}

		if filename == "" {
			cmdRemoveForm := generateCommandRemoveForms()

			err = cmdRemoveForm.Run()
			if err != nil {
				fmt.Println("Error running form:", err)
				return
			}

			filename = cogName
		}

		for i, cog := range config.Cogs {
			if cog.Name == filename {
				cogRemove = cog
				config.Cogs = slices.Delete(config.Cogs, i, i+1)
			}
		}

		jsonData, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			fmt.Println("failed to marshal config to JSON: %w", err)
			return
		}

		err = os.Remove(rootDir + "/src/cogs/" + cogRemove.File + ".py")

		err = os.WriteFile(configPath, jsonData, 0644)
		if err != nil {
			fmt.Println("failed to write updated botbox.conf: %w", err)
			return
		}

	},
}

func generateCommandRemoveForms() *huh.Form {
	cmdRemoveForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(&cogName).
				Height(8).
				Title("Select a cog to remove").
				OptionsFunc(func() []huh.Option[string] {
					return huh.NewOptions(cogList...)
				}, &cogName),
		),
	)

	return cmdRemoveForm
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
