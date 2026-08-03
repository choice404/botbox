/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"encoding/json"
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
		idxMultiPage
		idxPageInfo
		idxPageFieldStart
		idxPageFieldInfo
		idxBranchStart
		idxBranchInfo
		idxPageNext
		idxResponseStart
		idxResponseInfo
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
				// Every new command starts with clean page and response state
				allForms[idxPageInfo].Values.Map["pageName"] = new(string)
				allForms[idxPageInfo].Values.Map["pageTitle"] = new(string)
				allForms[idxPageFieldInfo].Values.Map["pageFields"] = new(string)
				allForms[idxPageNext].Values.Map["pages"] = new(string)
				allForms[idxResponseInfo].Values.Map["responses"] = new(string)
				emptyPages := "[]"
				modelValues.Map["pages"] = &emptyPages
				modelValues.Map["currentPage"] = new(string)
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
					Pages:       []PageInfo{},
					Responses:   []ResponseInfo{},
					ReturnType:  returnType,
				}
				commandString, _ := command.ToJSON()
				modelValues.Map["currentCommand"] = &commandString
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				// Modal commands first decide between a single page and a multi page flow
				if *formValues.Map["cmdType"] == "modal" {
					return idxMultiPage
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
				return idxResponseStart
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
				return idxResponseStart
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
					return idxResponseStart
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
	{ // NOTE: idxMultiPage
		values := map[string]*string{
			"multiPageConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Multi Page",
			Form: addMultiPageFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addMultiPageValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				// A fresh modal command always starts with an empty page collection
				emptyPages := "[]"
				modelValues.Map["pages"] = &emptyPages
				mirror := "[]"
				allForms[idxPageNext].Values.Map["pages"] = &mirror
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["multiPageConfirm"] == "yes" {
					return idxPageInfo
				}
				// A single page modal keeps the existing field loop
				return idxFieldStart
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxPageInfo
		values := map[string]*string{
			"pageName":  new(string),
			"pageTitle": new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Page Info",
			Form: addPageInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addPageInfoValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				page := PageInfo{
					Name:     *formValues.Map["pageName"],
					Title:    *formValues.Map["pageTitle"],
					Fields:   []FieldInfo{},
					Branches: []BranchRule{},
				}
				pageString, _ := pageToJSON(page)
				modelValues.Map["currentPage"] = &pageString
				// Each page collects its own fields, so the accumulator starts empty
				allForms[idxPageFieldInfo].Values.Map["pageFields"] = new(string)
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxPageFieldStart
		values := map[string]*string{
			"pageFieldStartConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Page Field Start",
			Form: addPageFieldStartFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addPageFieldStartValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				allForms[idxPageFieldInfo].Values.Map["fieldName"] = new(string)
				allForms[idxPageFieldInfo].Values.Map["fieldLabel"] = new(string)
				allForms[idxPageFieldInfo].Values.Map["fieldStyle"] = new(string)
				allForms[idxPageFieldInfo].Values.Map["fieldRequired"] = new(string)
				allForms[idxPageFieldInfo].Values.Map["fieldPlaceholder"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["pageFieldStartConfirm"] == "yes" {
					return -1
				}
				return idxBranchStart
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxPageFieldInfo
		values := map[string]*string{
			"pageFields":       new(string),
			"fieldName":        new(string),
			"fieldLabel":       new(string),
			"fieldStyle":       new(string),
			"fieldRequired":    new(string),
			"fieldPlaceholder": new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Page Field Info",
			Form: addPageFieldInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addPageFieldInfoValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				currentPage, _ := jsonToPage(*modelValues.Map["currentPage"])

				currentPage.Fields = append(currentPage.Fields, FieldInfo{
					Name:        *formValues.Map["fieldName"],
					Label:       *formValues.Map["fieldLabel"],
					Style:       *formValues.Map["fieldStyle"],
					Required:    *formValues.Map["fieldRequired"] == "yes",
					Placeholder: *formValues.Map["fieldPlaceholder"],
				})
				fieldString, _ := FieldInfoSliceToJSON(currentPage.Fields)
				formValues.Map["pageFields"] = &fieldString
				pageString, _ := pageToJSON(currentPage)
				modelValues.Map["currentPage"] = &pageString
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				// A modal page cannot hold more inputs than Discord allows
				fields, _ := JSONToFieldInfoSlice(*formValues.Map["pageFields"])
				if len(fields) >= MaxModalFields {
					return idxBranchStart
				}
				return idxPageFieldStart
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxBranchStart
		values := map[string]*string{
			"branchStartConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Branch Start",
			Form: addBranchStartFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addBranchStartValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				allForms[idxBranchInfo].Values.Map["branchField"] = new(string)
				allForms[idxBranchInfo].Values.Map["branchEquals"] = new(string)
				allForms[idxBranchInfo].Values.Map["branchGoto"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["branchStartConfirm"] == "yes" {
					// Branch rules test a field on this page, so an empty page has nothing to branch on
					fields, _ := JSONToFieldInfoSlice(*allForms[idxPageFieldInfo].Values.Map["pageFields"])
					if len(fields) == 0 {
						return idxPageNext
					}
					return -1
				}
				return idxPageNext
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxBranchInfo
		values := map[string]*string{
			"branchField":  new(string),
			"branchEquals": new(string),
			"branchGoto":   new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Branch Info",
			Form: addBranchInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addBranchInfoValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				currentPage, _ := jsonToPage(*modelValues.Map["currentPage"])

				currentPage.Branches = append(currentPage.Branches, BranchRule{
					Field:  *formValues.Map["branchField"],
					Equals: *formValues.Map["branchEquals"],
					Goto:   *formValues.Map["branchGoto"],
				})
				pageString, _ := pageToJSON(currentPage)
				modelValues.Map["currentPage"] = &pageString
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				return idxBranchStart
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxPageNext
		values := map[string]*string{
			"pages":              new(string),
			"pageNext":           new(string),
			"pageAnotherConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Page Next",
			Form: addPageNextFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addPageNextValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				currentPage, _ := jsonToPage(*modelValues.Map["currentPage"])
				currentPage.Next = *formValues.Map["pageNext"]

				pages, _ := JSONToPageInfoSlice(*modelValues.Map["pages"])
				pages = append(pages, currentPage)
				pagesString, _ := PageInfoSliceToJSON(pages)
				modelValues.Map["pages"] = &pagesString
				// The branch callback cannot see the model values, so the pages ride on the form too
				mirror := pagesString
				formValues.Map["pages"] = &mirror

				// The accepted pages live on the command so the accept summary and save path see them
				currentCommand, _ := JSONToCmdInfo(*modelValues.Map["currentCommand"])
				currentCommand.Pages = pages
				commandString, _ := currentCommand.ToJSON()
				modelValues.Map["currentCommand"] = &commandString

				// The next page starts with clean page forms
				allForms[idxPageInfo].Values.Map["pageName"] = new(string)
				allForms[idxPageInfo].Values.Map["pageTitle"] = new(string)
				formValues.Map["pageNext"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				pages, _ := JSONToPageInfoSlice(*formValues.Map["pages"])
				// The flow cannot chain more pages than the cap allows
				if *formValues.Map["pageAnotherConfirm"] == "yes" && len(pages) < MaxFlowPages {
					return idxPageInfo
				}
				return idxResponseStart
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxResponseStart
		values := map[string]*string{
			"responseStartConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Response Start",
			Form: addResponseStartFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addResponseStartValues",
			},
			ShowStatus: false,
			FormGroup:  "response",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				allForms[idxResponseInfo].Values.Map["responseContent"] = new(string)
				allForms[idxResponseInfo].Values.Map["responseEphemeral"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["responseStartConfirm"] == "yes" {
					return -1
				}
				return idxAccept
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxResponseInfo
		values := map[string]*string{
			"responses":         new(string),
			"responseContent":   new(string),
			"responseEphemeral": new(string),
		}
		wrapper := FormWrapper{
			Name: "Add Response Info",
			Form: addResponseInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "addResponseInfoValues",
			},
			ShowStatus: false,
			FormGroup:  "response",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				currentCommand, _ := JSONToCmdInfo(*modelValues.Map["currentCommand"])

				// Plain messages are the only response type that exists today
				currentCommand.Responses = append(currentCommand.Responses, ResponseInfo{
					Type:      "message",
					Content:   *formValues.Map["responseContent"],
					Ephemeral: *formValues.Map["responseEphemeral"] == "yes",
				})
				responseString, _ := ResponseInfoSliceToJSON(currentCommand.Responses)
				formValues.Map["responses"] = &responseString
				commandString, _ := currentCommand.ToJSON()
				modelValues.Map["currentCommand"] = &commandString
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				// A command cannot declare more responses than the cap allows
				responses, _ := JSONToResponseInfoSlice(*formValues.Map["responses"])
				if len(responses) >= MaxCommandResponses {
					return idxAccept
				}
				return idxResponseStart
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

// buildCommandSummary renders the accept screen text for a fully collected command
func buildCommandSummary(command CommandInfo) string {
	commandArgs := "None"
	if len(command.Args) > 0 {
		argNames := make([]string, len(command.Args))
		for i, arg := range command.Args {
			argNames[i] = fmt.Sprintf("%s (%s)", arg.Name, arg.Type)
		}
		commandArgs = strings.Join(argNames, ", ")
	}

	commandFields := "None"
	if len(command.Fields) > 0 {
		fieldNames := make([]string, len(command.Fields))
		for i, field := range command.Fields {
			fieldNames[i] = fmt.Sprintf("%s (%s)", field.Name, field.Style)
		}
		commandFields = strings.Join(fieldNames, ", ")
	}

	summary := fmt.Sprintf("Command Name: %s\nCommand Type: %s\nDescription: %s\nReturn Type: %s\nArguments: %v",
		command.Name, command.Type, command.Description, command.ReturnType, commandArgs)
	if command.Type == "modal" {
		summary = fmt.Sprintf("Command Name: %s\nCommand Type: %s\nDescription: %s\nReturn Type: %s\nFields: %v",
			command.Name, command.Type, command.Description, command.ReturnType, commandFields)
		// A multi page modal lists its pages instead of a flat field list
		if len(command.Pages) > 0 {
			pageLines := make([]string, len(command.Pages))
			for i, page := range command.Pages {
				next := page.Next
				if next == "" {
					next = "end"
				}
				pageLines[i] = fmt.Sprintf("  %s (%d fields, %d branches, next: %s)", page.Name, len(page.Fields), len(page.Branches), next)
			}
			summary = fmt.Sprintf("Command Name: %s\nCommand Type: %s\nDescription: %s\nReturn Type: %s\nPages:\n%s",
				command.Name, command.Type, command.Description, command.ReturnType, strings.Join(pageLines, "\n"))
		}
	}

	if len(command.Responses) > 0 {
		responseLines := make([]string, len(command.Responses))
		for i, response := range command.Responses {
			marker := ""
			if response.Ephemeral {
				marker = ", ephemeral"
			}
			responseLines[i] = fmt.Sprintf("  %s: %s%s", response.Type, response.Content, marker)
		}
		summary += "\nResponses:\n" + strings.Join(responseLines, "\n")
	}

	return summary
}

func addCmdAcceptFormGenerator(values Values, modelValues Values) *huh.Form {
	var summary string
	if modelValues.Map["currentCommand"] != nil && *modelValues.Map["currentCommand"] != "" {
		currentCommand, err := JSONToCmdInfo(*modelValues.Map["currentCommand"])
		if err == nil {
			summary = buildCommandSummary(*currentCommand)
		}
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
					// Accepting runs the full validation, a rejected command routes back like a plain no
					if b {
						if err := validateAcceptedCommand(modelValues); err != nil {
							return err
						}
					}
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
	return fieldInfoFormGenerator(values, "fields")
}

func addPageFieldInfoFormGenerator(values Values, modelValues Values) *huh.Form {
	return fieldInfoFormGenerator(values, "pageFields")
}

// fieldInfoFormGenerator builds the shared field form, fieldsKey names the accumulator on the wrapper
func fieldInfoFormGenerator(values Values, fieldsKey string) *huh.Form {
	fieldInfoForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(values.Map["fieldName"]).
				Title("Enter the field name").
				Prompt("> ").
				Validate(func(s string) error {
					fields, _ := JSONToFieldInfoSlice(*values.Map[fieldsKey])
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

// pageToJSON moves a single page onto the string value bus
func pageToJSON(page PageInfo) (string, error) {
	jsonData, err := json.Marshal(page)
	if err != nil {
		return "", fmt.Errorf("failed to marshal PageInfo to JSON: %w", err)
	}
	return string(jsonData), nil
}

// jsonToPage reads a single page back off the string value bus
func jsonToPage(jsonString string) (PageInfo, error) {
	var page PageInfo
	err := json.Unmarshal([]byte(jsonString), &page)
	if err != nil {
		return PageInfo{}, fmt.Errorf("failed to unmarshal JSON to PageInfo: %w", err)
	}
	return page, nil
}

// validatePageName checks a new page name against the pages collected so far
func validatePageName(s string, pages []PageInfo) error {
	if s == "" {
		return fmt.Errorf("page name cannot be empty")
	}
	if strings.Contains(s, " ") {
		return fmt.Errorf("page name cannot contain spaces")
	}
	for _, page := range pages {
		if page.Name == s {
			return fmt.Errorf("page name already exists")
		}
	}
	return nil
}

// validatePageReference checks a page name typed as a branch goto or default next target
func validatePageReference(s string) error {
	if strings.Contains(s, " ") {
		return fmt.Errorf("page name cannot contain spaces")
	}
	return nil
}

// validateResponseContent mirrors the ValidateResponses rules for a single content string
func validateResponseContent(s string) error {
	if s == "" {
		return fmt.Errorf("response content cannot be empty")
	}
	// The content lands inside a python string literal so these characters would break the generated file
	if strings.ContainsAny(s, "\"\\\n") {
		return fmt.Errorf("response content cannot contain double quotes, backslashes, or newlines")
	}
	return nil
}

// validateAcceptedCommand runs the full command validation against the commands accepted so far
func validateAcceptedCommand(modelValues Values) error {
	if modelValues.Map["currentCommand"] == nil || *modelValues.Map["currentCommand"] == "" {
		return fmt.Errorf("there is no command to validate")
	}
	command, err := JSONToCmdInfo(*modelValues.Map["currentCommand"])
	if err != nil {
		return fmt.Errorf("failed to read the current command: %w", err)
	}

	var existing []CommandInfo
	if modelValues.Map["slashCommands"] != nil {
		slashCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["slashCommands"])
		existing = append(existing, slashCommandList...)
	}
	if modelValues.Map["prefixCommands"] != nil {
		prefixCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["prefixCommands"])
		existing = append(existing, prefixCommandList...)
	}
	return ValidateCommand(*command, existing)
}

func addMultiPageFormGenerator(values Values, modelValues Values) *huh.Form {
	multiPageForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want the modal to use multiple pages?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["multiPageConfirm"] = &s
					return nil
				}),
		),
	)
	return multiPageForm
}

func addPageInfoFormGenerator(values Values, modelValues Values) *huh.Form {
	pageInfoForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(values.Map["pageName"]).
				Title("Enter the page name").
				Prompt("> ").
				Validate(func(s string) error {
					var pages []PageInfo
					if modelValues.Map["pages"] != nil {
						pages, _ = JSONToPageInfoSlice(*modelValues.Map["pages"])
					}
					return validatePageName(s, pages)
				}),
			huh.NewInput().
				Value(values.Map["pageTitle"]).
				Title("Enter the page title").
				Prompt("> ").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("page title cannot be empty")
					}
					return nil
				}),
		),
	)
	return pageInfoForm
}

func addPageFieldStartFormGenerator(values Values, modelValues Values) *huh.Form {
	pageFieldStartForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to add a field to this page?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["pageFieldStartConfirm"] = &s
					return nil
				}),
		),
	)
	return pageFieldStartForm
}

func addBranchStartFormGenerator(values Values, modelValues Values) *huh.Form {
	branchStartForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to add a branch rule to this page?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["branchStartConfirm"] = &s
					return nil
				}),
		),
	)
	return branchStartForm
}

