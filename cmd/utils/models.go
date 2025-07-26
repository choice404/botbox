/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func CreateModel(callback func(*Model) []error) Model {
	m := Model{width: maxWidth}
	m.title = "Create a BotBox Project"
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.currentFormPtr = 0
	m.state = 0
	m.callback = callback
	createForms := CreateFormWrapperGenerator()
	m.forms = createForms
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

	m.ModelValues = Values{
		Map:  values,
		Name: "ModelValues",
	}

	m.displayCallback = func() string {
		s := m.styles
		var display strings.Builder
		display.WriteString(s.HeaderText.Render("Project Created") + "\n\n")
		display.WriteString(s.KeyText.Render("Bot Name: ") + s.ValueText.Render(*m.ModelValues.Map["botName"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Bot Description: ") + s.ValueText.Render(*m.ModelValues.Map["botDescription"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Project Author: ") + s.ValueText.Render(*m.ModelValues.Map["botAuthor"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Bot Prefix: ") + s.ValueText.Render(*m.ModelValues.Map["botPrefix"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Environment: ") + s.ValueText.Render(*m.ModelValues.Map["envChoice"]) + "\n")
		return display.String()
	}

	return m
}

func InitModel(callback func(*Model) []error) Model {
	m := Model{width: maxWidth}
	m.title = "Initialize a BotBox Project in the Current Directory"
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.currentFormPtr = 0
	m.state = 0
	m.callback = callback
	createForms := AddFormWrapperGenerator()
	m.forms = createForms

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

	m.ModelValues = Values{
		Map:  values,
		Name: "ModelValues",
	}

	m.displayCallback = func() string {
		s := m.styles
		var display strings.Builder
		display.WriteString(s.HeaderText.Render("Project Created") + "\n\n")
		display.WriteString(s.KeyText.Render("Bot Name: ") + s.ValueText.Render(*m.ModelValues.Map["botName"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Bot Description: ") + s.ValueText.Render(*m.ModelValues.Map["botDescription"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Project Author: ") + s.ValueText.Render(*m.ModelValues.Map["botAuthor"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Bot Prefix: ") + s.ValueText.Render(*m.ModelValues.Map["botPrefix"]) + "\n")
		display.WriteString("  - " + s.KeyText.Render("Environment: ") + s.ValueText.Render(*m.ModelValues.Map["envChoice"]) + "\n")
		return display.String()
	}
	return m
}

func AddModel(callback func(*Model) []error, initCallback func(*Model, []Values)) Model {
	m := Model{width: maxWidth}
	m.title = "Add a Cog"
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.currentFormPtr = 0
	m.callback = callback
	m.initCallback = initCallback
	m.ModelValues.Map = map[string]*string{
		"filename":       new(string),
		"currentCommand": new(string),
		"slashCommands":  new(string),
		"prefixCommands": new(string),
	}

	emptySlice := "[]"
	m.ModelValues.Map["slashCommands"] = &emptySlice
	m.ModelValues.Map["prefixCommands"] = &emptySlice

	addForms := AddFormWrapperGenerator()

	m.displayCallback = func() string {
		s := m.styles
		var display strings.Builder
		display.WriteString(s.HeaderText.Render("Cog Added") + "\n\n")
		display.WriteString(s.KeyText.Render("Cog Name: ") + s.ValueText.Render(*m.ModelValues.Map["filename"]) + "\n")
		slashCommands, _ := JSONToCmdInfoSlice(*m.ModelValues.Map["slashCommands"])
		prefixCommands, _ := JSONToCmdInfoSlice(*m.ModelValues.Map["prefixCommands"])
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

func RemoveModel(callback func(*Model) []error, initCallback func(*Model, []Values)) Model {
	m := Model{width: maxWidth}
	m.title = "Remove a Cog"
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.currentFormPtr = 0
	m.state = 0
	m.callback = callback
	removeForms := RemoveFormWrapperGenerator()
	m.forms = removeForms
	m.initCallback = initCallback

	values := map[string]*string{
		"cogName": new(string),
	}

	m.ModelValues = Values{
		Map:  values,
		Name: "ModelValues",
	}

	m.displayCallback = func() string {
		s := m.styles
		var display strings.Builder
		display.WriteString(s.HeaderText.Render("Cog queued to remove") + "\n\n")
		display.WriteString(s.KeyText.Render("  - Name: ") + s.ValueText.Render(*m.ModelValues.Map["cogName"]) + "\n")
		return display.String()
	}
	return m
}

func LocalConfigModel(callback func(*Model) []error, initCallback func(*Model, []Values)) Model {
	m := Model{width: maxWidth}
	m.title = "BotBox Project Configuration"
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.currentFormPtr = 0
	m.state = 0
	m.callback = callback
	localConfigForms := ConfigFormWrapperGenerator()
	m.forms = localConfigForms
	m.initCallback = initCallback

	m.ModelValues.Map = map[string]*string{
		"name":        new(string),
		"author":      new(string),
		"description": new(string),
		"prefix":      new(string),
		"cogs":        new(string),
	}

	m.displayCallback = func() string {
		s := m.styles
		var display strings.Builder
		if *m.ModelValues.Map["name"] != "" {
			display.WriteString(s.KeyText.Render("Bot Name: ") + s.ValueText.Render(*m.ModelValues.Map["name"]) + "\n")
		}
		if *m.ModelValues.Map["author"] != "" {
			display.WriteString(s.KeyText.Render("Author: ") + s.ValueText.Render(*m.ModelValues.Map["author"]) + "\n")
		}
		if *m.ModelValues.Map["description"] != "" {
			display.WriteString(s.KeyText.Render("Description: ") + s.ValueText.Render(*m.ModelValues.Map["description"]) + "\n")
		}
		if *m.ModelValues.Map["prefix"] != "" {
			display.WriteString(s.KeyText.Render("Command Prefix: ") + s.ValueText.Render(*m.ModelValues.Map["prefix"]) + "\n")
		}
		if *m.ModelValues.Map["cogs"] != "" {
			cogsString := *m.ModelValues.Map["cogs"]
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

func GlobalConfigModel(callback func(*Model) []error, initCallback func(*Model, []Values)) Model {
	m := Model{width: maxWidth}
	m.title = "BotBox CLI Configuration"
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.currentFormPtr = 0
	m.state = 0
	m.callback = callback
	localConfigForms := ConfigFormWrapperGenerator()
	m.forms = localConfigForms
	m.initCallback = initCallback

	m.ModelValues.Map = map[string]*string{
		"version":                new(string),
		"check_updates":          new(string),
		"auto_update":            new(string),
		"default_user":           new(string),
		"github_username":        new(string),
		"scroll_enabled":         new(string),
		"color_scheme":           new(string),
		"default_command_prefix": new(string),
		"default_python_version": new(string),
		"auto_git_init":          new(string),
		"editor":                 new(string),
	}

	m.displayCallback = func() string {
		s := m.styles
		var display strings.Builder
		if *m.ModelValues.Map["version"] != "" {
			display.WriteString(s.KeyText.Render("Version: ") + s.ValueText.Render(*m.ModelValues.Map["version"]) + "\n")
		}
		if *m.ModelValues.Map["check_updates"] != "" {
			display.WriteString(s.KeyText.Render("Check Updates: ") + s.ValueText.Render(*m.ModelValues.Map["check_updates"]) + "\n")
		}
		if *m.ModelValues.Map["auto_update"] != "" {
			display.WriteString(s.KeyText.Render("Auto Update: ") + s.ValueText.Render(*m.ModelValues.Map["auto_update"]) + "\n")
		}
		if *m.ModelValues.Map["default_user"] != "" {
			display.WriteString(s.KeyText.Render("Default User: ") + s.ValueText.Render(*m.ModelValues.Map["default_user"]) + "\n")
		}
		if *m.ModelValues.Map["github_username"] != "" {
			display.WriteString(s.KeyText.Render("GitHub Username: ") + s.ValueText.Render(*m.ModelValues.Map["github_username"]) + "\n")
		}
		if *m.ModelValues.Map["scroll_enabled"] != "" {
			display.WriteString(s.KeyText.Render("Scroll Enabled: ") + s.ValueText.Render(*m.ModelValues.Map["scroll_enabled"]) + "\n")
		}
		if *m.ModelValues.Map["color_scheme"] != "" {
			display.WriteString(s.KeyText.Render("Color Scheme: ") + s.ValueText.Render(*m.ModelValues.Map["color_scheme"]) + "\n")
		}
		if *m.ModelValues.Map["default_command_prefix"] != "" {
			display.WriteString(s.KeyText.Render("Default Command Prefix: ") + s.ValueText.Render(*m.ModelValues.Map["default_command_prefix"]) + "\n")
		}
		if *m.ModelValues.Map["default_python_version"] != "" {
			display.WriteString(s.KeyText.Render("Default Python Version: ") + s.ValueText.Render(*m.ModelValues.Map["default_python_version"]) + "\n")
		}
		if *m.ModelValues.Map["auto_git_init"] != "" {
			display.WriteString(s.KeyText.Render("Auto Git Init: ") + s.ValueText.Render(*m.ModelValues.Map["auto_git_init"]) + "\n")
		}
		if *m.ModelValues.Map["editor"] != "" {
			display.WriteString(s.KeyText.Render("Editor: ") + s.ValueText.Render(*m.ModelValues.Map["editor"]) + "\n")
		}

		return display.String()
	}

	return m
}

func ConfigSyncModel(callback func(*Model) []error, initCallback func(*Model, []Values)) Model {
	m := Model{width: maxWidth}
	m.title = "Sync BotBox Project Configuration"
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.currentFormPtr = 0
	m.state = 0
	m.callback = callback
	localConfigForms := ConfigSyncFormWrapperGenerator()
	m.forms = localConfigForms
	m.initCallback = initCallback

	m.ModelValues = Values{
		Map: map[string]*string{
			"updatedCogs":  new(string),
			"addedCogs":    new(string),
			"removedCogs":  new(string),
			"headerIssues": new(string),
			"noChanges":    new(string),
		},
		Name: "ModelValues",
	}

	m.displayCallback = func() string {
		s := m.styles
		var display strings.Builder
		if *m.ModelValues.Map["noChanges"] != "" {
			display.WriteString(s.HeaderText.Render("No changes detected") + "\n\n")
		} else {
			display.WriteString(s.HeaderText.Render("Updated config") + "\n\n")
			if *m.ModelValues.Map["addedCogs"] != "" {
				display.WriteString(s.KeyText.Render("Added Cogs: ") + s.ValueText.Render(*m.ModelValues.Map["addedCogs"]) + "\n")
			}
			if *m.ModelValues.Map["updatedCogs"] != "" {
				display.WriteString(s.KeyText.Render("Updated Cogs: ") + s.ValueText.Render(*m.ModelValues.Map["updatedCogs"]) + "\n")
			}
			if *m.ModelValues.Map["removedCogs"] != "" {
				display.WriteString(s.KeyText.Render("Removed Cogs: ") + s.ValueText.Render(*m.ModelValues.Map["removedCogs"]) + "\n")
			}
			if *m.ModelValues.Map["headerIssues"] != "" {
				display.WriteString(s.KeyText.Render("Header Issues: ") + s.ValueText.Render(*m.ModelValues.Map["headerIssues"]) + "\n")
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
