/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package cmd

import (
	"fmt"

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configIsOptional bool
	allConfigs       bool
	configArgs       []string

	configCmd = &cobra.Command{
		Use:   "config [sections...]",
		Short: "Disoplay configuration",
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
)

func localConfigInitCallback(model *utils.Model, allFormsModels []utils.Values) {
	var errors []error
	modelValues := model.ModelValues

	config, err := utils.LoadConfig()
	if err != nil {
		errors = append(errors, fmt.Errorf("error finding root directory: %w", err))
		model.HandleError(errors)
		return
	}
	if allConfigs {
		modelValues.Map["name"] = &config.BotInfo.Name
		modelValues.Map["author"] = &config.BotInfo.Author
		modelValues.Map["description"] = &config.BotInfo.Description
		modelValues.Map["prefix"] = &config.BotInfo.CommandPrefix
		var cogConfigs string
		if len(config.Cogs) > 0 {
			cogConfigs, err = utils.CogConfigSliceToJSON(config.Cogs)
			if err != nil {
				errors = append(errors, fmt.Errorf("error converting cogs to JSON: %w", err))
				model.HandleError(errors)
				return
			}
		} else {
			cogConfigs = "[]"
		}
		modelValues.Map["cogs"] = &cogConfigs
	} else {
		var invalidArgs []string
		for _, arg := range configArgs {
			switch arg {
			case "name":
				modelValues.Map["name"] = &config.BotInfo.Name
			case "author":
				modelValues.Map["author"] = &config.BotInfo.Author
			case "description":
				modelValues.Map["description"] = &config.BotInfo.Description
			case "prefix":
				modelValues.Map["prefix"] = &config.BotInfo.CommandPrefix
			case "cogs":
				var cogConfigs string
				if len(config.Cogs) > 0 {
					cogConfigs, err = utils.CogConfigSliceToJSON(config.Cogs)
					if err != nil {
						errors = append(errors, fmt.Errorf("error converting cogs to JSON: %w", err))
						model.HandleError(errors)
						return
					}
				} else {
					cogConfigs = "[]"
				}
				modelValues.Map["cogs"] = &cogConfigs
			default:
				invalidArgs = append(invalidArgs, arg)
			}
		}
		if len(invalidArgs) > 0 {
			errors = append(errors, fmt.Errorf("invalid configuration sectins: %v", invalidArgs))
			model.HandleError(errors)
			return
		}
	}
}

func globalConfigInitCallback(model *utils.Model, allFormsModels []utils.Values) {
	var errors []error
	modelValues := model.ModelValues
	if allConfigs {
		*modelValues.Map["version"] = viper.GetString("cli.version")
		if viper.GetBool("cli.check_updates") {
			*modelValues.Map["check_updates"] = "true"
		} else {
			*modelValues.Map["check_updates"] = "false"
		}
		if viper.GetBool("cli.auto_update") {
			*modelValues.Map["auto_update"] = "true"
		} else {
			*modelValues.Map["auto_update"] = "false"
		}
		*modelValues.Map["default_user"] = viper.GetString("user.default_user")
		*modelValues.Map["github_username"] = viper.GetString("user.github_username")
		if viper.GetBool("display.scroll_enabled") {
			*modelValues.Map["scroll_enabled"] = "true"
		} else {
			*modelValues.Map["scroll_enabled"] = "false"
		}
		*modelValues.Map["color_scheme"] = viper.GetString("display.color_scheme")
		*modelValues.Map["default_command_prefix"] = viper.GetString("defaults.command_prefix")
		*modelValues.Map["default_python_version"] = viper.GetString("defaults.python_version")
		if viper.GetBool("defaults.auto_git_init") {
			*modelValues.Map["auto_git_init"] = "true"
		} else {
			*modelValues.Map["auto_git_init"] = "false"
		}
		*modelValues.Map["editor"] = viper.GetString("dev.editor")
	} else {
		var invalidArgs []string
		for _, arg := range configArgs {
			switch arg {
			case "version":
				*modelValues.Map["version"] = viper.GetString("cli.version")
			case "check-updates":
				if viper.GetBool("cli.check_updates") {
					*modelValues.Map["check_updates"] = "true"
				} else {
					*modelValues.Map["check_updates"] = "false"
				}
			case "auto-update":
				if viper.GetBool("cli.auto_update") {
					*modelValues.Map["auto_update"] = "true"
				} else {
					*modelValues.Map["auto_update"] = "false"
				}
			case "user":
				*modelValues.Map["default_user"] = viper.GetString("user.default_user")
			case "github":
				*modelValues.Map["github_username"] = viper.GetString("user.github_username")
			case "scroll":
				if viper.GetBool("display.scroll_enabled") {
					*modelValues.Map["scroll_enabled"] = "true"
				} else {
					*modelValues.Map["scroll_enabled"] = "false"
				}
			case "colorscheme":
				*modelValues.Map["color_scheme"] = viper.GetString("display.color_scheme")
			case "command-prefix":
				*modelValues.Map["default_command_prefix"] = viper.GetString("defaults.command_prefix")
			case "git":
				if viper.GetBool("defaults.auto_git_init") {
					*modelValues.Map["auto_git_init"] = "true"
				} else {
					*modelValues.Map["auto_git_init"] = "false"
				}
			case "python":
				*modelValues.Map["default_python_version"] = viper.GetString("defaults.python_version")
			case "editor":
				*modelValues.Map["editor"] = viper.GetString("dev.editor")
			default:
				invalidArgs = append(invalidArgs, arg)
			}
		}
		if len(invalidArgs) > 0 {
			errors = append(errors, fmt.Errorf("invalid configuration sectins: %v", invalidArgs))
			model.HandleError(errors)
			return
		}
	}
}

func configCallback(model *utils.Model) []error { return nil }

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolP("global", "g", false, "Show global CLI configuration")
	configCmd.Flags().Bool("local", false, "Show local project configuration (default)")

	configCmd.Flags().StringP("format", "f", "default", "Output format (default, json, yaml)")
	configCmd.Flags().BoolP("keys-only", "k", false, "Show only configuration keys")

	configCmd.MarkFlagsMutuallyExclusive("global", "local")

	configCmd.Flags().BoolP("help", "h", false, "Help for config")
}

/*
Copyright © 2025 2025 Austin "Choice404" Choi

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
