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

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize cog configuration with actual cog files",
	Long: `Synchronize the botbox.conf file with the actual cog files in your project.

This command analyzes your src/cogs/ directory and updates the configuration to:
  - Add newly created cog files that aren't in the configuration
  - Remove configuration entries for deleted cog files  
  - Update command information by parsing cog file contents
  - Fix inconsistencies between code and configuration

Use this command when you've manually added/removed cogs or when the 
configuration seems out of sync with your actual project structure. 
It ensures your bot will load all available cogs correctly.`,
	Run: func(cmd *cobra.Command, args []string) {
		model := utils.ConfigSyncModel(configCallback, configSyncInitCallback)
		utils.CupSleeve(model)
	},
}

func configSyncInitCallback(model *utils.Model, allFormsModels []utils.Values) {
	modelValues := model.ModelValues

	if modelValues.Map == nil {
		model.HandleError([]error{fmt.Errorf("model values not properly initialized")})
		return
	}

	result, err := utils.SyncCogsWithConfig()
	if err != nil {
		errors := []error{fmt.Errorf("failed to sync cogs with config: %w", err)}
		model.HandleError(errors)
		return
	}

	emptyString := ""
	for key := range modelValues.Map {
		if modelValues.Map[key] == nil {
			modelValues.Map[key] = &emptyString
		}
	}

	if len(result.AddedCogs) > 0 {
		addedCogs := strings.Join(result.AddedCogs, ", ")
		*modelValues.Map["addedCogs"] = addedCogs
	}

	if len(result.UpdatedCogs) > 0 {
		updatedCogs := strings.Join(result.UpdatedCogs, ", ")
		*modelValues.Map["updatedCogs"] = updatedCogs
	}

	if len(result.RemovedCogs) > 0 {
		removedCogs := strings.Join(result.RemovedCogs, ", ")
		*modelValues.Map["removedCogs"] = removedCogs
	}

	if len(result.HeaderIssues) > 0 {
		headerIssues := strings.Join(result.HeaderIssues, ", ")
		*modelValues.Map["headerIssues"] = headerIssues
	}

	if len(result.Errors) > 0 {
		var errors []error
		for _, err := range result.Errors {
			errors = append(errors, fmt.Errorf("sync error: %s", err))
		}
		model.HandleError(errors)
		return
	}

	if len(result.AddedCogs) == 0 && len(result.UpdatedCogs) == 0 && len(result.RemovedCogs) == 0 {
		*modelValues.Map["noChanges"] = "No changes detected in cogs."
	}
}

func init() {
	configCmd.AddCommand(syncCmd)
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
