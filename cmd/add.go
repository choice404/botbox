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
			fmt.Println("Current directory is not in a botbox project.")
			return
		}

		if len(args) > 0 {
			addCogName = args[0]
		} else {
			addCogName = ""
		}

		model := utils.AddModel(addCallback, addInitCallback)
		utils.CupSleeve(model)
	},
}

func addCallback(model *utils.Model) []error {
	values := model.ModelValues
	rootDir, err := utils.FindBotConf()
	config, err := utils.LoadConfig()
	var errors []error
	if err != nil {
		errors = append(errors, fmt.Errorf("error finding root directory: %w", err))
		return errors
	}

	filename := *values.Map["filename"]
	slashCommandList, _ := utils.JSONToCmdInfoSlice(*values.Map["slashCommands"])
	prefixCommandList, _ := utils.JSONToCmdInfoSlice(*values.Map["prefixCommands"])

	filePath := filepath.Join(rootDir, "src", "cogs", filename+".py")
	file, err := os.Create(filePath)
	if err != nil {
		errors = append(errors, fmt.Errorf("error creating file: %w", err))
		return errors
	}
	defer file.Close()

	className := strings.ToUpper(string(filename[0])) + filename[1:]

	var cogContent strings.Builder

	fmt.Fprintf(&cogContent, `"""
Bot Author: %s

%s
%s
"""

import discord
from discord import app_commands
from discord.ext import commands
from dotenv import load_dotenv
import os

load_dotenv()

GUILD_ID = int(os.getenv("GUILD_ID", 0))
GUILD = discord.Object(id=GUILD_ID) if GUILD_ID else None

class %s(commands.Cog, name="%s"):
    def __init__(self, bot) -> None:
        self.bot = bot
        print("%s cog loaded")
`, config.BotInfo.Author, config.BotInfo.Name, config.BotInfo.Description, className, className, filename)

	for _, command := range slashCommandList {
		fullArgStr := buildArgString(command.Args)

		fmt.Fprintf(&cogContent, `
    @app_commands.command(name="%s", description="%s")`, command.Name, command.Description)

		if fullArgStr != "" {
			cogContent.WriteString(`
    @app_commands.describe(`)

			for _, arg := range command.Args {
				fmt.Fprintf(&cogContent, `
        %s="%s",`, arg.Name, arg.Description)
			}

			fmt.Fprintf(&cogContent, `
    )`)
		}

		if command.Scope == "guild" {
			fmt.Fprintf(&cogContent, `
    @app_commands.guilds(GUILD) if GUILD else app_commands.default_permissions`)
		}

		fmt.Fprintf(&cogContent, `
    async def %s(self, interaction: discord.Interaction, %s) -> %s:
        """
        %s when the user types "/%s"

            Parameters:`, command.Name, fullArgStr, command.ReturnType, command.Description, command.Name)

		for _, arg := range command.Args {
			fmt.Fprintf(&cogContent, `
                    %s (%s): %s`, arg.Name, arg.Type, arg.Description)
		}

		fmt.Fprintf(&cogContent, `

            Returns:
                    %s
        """

        try:
            await interaction.response.send_message(f"%s", ephemeral=True)
        except Exception as e:
            print(f"Error: {e}")
            await interaction.response.send_message(f"Error: {e}", ephemeral=True)

        return %s
`, command.ReturnType, command.Name, getReturnValue(command.ReturnType))
	}

	for _, command := range prefixCommandList {
		fullArgStr := buildArgString(command.Args)

		fmt.Fprintf(&cogContent, `
    @commands.command()
    async def %s(self, ctx: commands.Context, %s) -> %s:
        """
        %s when the user types "/%s"

            Parameters:
`, command.Name, fullArgStr, command.ReturnType, command.Description, command.Name)

		for _, arg := range command.Args {
			fmt.Fprintf(&cogContent, `
                    %s (%s): %s`, arg.Name, arg.Type, arg.Description)
		}

		fmt.Fprintf(&cogContent, `

            Returns:
                    %s
        """

        try:
            await ctx.send(f"%s", ephemeral=True)
        except Exception as e:
            print(f"Error: {e}")
            await ctx.send(f"Error: {e}", ephemeral=True)

        return %s
`, command.ReturnType, command.Name, getReturnValue(command.ReturnType))
	}

	fmt.Fprintf(&cogContent, `

async def setup(bot):
    await bot.add_cog(%s(bot))

"""
File generated by BotBox - https://github.com/choice404/botbox
"""`, className)

	_, err = file.WriteString(cogContent.String())
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
		File:           strings.ToLower(string(filename[0])) + filename[1:],
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
		*modelValues.Map["filename"] = addCogName
		configs, err := utils.LoadConfig()
		if err != nil {
			errors = append(errors, fmt.Errorf("error loading configuration: %w", err))
			model.HandleError(errors)
			return
		}
		if len(configs.Cogs) == 0 {
			errors = append(errors, fmt.Errorf("error: no cogs available to add"))
			model.HandleError(errors)
			return
		}

		cogExists := false
		for _, cog := range configs.Cogs {
			if cog.Name == addCogName {
				cogExists = true
				break
			}
		}
		if cogExists {
			errors = append(errors, fmt.Errorf("cog '%s' already exists in the project", addCogName))
			model.HandleError(errors)
			return
		}
	}
}

func getReturnValue(returnType string) string {
	switch returnType {
	case "str":
		return `""`
	case "int":
		return "0"
	case "float":
		return "0.0"
	case "bool":
		return "False"
	default:
		return "None"
	}
}

func buildArgString(args []utils.ArgInfo) string {
	if len(args) == 0 {
		return ""
	}
	var argBuilder strings.Builder
	for i, arg := range args {
		fmt.Fprintf(&argBuilder, "%s: %s", arg.Name, arg.Type)
		if i < len(args)-1 {
			argBuilder.WriteString(", ")
		}
	}
	return argBuilder.String()
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
