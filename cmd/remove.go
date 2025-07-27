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

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	removeCogName string
	cogRemove     utils.CogConfig
)

var removeCmd = &cobra.Command{
	Use:   "remove [cog-name]",
	Short: "Remove a cog from the current Bot Box project",
	Long: `Remove an existing cog from your Bot Box project.

This command will:
  - Display a list of available cogs to choose from
  - Remove the selected cog file from the src/cogs/ directory
  - Update the botbox.conf configuration to remove the cog entry
  - Maintain project consistency by cleaning up all references

You can specify the cog name as an argument or select from an interactive list. 
The command ensures safe removal without breaking your project configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := utils.FindBotConf()
		if err != nil {
			fmt.Println("Current directory is not in a botbox project.")
			return
		}

		if len(args) > 0 {
			removeCogName = args[0]
		} else {
			removeCogName = ""
		}
		model := utils.RemoveModel(removeCallback, removeInitCallback)
		utils.CupSleeve(model)
	},
}

func removeCallback(model *utils.Model) []error {
	var errors []error
	values := model.ModelValues
	filename := values.Map["cogName"]

	rootDir, err := utils.FindBotConf()
	if err != nil {
		errors = append(errors, fmt.Errorf("error finding root directory: %w", err))
		return errors
	}

	configPath := filepath.Join(rootDir, "botbox.conf")

	config, err := utils.LoadConfig()
	if err != nil {
		errors = append(errors, fmt.Errorf("error loading config: %w", err))
		return errors
	}

	for i, cog := range config.Cogs {
		if cog.Name == *filename {
			cogRemove = cog
			config.Cogs = slices.Delete(config.Cogs, i, i+1)
		}
	}

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		errors = append(errors, fmt.Errorf("error marshalling config to JSON: %w", err))
		return errors
	}

	err = os.Remove(rootDir + "/src/cogs/" + cogRemove.File + ".py")

	err = os.WriteFile(configPath, jsonData, 0644)
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to write updated botbox.conf: %w", err))
		return errors
	}
	return nil
}

func removeInitCallback(model *utils.Model, allFormsModels []utils.Values) {
	var errors []error
	modelValues := model.ModelValues
	if removeCogName != "" {
		*modelValues.Map["cogName"] = removeCogName
		configs, err := utils.LoadConfig()
		if err != nil {
			errors = append(errors, fmt.Errorf("error loading configuration: %w", err))
			model.HandleError(errors)
			return
		}
		if len(configs.Cogs) == 0 {
			errors = append(errors, fmt.Errorf("error: no cogs available to remove"))
			model.HandleError(errors)
			return
		}
		cogExists := false
		for _, cog := range configs.Cogs {
			if cog.Name == removeCogName {
				cogExists = true
				break
			}
		}
		if !cogExists {
			errors = append(errors, fmt.Errorf("cog '%s' does not exist in the project", removeCogName))
			model.HandleError(errors)
			return
		}
	}
}

func init() {
	rootCmd.AddCommand(removeCmd)
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
