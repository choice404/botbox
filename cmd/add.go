/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	addCogName string
)

var addCmd = &cobra.Command{
	Use:   "add [cog-name]",
	Short: "Add a new cog to the current Bot Box project",
	Long: `Add a new cog (command module) to your Bot Box project with an interactive setup.

This command guides you through creating a new cog by specifying:
  - Cog name and file structure
  - Slash commands with descriptions and arguments
  - Prefix commands for traditional bot interactions
  - Command argument types and return values
  - Command scopes (guild or global)

The generated cog will be automatically registered in botbox.conf and include 
proper Discord.py boilerplate code. It's recommended to use this command instead 
of manually creating cogs to ensure proper integration.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := utils.FindBotConf()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Current directory is not in a botbox project.")
			return
		}

		if len(args) > 0 {
			addCogName = args[0]
		} else {
			addCogName = ""
		}

		if isHeadless(cmd, []string{"commands"}) {
			runAddHeadless(cmd, args)
			return
		}

		model := utils.AddModel(addCallback, addInitCallback)
		utils.CupSleeve(model)
	},
}

func runAddHeadless(cmd *cobra.Command, args []string) {
	if addCogName == "" {
		fmt.Fprintln(os.Stderr, "Error: a cog name is required when running without the TUI")
		os.Exit(1)
	}

	rawCommands, _ := cmd.Flags().GetString("commands")
	commands, err := parseCommandsInput(rawCommands)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	// Validate each command against the ones accepted before it
	var slashCommands, prefixCommands []utils.CommandInfo
	for i, command := range commands {
		if err := utils.ValidateCommand(command, commands[:i]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: command '%s': %v\n", command.Name, err)
			os.Exit(1)
		}
		if command.Type == "slash" {
			slashCommands = append(slashCommands, command)
		} else {
			prefixCommands = append(prefixCommands, command)
		}
	}

	slashJSON, err := utils.CmdInfoSliceToJSON(slashCommands)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	prefixJSON, err := utils.CmdInfoSliceToJSON(prefixCommands)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	model := utils.AddModel(addCallback, addInitCallback)
	model.ModelValues.Map["slashCommands"] = &slashJSON
	model.ModelValues.Map["prefixCommands"] = &prefixJSON

	if utils.PrintErrors(utils.RunHeadless(model)) {
		os.Exit(1)
	}

	rootDir, err := utils.FindBotConf()
	if err == nil {
		fileBase := strings.ToLower(string(addCogName[0])) + addCogName[1:]
		fmt.Println(filepath.Join(rootDir, "src", "cogs", fileBase+".py"))
	}
}

/**
 * parseCommandsInput
 * Parses the --commands flag which accepts inline JSON, @path/to/file.json, or - for stdin
 * @param raw {string} - the raw flag value
 * @return []utils.CommandInfo - the parsed commands
 * @return error - any read or parse failure
 **/
func parseCommandsInput(raw string) ([]utils.CommandInfo, error) {
	if raw == "" {
		return nil, nil
	}

	var data []byte
	switch {
	case raw == "-":
		stdinData, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("error reading commands from stdin: %w", err)
		}
		data = stdinData
	case strings.HasPrefix(raw, "@"):
		fileData, err := os.ReadFile(strings.TrimPrefix(raw, "@"))
		if err != nil {
			return nil, fmt.Errorf("error reading commands file: %w", err)
		}
		data = fileData
	default:
		data = []byte(raw)
	}

	var commands []utils.CommandInfo
	if err := json.Unmarshal(data, &commands); err != nil {
		return nil, fmt.Errorf("error parsing commands JSON: %w", err)
	}
	return commands, nil
}

func addCallback(model *utils.Model) []error {
	values := model.ModelValues
	var errors []error
	rootDir, err := utils.FindBotConf()
	if err != nil {
		errors = append(errors, fmt.Errorf("error finding root directory: %w", err))
		return errors
	}
	config, err := utils.LoadConfig()
	if err != nil {
		errors = append(errors, fmt.Errorf("error loading configuration: %w", err))
		return errors
	}

	filename := *values.Map["filename"]
	if filename == "" {
		errors = append(errors, fmt.Errorf("filename cannot be empty"))
		return errors
	}
	fileBase := strings.ToLower(string(filename[0])) + filename[1:]
	slashCommandList, _ := utils.JSONToCmdInfoSlice(*values.Map["slashCommands"])
	prefixCommandList, _ := utils.JSONToCmdInfoSlice(*values.Map["prefixCommands"])

	filePath := filepath.Join(rootDir, "src", "cogs", fileBase+".py")
	file, err := os.Create(filePath)
	if err != nil {
		errors = append(errors, fmt.Errorf("error creating file: %w", err))
		return errors
	}
	defer file.Close()

	className := strings.ToUpper(string(filename[0])) + filename[1:]

	cogContent, err := utils.RenderTemplate("cog.py.tmpl", utils.CogTemplateData{
		Author:         config.BotInfo.Author,
		BotName:        config.BotInfo.Name,
		BotDescription: config.BotInfo.Description,
		ClassName:      className,
		Filename:       filename,
		SlashCommands:  slashCommandList,
		PrefixCommands: prefixCommandList,
	})
	if err != nil {
		errors = append(errors, fmt.Errorf("error rendering cog template: %w", err))
		return errors
	}

	_, err = file.WriteString(cogContent)
	if err != nil {
		errors = append(errors, fmt.Errorf("error writing to file: %w", err))
		return errors
	}

	err = file.Sync()
	if err != nil {
		errors = append(errors, fmt.Errorf("error syncing file: %w", err))
		return errors
	}

	cog := utils.CogConfig{
		Name:           strings.ToUpper(string(filename[0])) + filename[1:],
		File:           fileBase,
		SlashCommands:  []utils.CommandInfo{},
		PrefixCommands: []utils.CommandInfo{},
	}

	for _, slashCommand := range slashCommandList {
		cog.SlashCommands = append(cog.SlashCommands, slashCommand)
	}

	for _, prefixCommand := range prefixCommandList {
		cog.PrefixCommands = append(cog.PrefixCommands, prefixCommand)
	}

	cog.Env = "development"
	config.Cogs = append(config.Cogs, cog)

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to marshal config to JSON: %w", err))
		return errors
	}

	confDir, err := utils.FindBotConf()
	if err != nil {
		errors = append(errors, fmt.Errorf("error finding botbox config directory: %w", err))
		return errors
	}

	confPath := filepath.Join(confDir, "botbox.conf")

	err = os.WriteFile(confPath, jsonData, 0644)
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to write updated botbox.conf: %w", err))
		return errors
	}

	return nil
}

func addInitCallback(model *utils.Model, allFormsModels []utils.Values) {
	modelValues := model.ModelValues
	var errors []error
	if addCogName != "" {
		if err := utils.ValidateFileName(addCogName); err != nil {
			errors = append(errors, fmt.Errorf("invalid cog name: %w", err))
			model.HandleError(errors)
			return
		}
		configs, err := utils.LoadConfig()
		if err != nil {
			errors = append(errors, fmt.Errorf("error loading configuration: %w", err))
			model.HandleError(errors)
			return
		}

		cogFile := strings.ToLower(string(addCogName[0])) + addCogName[1:]
		cogExists := false
		for _, cog := range configs.Cogs {
			if cog.File == cogFile {
				cogExists = true
				break
			}
		}
		if cogExists {
			errors = append(errors, fmt.Errorf("cog '%s' already exists in the project", addCogName))
			model.HandleError(errors)
			return
		}

		*modelValues.Map["filename"] = addCogName
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().String("commands", "", "JSON array of commands to generate, accepts inline JSON, @path/to/file.json, or - for stdin")
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
