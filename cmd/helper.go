/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"unicode"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

var (
	botName                string
	botDescription         string
	botAuthor              string
	botPrefix              string
	botTokenDopplerProject string
	botGuildDopplerEnv     string
	envChoice              string
	licenseType            string
	licenseText            string
)

func Banner() {
	fmt.Println(`
    ____        __     ____            
   / __ )____  / /_   / __ )____  _  __
  / __  / __ \/ __/  / __  / __ \| |/_/
 / /_/ / /_/ / /_   / /_/ / /_/ />  <  
/_____/\____/\__/  /_____/\____/_/|_|  
  `)
}

func FindBotConf() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	originalDir := currentDir

	for {
		confDir := filepath.Join(currentDir)

		_, err := os.Stat(filepath.Join(confDir, "botbox.conf"))
		if err == nil {
			confPath, err := filepath.Abs(filepath.Join(confDir, "botbox.conf"))
			if err != nil {
				return "", fmt.Errorf("failed to get absolute path of %s: %w", confPath, err)
			}

			return confDir, nil
		}

		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("error checking file %s: %w", confDir, err)
		}

		parentDir := filepath.Dir(currentDir)

		if parentDir == currentDir {
			break
		}

		currentDir = parentDir
	}

	return "", fmt.Errorf("Not a botbox project: %s", originalDir)
}

