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

	slashCommandList  []CommandInfo
	prefixCommandList []CommandInfo
	command           string
	commandConfirm    bool
	description       string
	returnType        string
	cmdType           string
	args              []ArgInfo
	argsInput         string
	argType           string
	argDescription    string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds new cog to the project",
	Long: `This command allows you to add a new cog to your bot project.
You can specify the cog structure as well as create commands for the cog.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fileNameForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("filename").
					Title("Enter the filename").
					Prompt("> ").
					Validate(func(s string) error {
						err := validateFileName(s)
						if err != nil {
							return err
						}
						return nil
					}),
			),
		)

		Banner()
		filename := ""
		if len(args) > 0 {
			filename = args[0]
			err := validateFileName(filename)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		if filename == "" {
			err := fileNameForm.Run()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			filename = fileNameForm.GetString("filename")
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
		cmdConfirmForm, cmdInfoForm, cmdAcceptForm := generateCmdForms()
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
		cmdType = cmdInfoForm.GetString("cmdType")
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
			argsInput = ""
			argType = ""
			argDescription = ""
		}
		err = cmdAcceptForm.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if cmdAcceptForm.GetBool("confirm") {
			if cmdType == "slash" {
				slashCommandList = append(slashCommandList, CommandInfo{Name: command, Description: description, Args: args, ReturnType: returnType})
			} else if cmdType == "prefix" {
				prefixCommandList = append(prefixCommandList, CommandInfo{Name: command, Description: description, Args: args, ReturnType: returnType})
			}
		}
		command = ""
		description = ""
		returnType = ""
		cmdType = ""
		args = []ArgInfo{}
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

class %s(commands.Cog, name="%s"):
    def __init__(self, bot) -> None:
        self.bot = bot
        print("%s cog loaded")
`, config.BotInfo.Author, config.BotInfo.Name, config.BotInfo.Description, className, className, filename)
	if err != nil {
		return
	}

	for _, command := range slashCommandList {
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
			fmt.Println("Error writing return statement:", err)
			return
		}
	}

	for _, command := range prefixCommandList {
		var argBuilder strings.Builder
		for i, arg := range command.Args {
			fmt.Fprintf(&argBuilder, "%s: %s", arg.Name, arg.Type)
			if i < len(command.Args)-1 {
				argBuilder.WriteString(", ")
			}
		}
		fullArgStr := argBuilder.String()

		_, err = fmt.Fprintf(file, `
    @commands.command()
    async def %s(self, interaction: discord.Interaction, %s) -> %s:
        """
        %s when the user types "/%s"

            Parameters:
`, command.Name, fullArgStr, command.ReturnType, command.Description, command.Name)
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
			fmt.Println("Error writing return statement:", err)
			return
		}
	}

	_, err = fmt.Fprintf(file, `

async def setup(bot):
    await bot.add_cog(%s(bot))
  `, className)

	err = file.Sync()
	if err != nil {
		fmt.Println("Error syncing file:", err)
		return
	}

	cog := CogConfig{
		Name:           strings.ToUpper(string(filename[0])) + filename[1:],
		File:           strings.ToLower(string(filename[0])) + filename[1:],
		SlashCommands:  []string{},
		PrefixCommands: []string{},
	}

	for _, slashCommand := range slashCommandList {
		cog.SlashCommands = append(cog.SlashCommands, slashCommand.Name)
	}

	for _, prefixCommand := range prefixCommandList {
		cog.PrefixCommands = append(cog.PrefixCommands, prefixCommand.Name)
	}

	config.Cogs = append(config.Cogs, cog)

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("failed to marshal config to JSON: %w", err)
		return
	}

	confDir, err := FindBotConf()
	if err != nil {
		fmt.Println("Error: %w", err)
		return
	}

	confPath := filepath.Join(confDir, "botbox.conf")

	err = os.WriteFile(confPath, jsonData, 0644)
	if err != nil {
		fmt.Println("failed to write updated botbox.conf: %w", err)
		return
	}

}

func generateCmdForms() (*huh.Form, *huh.Form, *huh.Form) {
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
					if commandExists(s, append(slashCommandList, prefixCommandList...)) {
						return fmt.Errorf("command name already exists")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Key("cmdType").
				Title("Select the command type").
				Options(
					huh.NewOption("slash", "slash"),
					huh.NewOption("prefix", "prefix"),
				).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("return type cannot be empty")
					}
					if s != "slash" && s != "prefix" {
						return fmt.Errorf("return type must either slash or prefix")
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
	cmdAcceptForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Command Info").
				Description(fmt.Sprintf("Command Name: %s\nCommand Type: %s\nDescription: %s\nReturn Type: %s\nArguments: %v", command, cmdType, description, returnType, args)),
			huh.NewConfirm().
				Title("Does everything look correct?").
				Key("confirm").
				Affirmative("yes").
				Negative("no"),
		),
	)

	return cmdStartForm, cmdInfoForm, cmdAcceptForm
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
						return fmt.Errorf("Argument name cannot be empty")
					}
					if strings.Contains(s, " ") {
						return fmt.Errorf("Argument name cannot contain spaces")
					}
					if argExists(s, args) {
						return fmt.Errorf("Argument name already exists")
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

func commandExists(commandName string, commandList []CommandInfo) bool {
	for _, cmd := range commandList {
		if cmd.Name == commandName {
			return true
		}
	}
	return false
}

func argExists(argName string, args []ArgInfo) bool {
	for _, arg := range args {
		if arg.Name == argName {
			return true
		}
	}
	return false
}

func fileExists(fileName string) bool {
	rootDir, err := FindBotConf()
	filePath := filepath.Join(rootDir, "src", "cogs", fileName+".py")
	_, err = os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func validateFileName(fileName string) error {
	if fileExists(fileName) {
		return fmt.Errorf("file with name '%s' already exists", fileName)
	}
	if fileName == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	if strings.Contains(fileName, " ") {
		return fmt.Errorf("filename cannot contain spaces")
	}
	if strings.Contains(fileName, ".") || strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		return fmt.Errorf("filename cannot contain '.' or '/' or '\\'")
	}
	if strings.Contains(fileName, "-") || strings.Contains(fileName, ":") || strings.Contains(fileName, "*") || strings.Contains(fileName, "?") || strings.Contains(fileName, "\"") {
		return fmt.Errorf("filename cannot contain '-', ':', '*', '?', or '\"'")
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

MIT License

Copyright (c) 2025 Austin

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
