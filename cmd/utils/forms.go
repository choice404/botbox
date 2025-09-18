/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/charmbracelet/huh"
)

/**
 * Create Forms and Model Generators
 */
func CreateFormWrapperGenerator() []FormWrapper {
	values := map[string]*string{
		"botName":                new(string),
		"botDescription":         new(string),
		"botAuthor":              new(string),
		"botPrefix":              new(string),
		"envChoice":              new(string),
		"botTokenDopplerProject": new(string),
		"botGuildDopplerEnv":     new(string),
		"licenseType":            new(string),
	}

	wrapper := FormWrapper{
		Name: "Create Bot",
		Form: createFormGenerator,
		Values: Values{
			Map:  values,
			Name: "createBotValues",
		},
		ShowStatus: true,
		Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
			if formValues.Map["botName"] != nil {
				*modelValues.Map["botName"] = *formValues.Map["botName"]
			}
			if formValues.Map["botDescription"] != nil {
				*modelValues.Map["botDescription"] = *formValues.Map["botDescription"]
			}
			if formValues.Map["botAuthor"] != nil {
				*modelValues.Map["botAuthor"] = *formValues.Map["botAuthor"]
			}
			if formValues.Map["botPrefix"] != nil {
				*modelValues.Map["botPrefix"] = *formValues.Map["botPrefix"]
			}
			if formValues.Map["envChoice"] != nil {
				*modelValues.Map["envChoice"] = *formValues.Map["envChoice"]
			}
			if formValues.Map["botTokenDopplerProject"] != nil {
				*modelValues.Map["botTokenDopplerProject"] = *formValues.Map["botTokenDopplerProject"]
			}
			if formValues.Map["botGuildDopplerEnv"] != nil {
				*modelValues.Map["botGuildDopplerEnv"] = *formValues.Map["botGuildDopplerEnv"]
			}
			if formValues.Map["licenseType"] != nil {
				*modelValues.Map["licenseType"] = *formValues.Map["licenseType"]
			}
		},
	}
	return []FormWrapper{wrapper}
}