func addBranchInfoFormGenerator(values Values, modelValues Values) *huh.Form {
	var currentPage PageInfo
	if modelValues.Map["currentPage"] != nil {
		currentPage, _ = jsonToPage(*modelValues.Map["currentPage"])
	}

	// The branch can only test a field the user already added to this page
	fieldOptions := make([]huh.Option[string], 0, len(currentPage.Fields))
	for _, field := range currentPage.Fields {
		fieldOptions = append(fieldOptions, huh.NewOption(field.Name, field.Name))
	}

	var fieldPicker huh.Field
	if len(fieldOptions) > 0 {
		fieldPicker = huh.NewSelect[string]().
			Value(values.Map["branchField"]).
			Title("Select the field this branch tests").
			Options(fieldOptions...)
	} else {
		// The branch start guard makes this unreachable in the normal flow, keep a fallback anyway
		fieldPicker = huh.NewInput().
			Value(values.Map["branchField"]).
			Title("Enter the field this branch tests").
			Prompt("> ")
	}

	branchInfoForm := huh.NewForm(
		huh.NewGroup(
			fieldPicker,
			huh.NewInput().
				Value(values.Map["branchEquals"]).
				Title("Enter the value that triggers this branch").
				Prompt("> "),
			huh.NewInput().
				Value(values.Map["branchGoto"]).
				Title("Enter the page to jump to when the value matches").
				Prompt("> ").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("branch goto cannot be empty")
					}
					// The target page may not exist yet, the accept step validates the full flow
					return validatePageReference(s)
				}),
		),
	)
	return branchInfoForm
}

