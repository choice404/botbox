/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

// Flags that carry project values, providing any of them implies headless mode
var projectValueFlags = []string{"name", "description", "author", "prefix", "env", "token", "doppler-project", "guild", "doppler-env", "license", "help-style", "docker"}

/**
 * registerProjectFlags
 * Registers the shared create/init flag set on a command
 * @param cmd {*cobra.Command} - the command to register the flags on
 * @return ...
 **/
func registerProjectFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Bot name")
	cmd.Flags().String("description", "", "Bot description")
	cmd.Flags().String("author", "", "Bot author (defaults to user.default_user from the global config)")
	cmd.Flags().String("prefix", "", "Single non alphanumeric character command prefix (defaults to defaults.command_prefix)")
	cmd.Flags().String("env", "env", "How to handle environment variables, env or doppler")
	cmd.Flags().String("token", "", "Bot token, prefer the BOTBOX_TOKEN environment variable so the token stays out of shell history")
	cmd.Flags().String("doppler-project", "", "Doppler project name when using --env doppler")
	cmd.Flags().String("guild", "", "Guild ID when using --env env")
	cmd.Flags().String("doppler-env", "", "Doppler environment name when using --env doppler")
	cmd.Flags().String("license", "mit", "License type: mit, apache-2.0, gpl-3.0, bsd-3-clause, unlicense, no-license")
	cmd.Flags().String("help-style", "compact", "How the generated help command formats its output: compact or detailed")
	cmd.Flags().Bool("docker", false, "Generate Docker files (Dockerfile, docker-compose.yml, .dockerignore)")
	cmd.Flags().Bool("force", false, "Overwrite existing files without prompting")
}

/**
 * isHeadless
 * Decides if a command should run without the tui
 * @param cmd {*cobra.Command} - the command being run
 * @param valueFlags {[]string} - flags that imply headless mode when set
 * @return bool - true when the command should skip the tui
 **/
func isHeadless(cmd *cobra.Command, valueFlags []string) bool {
	if headless, err := cmd.Flags().GetBool("headless"); err == nil && headless {
		return true
	}
	for _, name := range valueFlags {
		if cmd.Flags().Changed(name) {
			return true
		}
	}
	return false
}

/**
 * collectProjectValues
 * Reads the project flags, applies global config defaults, and validates everything
 * @param cmd {*cobra.Command} - the command holding the flags
 * @param args {[]string} - positional args, the first one is the bot name
 * @return map[string]string - the validated ModelValues entries
 * @return error - the first validation failure
 **/
func collectProjectValues(cmd *cobra.Command, args []string) (map[string]string, error) {
	flags := cmd.Flags()

	name, _ := flags.GetString("name")
	if name == "" && len(args) > 0 {
		name = args[0]
	}
	if err := utils.ValidateBotName(name); err != nil {
		return nil, err
	}

	description, _ := flags.GetString("description")
	if err := utils.ValidateBotDescription(description); err != nil {
		return nil, err
	}

	author, _ := flags.GetString("author")
	if author == "" && GlobalConfig != nil {
		author = GlobalConfig.User.DefaultUser
	}
	if err := utils.ValidateBotAuthor(author); err != nil {
		return nil, fmt.Errorf("%w (set --author or user.default_user in the global config)", err)
	}

	prefix, _ := flags.GetString("prefix")
	if prefix == "" && GlobalConfig != nil {
		prefix = GlobalConfig.Defaults.CommandPrefix
	}
	if prefix == "" {
		prefix = "!"
	}
	if err := utils.ValidateBotPrefix(prefix); err != nil {
		return nil, err
	}

	envChoice, _ := flags.GetString("env")
	if err := utils.ValidateEnvChoice(envChoice); err != nil {
		return nil, err
	}

	// The token and doppler project share one value key, same for guild and doppler env
	tokenOrProject := ""
	guildOrEnv := ""
	if envChoice == "env" {
		tokenOrProject, _ = flags.GetString("token")
		if tokenOrProject == "" {
			tokenOrProject = os.Getenv("BOTBOX_TOKEN")
		}
		if err := utils.ValidateToken(envChoice, tokenOrProject); err != nil {
			return nil, fmt.Errorf("%w (set --token or BOTBOX_TOKEN)", err)
		}
		guildOrEnv, _ = flags.GetString("guild")
	} else {
		tokenOrProject, _ = flags.GetString("doppler-project")
		if tokenOrProject == "" {
			return nil, fmt.Errorf("doppler project name cannot be empty, set --doppler-project")
		}
		guildOrEnv, _ = flags.GetString("doppler-env")
	}

	license, _ := flags.GetString("license")
	if err := utils.ValidateLicense(license); err != nil {
		return nil, err
	}

	// An unset help style falls back to the default instead of failing validation
	helpStyle, _ := flags.GetString("help-style")
	if helpStyle == "" {
		helpStyle = utils.DefaultHelpStyle
	}
	if err := utils.ValidateHelpStyle(helpStyle); err != nil {
		return nil, err
	}

	// The docker flag rides the values bus as yes or no like the force flag does
	docker, _ := flags.GetBool("docker")
	dockerize := "no"
	if docker {
		dockerize = "yes"
	}

	return map[string]string{
		"botName":                name,
		"botDescription":         description,
		"botAuthor":              author,
		"botPrefix":              prefix,
		"envChoice":              envChoice,
		"botTokenDopplerProject": tokenOrProject,
		"botGuildDopplerEnv":     guildOrEnv,
		"licenseType":            license,
		"helpStyle":              helpStyle,
		"dockerize":              dockerize,
	}, nil
}

/**
 * applyModelValues
 * Copies collected values into a model's ModelValues map
 * @param model {*utils.Model} - the model to fill
 * @param values {map[string]string} - the values to apply
 * @return ...
 **/
func applyModelValues(model *utils.Model, values map[string]string) {
	for key, value := range values {
		v := value
		if model.ModelValues.Map[key] == nil {
			model.ModelValues.Map[key] = &v
		} else {
			*model.ModelValues.Map[key] = v
		}
	}
}

/**
 * setForceValue
 * Stores the force flag on the model values bus as yes or no
 * @param model {*utils.Model} - the model to store the flag on
 * @param force {bool} - the force flag value
 * @return ...
 **/
func setForceValue(model *utils.Model, force bool) {
	s := "no"
	if force {
		s = "yes"
	}
	model.ModelValues.Map["force"] = &s
}

/**
 * readForceValue
 * Reads the force flag back off the model values bus
 * @param values {utils.Values} - the model values
 * @return bool - true when force was set
 **/
func readForceValue(values utils.Values) bool {
	if f, ok := values.Map["force"]; ok && f != nil {
		return *f == "yes"
	}
	return false
}

/*
Copyright © 2025 Austin "Choice404" Choi

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