func createFormGenerator(values Values, modelValues Values) *huh.Form {
	createForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the name of your bot").
				Prompt("> ").
				Value(values.Map["botName"]).
				Validate(func(s string) error {
					if *values.Map["botName"] == "" {
						return fmt.Errorf("Bot name cannot be empty")
					}
					if len(s) > 20 {
						return fmt.Errorf("Bot name is too long")
					}
					r := []rune(s)[0]
					if !unicode.IsLetter(r) {
						return fmt.Errorf("Bot name must start with a letter")
					}
					if strings.ContainsRune(s, ' ') || strings.ContainsRune(s, '\t') {
						return fmt.Errorf("Bot name cannot contain whitespace")
					}
					if strings.ContainsAny(s, "!@#$%^&*()_+={}[]|\\:;\"'<>,.?/~`") {
						return fmt.Errorf("Bot name cannot contain special characters")
					}
					return nil
				}),

			huh.NewText().
				Title("Enter a description of your bot").
				Value(values.Map["botDescription"]).
				CharLimit(100).
				Validate(func(s string) error {
					if *values.Map["botDescription"] == "" {
						return fmt.Errorf("Description cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Title("Enter the author of your bot").
				Prompt("> ").
				Value(values.Map["botAuthor"]).
				Validate(func(s string) error {
					if *values.Map["botAuthor"] == "" {
						return fmt.Errorf("Author name cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Title("Enter the command prefix for your bot (default: '!')").
				Prompt("> ").
				Value(values.Map["botPrefix"]).
				Validate(func(s string) error {
					if s == "" {
						*values.Map["botPrefix"] = "!"
						return nil
					}
					if len(s) > 1 {
						return fmt.Errorf("Command prefix must be a single character")
					}
					r := []rune(s)[0]

					if unicode.IsLetter(r) || unicode.IsDigit(r) {
						return fmt.Errorf("Command prefix can not be an alphanumeric character")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How do you want to handle environment variables?").
				Options(
					huh.NewOption("Create a .env file", "env"),
					huh.NewOption("Use Doppler", "doppler"),
				).
				Value(values.Map["envChoice"]).
				Validate(func(s string) error {
					if s != "env" && s != "doppler" {
						return fmt.Errorf("Please select either 'env' or 'doppler'")
					}
					return nil
				}),

			huh.NewInput().
				TitleFunc(func() string {
					if *values.Map["envChoice"] == "env" {
						return "Enter the bot token"
					}
					return "Enter the Doppler project name"
				}, values.Map["envChoice"]).
				Prompt("> ").
				EchoMode(huh.EchoModePassword).
				Validate(func(s string) error {
					if *values.Map["envChoice"] == "env" {
						if s == "" {
							return fmt.Errorf("Token cannot be empty")
						}
						if len(s) < 9 {
							return fmt.Errorf("Token is too short")
						}
					}
					return nil
				}).
				Value(values.Map["botTokenDopplerProject"]),

			huh.NewInput().
				TitleFunc(func() string {
					if *values.Map["envChoice"] == "env" {
						return "Enter the bot guild ID"
					}
					return "Enter the Doppler environment name"
				}, values.Map["envChoice"]).
				Prompt("> ").
				Value(values.Map["botGuildDopplerEnv"]).
				Validate(func(s string) error {
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What license do you want to use?").
				Options(
					huh.NewOption("MIT", "mit"),
					huh.NewOption("Apache 2.0", "apache-2.0"),
					huh.NewOption("GPLv3", "gpl-3.0"),
					huh.NewOption("BSD 3-Clause", "bsd-3-clause"),
					huh.NewOption("Unlicense", "unlicense"),
					huh.NewOption("No license", "no-license"),
				).
				Value(values.Map["licenseType"]).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("Please select a license type")
					}
					return nil
				}),
		),
	).
		WithWidth(100).
		WithShowHelp(false).
		WithShowErrors(false)

	return createForm
}

/**
 * Add forms and model generators for other functionalities
 */
func AddFormWrapperGenerator() []FormWrapper {
	forms := []FormWrapper{}
	{ // NOTE: 0
		values := map[string]*string{
			"filename": new(string),
		}

		wrapper := FormWrapper{
			Name: "Add File Name",
			Form: addFileNameFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addFileNameValues",
			},
			ShowStatus: false,
			FormGroup:  "filename",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				if *formValues.Map["filename"] == "" {
					*formValues.Map["filename"] = "cog.py"
				}
				*modelValues.Map["filename"] = *formValues.Map["filename"]
			},

			SkipCondition: func(modelValues Values, allForms []FormWrapper, currentIndex int) bool {
				if modelValues.Map["filename"] != nil && *modelValues.Map["filename"] != "" {
					return true
				}
				return false
			},
			SkipCallback: func(modelValues Values, allForms []FormWrapper, currentIndex int) {
				if modelValues.Map["filename"] != nil && *modelValues.Map["filename"] != "" {
					filename := *modelValues.Map["filename"]
					allForms[currentIndex].Values.Map["filename"] = &filename
				}
			},
		}

		forms = append(forms, wrapper)
	}
	{ // NOTE: 1
		values := map[string]*string{
			"cmdStartConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Command Start",
			Form: addCmdStartFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addCommandStartValues",
			},
			ShowStatus: false,
			FormGroup:  "command",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				allForms[2].Values.Map["cmdName"] = new(string)
				allForms[2].Values.Map["cmdType"] = new(string)
				allForms[2].Values.Map["cmdScope"] = new(string)
				allForms[2].Values.Map["cmdDescription"] = new(string)
				allForms[2].Values.Map["cmdReturnType"] = new(string)
				allForms[4].Values.Map["args"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["cmdStartConfirm"] == "yes" {
					return -1
				} else {
					return -2
				}
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: 2
		values := map[string]*string{
			"cmdName":        new(string),
			"cmdType":        new(string),
			"cmdScope":       new(string),
			"cmdDescription": new(string),
			"cmdReturnType":  new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Command Info",
			Form: addCmdInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addCommandInfoValues",
			},
			ShowStatus: false,
			FormGroup:  "command",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				command := CommandInfo{
					Name:        *formValues.Map["cmdName"],
					Type:        *formValues.Map["cmdType"],
					Scope:       *formValues.Map["cmdScope"],
					Description: *formValues.Map["cmdDescription"],
					Args:        []ArgInfo{},
					ReturnType:  *formValues.Map["cmdReturnType"],
				}
				commandString, _ := command.ToJSON()
				modelValues.Map["currentCommand"] = &commandString
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: 3
		values := map[string]*string{
			"argStartConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Argument Start",
			Form: addArgStartFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addArgumentStartValues",
			},
			ShowStatus: false,
			FormGroup:  "argument",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				allForms[4].Values.Map["argName"] = new(string)
				allForms[4].Values.Map["argDescription"] = new(string)
				allForms[4].Values.Map["argType"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["argStartConfirm"] == "yes" {
					return -1
				}
				return 5
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: 4
		values := map[string]*string{
			"args":           new(string),
			"argName":        new(string),
			"argDescription": new(string),
			"argType":        new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Argument Info",
			Form: addArgInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addArgumentInfoValues",
			},
			ShowStatus: false,
			FormGroup:  "argument",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				currentCommand, _ := JSONToCmdInfo(*modelValues.Map["currentCommand"])

				currentCommand.Args = append(currentCommand.Args, ArgInfo{
					Name:        *values["argName"],
					Type:        *values["argType"],
					Description: *values["argDescription"],
				})
				argString, _ := ArgInfoSliceToJSON(currentCommand.Args)
				formValues.Map["args"] = &argString
				commandString, _ := currentCommand.ToJSON()
				modelValues.Map["currentCommand"] = &commandString
			},
			BranchCallback: func(values Values, allForms []FormWrapper) int {
				return 3
			},
			BranchValueHandler: func(targetFormIndex int, targetValues Values) {
				if targetFormIndex == 1 {
					ResetFormValues(targetValues)
				}
				if targetFormIndex == 2 {
					ResetFormValues(targetValues)
				}
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: 5
		values := map[string]*string{
			"cmdAcceptConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Command Accept",
			Form: addCmdAcceptFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addCommandAcceptValues",
			},
			ShowStatus: false,
			FormGroup:  "command",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				if *formValues.Map["cmdAcceptConfirm"] == "yes" {
					command, _ := JSONToCmdInfo(*modelValues.Map["currentCommand"])
					if command.Type == "slash" {
						slashCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["slashCommands"])
						slashCommandList = append(slashCommandList, *command)
						jsonData, _ := CmdInfoSliceToJSON(slashCommandList)
						modelValues.Map["slashCommands"] = &jsonData
					} else if command.Type == "prefix" {
						prefixCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["prefixCommands"])
						prefixCommandList = append(prefixCommandList, *command)
						jsonData, _ := CmdInfoSliceToJSON(prefixCommandList)
						modelValues.Map["prefixCommands"] = &jsonData
					}
				}
			},
			BranchCallback: func(values Values, allForms []FormWrapper) int {
				return 1
			},
			BranchValueHandler: func(targetFormIndex int, targetValues Values) {
				if targetFormIndex == 1 {
					ResetFormValues(targetValues)
				}
				if targetFormIndex == 2 {
					ResetFormValues(targetValues)
				}
			},
		}
		forms = append(forms, wrapper)
	}

	return forms
}

func addFileNameFormGenerator(values Values, modelValues Values) *huh.Form {
	fileNameForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(values.Map["filename"]).
				Title("Enter the filename").
				Prompt("> ").
				Validate(func(s string) error {
					err := validateFileName(s)
					if err != nil {
						return err
					}
					return nil
				}),
		),
	)
	return fileNameForm
}

func addCmdStartFormGenerator(values Values, modelValues Values) *huh.Form {
	cmdStartForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to add a command?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["cmdStartConfirm"] = &s
					return nil
				}),
		),
	)
	return cmdStartForm
}

