/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"fmt"
	"strings"

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
		"helpStyle":              new(string),
		"dockerize":              new(string),
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
			if formValues.Map["helpStyle"] != nil {
				*modelValues.Map["helpStyle"] = *formValues.Map["helpStyle"]
			}
			if formValues.Map["dockerize"] != nil {
				*modelValues.Map["dockerize"] = *formValues.Map["dockerize"]
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
				Validate(ValidateBotName),

			huh.NewText().
				Title("Enter a description of your bot").
				Value(values.Map["botDescription"]).
				CharLimit(100).
				Validate(ValidateBotDescription),

			huh.NewInput().
				Title("Enter the author of your bot").
				Prompt("> ").
				Value(values.Map["botAuthor"]).
				Validate(ValidateBotAuthor),

			huh.NewInput().
				Title("Enter the command prefix for your bot (default: '!')").
				Prompt("> ").
				Value(values.Map["botPrefix"]).
				Validate(func(s string) error {
					if s == "" {
						*values.Map["botPrefix"] = "!"
						return nil
					}
					return ValidateBotPrefix(s)
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
				Validate(ValidateEnvChoice),

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
					return ValidateToken(*values.Map["envChoice"], s)
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
				Validate(ValidateLicense),

			huh.NewSelect[string]().
				Title("How should the generated help command format its output?").
				Options(
					huh.NewOption("Compact", "compact"),
					huh.NewOption("Detailed", "detailed"),
				).
				Value(values.Map["helpStyle"]).
				Validate(ValidateHelpStyle),

			huh.NewConfirm().
				Title("Generate Docker files for this project?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					// Store the confirm result on the values bus as yes or no
					v := "no"
					if b {
						v = "yes"
					}
					*values.Map["dockerize"] = v
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
	// Positions of each wrapper in the returned slice, referenced by callbacks and branches
	const (
		idxFileName = iota
		idxCmdStart
		idxCmdInfo
		idxArgStart
		idxArgInfo
		idxFieldStart
		idxFieldInfo
		idxAccept
	)

	forms := []FormWrapper{}
	{ // NOTE: idxFileName
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
	{ // NOTE: idxCmdStart
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
				allForms[idxCmdInfo].Values.Map["cmdName"] = new(string)
				allForms[idxCmdInfo].Values.Map["cmdType"] = new(string)
				allForms[idxCmdInfo].Values.Map["cmdScope"] = new(string)
				allForms[idxCmdInfo].Values.Map["cmdDescription"] = new(string)
				allForms[idxCmdInfo].Values.Map["cmdReturnType"] = new(string)
				allForms[idxArgInfo].Values.Map["args"] = new(string)
				allForms[idxFieldInfo].Values.Map["fields"] = new(string)
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
	{ // NOTE: idxCmdInfo
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
				// Modal commands only respond through the modal, so their return type is fixed
				returnType := *formValues.Map["cmdReturnType"]
				if *formValues.Map["cmdType"] == "modal" {
					returnType = "None"
				}
				command := CommandInfo{
					Name:        *formValues.Map["cmdName"],
					Type:        *formValues.Map["cmdType"],
					Scope:       *formValues.Map["cmdScope"],
					Description: *formValues.Map["cmdDescription"],
					Args:        []ArgInfo{},
					Fields:      []FieldInfo{},
					ReturnType:  returnType,
				}
				commandString, _ := command.ToJSON()
				modelValues.Map["currentCommand"] = &commandString
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				// Modal commands collect fields instead of arguments
				if *formValues.Map["cmdType"] == "modal" {
					return idxFieldStart
				}
				return -1
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxArgStart
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
				allForms[idxArgInfo].Values.Map["argName"] = new(string)
				allForms[idxArgInfo].Values.Map["argDescription"] = new(string)
				allForms[idxArgInfo].Values.Map["argType"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["argStartConfirm"] == "yes" {
					return -1
				}
				return idxAccept
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxArgInfo
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
				return idxArgStart
			},
			BranchValueHandler: func(targetFormIndex int, targetValues Values) {
				if targetFormIndex == idxCmdStart {
					ResetFormValues(targetValues)
				}
				if targetFormIndex == idxCmdInfo {
					ResetFormValues(targetValues)
				}
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxFieldStart
		values := map[string]*string{
			"fieldStartConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Field Start",
			Form: addFieldStartFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addFieldStartValues",
			},
			ShowStatus: false,
			FormGroup:  "field",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				allForms[idxFieldInfo].Values.Map["fieldName"] = new(string)
				allForms[idxFieldInfo].Values.Map["fieldLabel"] = new(string)
				allForms[idxFieldInfo].Values.Map["fieldStyle"] = new(string)
				allForms[idxFieldInfo].Values.Map["fieldRequired"] = new(string)
				allForms[idxFieldInfo].Values.Map["fieldPlaceholder"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["fieldStartConfirm"] == "yes" {
					return -1
				}
				return idxAccept
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxFieldInfo
		values := map[string]*string{
			"fields":           new(string),
			"fieldName":        new(string),
			"fieldLabel":       new(string),
			"fieldStyle":       new(string),
			"fieldRequired":    new(string),
			"fieldPlaceholder": new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Field Info",
			Form: addFieldInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addFieldInfoValues",
			},
			ShowStatus: false,
			FormGroup:  "field",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				currentCommand, _ := JSONToCmdInfo(*modelValues.Map["currentCommand"])

				currentCommand.Fields = append(currentCommand.Fields, FieldInfo{
					Name:        *values["fieldName"],
					Label:       *values["fieldLabel"],
					Style:       *values["fieldStyle"],
					Required:    *values["fieldRequired"] == "yes",
					Placeholder: *values["fieldPlaceholder"],
				})
				fieldString, _ := FieldInfoSliceToJSON(currentCommand.Fields)
				formValues.Map["fields"] = &fieldString
				commandString, _ := currentCommand.ToJSON()
				modelValues.Map["currentCommand"] = &commandString
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				// A modal page cannot hold more inputs than Discord allows
				fields, _ := JSONToFieldInfoSlice(*formValues.Map["fields"])
				if len(fields) >= MaxModalFields {
					return idxAccept
				}
				return idxFieldStart
			},
			BranchValueHandler: func(targetFormIndex int, targetValues Values) {
				if targetFormIndex == idxCmdStart {
					ResetFormValues(targetValues)
				}
				if targetFormIndex == idxCmdInfo {
					ResetFormValues(targetValues)
				}
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxAccept
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
					// Modal commands are app commands, so they live with the slash commands
					if command.Type == "slash" || command.Type == "modal" {
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
				return idxCmdStart
			},
			BranchValueHandler: func(targetFormIndex int, targetValues Values) {
				if targetFormIndex == idxCmdStart {
					ResetFormValues(targetValues)
				}
				if targetFormIndex == idxCmdInfo {
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
					err := ValidateFileName(s)
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
					slashCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["slashCommands"])
					prefixCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["prefixCommands"])
					return ValidateCommandName(s, append(slashCommandList, prefixCommandList...))
				}),
			huh.NewSelect[string]().
				Value(values.Map["cmdType"]).
				Title("Select the command type").
				Options(
					huh.NewOption("slash", "slash"),
					huh.NewOption("prefix", "prefix"),
					huh.NewOption("modal", "modal"),
				).
				Validate(ValidateCommandType),
			huh.NewSelect[string]().
				Value(values.Map["cmdScope"]).
				Title("Select the command scope").
				Options(
					huh.NewOption("guild", "guild"),
					huh.NewOption("global", "global"),
				).
				Validate(ValidateCommandScope),
			huh.NewText().
				Value(values.Map["cmdDescription"]).
				Title("Enter the command description").
				CharLimit(400).
				Validate(ValidateCommandDescription),
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
				Validate(ValidateReturnType),
		),
	)
	return cmdInfoForm
}

func addCmdAcceptFormGenerator(values Values, modelValues Values) *huh.Form {
	var commandName, commandType, commandDesc, commandReturn, commandArgs, commandFields string

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

			if len(currentCommand.Fields) > 0 {
				fieldNames := make([]string, len(currentCommand.Fields))
				for i, field := range currentCommand.Fields {
					fieldNames[i] = fmt.Sprintf("%s (%s)", field.Name, field.Style)
				}
				commandFields = strings.Join(fieldNames, ", ")
			} else {
				commandFields = "None"
			}
		}
	}

	summary := fmt.Sprintf("Command Name: %s\nCommand Type: %s\nDescription: %s\nReturn Type: %s\nArguments: %v",
		commandName, commandType, commandDesc, commandReturn, commandArgs)
	if commandType == "modal" {
		summary = fmt.Sprintf("Command Name: %s\nCommand Type: %s\nDescription: %s\nReturn Type: %s\nFields: %v",
			commandName, commandType, commandDesc, commandReturn, commandFields)
	}

	cmdAcceptForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Command Info").
				Description(summary),
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
					args, _ := JSONToArgInfoSlice(*values.Map["args"])
					return ValidateArgName(s, args)
				}),
			huh.NewText().
				Value(values.Map["argDescription"]).
				Title("Enter the argument description").
				CharLimit(200).
				Validate(ValidateArgDescription),
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
				Validate(ValidateArgType),
		),
	)
	return argInfoForm
}

