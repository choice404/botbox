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

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	editCogName     string
	editWrittenPath string
	editBackupPath  string
)

// Flags that carry edit values, providing any of them implies headless mode
var editValueFlags = []string{"remove-command", "add-commands", "replace-commands", "env", "no-backup"}

// editOptions carries the headless edit flags after they leave cobra
type editOptions struct {
	removeCommands  []string
	addCommands     string
	replaceCommands string
	replaceSet      bool
	env             string
	noBackup        bool
}

var editCmd = &cobra.Command{
	Use:   "edit [cog-name]",
	Short: "Edit an existing cog in the current Bot Box project",
	Long: `Edit an existing cog (command module) in your Bot Box project.

This command lets you change a cog without recreating it:
  - Add new slash, prefix, or modal commands
  - Edit the info, arguments, fields, pages, and responses of a command
  - Remove commands from the cog
  - Switch the cog between the development and production environments

Applying the changes regenerates the cog's .py file from its definition, so
custom code written inside command bodies is not preserved. A .py.bak copy of
the previous file is written first unless --no-backup is given.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := utils.FindBotConf()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Current directory is not in a botbox project.")
			return
		}

		config, err := utils.LoadConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return
		}
		if len(config.Cogs) == 0 {
			fmt.Fprintln(os.Stderr, "Error: no cogs available to edit")
			return
		}

		if len(args) > 0 {
			editCogName = args[0]
		} else {
			editCogName = ""
		}

		if isHeadless(cmd, editValueFlags) {
			runEditHeadless(cmd)
			return
		}

		model := utils.EditModel(editCallback, editInitCallback)
		utils.CupSleeve(model)
		printEditResult()
	},
}

func runEditHeadless(cmd *cobra.Command) {
	model, err := buildEditHeadlessModel(collectEditOptions(cmd))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if utils.PrintErrors(utils.RunHeadless(*model)) {
		os.Exit(1)
	}

	printEditResult()
}

/**
 * collectEditOptions
 * Reads the headless edit flags off the command
 * @param cmd {*cobra.Command} - the command holding the flags
 * @return editOptions - the collected flag values
 **/
func collectEditOptions(cmd *cobra.Command) editOptions {
	flags := cmd.Flags()
	removeCommands, _ := flags.GetStringArray("remove-command")
	addCommands, _ := flags.GetString("add-commands")
	replaceCommands, _ := flags.GetString("replace-commands")
	env, _ := flags.GetString("env")
	noBackup, _ := flags.GetBool("no-backup")

	return editOptions{
		removeCommands:  removeCommands,
		addCommands:     addCommands,
		replaceCommands: replaceCommands,
		replaceSet:      flags.Changed("replace-commands"),
		env:             env,
		noBackup:        noBackup,
	}
}

/**
 * buildEditHeadlessModel
 * Applies the edit operations to the cog's command set and builds the model to run
 * Operations run in a fixed order, replace first, then removes, then adds
 * @param opts {editOptions} - the collected headless flags
 * @return *utils.Model - the model ready for RunHeadless
 * @return error - the first validation or parse failure
 **/
func buildEditHeadlessModel(opts editOptions) (*utils.Model, error) {
	if editCogName == "" {
		return nil, fmt.Errorf("a cog name is required when running without the TUI")
	}

	config, err := utils.LoadConfig()
	if err != nil {
		return nil, err
	}

	cogIndex, found := findProjectCog(config, editCogName)
	if !found {
		return nil, fmt.Errorf("cog '%s' does not exist in the project", editCogName)
	}
	cog := config.Cogs[cogIndex]

	// The working set holds slash commands first so validation sees the config order
	commands := append([]utils.CommandInfo{}, cog.SlashCommands...)
	commands = append(commands, cog.PrefixCommands...)

	// Replace swaps the whole command set before removes and adds run
	if opts.replaceSet {
		replaced, err := parseCommandsInput(opts.replaceCommands)
		if err != nil {
			return nil, err
		}
		commands = normalizeModalReturns(replaced)
	}

	// Removes must name a command that exists in the working set
	for _, name := range opts.removeCommands {
		removed := false
		for i, command := range commands {
			if command.Name == name {
				commands = append(commands[:i], commands[i+1:]...)
				removed = true
				break
			}
		}
		if !removed {
			return nil, fmt.Errorf("command '%s' does not exist in cog '%s'", name, cog.Name)
		}
	}

	if opts.addCommands != "" {
		added, err := parseCommandsInput(opts.addCommands)
		if err != nil {
			return nil, err
		}
		commands = append(commands, normalizeModalReturns(added)...)
	}

	// Every command in the final set validates against the ones before it
	for i, command := range commands {
		if err := utils.ValidateCommand(command, commands[:i]); err != nil {
			return nil, fmt.Errorf("command '%s': %v", command.Name, err)
		}
	}

	if opts.env != "" && opts.env != "development" && opts.env != "production" {
		return nil, fmt.Errorf("env must be development or production")
	}

	// Modal commands are app commands, so they live with the slash commands
	slashCommands := []utils.CommandInfo{}
	prefixCommands := []utils.CommandInfo{}
	for _, command := range commands {
		if command.Type == "prefix" {
			prefixCommands = append(prefixCommands, command)
		} else {
			slashCommands = append(slashCommands, command)
		}
	}

	slashJSON, err := utils.CmdInfoSliceToJSON(slashCommands)
	if err != nil {
		return nil, err
	}
	prefixJSON, err := utils.CmdInfoSliceToJSON(prefixCommands)
	if err != nil {
		return nil, err
	}

	model := utils.EditModel(editCallback, editInitCallback)
	*model.ModelValues.Map["cogName"] = cog.Name
	model.ModelValues.Map["slashCommands"] = &slashJSON
	model.ModelValues.Map["prefixCommands"] = &prefixJSON
	*model.ModelValues.Map["cogEnv"] = opts.env
	if opts.noBackup {
		*model.ModelValues.Map["backup"] = "no"
	}

	return &model, nil
}

/**
 * normalizeModalReturns
 * Returns a copy of the commands with modal return types pinned to None
 * @param commands {[]utils.CommandInfo} - the commands to normalize
 * @return []utils.CommandInfo - the normalized copy
 **/
func normalizeModalReturns(commands []utils.CommandInfo) []utils.CommandInfo {
	normalized := make([]utils.CommandInfo, 0, len(commands))
	for _, command := range commands {
		// Modal commands only respond through the modal, so their return type is fixed
		if command.Type == "modal" {
			command.ReturnType = "None"
		}
		normalized = append(normalized, command)
	}
	return normalized
}

/**
 * findProjectCog
 * Finds a cog in the config by its display name or file base
 * @param config {utils.Config} - the loaded project config
 * @param name {string} - the name or file base to match
 * @return int - the index of the matching cog
 * @return bool - true when a cog matched
 **/
func findProjectCog(config utils.Config, name string) (int, bool) {
	for i, cog := range config.Cogs {
		if cog.Name == name || cog.File == name {
			return i, true
		}
	}
	return -1, false
}

/**
 * printEditResult
 * Prints the regeneration warning to stderr and the written file path to stdout
 * @return ...
 **/
func printEditResult() {
	if editWrittenPath == "" {
		return
	}
	warning := "regeneration rewrites the cog file from its definition, custom code in command bodies is not preserved"
	if editBackupPath != "" {
		warning += ", previous version saved to " + editBackupPath
	}
	fmt.Fprintln(os.Stderr, "Warning:", warning)
	fmt.Println(editWrittenPath)
}

func editCallback(model *utils.Model) []error {
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

	cogName := *values.Map["cogName"]
	if cogName == "" {
		errors = append(errors, fmt.Errorf("no cog selected to edit"))
		return errors
	}
	cogIndex := -1
	for i, cog := range config.Cogs {
		if cog.Name == cogName {
			cogIndex = i
			break
		}
	}
	if cogIndex < 0 {
		errors = append(errors, fmt.Errorf("cog '%s' does not exist in the project", cogName))
		return errors
	}

	slashCommands, err := utils.JSONToCmdInfoSlice(*values.Map["slashCommands"])
	if err != nil {
		errors = append(errors, fmt.Errorf("error reading slash commands: %w", err))
		return errors
	}
	prefixCommands, err := utils.JSONToCmdInfoSlice(*values.Map["prefixCommands"])
	if err != nil {
		errors = append(errors, fmt.Errorf("error reading prefix commands: %w", err))
		return errors
	}
	if slashCommands == nil {
		slashCommands = []utils.CommandInfo{}
	}
	if prefixCommands == nil {
		prefixCommands = []utils.CommandInfo{}
	}

	// Prefix commands carry no scope in the generated file, the parser always reads them as global
	normalizedPrefix := make([]utils.CommandInfo, 0, len(prefixCommands))
	for _, command := range prefixCommands {
		command.Scope = "global"
		normalizedPrefix = append(normalizedPrefix, command)
	}
	prefixCommands = normalizedPrefix

	cog := config.Cogs[cogIndex]
	cog.SlashCommands = slashCommands
	cog.PrefixCommands = prefixCommands
	if env := *values.Map["cogEnv"]; env != "" {
		cog.Env = env
	}

	backup := *values.Map["backup"] != "no"
	if err := utils.RegenerateCogFile(rootDir, config, cog, backup); err != nil {
		errors = append(errors, fmt.Errorf("error regenerating cog file: %w", err))
		return errors
	}

	config.Cogs[cogIndex] = cog

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to marshal config to JSON: %w", err))
		return errors
	}

	confPath := filepath.Join(rootDir, "botbox.conf")
	err = os.WriteFile(confPath, jsonData, 0644)
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to write updated botbox.conf: %w", err))
		return errors
	}

	// The run function prints the result after the tui or headless run finishes
	editWrittenPath = filepath.Join(rootDir, "src", "cogs", cog.File+".py")
	editBackupPath = ""
	if backup {
		if _, err := os.Stat(editWrittenPath + ".bak"); err == nil {
			editBackupPath = editWrittenPath + ".bak"
		}
	}

	return nil
}

func editInitCallback(model *utils.Model, allFormsModels []utils.Values) {
	if editCogName == "" {
		return
	}

	modelValues := model.ModelValues
	var errors []error

	config, err := utils.LoadConfig()
	if err != nil {
		errors = append(errors, fmt.Errorf("error loading configuration: %w", err))
		model.HandleError(errors)
		return
	}

	cogIndex, found := findProjectCog(config, editCogName)
	if !found {
		errors = append(errors, fmt.Errorf("cog '%s' does not exist in the project", editCogName))
		model.HandleError(errors)
		return
	}
	cog := config.Cogs[cogIndex]

	*modelValues.Map["cogName"] = cog.Name

	// The headless path prefills the command lists, the tui arg path loads them here
	if *modelValues.Map["slashCommands"] == "" {
		slashCommands := cog.SlashCommands
		if slashCommands == nil {
			slashCommands = []utils.CommandInfo{}
		}
		slashJSON, err := utils.CmdInfoSliceToJSON(slashCommands)
		if err != nil {
			errors = append(errors, fmt.Errorf("error reading slash commands: %w", err))
			model.HandleError(errors)
			return
		}
		*modelValues.Map["slashCommands"] = slashJSON
	}
	if *modelValues.Map["prefixCommands"] == "" {
		prefixCommands := cog.PrefixCommands
		if prefixCommands == nil {
			prefixCommands = []utils.CommandInfo{}
		}
		prefixJSON, err := utils.CmdInfoSliceToJSON(prefixCommands)
		if err != nil {
			errors = append(errors, fmt.Errorf("error reading prefix commands: %w", err))
			model.HandleError(errors)
			return
		}
		*modelValues.Map["prefixCommands"] = prefixJSON
	}
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringArray("remove-command", nil, "Name of a command to remove from the cog, repeatable")
	editCmd.Flags().String("add-commands", "", "JSON array of commands to add, accepts inline JSON, @path/to/file.json, or - for stdin")
	editCmd.Flags().String("replace-commands", "", "JSON array that replaces every command on the cog, accepts inline JSON, @path/to/file.json, or - for stdin")
	editCmd.Flags().String("env", "", "Cog environment: development or production")
	editCmd.Flags().Bool("no-backup", false, "Skip the .py.bak backup written before the cog file is regenerated")
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