func addCmdInfoFormGenerator(values Values, modelValues Values) *huh.Form {
	cmdInfoForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(values.Map["cmdName"]).
				Title("Enter the command name").
				Prompt("> ").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("command name cannot be empty")
					}
					if strings.Contains(s, " ") {
						return fmt.Errorf("command name cannot contain spaces")
					}
					slashCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["slashCommands"])
					prefixCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["prefixCommands"])
					if commandExists(s, append(slashCommandList, prefixCommandList...)) {
						return fmt.Errorf("command name already exists")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Value(values.Map["cmdType"]).
				Title("Select the command type").
				Options(
					huh.NewOption("slash", "slash"),
					huh.NewOption("prefix", "prefix"),
				).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("return type cannot be empty")
					}
					if s != "slash" && s != "prefix" {
						return fmt.Errorf("return type must either slash or prefix")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Value(values.Map["cmdScope"]).
				Title("Select the command scope").
				Options(
					huh.NewOption("guild", "guild"),
					huh.NewOption("global", "global"),
				).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("command scope cannot be empty")
					}
					if s != "guild" && s != "global" {
						return fmt.Errorf("command scope must either guild or global")
					}
					return nil
				}),
			huh.NewText().
				Value(values.Map["cmdDescription"]).
				Title("Enter the command description").
				CharLimit(400).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("description cannot be empty")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Value(values.Map["cmdReturnType"]).
				Title("Enter the command return type").
				Options(
					huh.NewOption("str", "str"),
					huh.NewOption("int", "int"),
					huh.NewOption("float", "float"),
					huh.NewOption("bool", "bool"),
					huh.NewOption("None", "None"),
				).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("return type cannot be empty")
					}
					if s != "str" && s != "int" && s != "float" && s != "bool" && s != "None" {
						return fmt.Errorf("return type must be one of str, int, float, bool, None")
					}
					return nil
				}),
		),
	)
	return cmdInfoForm
}

