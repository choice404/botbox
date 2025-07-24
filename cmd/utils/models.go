/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func CreateModel(callback func(map[string]string)) Model {
	m := Model{width: maxWidth}
	m.title = "Create a BotBox Project"
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.currentFormPtr = 0
	m.state = 0
	m.callback = callback
	createForms := CreateFormWrapperGenerator()
	m.forms = createForms
	m.modelValues = map[string]*string{
		"botName":                new(string),
		"botDescription":         new(string),
		"botAuthor":              new(string),
		"botPrefix":              new(string),
		"envChoice":              new(string),
		"botTokenDopplerProject": new(string),
		"botGuildDopplerEnv":     new(string),
		"licenseType":            new(string),
	}
	m.displayCallback = func() string {
		s := m.styles
		var display strings.Builder
		display.WriteString(s.HeaderText.Render("Project Created") + "\n\n")
		display.WriteString(s.KeyText.Render("Bot Name: ") + s.ValueText.Render(*m.modelValues["botName"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Bot Description: ") + s.ValueText.Render(*m.modelValues["botDescription"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Project Author: ") + s.ValueText.Render(*m.modelValues["botAuthor"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Bot Prefix: ") + s.ValueText.Render(*m.modelValues["botPrefix"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Environment: ") + s.ValueText.Render(*m.modelValues["envChoice"]) + "\n")
		return display.String()
	}

	return m
}

func InitModel(callback func(map[string]string)) Model {
	m := Model{width: maxWidth}
	m.title = "Initialize a BotBox Project in the Current Directory"
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.currentFormPtr = 0
	m.state = 0
	m.callback = callback
	createForms := AddFormWrapperGenerator()
	m.forms = createForms
	m.displayCallback = func() string {
		s := m.styles
		var display strings.Builder
		display.WriteString(s.HeaderText.Render("Project Created") + "\n\n")
		display.WriteString(s.KeyText.Render("Bot Name: ") + s.ValueText.Render(*m.modelValues["botName"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Bot Description: ") + s.ValueText.Render(*m.modelValues["botDescription"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Project Author: ") + s.ValueText.Render(*m.modelValues["botAuthor"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Bot Prefix: ") + s.ValueText.Render(*m.modelValues["botPrefix"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Environment: ") + s.ValueText.Render(*m.modelValues["envChoice"]) + "\n")
		return display.String()
	}
	return m
}

func AddModel(callback func(map[string]string)) Model {
	m := Model{width: maxWidth}
	m.title = "Add a Cog"
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.currentFormPtr = 0
	m.callback = callback
	m.modelValues = map[string]*string{
		"filename":       new(string),
		"currentCommand": new(string),
		"slashCommands":  new(string),
		"prefixCommands": new(string),
	}

	emptySlice := "[]"
	m.modelValues["slashCommands"] = &emptySlice
	m.modelValues["prefixCommands"] = &emptySlice

	addForms := AddFormWrapperGenerator()

	m.displayCallback = func() string {
		s := m.styles
		var display strings.Builder
		display.WriteString(s.HeaderText.Render("Cog Added") + "\n\n")
		display.WriteString(s.KeyText.Render("Cog Name: ") + s.ValueText.Render(*m.modelValues["filename"]) + "\n")
		slashCommands, _ := JSONToCmdInfoSlice(*m.modelValues["slashCommands"])
		prefixCommands, _ := JSONToCmdInfoSlice(*m.modelValues["prefixCommands"])
		if len(slashCommands) > 0 {
			display.WriteString(s.KeyText.Render("Slash Commands:") + "\n")
			for _, slashCommand := range slashCommands {
				var args []string
				for _, command := range slashCommand.Args {
					args = append(args, command.Name+": "+command.Type)
				}
				argsStr := strings.Join(args, ", ")

				commandLine := slashCommand.Name + "(" + argsStr + ") -> " + slashCommand.ReturnType
				display.WriteString("    - " + s.ValueText.Render(commandLine) + "\n")
			}
		}
		if len(prefixCommands) > 0 {
			display.WriteString(s.KeyText.Render("Prefix Commands:") + "\n")
			for _, prefixCommand := range prefixCommands {
				var args []string
				for _, command := range prefixCommand.Args {
					args = append(args, command.Name+": "+command.Type)
				}
				argsStr := strings.Join(args, ", ")

				commandLine := prefixCommand.Name + "(" + argsStr + ") -> " + prefixCommand.ReturnType
				display.WriteString("    - " + s.ValueText.Render(commandLine) + "\n")
			}
		}
		return display.String()
	}

	m.forms = addForms

	return m
}

func RemoveModel(callback func(map[string]string)) Model {
	m := Model{width: maxWidth}
	m.title = "Remove a Cog"
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.currentFormPtr = 0
	m.state = 0
	m.callback = callback
	removeForms := RemoveFormWrapperGenerator()
	m.forms = removeForms

	m.modelValues = map[string]*string{
		"cogName": new(string),
	}

	m.displayCallback = func() string {
		s := m.styles
		var display strings.Builder
		display.WriteString(s.HeaderText.Render("Cog queued to remove") + "\n\n")
		display.WriteString(s.KeyText.Render("  - Name: ") + s.ValueText.Render(*m.modelValues["cogName"]) + "\n")
		return display.String()
	}

	return m
}

func ConfigModel(callback func(map[string]string), initCallback func(map[string]*string, []map[string]*string)) Model {
	m := Model{width: maxWidth}
	m.title = "BotBox Configuration"
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.currentFormPtr = 0
	m.state = 0
	m.callback = callback
	configForms := ConfigFormWrapperGenerator()
	m.forms = configForms
	m.initCallback = initCallback

	m.modelValues = map[string]*string{
		"botName":        new(string),
		"rootDir":        new(string),
		"botAuthor":      new(string),
		"botDescription": new(string),
		"botPrefix":      new(string),
		"cogs":           new(string),
	}

	m.displayCallback = func() string {
		s := m.styles
		var display strings.Builder
		if *m.modelValues["botName"] != "" {
			display.WriteString(s.KeyText.Render("Bot Name: ") + s.ValueText.Render(*m.modelValues["botName"]) + "\n")
		}
		if *m.modelValues["rootDir"] != "" {
			display.WriteString(s.KeyText.Render("Root Directory: ") + s.ValueText.Render(*m.modelValues["rootDir"]) + "\n")
		}
		if *m.modelValues["botAuthor"] != "" {
			display.WriteString(s.KeyText.Render("Author: ") + s.ValueText.Render(*m.modelValues["botAuthor"]) + "\n")
		}
		if *m.modelValues["botDescription"] != "" {
			display.WriteString(s.KeyText.Render("Description: ") + s.ValueText.Render(*m.modelValues["botDescription"]) + "\n")
		}
		if *m.modelValues["botPrefix"] != "" {
			display.WriteString(s.KeyText.Render("Command Prefix: ") + s.ValueText.Render(*m.modelValues["botPrefix"]) + "\n")
		}
		if *m.modelValues["cogs"] != "" {
			cogsString := *m.modelValues["cogs"]
			cogs, _ := JSONToCogConfigSlice(cogsString)

			display.WriteString(s.KeyText.Render("Cogs:") + "\n")
			for _, cog := range cogs {
				display.WriteString("  - " + s.ValueText.Render(cog.File) + "(" + s.ValueText.Render(cog.Name) + ")\n")

				if cog.Env != "" {
					display.WriteString("    " + s.KeyText.Render("Environment: ") + s.ValueText.Render(cog.Env) + "\n")
				} else {
					display.WriteString("    " + s.KeyText.Render("Environment: ") + s.ValueText.Render("None") + "\n")
				}

				slashCommands := cog.SlashCommands
				prefixCommands := cog.PrefixCommands

				if len(slashCommands) > 0 {
					display.WriteString(s.KeyText.Render("    Slash Commands:") + "\n")
					for _, slashCommand := range slashCommands {
						var args []string
						for _, command := range slashCommand.Args {
							args = append(args, command.Name+": "+command.Type)
						}
						argsStr := strings.Join(args, ", ")

						commandLine := slashCommand.Name + "(" + argsStr + ") -> " + slashCommand.ReturnType
						display.WriteString("      - " + s.ValueText.Render(commandLine) + "\n")
					}
				}
				if len(prefixCommands) > 0 {
					display.WriteString(s.KeyText.Render("    Prefix Commands:") + "\n")
					for _, prefixCommand := range prefixCommands {
						var args []string
						for _, command := range prefixCommand.Args {
							args = append(args, command.Name+": "+command.Type)
						}
						argsStr := strings.Join(args, ", ")

						commandLine := prefixCommand.Name + "(" + argsStr + ") -> " + prefixCommand.ReturnType
						display.WriteString("      - " + s.ValueText.Render(commandLine) + "\n")
					}
				}
			}
		}
		return display.String()
	}

	return m
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