func BotBoxCreate(actionCallback func()) {
	Banner()
	botBoxConfigForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the name of your bot").
				Prompt("> ").
				Value(&botName).
				Validate(func(s string) error {
					if botName == "" {
						return fmt.Errorf("bot name cannot be empty")
					}
					if len(s) > 20 {
						return fmt.Errorf("bot name is too long")
					}
					r := []rune(s)[0]
					if !unicode.IsLetter(r) {
						return fmt.Errorf("bot name must start with a letter")
					}
					return nil
				}),

			huh.NewText().
				Title("Enter a description of your bot").
				Value(&botDescription).
				CharLimit(100),

			huh.NewInput().
				Title("Enter the author of your bot").
				Prompt("> ").
				Value(&botAuthor).
				Validate(func(s string) error {
					if botAuthor == "" {
						return fmt.Errorf("author name cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Title("Enter the command prefix for your bot (default: '!')").
				Prompt("> ").
				Value(&botPrefix).
				Validate(func(s string) error {
					if s == "" {
						botPrefix = "!"
						return nil
					}
					if len(s) > 1 {
						return fmt.Errorf("command prefix must be a single character")
					}
					r := []rune(s)[0]

					if unicode.IsLetter(r) || unicode.IsDigit(r) {
						return fmt.Errorf("command prefix can not be an alphanumeric character")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How do you want to handle environment variables?").
				Options(
					huh.NewOption("Create a .env file", "env"),
					huh.NewOption("Use Doppler", "doppler"),
				).
				Value(&envChoice),
			huh.NewInput().
				TitleFunc(func() string {
					if envChoice == "env" {
						return "Enter the bot token"
					}
					return "Enter the Doppler project name"
				}, &envChoice).
				Prompt("> ").
				EchoMode(huh.EchoModePassword).
				Validate(func(s string) error {
					if envChoice == "env" {
						if s == "" {
							return fmt.Errorf("token cannot be empty")
						}
						if len(s) < 10 {
							return fmt.Errorf("token is too short")
						}
					}
					return nil
				}).
				Value(&botTokenDopplerProject),
			huh.NewInput().
				TitleFunc(func() string {
					if envChoice == "env" {
						return "Enter the bot guild ID"
					}
					return "Enter the Doppler environment name"
				}, &envChoice).
				Prompt("> ").
				Value(&botGuildDopplerEnv),
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
				Value(&licenseType),
		),
	)

	botBoxConfigForm.Run()

	err := spinner.New().
		Title("Creating project structure...").
		Action(actionCallback).
		Run()

	if err != nil {
		fmt.Println("Error creating project structure:", err)
		return
	}
}

func FetchLicense(licenseKey string) (string, error) {
	if licenseKey == "" || licenseKey == "none" {
		return "", fmt.Errorf("no license key provided or selected 'none'")
	}

	apiURL := fmt.Sprintf("https://api.github.com/licenses/%s", licenseKey)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request to %s: %w", apiURL, err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "bot-box")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch license %s: %w", licenseKey, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to fetch license %s: status %s, body: %s",
			licenseKey, resp.Status, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body for %s: %w", licenseKey, err)
	}

	var licenseResp LicenseResponse
	err = json.Unmarshal(bodyBytes, &licenseResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response for %s: %w", licenseKey, err)
	}

	if licenseResp.Body == "" {
		return "", fmt.Errorf("no license body found in response for %s", licenseKey)
	}

	return licenseResp.Body, nil
}

func CreateProject(rootDir string) {
	if rootDir[len(rootDir)-1:] != "/" {
		rootDir += "/"
	}

	directories := []string{
		"src",
		"src/cogs",
	}

	if confOpt, err := CreateFileOption(rootDir + "botbox.conf"); err == nil && confOpt {
		confFile, err := os.Create(rootDir + "botbox.conf")
		if err != nil {
			fmt.Printf("Error creating botbox.conf file: %v\n", err)
			return
		}
		defer confFile.Close()
		_, err = fmt.Fprintf(confFile, `{
  "bot": {
    "name": "%s",
    "command_prefix": "%s",
    "author": "%s",
    "description": "%s" 
  },
  "cogs": [
    {
      "name": "HelloWorld",
      "file": "helloWorld",
      "commands": [
          "hello"
      ]
    }
  ]
}`, botName, botPrefix, botAuthor, botDescription)
	} else if err == nil && !confOpt {
		fmt.Println("Not overriding botbox.conf file.")
	} else {
		fmt.Printf("Error creating botbox.conf file: %v\n", err)
		return
	}

	if readmeOpt, err := CreateFileOption(rootDir + "README.md"); err == nil && readmeOpt {
		readmeFile, err := os.Create(rootDir + "README.md")
		if err != nil {
			fmt.Printf("Error creating README.md file: %v\n", err)
			return
		}
		defer readmeFile.Close()
		_, err = fmt.Fprintf(readmeFile, `# %s

## Table of Contents

- [About](#about)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [License](#license)
- [Contributors](#contributors)

## About
%s

### Author
%s

## Installation
1. Clone the repository
2. Install the dependencies
3. Run the bot
4. Enjoy!

## Usage
1. Run the bot with the command `+"`./run.sh`"+`
2. Use the command prefix to run commands
3. Use the command name to run commands
4. Use the command arguments to run commands

## License
`, botName, botDescription, botAuthor)
		if err != nil {
			fmt.Printf("Error writing to README.md file: %v\n", err)
			return
		}

		if licenseType != "no-license" && licenseType != "" {
			_, err = fmt.Fprintf(readmeFile, `This project is licensed under the %s License - see the [LICENSE](LICENSE) file for details.
    `, licenseType)
			if err != nil {
				fmt.Printf("Error writing to README.md file: %v\n", err)
				return
			}
		} else {
			_, err = fmt.Fprint(readmeFile, `All rights reserved.`)
		}

		_, err = fmt.Fprintf(readmeFile, ` ## Contributors

- %s`, botAuthor)
		if err != nil {
			fmt.Printf("Error writing to README.md file: %v\n", err)
			return
		}

	} else if err == nil && !readmeOpt {
		fmt.Println("Not overriding README.md file.")
	} else {
		fmt.Printf("Error creating README.md file: %v\n", err)
		return
	}

	if licenseType != "no-license" && licenseType != "" {
		if licenseOpt, err := CreateFileOption(rootDir + "LICENSE"); err == nil && licenseOpt {
			licenseFile, err := os.Create(rootDir + "LICENSE")
			if err != nil {
				fmt.Printf("Error creating LICENSE file: %v\n", err)
				return
			}
			defer licenseFile.Close()
			licenseText, err = FetchLicense(licenseType)

			if err != nil {
				fmt.Printf("Error fetching license text: %v\n", err)
				return
			}

			_, err = fmt.Fprint(licenseFile, licenseText)
			if err != nil {
				fmt.Printf("Error writing to LICENSE file: %v\n", err)
				return
			}
		} else if err == nil && !licenseOpt {
			fmt.Println("Not overriding LICENSE file.")
		} else {
			fmt.Printf("Error creating LICENSE file: %v\n", err)
			return
		}
	}

	if envChoice == "doppler" {
		if dopplerOpt, err := CreateFileOption(rootDir + "doppler.yaml"); err == nil && dopplerOpt {
			dopplerFile, err := os.Create(rootDir + "doppler.yaml")
			if err != nil {
				fmt.Printf("Error creating doppler.yaml file: %v\n", err)
				return
			}
			defer dopplerFile.Close()
			_, err = fmt.Fprintf(dopplerFile, `setup:
  - project: %s
    config: %s
`, botTokenDopplerProject, botGuildDopplerEnv)
		} else if err == nil && !dopplerOpt {
			fmt.Println("Not overriding doppler.yaml file.")
		} else {
			fmt.Printf("Error creating doppler.yaml file: %v\n", err)
			return
		}
	} else if envChoice == "env" {
		if envOpt, err := CreateFileOption(rootDir + ".env"); err == nil && envOpt {
			envFile, err := os.Create(rootDir + ".env")
			if err != nil {
				fmt.Printf("Error creating .env file: %v\n", err)
				return
			}
			defer envFile.Close()
			_, err = fmt.Fprintf(envFile, `DISCORD_TOKEN=%s
DISCORD_GUILD=%s
`, botTokenDopplerProject, botGuildDopplerEnv)
			if err != nil {
				fmt.Printf("Error writing to .env file: %v\n", err)
				return
			}
		} else if err == nil && !envOpt {
			fmt.Println("Not overriding .env file.")
		} else {
			fmt.Printf("Error creating .env file: %v\n", err)
			return
		}
	} else if envChoice == "none" {
		fmt.Println("No environment file will be created.")
	} else {
		fmt.Println("Invalid environment choice.")
		return
	}

	if runOpt, err := CreateFileOption(rootDir + "run.sh"); err == nil && runOpt {
		runFile, err := os.Create(rootDir + "run.sh")
		if err != nil {
			fmt.Printf("Error creating run.sh file: %v\n", err)
			return
		}
		defer runFile.Close()
		_, err = fmt.Fprint(runFile, `#!/bin/bash
# This script will run the bot
# Make sure to give it execute permissions with chmod +x run.sh
# and run it with ./run.sh
# If you are using Doppler, make sure to run doppler run -- ./run.sh
# If you are using a .env file, make sure to run source .env before running this script

python3 src/main.py
`)
		if err != nil {
			fmt.Printf("Error writing to run.sh file: %v\n", err)
			return
		}
		err = os.Chmod(rootDir+"run.sh", 0755)
		if err != nil {
			fmt.Printf("Error setting permissions for run.sh file: %v\n", err)
			return
		}
	} else if err == nil && !runOpt {
		fmt.Println("Not overriding run.sh file.")
	} else {
		fmt.Printf("Error creating run.sh file: %v\n", err)
		return
	}

	for _, dir := range directories {
		err := os.MkdirAll(rootDir+dir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	if mainOpt, err := CreateFileOption(rootDir + "src/main.py"); err == nil && mainOpt {
		mainFile, err := os.Create(rootDir + "src/main.py")
		if err != nil {
			fmt.Printf("Error creating main.py file: %v\n", err)
			return
		}
		defer mainFile.Close()
		_, err = fmt.Fprint(mainFile, `"""
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information for BotBox
"""

import discord
from discord.ext import commands, tasks
from dotenv import load_dotenv
import os
from cogs import helloWorld
import json

class Bot(commands.Bot):
    def __init__(self):
        with open('config.json') as f:
            config = json.load(f)
        self.name =  config['bot']['name']
        intents = discord.Intents.all()
        intents.message_content = True
        super().__init__(command_prefix = config['bot']['command_prefix'], intents=intents, help_command = None)
        self.synced = False

    async def syncing(self):
        if not self.synced:
            await self.tree.sync()
            self.synced = True
        print(f"Synced slash commands for {self.user}")

    async def on_command_error(self, ctx, error):
        await ctx.reply(error, ephemeral = True)

load_dotenv()
bot = Bot()
TOKEN = str(os.getenv('DISCORD_TOKEN'))
GUILD = os.getenv('DISCORD_GUILD')

@bot.event
async def on_ready():
    print(f'{bot.user} has connected to Discord!')
    print(f'Connected to guild: {GUILD}')
    await bot.add_cog(helloWorld.HelloWorld(bot))
    await bot.syncing()

def main():
    print(f"{bot.name} is starting up...")
    bot.run(TOKEN)

if __name__ == '__main__':
    main()

"""
Copyright © 2025 Austin "Choice404" Choi

https://github.com/choice404/botbox

Bot Box
The template entry point for the bot.

This code and BotBox is licensed under the MIT License.
https://github.com/choice404/botbox/license
"""
`)
		if err != nil {
			fmt.Printf("Error writing to main.py file: %v\n", err)
			return
		}
		err = os.Chmod(rootDir+"src/main.py", 0755)
		if err != nil {
			fmt.Printf("Error setting permissions for main.py file: %v\n", err)
			return
		}
	}

	if helloWorldOpt, err := CreateFileOption(rootDir + "src/cogs/helloWorld.py"); err == nil && helloWorldOpt {
		helloWorldFile, err := os.Create(rootDir + "src/cogs/helloWorld.py")
		if err != nil {
			fmt.Printf("Error creating helloWorld.py file: %v\n", err)
			return
		}
		defer helloWorldFile.Close()
		_, err = fmt.Fprintf(helloWorldFile, `"""
Bot Author %s

%s
%s
"""

import discord
from discord import app_commands
from discord.ext import commands
from dotenv import load_dotenv

load_dotenv()

class HelloWorld(commands.Cog):
    def __init__(self, bot) -> None:
        self.bot = bot
        print("HelloWorld cog loaded")

    @app_commands.command(name="hello", description="Bot responds with world")
    async def hello(self, interaction: discord.Interaction) -> None:
        """
        Bot responds with "world" when the user types "/hello"

            Parameters:
                    interaction (discord.Interaction): The interaction object that triggered the command

            Returns:
                    None
        """

        try:
            await interaction.response.send_message(f"world", ephemeral=True)
        except Exception as e:
            print(f"Error: {e}")
            await interaction.response.send_message(f"Error: {e}", ephemeral=True)
  `, botAuthor, botName, botDescription)

		if err != nil {
			fmt.Printf("Error writing to helloWorld.py file: %v\n", err)
			return
		}
		err = os.Chmod(rootDir+"src/cogs/helloWorld.py", 0755)
		if err != nil {
			fmt.Printf("Error setting permissions for helloWorld.py file: %v\n", err)
			return
		}
	}

	if initOpt, err := CreateFileOption(rootDir + "src/cogs/__init__.py"); err == nil && initOpt {
		_, err := os.Create(rootDir + "src/cogs/__init__.py")
		if err != nil {
			fmt.Printf("Error creating __init__.py file: %v\n", err)
			return
		}
	}

	fmt.Println("Project structure created successfully!")

}

func CreateFileOption(filename string) (bool, error) {
	var override bool
	formTitle := fmt.Sprintf("The file %s already exists. Do you want to override it?", filename)
	overrideForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(formTitle).
				Affirmative("yes").
				Negative("no").
				Value(&override),
		),
	)
	if _, err := os.Stat(filename); err == nil {
		overrideForm.Run()
		if override {
			return true, nil
		} else {
			return false, nil
		}
	}
	return true, nil
}

func LoadConfig() (Config, error) {
	var cfg Config

	confDir, err := FindBotConf()
	if err != nil {
		return cfg, fmt.Errorf("failed to find config directory: %w", err)
	}

	confPath := filepath.Join(confDir, "botbox.conf")

	jsonData, err := os.ReadFile(confPath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file %s: %w", confPath, err)
	}

	err = json.Unmarshal(jsonData, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to parse config JSON from %s: %w", confPath, err)
	}

	return cfg, nil
}

/*
Copyright © 2025 Austin "Choice404" Choi

https://github.com/choice404/botbox

Bot Box
A discord bot template generator to help create discord bots quickly and easily

This code is licensed under the MIT License.
Please see the LICENSE file in the root directory of this project for the full license details.
*/