func addCmdAcceptFormGenerator(values Values, modelValues Values) *huh.Form {
	var commandName, commandType, commandDesc, commandReturn, commandArgs string

	if modelValues.Map["currentCommand"] != nil && *modelValues.Map["currentCommand"] != "" {
		currentCommand, err := JSONToCmdInfo(*modelValues.Map["currentCommand"])
		if err == nil {
			commandName = currentCommand.Name
			commandType = currentCommand.Type
			commandDesc = currentCommand.Description
			commandReturn = currentCommand.ReturnType

			if len(currentCommand.Args) > 0 {
				argNames := make([]string, len(currentCommand.Args))
				for i, arg := range currentCommand.Args {
					argNames[i] = fmt.Sprintf("%s (%s)", arg.Name, arg.Type)
				}
				commandArgs = strings.Join(argNames, ", ")
			} else {
				commandArgs = "None"
			}
		}
	}

	cmdAcceptForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Command Info").
				Description(fmt.Sprintf("Command Name: %s\nCommand Type: %s\nDescription: %s\nReturn Type: %s\nArguments: %v",
					commandName, commandType, commandDesc, commandReturn, commandArgs)),
			huh.NewConfirm().
				Title("Does everything look correct?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["cmdAcceptConfirm"] = &s
					return nil
				}),
		),
	)
	return cmdAcceptForm
}

func addArgStartFormGenerator(values Values, modelValues Values) *huh.Form {
	argStartForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to add an argument?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["argStartConfirm"] = &s
					return nil
				}),
		),
	)
	return argStartForm
}