func addFieldStartFormGenerator(values Values, modelValues Values) *huh.Form {
	fieldStartForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to add a modal field?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["fieldStartConfirm"] = &s
					return nil
				}),
		),
	)
	return fieldStartForm
}

func addFieldInfoFormGenerator(values Values, modelValues Values) *huh.Form {
	fieldInfoForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(values.Map["fieldName"]).
				Title("Enter the field name").
				Prompt("> ").
				Validate(func(s string) error {
					fields, _ := JSONToFieldInfoSlice(*values.Map["fields"])
					return ValidateFieldName(s, fields)
				}),
			huh.NewInput().
				Value(values.Map["fieldLabel"]).
				Title("Enter the field label").
				Prompt("> ").
				Validate(ValidateFieldLabel),
			huh.NewSelect[string]().
				Value(values.Map["fieldStyle"]).
				Title("Select the field style").
				Options(
					huh.NewOption("short", "short"),
					huh.NewOption("paragraph", "paragraph"),
				).
				Validate(ValidateFieldStyle),
			huh.NewConfirm().
				Title("Is the field required?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["fieldRequired"] = &s
					return nil
				}),
			huh.NewInput().
				Value(values.Map["fieldPlaceholder"]).
				Title("Enter the field placeholder (optional)").
				Prompt("> "),
		),
	)
	return fieldInfoForm
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