func addPageNextFormGenerator(values Values, modelValues Values) *huh.Form {
	pageNextForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(values.Map["pageNext"]).
				Title("Enter the default next page (leave empty to end the flow)").
				Prompt("> ").
				Validate(validatePageReference),
			huh.NewConfirm().
				Title("Do you want to add another page?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["pageAnotherConfirm"] = &s
					return nil
				}),
		),
	)
	return pageNextForm
}

func addResponseStartFormGenerator(values Values, modelValues Values) *huh.Form {
	responseStartForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to add a custom response?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["responseStartConfirm"] = &s
					return nil
				}),
		),
	)
	return responseStartForm
}

func addResponseInfoFormGenerator(values Values, modelValues Values) *huh.Form {
	responseInfoForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(values.Map["responseContent"]).
				Title("Enter the response message content").
				Prompt("> ").
				Validate(validateResponseContent),
			huh.NewConfirm().
				Title("Should the response be ephemeral?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["responseEphemeral"] = &s
					return nil
				}),
		),
	)
	return responseInfoForm
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
 * Edit Forms and Model Generators
 */
func EditFormWrapperGenerator() []FormWrapper {
	// Positions of each wrapper in the returned slice, referenced by callbacks and branches
	const (
		idxEditSelectCog = iota
		idxEditAction
		idxEditCmdInfo
		idxEditArgStart
		idxEditArgInfo
		idxEditFieldStart
		idxEditFieldInfo
		idxEditAccept
		idxEditMultiPage
		idxEditPageInfo
		idxEditPageFieldStart
		idxEditPageFieldInfo
		idxEditBranchStart
		idxEditBranchInfo
		idxEditPageNext
		idxEditModInfo
		idxEditRedefine
		idxEditRedefineResponses
		idxEditResponseStart
		idxEditResponseInfo
		idxEditPickCommand
		idxEditRemoveCommand
	)

	// resetCommandState clears every per command form so a new command flow starts clean
	resetCommandState := func(modelValues Values, allForms []FormWrapper) {
		allForms[idxEditCmdInfo].Values.Map["cmdName"] = new(string)
		allForms[idxEditCmdInfo].Values.Map["cmdType"] = new(string)
		allForms[idxEditCmdInfo].Values.Map["cmdScope"] = new(string)
		allForms[idxEditCmdInfo].Values.Map["cmdDescription"] = new(string)
		allForms[idxEditCmdInfo].Values.Map["cmdReturnType"] = new(string)
		allForms[idxEditArgInfo].Values.Map["args"] = new(string)
		allForms[idxEditFieldInfo].Values.Map["fields"] = new(string)
		allForms[idxEditPageInfo].Values.Map["pageName"] = new(string)
		allForms[idxEditPageInfo].Values.Map["pageTitle"] = new(string)
		allForms[idxEditPageFieldInfo].Values.Map["pageFields"] = new(string)
		allForms[idxEditPageNext].Values.Map["pages"] = new(string)
		allForms[idxEditResponseInfo].Values.Map["responses"] = new(string)
		emptyPages := "[]"
		modelValues.Map["pages"] = &emptyPages
		modelValues.Map["currentPage"] = new(string)
		modelValues.Map["currentCommand"] = new(string)
		modelValues.Map["editingOriginal"] = new(string)
	}

	forms := []FormWrapper{}
	{ // NOTE: idxEditSelectCog
		values := map[string]*string{
			"cogName": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Select Cog",
			Form: editSelectCogFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editSelectCogValues",
			},
			ShowStatus: false,
			FormGroup:  "cog",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				if formValues.Map["cogName"] == nil || *formValues.Map["cogName"] == "" {
					return
				}
				*modelValues.Map["cogName"] = *formValues.Map["cogName"]
				editLoadCogCommands(modelValues)
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
	{ // NOTE: idxEditAction
		values := map[string]*string{
			"editAction": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Action",
			Form: editActionFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editActionValues",
			},
			ShowStatus: false,
			FormGroup:  "action",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				// A fresh add starts with a clean command flow, edit prefills it later
				if *formValues.Map["editAction"] == "add" {
					resetCommandState(modelValues, allForms)
				}
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				switch *formValues.Map["editAction"] {
				case "add":
					return idxEditCmdInfo
				case "edit":
					return idxEditPickCommand
				case "remove":
					return idxEditRemoveCommand
				default:
					return -2
				}
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditCmdInfo
		values := map[string]*string{
			"cmdName":        new(string),
			"cmdType":        new(string),
			"cmdScope":       new(string),
			"cmdDescription": new(string),
			"cmdReturnType":  new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Add Command Info",
			Form: addCmdInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editAddCommandInfoValues",
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
					Pages:       []PageInfo{},
					Responses:   []ResponseInfo{},
					ReturnType:  returnType,
				}
				commandString, _ := command.ToJSON()
				modelValues.Map["currentCommand"] = &commandString
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				// Modal commands first decide between a single page and a multi page flow
				if *formValues.Map["cmdType"] == "modal" {
					return idxEditMultiPage
				}
				return -1
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditArgStart
		values := map[string]*string{
			"argStartConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Argument Start",
			Form: addArgStartFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editArgumentStartValues",
			},
			ShowStatus: false,
			FormGroup:  "argument",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				allForms[idxEditArgInfo].Values.Map["argName"] = new(string)
				allForms[idxEditArgInfo].Values.Map["argDescription"] = new(string)
				allForms[idxEditArgInfo].Values.Map["argType"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["argStartConfirm"] == "yes" {
					return -1
				}
				return idxEditRedefineResponses
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditArgInfo
		values := map[string]*string{
			"args":           new(string),
			"argName":        new(string),
			"argDescription": new(string),
			"argType":        new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Argument Info",
			Form: addArgInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editArgumentInfoValues",
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
				return idxEditArgStart
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditFieldStart
		values := map[string]*string{
			"fieldStartConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Field Start",
			Form: addFieldStartFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editFieldStartValues",
			},
			ShowStatus: false,
			FormGroup:  "field",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				allForms[idxEditFieldInfo].Values.Map["fieldName"] = new(string)
				allForms[idxEditFieldInfo].Values.Map["fieldLabel"] = new(string)
				allForms[idxEditFieldInfo].Values.Map["fieldStyle"] = new(string)
				allForms[idxEditFieldInfo].Values.Map["fieldRequired"] = new(string)
				allForms[idxEditFieldInfo].Values.Map["fieldPlaceholder"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["fieldStartConfirm"] == "yes" {
					return -1
				}
				return idxEditRedefineResponses
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditFieldInfo
		values := map[string]*string{
			"fields":           new(string),
			"fieldName":        new(string),
			"fieldLabel":       new(string),
			"fieldStyle":       new(string),
			"fieldRequired":    new(string),
			"fieldPlaceholder": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Field Info",
			Form: addFieldInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editFieldInfoValues",
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
					return idxEditRedefineResponses
				}
				return idxEditFieldStart
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditAccept
		values := map[string]*string{
			"cmdAcceptConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Command Accept",
			Form: addCmdAcceptFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editCommandAcceptValues",
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
				} else if modelValues.Map["editingOriginal"] != nil && *modelValues.Map["editingOriginal"] != "" {
					// A rejected edit restores the untouched original command
					original, err := JSONToCmdInfo(*modelValues.Map["editingOriginal"])
					if err == nil {
						if original.Type == "prefix" {
							prefixCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["prefixCommands"])
							prefixCommandList = append(prefixCommandList, *original)
							jsonData, _ := CmdInfoSliceToJSON(prefixCommandList)
							modelValues.Map["prefixCommands"] = &jsonData
						} else {
							slashCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["slashCommands"])
							slashCommandList = append(slashCommandList, *original)
							jsonData, _ := CmdInfoSliceToJSON(slashCommandList)
							modelValues.Map["slashCommands"] = &jsonData
						}
					}
				}
				// The editing marker never outlives the accept decision
				modelValues.Map["editingOriginal"] = new(string)
			},
			BranchCallback: func(values Values, allForms []FormWrapper) int {
				return idxEditAction
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditMultiPage
		values := map[string]*string{
			"multiPageConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Multi Page",
			Form: addMultiPageFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editMultiPageValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				// A fresh modal command always starts with an empty page collection
				emptyPages := "[]"
				modelValues.Map["pages"] = &emptyPages
				mirror := "[]"
				allForms[idxEditPageNext].Values.Map["pages"] = &mirror
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["multiPageConfirm"] == "yes" {
					return idxEditPageInfo
				}
				// A single page modal keeps the existing field loop
				return idxEditFieldStart
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditPageInfo
		values := map[string]*string{
			"pageName":  new(string),
			"pageTitle": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Page Info",
			Form: addPageInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editPageInfoValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				page := PageInfo{
					Name:     *formValues.Map["pageName"],
					Title:    *formValues.Map["pageTitle"],
					Fields:   []FieldInfo{},
					Branches: []BranchRule{},
				}
				pageString, _ := pageToJSON(page)
				modelValues.Map["currentPage"] = &pageString
				// Each page collects its own fields, so the accumulator starts empty
				allForms[idxEditPageFieldInfo].Values.Map["pageFields"] = new(string)
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditPageFieldStart
		values := map[string]*string{
			"pageFieldStartConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Page Field Start",
			Form: addPageFieldStartFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editPageFieldStartValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				allForms[idxEditPageFieldInfo].Values.Map["fieldName"] = new(string)
				allForms[idxEditPageFieldInfo].Values.Map["fieldLabel"] = new(string)
				allForms[idxEditPageFieldInfo].Values.Map["fieldStyle"] = new(string)
				allForms[idxEditPageFieldInfo].Values.Map["fieldRequired"] = new(string)
				allForms[idxEditPageFieldInfo].Values.Map["fieldPlaceholder"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["pageFieldStartConfirm"] == "yes" {
					return -1
				}
				return idxEditBranchStart
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditPageFieldInfo
		values := map[string]*string{
			"pageFields":       new(string),
			"fieldName":        new(string),
			"fieldLabel":       new(string),
			"fieldStyle":       new(string),
			"fieldRequired":    new(string),
			"fieldPlaceholder": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Page Field Info",
			Form: addPageFieldInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editPageFieldInfoValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				currentPage, _ := jsonToPage(*modelValues.Map["currentPage"])

				currentPage.Fields = append(currentPage.Fields, FieldInfo{
					Name:        *formValues.Map["fieldName"],
					Label:       *formValues.Map["fieldLabel"],
					Style:       *formValues.Map["fieldStyle"],
					Required:    *formValues.Map["fieldRequired"] == "yes",
					Placeholder: *formValues.Map["fieldPlaceholder"],
				})
				fieldString, _ := FieldInfoSliceToJSON(currentPage.Fields)
				formValues.Map["pageFields"] = &fieldString
				pageString, _ := pageToJSON(currentPage)
				modelValues.Map["currentPage"] = &pageString
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				// A modal page cannot hold more inputs than Discord allows
				fields, _ := JSONToFieldInfoSlice(*formValues.Map["pageFields"])
				if len(fields) >= MaxModalFields {
					return idxEditBranchStart
				}
				return idxEditPageFieldStart
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditBranchStart
		values := map[string]*string{
			"branchStartConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Branch Start",
			Form: addBranchStartFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editBranchStartValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				allForms[idxEditBranchInfo].Values.Map["branchField"] = new(string)
				allForms[idxEditBranchInfo].Values.Map["branchEquals"] = new(string)
				allForms[idxEditBranchInfo].Values.Map["branchGoto"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["branchStartConfirm"] == "yes" {
					// Branch rules test a field on this page, so an empty page has nothing to branch on
					fields, _ := JSONToFieldInfoSlice(*allForms[idxEditPageFieldInfo].Values.Map["pageFields"])
					if len(fields) == 0 {
						return idxEditPageNext
					}
					return -1
				}
				return idxEditPageNext
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditBranchInfo
		values := map[string]*string{
			"branchField":  new(string),
			"branchEquals": new(string),
			"branchGoto":   new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Branch Info",
			Form: addBranchInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editBranchInfoValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				currentPage, _ := jsonToPage(*modelValues.Map["currentPage"])

				currentPage.Branches = append(currentPage.Branches, BranchRule{
					Field:  *formValues.Map["branchField"],
					Equals: *formValues.Map["branchEquals"],
					Goto:   *formValues.Map["branchGoto"],
				})
				pageString, _ := pageToJSON(currentPage)
				modelValues.Map["currentPage"] = &pageString
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				return idxEditBranchStart
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditPageNext
		values := map[string]*string{
			"pages":              new(string),
			"pageNext":           new(string),
			"pageAnotherConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Page Next",
			Form: addPageNextFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editPageNextValues",
			},
			ShowStatus: false,
			FormGroup:  "page",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				currentPage, _ := jsonToPage(*modelValues.Map["currentPage"])
				currentPage.Next = *formValues.Map["pageNext"]

				pages, _ := JSONToPageInfoSlice(*modelValues.Map["pages"])
				pages = append(pages, currentPage)
				pagesString, _ := PageInfoSliceToJSON(pages)
				modelValues.Map["pages"] = &pagesString
				// The branch callback cannot see the model values, so the pages ride on the form too
				mirror := pagesString
				formValues.Map["pages"] = &mirror

				// The accepted pages live on the command so the accept summary and save path see them
				currentCommand, _ := JSONToCmdInfo(*modelValues.Map["currentCommand"])
				currentCommand.Pages = pages
				commandString, _ := currentCommand.ToJSON()
				modelValues.Map["currentCommand"] = &commandString

				// The next page starts with clean page forms
				allForms[idxEditPageInfo].Values.Map["pageName"] = new(string)
				allForms[idxEditPageInfo].Values.Map["pageTitle"] = new(string)
				formValues.Map["pageNext"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				pages, _ := JSONToPageInfoSlice(*formValues.Map["pages"])
				// The flow cannot chain more pages than the cap allows
				if *formValues.Map["pageAnotherConfirm"] == "yes" && len(pages) < MaxFlowPages {
					return idxEditPageInfo
				}
				return idxEditRedefineResponses
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditModInfo
		values := map[string]*string{
			"cmdName":        new(string),
			"cmdType":        new(string),
			"cmdScope":       new(string),
			"cmdDescription": new(string),
			"cmdReturnType":  new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Command Info",
			Form: addCmdInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editCommandInfoValues",
			},
			ShowStatus: false,
			FormGroup:  "command",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				currentCommand, err := JSONToCmdInfo(*modelValues.Map["currentCommand"])
				if err != nil {
					return
				}
				// Only the info fields change here, the collected sets stay until redefined
				currentCommand.Name = *formValues.Map["cmdName"]
				currentCommand.Type = *formValues.Map["cmdType"]
				currentCommand.Scope = *formValues.Map["cmdScope"]
				currentCommand.Description = *formValues.Map["cmdDescription"]
				// Modal commands only respond through the modal, so their return type is fixed
				returnType := *formValues.Map["cmdReturnType"]
				if currentCommand.Type == "modal" {
					returnType = "None"
				}
				currentCommand.ReturnType = returnType
				commandString, _ := currentCommand.ToJSON()
				modelValues.Map["currentCommand"] = &commandString
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				return idxEditRedefine
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditRedefine
		values := map[string]*string{
			"redefineConfirm": new(string),
			"cmdType":         new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Redefine",
			Form: editRedefineFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editRedefineValues",
			},
			ShowStatus: false,
			FormGroup:  "command",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				currentCommand, err := JSONToCmdInfo(*modelValues.Map["currentCommand"])
				if err != nil {
					return
				}
				// The branch callback cannot see the model values, so the type rides on the form
				commandType := currentCommand.Type
				formValues.Map["cmdType"] = &commandType

				if *formValues.Map["redefineConfirm"] != "yes" {
					return
				}
				// Redefining drops the collected sets and clears every accumulator
				currentCommand.Args = []ArgInfo{}
				currentCommand.Fields = []FieldInfo{}
				currentCommand.Pages = []PageInfo{}
				commandString, _ := currentCommand.ToJSON()
				modelValues.Map["currentCommand"] = &commandString

				allForms[idxEditArgInfo].Values.Map["args"] = new(string)
				allForms[idxEditFieldInfo].Values.Map["fields"] = new(string)
				allForms[idxEditPageInfo].Values.Map["pageName"] = new(string)
				allForms[idxEditPageInfo].Values.Map["pageTitle"] = new(string)
				allForms[idxEditPageFieldInfo].Values.Map["pageFields"] = new(string)
				allForms[idxEditPageNext].Values.Map["pages"] = new(string)
				emptyPages := "[]"
				modelValues.Map["pages"] = &emptyPages
				modelValues.Map["currentPage"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["redefineConfirm"] == "yes" {
					if formValues.Map["cmdType"] != nil && *formValues.Map["cmdType"] == "modal" {
						return idxEditMultiPage
					}
					return idxEditArgStart
				}
				return idxEditRedefineResponses
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditRedefineResponses
		values := map[string]*string{
			"redefineResponsesConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Redefine Responses",
			Form: editRedefineResponsesFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editRedefineResponsesValues",
			},
			ShowStatus: false,
			FormGroup:  "response",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				if *formValues.Map["redefineResponsesConfirm"] != "yes" {
					return
				}
				currentCommand, err := JSONToCmdInfo(*modelValues.Map["currentCommand"])
				if err != nil {
					return
				}
				currentCommand.Responses = []ResponseInfo{}
				commandString, _ := currentCommand.ToJSON()
				modelValues.Map["currentCommand"] = &commandString
				allForms[idxEditResponseInfo].Values.Map["responses"] = new(string)
			},
			// Only an edited command with kept responses has anything to redefine
			SkipCondition: func(modelValues Values, allForms []FormWrapper, currentIndex int) bool {
				if modelValues.Map["editingOriginal"] == nil || *modelValues.Map["editingOriginal"] == "" {
					return true
				}
				if modelValues.Map["currentCommand"] == nil || *modelValues.Map["currentCommand"] == "" {
					return true
				}
				currentCommand, err := JSONToCmdInfo(*modelValues.Map["currentCommand"])
				if err != nil {
					return true
				}
				return len(currentCommand.Responses) == 0
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["redefineResponsesConfirm"] == "yes" {
					return -1
				}
				return idxEditAccept
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditResponseStart
		values := map[string]*string{
			"responseStartConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Response Start",
			Form: addResponseStartFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editResponseStartValues",
			},
			ShowStatus: false,
			FormGroup:  "response",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				allForms[idxEditResponseInfo].Values.Map["responseContent"] = new(string)
				allForms[idxEditResponseInfo].Values.Map["responseEphemeral"] = new(string)
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if *formValues.Map["responseStartConfirm"] == "yes" {
					return -1
				}
				return idxEditAccept
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditResponseInfo
		values := map[string]*string{
			"responses":         new(string),
			"responseContent":   new(string),
			"responseEphemeral": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Response Info",
			Form: addResponseInfoFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editResponseInfoValues",
			},
			ShowStatus: false,
			FormGroup:  "response",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				currentCommand, _ := JSONToCmdInfo(*modelValues.Map["currentCommand"])

				// Plain messages are the only response type that exists today
				currentCommand.Responses = append(currentCommand.Responses, ResponseInfo{
					Type:      "message",
					Content:   *formValues.Map["responseContent"],
					Ephemeral: *formValues.Map["responseEphemeral"] == "yes",
				})
				responseString, _ := ResponseInfoSliceToJSON(currentCommand.Responses)
				formValues.Map["responses"] = &responseString
				commandString, _ := currentCommand.ToJSON()
				modelValues.Map["currentCommand"] = &commandString
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				// A command cannot declare more responses than the cap allows
				responses, _ := JSONToResponseInfoSlice(*formValues.Map["responses"])
				if len(responses) >= MaxCommandResponses {
					return idxEditAccept
				}
				return idxEditResponseStart
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditPickCommand
		values := map[string]*string{
			"editCmdName": new(string),
			"editFound":   new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Pick Command",
			Form: editPickCommandFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editPickCommandValues",
			},
			ShowStatus: false,
			FormGroup:  "command",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				found := "no"
				formValues.Map["editFound"] = &found

				name := ""
				if formValues.Map["editCmdName"] != nil {
					name = *formValues.Map["editCmdName"]
				}
				if name == "" {
					return
				}

				slashCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["slashCommands"])
				prefixCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["prefixCommands"])

				// The picked command leaves its list so name checks run against the rest
				var command *CommandInfo
				for i, candidate := range slashCommandList {
					if candidate.Name == name {
						picked := candidate
						command = &picked
						slashCommandList = append(slashCommandList[:i], slashCommandList[i+1:]...)
						break
					}
				}
				if command == nil {
					for i, candidate := range prefixCommandList {
						if candidate.Name == name {
							picked := candidate
							command = &picked
							prefixCommandList = append(prefixCommandList[:i], prefixCommandList[i+1:]...)
							break
						}
					}
				}
				if command == nil {
					return
				}

				slashJSON, _ := CmdInfoSliceToJSON(slashCommandList)
				modelValues.Map["slashCommands"] = &slashJSON
				prefixJSON, _ := CmdInfoSliceToJSON(prefixCommandList)
				modelValues.Map["prefixCommands"] = &prefixJSON

				commandString, _ := command.ToJSON()
				modelValues.Map["currentCommand"] = &commandString
				originalString := commandString
				modelValues.Map["editingOriginal"] = &originalString

				// The info form starts prefilled with the picked command's values
				cmdName := command.Name
				cmdType := command.Type
				cmdScope := command.Scope
				cmdDescription := command.Description
				cmdReturnType := command.ReturnType
				allForms[idxEditModInfo].Values.Map["cmdName"] = &cmdName
				allForms[idxEditModInfo].Values.Map["cmdType"] = &cmdType
				allForms[idxEditModInfo].Values.Map["cmdScope"] = &cmdScope
				allForms[idxEditModInfo].Values.Map["cmdDescription"] = &cmdDescription
				allForms[idxEditModInfo].Values.Map["cmdReturnType"] = &cmdReturnType

				yes := "yes"
				formValues.Map["editFound"] = &yes
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				if formValues.Map["editFound"] != nil && *formValues.Map["editFound"] == "yes" {
					return idxEditModInfo
				}
				return idxEditAction
			},
		}
		forms = append(forms, wrapper)
	}
	{ // NOTE: idxEditRemoveCommand
		values := map[string]*string{
			"removeCmdName": new(string),
			"removeConfirm": new(string),
		}
		wrapper := FormWrapper{
			Name: "Edit Remove Command",
			Form: editRemoveCommandFormGenerator,
			Values: Values{
				Map:  values,
				Name: "editRemoveCommandValues",
			},
			ShowStatus: false,
			FormGroup:  "command",
			Callback: func(formValues Values, modelValues Values, allForms []FormWrapper) {
				if *formValues.Map["removeConfirm"] != "yes" {
					return
				}
				name := ""
				if formValues.Map["removeCmdName"] != nil {
					name = *formValues.Map["removeCmdName"]
				}
				if name == "" {
					return
				}

				slashCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["slashCommands"])
				for i, candidate := range slashCommandList {
					if candidate.Name == name {
						slashCommandList = append(slashCommandList[:i], slashCommandList[i+1:]...)
						slashJSON, _ := CmdInfoSliceToJSON(slashCommandList)
						modelValues.Map["slashCommands"] = &slashJSON
						return
					}
				}
				prefixCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["prefixCommands"])
				for i, candidate := range prefixCommandList {
					if candidate.Name == name {
						prefixCommandList = append(prefixCommandList[:i], prefixCommandList[i+1:]...)
						prefixJSON, _ := CmdInfoSliceToJSON(prefixCommandList)
						modelValues.Map["prefixCommands"] = &prefixJSON
						return
					}
				}
			},
			BranchCallback: func(formValues Values, allForms []FormWrapper) int {
				return idxEditAction
			},
		}
		forms = append(forms, wrapper)
	}

	return forms
}

// editLoadCogCommands fills the model value bus with the selected cog's command lists
func editLoadCogCommands(modelValues Values) {
	config, err := LoadConfig()
	if err != nil {
		return
	}
	for _, cog := range config.Cogs {
		if cog.Name != *modelValues.Map["cogName"] {
			continue
		}
		slashCommands := cog.SlashCommands
		if slashCommands == nil {
			slashCommands = []CommandInfo{}
		}
		prefixCommands := cog.PrefixCommands
		if prefixCommands == nil {
			prefixCommands = []CommandInfo{}
		}
		slashJSON, _ := CmdInfoSliceToJSON(slashCommands)
		modelValues.Map["slashCommands"] = &slashJSON
		prefixJSON, _ := CmdInfoSliceToJSON(prefixCommands)
		modelValues.Map["prefixCommands"] = &prefixJSON
		return
	}
}

// editCommandNames lists every command name on the model value bus, slash commands first
func editCommandNames(modelValues Values) []string {
	var names []string
	if modelValues.Map["slashCommands"] != nil {
		slashCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["slashCommands"])
		for _, command := range slashCommandList {
			names = append(names, command.Name)
		}
	}
	if modelValues.Map["prefixCommands"] != nil {
		prefixCommandList, _ := JSONToCmdInfoSlice(*modelValues.Map["prefixCommands"])
		for _, command := range prefixCommandList {
			names = append(names, command.Name)
		}
	}
	return names
}

// noCommandsFormGenerator builds the note shown when a cog has no commands to pick from
func noCommandsFormGenerator() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("No Commands").
				Description("This cog has no commands yet. Add a command first."),
		),
	)
}

func editSelectCogFormGenerator(values Values, modelValues Values) *huh.Form {
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
					Description("There are no cogs to edit. Please add some cogs first."),
			),
		)
		return noCogForm
	}

	cogSelectForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(values.Map["cogName"]).
				Height(8).
				Title("Select a cog to edit").
				Options(huh.NewOptions(cogList...)...),
		),
	)
	return cogSelectForm
}