func addArgInfoFormGenerator(values Values, modelValues Values) *huh.Form {
	argInfoForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(values.Map["argName"]).
				Title("Enter the argument name").
				Prompt("> ").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("Argument name cannot be empty")
					}
					if strings.Contains(s, " ") {
						return fmt.Errorf("Argument name cannot contain spaces")
					}
					if strings.Contains(s, "-") {
						return fmt.Errorf("Argument name cannot contain dashes")
					}
					args, _ := JSONToArgInfoSlice(*values.Map["args"])
					if argExists(s, args) {
						return fmt.Errorf("Argument name already exists")
					}
					return nil
				}),
			huh.NewText().
				Value(values.Map["argDescription"]).
				Title("Enter the argument description").
				CharLimit(200).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("argument description cannot be empty")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Value(values.Map["argType"]).
				Title("Enter the argument type").
				Options(
					huh.NewOption("str", "str"),
					huh.NewOption("int", "int"),
					huh.NewOption("float", "float"),
					huh.NewOption("bool", "bool"),
					huh.NewOption("discord.Member", "discord.Member"),
					huh.NewOption("discord.Role", "discord.Role"),
				).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("argument type cannot be empty")
					}
					if s != "str" && s != "int" && s != "float" && s != "bool" && s != "discord.Member" && s != "discord.Role" {
						return fmt.Errorf("argument type must be one of str, int, float, bool, discord.Member, discord.Role")
					}
					return nil
				}),
		),
	)
	return argInfoForm
}

/**
 * Remove Forms and Model Generators
 */
func RemoveFormWrapperGenerator() []FormWrapper {
	forms := []FormWrapper{}
	{ // NOTE: 0
		values := map[string]*string{
			"cogName": new(string),
		}
		wrapper := FormWrapper{
			Name: "Remove Cog",
			Form: removeCogFormGenerator,
			Values: Values{
				Map:  values,
				Name: "removeCogValues",
			},
			ShowStatus: true,
			FormGroup:  "cog",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				if formValues.Map["cogName"] != nil {
					*modelValues.Map["cogName"] = *formValues.Map["cogName"]
				} else {
					*modelValues.Map["cogName"] = ""
				}
			},
			SkipCondition: func(modelValues Values, allForms []FormWrapper, currentIndex int) bool {
				if modelValues.Map["cogName"] != nil && *modelValues.Map["cogName"] != "" {
					return true
				}
				return false
			},
			SkipCallback: func(modelValues Values, allForms []FormWrapper, currentIndex int) {
				if modelValues.Map["cogName"] != nil && *modelValues.Map["cogName"] != "" {
					cogName := *modelValues.Map["cogName"]
					allForms[currentIndex].Values.Map["cogName"] = &cogName
				}
			},
		}
		forms = append(forms, wrapper)
	}
	return forms
}

func removeCogFormGenerator(values Values, modelValues Values) *huh.Form {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		errorForm := huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("Error").
					Description("Failed to load configuration file."),
			),
		)
		errorForm.State = huh.StateCompleted
		return errorForm
	}

	var cogList []string
	for _, cog := range config.Cogs {
		cogList = append(cogList, cog.Name)
	}

	if len(cogList) == 0 {
		noCogForm := huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("No Cogs Available").
					Description("There are no cogs to remove. Please add some cogs first."),
			),
		)
		return noCogForm
	}

	cmdRemoveForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(values.Map["cogName"]).
				Height(8).
				Title("Select a cog to remove").
				Options(huh.NewOptions(cogList...)...),
		),
	)
	return cmdRemoveForm
}

/**
 * Configuration Forms and Model Generators
 */
func ConfigFormWrapperGenerator() []FormWrapper {
	return []FormWrapper{}
}

/**
 * Nont specific forms for special use cases
 */

func finalCompleteFormGenerator(values Values, modelValues Values) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("").
				Description(""),
		),
	)

	form.State = huh.StateCompleted

	return form
}

func ConfigSyncFormWrapperGenerator() []FormWrapper {
	return []FormWrapper{
		{
			Name: "Sync Config",
			Form: finalCompleteFormGenerator,
			Values: Values{
				Map:  map[string]*string{},
				Name: "configSyncValues",
			},
			ShowStatus: false,
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
			},
		},
	}
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
