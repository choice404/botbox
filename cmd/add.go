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

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	cogName     string
	cogClass    string
	cogFunction string

	commandList    []CommandInfo
	command        string
	commandConfirm bool
	description    string
	returnType     string
	args           []ArgInfo
	argsInput      string
	argType        string
	argDescription string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds new files to the project",
	Long:  `Adds a new file to the project. By default it will add a cog however the user can specify to add a custom command file.`,
	Run: func(cmd *cobra.Command, args []string) {
		Banner()
		var filename string
		if len(args) > 0 {
			filename = args[0]
		}
		addCogs(filename)
	},
}

func addCogs(filename string) {
	rootDir, err := FindBotConf()
	if err != nil {
		fmt.Println("Error finding bot config:", err)
		return
	}
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	for {
		cmdConfirmForm, cmdInfoForm := generateCmdForms()
		err := cmdConfirmForm.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if !cmdConfirmForm.GetBool("confirm") {
			break
		}

		err = cmdInfoForm.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		command = cmdInfoForm.GetString("name")
		description = cmdInfoForm.GetString("description")
		returnType = cmdInfoForm.GetString("returnType")
		for {
			argConfirmForm, argInfoForm := generateArgForms()
			err := argConfirmForm.Run()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			if !argConfirmForm.GetBool("confirm") {
				break
			}
			err = argInfoForm.Run()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			argsInput = argInfoForm.GetString("name")
			argType = argInfoForm.GetString("type")
			argDescription = argInfoForm.GetString("description")
			args = append(args, ArgInfo{Name: argsInput, Type: argType, Description: argDescription})
		}
		commandList = append(commandList, CommandInfo{Name: command, Description: description, Args: args, ReturnType: returnType})
	}

	filePath := filepath.Join(rootDir, "src", "cogs", filename+".py")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

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
		return
	}

	for _, command := range commandList {
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
			return
		}

		for _, arg := range command.Args {
			_, err = fmt.Fprintf(file, `                    %s (%s): %s
`, arg.Name, arg.Type, arg.Description)
			if err != nil {
				return
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
			return
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
			return
		}
	}

	err = file.Sync()
	if err != nil {
		return
	}
}

func oldAddCogs(filename string) error {
	rootDir, err := FindBotConf()
	if err != nil {
		fmt.Println("Error finding bot config:", err)
		return err
	}
	config, err := LoadConfig()
	Banner()

	if filename == "" {
		fmt.Println("Eenter the cog name:")
		fmt.Scanln(&filename)
	}
	filePath := filepath.Join(rootDir, "src", "cogs", filename+".py")
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

func generateCmdForms() (*huh.Form, *huh.Form) {
	cmdStartForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("confirm").
				Title("Do you want to add a command?").
				Affirmative("yes").
				Negative("no"),
		),
	)
	cmdInfoForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Enter the command name").
				Prompt("> ").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("command name cannot be empty")
					}
					if strings.Contains(s, " ") {
						return fmt.Errorf("command name cannot contain spaces")
					}
					return nil
				}),

			huh.NewText().
				Key("description").
				Title("Enter the command description").
				CharLimit(400).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("description cannot be empty")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Key("returnType").
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
	return cmdStartForm, cmdInfoForm
}

func generateArgForms() (*huh.Form, *huh.Form) {

	argStartForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("confirm").
				Title("Do you want to add an argument?").
				Affirmative("yes").
				Negative("no"),
		),
	)

	argInfoForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Enter the argument name").
				Prompt("> ").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("argument name cannot be empty")
					}
					if strings.Contains(s, " ") {
						return fmt.Errorf("argument name cannot contain spaces")
					}
					return nil
				}),

			huh.NewText().
				Key("description").
				Title("Enter the argument description").
				CharLimit(200).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("argument description cannot be empty")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Key("type").
				Title("Enter the argument type").
				Options(
					huh.NewOption("str", "str"),
					huh.NewOption("int", "int"),
					huh.NewOption("float", "float"),
					huh.NewOption("bool", "bool"),
					huh.NewOption("None", "None"),
				).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("argument type cannot be empty")
					}
					if s != "str" && s != "int" && s != "float" && s != "bool" && s != "None" {
						return fmt.Errorf("argument type must be one of str, int, float, bool, None")
					}
					return nil
				}),
		),
	)

	return argStartForm, argInfoForm
}

func confirmCommand() *huh.Form {
	cmdForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Does everything look correct?").
				Affirmative("yes").
				Negative("no"),
		),
	)

	return cmdForm
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