func editActionFormGenerator(values Values, modelValues Values) *huh.Form {
	actionForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(values.Map["editAction"]).
				Title("What do you want to do?").
				Options(
					huh.NewOption("Add a command", "add"),
					huh.NewOption("Edit a command", "edit"),
					huh.NewOption("Remove a command", "remove"),
					huh.NewOption("Apply changes", "apply"),
				),
		),
	)
	return actionForm
}

func editPickCommandFormGenerator(values Values, modelValues Values) *huh.Form {
	names := editCommandNames(modelValues)
	if len(names) == 0 {
		return noCommandsFormGenerator()
	}

	pickForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(values.Map["editCmdName"]).
				Height(8).
				Title("Select a command to edit").
				Options(huh.NewOptions(names...)...),
		),
	)
	return pickForm
}

func editRemoveCommandFormGenerator(values Values, modelValues Values) *huh.Form {
	names := editCommandNames(modelValues)
	if len(names) == 0 {
		return noCommandsFormGenerator()
	}

	removeForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(values.Map["removeCmdName"]).
				Height(8).
				Title("Select a command to remove").
				Options(huh.NewOptions(names...)...),
			huh.NewConfirm().
				Title("Remove this command?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["removeConfirm"] = &s
					return nil
				}),
		),
	)
	return removeForm
}

func editRedefineFormGenerator(values Values, modelValues Values) *huh.Form {
	redefineForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Redefine the arguments, fields, and pages?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["redefineConfirm"] = &s
					return nil
				}),
		),
	)
	return redefineForm
}

func editRedefineResponsesFormGenerator(values Values, modelValues Values) *huh.Form {
	redefineResponsesForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Redefine the responses?").
				Affirmative("yes").
				Negative("no").
				Validate(func(b bool) error {
					var s string
					if b {
						s = "yes"
					} else {
						s = "no"
					}
					values.Map["redefineResponsesConfirm"] = &s
					return nil
				}),
		),
	)
	return redefineResponsesForm
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
