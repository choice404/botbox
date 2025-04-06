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

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds new files to the project",
	Long:  `Adds a new file to the project. By default it will add a cog however the user can specify to add a custom command file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
		filename := ""
		if len(args) > 0 {
			filename = args[0]
		}
		addCogs(filename)
	},
}

func addCogs(filename string) error {
	rootDir, err := FindBotConf()
	config, err := LoadConfig()
	Banner()

	if filename == "" {
		fmt.Println("Eenter the cog name:")
		fmt.Scanln(&filename)
	}
	filePath := filepath.Join(rootDir, "cogs", filename+".py")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to create cog file %s: %w", filename, err)
	}
	defer file.Close()

	var commands []CommandInfo
	var command string
	var description string
	var returnType string
	for {
		fmt.Println("Enter a command name (or '!!' to finish):")
		fmt.Scanln(&command)
		if command == "!!" {
			break
		}
		fmt.Println("Enter a description for the command:")
		fmt.Scanln(&description)

		var args []ArgInfo
		var argsInput string
		var argType string
		var argDescription string
		for {
			fmt.Println("Enter an argument name (or '!!' to finish):")
			fmt.Scanln(&argsInput)
			if argsInput == "!!" {
				break
			}
			fmt.Println("Enter the argument type:")
			fmt.Scanln(&argType)
			fmt.Println("Enter a description for the argument:")
			fmt.Scanln(&argDescription)
			args = append(args, ArgInfo{Name: argsInput, Type: argType, Description: argDescription})
		}

		fmt.Println("Enter the return type:")
		fmt.Scanln(&returnType)
		commands = append(commands, CommandInfo{Name: command, Description: description, Args: args, ReturnType: returnType})
	}

	className := strings.ToUpper(string(filename[0])) + filename[1:]

	_, err = fmt.Fprintf(file, `"""
Author %s

%s
%s
"""

import discord
from discord import app_commands
from discord.ext import commands
from dotenv import load_dotenv

load_dotenv()

class %s(commands.Cog):
    def __init__(self, bot) -> None:
        self.bot = bot
        print("%s cog loaded")
`, config.BotInfo.Author, config.BotInfo.Name, config.BotInfo.Description, className, filename)
	if err != nil {
		return fmt.Errorf("failed to write initial part of cog file %s: %w", filePath, err)
	}

	for _, command := range commands {
		var argBuilder strings.Builder
		for i, arg := range command.Args {
			fmt.Fprintf(&argBuilder, "%s: %s", arg.Name, arg.Type)
			if i < len(command.Args)-1 {
				argBuilder.WriteString(", ")
			}
		}
		fullArgStr := argBuilder.String()

		_, err = fmt.Fprintf(file, `
    @app_commands.command(name="%s", description="%s")
    async def %s(self, interaction: discord.Interaction, %s) -> %s:
        """
        %s when the user types "/%s"

            Parameters:
`, command.Name, command.Description, command.Name, fullArgStr, command.ReturnType, command.Description, command.Name)
		if err != nil {
			return fmt.Errorf("writing command signature for %s: %w", command.Name, err)
		}

		for _, arg := range command.Args {
			_, err = fmt.Fprintf(file, `                    %s (%s): %s
`, arg.Name, arg.Type, arg.Description)
			if err != nil {
				return fmt.Errorf("writing arg description for %s in %s: %w", arg.Name, command.Name, err)
			}
		}

		_, err = fmt.Fprintf(file, `
            Returns:
                    %s
        """

        try:
            await interaction.response.send_message(f"%s", ephemeral=True)
        except Exception as e:
            print(f"Error: {e}")
            await interaction.response.send_message(f"Error: {e}", ephemeral=True)
`, command.ReturnType, command.Name)
		if err != nil {
			return fmt.Errorf("writing return/try block for %s: %w", command.Name, err)
		}

		var retValue any
		switch command.ReturnType {
		case "str":
			retValue = `""`
		case "int":
			retValue = 0
		case "float":
			retValue = 0.0
		case "bool":
			retValue = "False"
		default:
			retValue = "None"
		}

		_, err = fmt.Fprintf(file, `
        return %v
`, retValue)
		if err != nil {
			return fmt.Errorf("writing return statement for %s: %w", command.Name, err)
		}
	}

	err = file.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync cog file %s: %w", filePath, err)
	}
	cog := CogConfig{
		Name:     strings.ToUpper(string(filename[0])) + filename[1:],
		File:     strings.ToLower(string(filename[0])) + filename[1:],
		Commands: []string{},
	}

	for _, command := range commands {
		cog.Commands = append(cog.Commands, command.Name)
	}

	config.Cogs = append(config.Cogs, cog)

	jsonData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	err = os.WriteFile("config.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated config.json: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2025 Austin "Choice404" Choi

https://github.com/choice404/botbox

Bot Box
A discord bot template generator to help create discord bots quickly and easily

This code is licensed under the MIT License.
Please see the LICENSE file in the root directory of this project for the full license details.
*/
