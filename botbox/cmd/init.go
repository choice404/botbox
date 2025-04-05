/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

type LicenseResponse struct {
	Body string `json:"body"`
}

var (
	botName        string
	botDescription string
	botAuthor      string
	botPrefix      string
	botToken       string
	botGuild       string
	envChoice      string
	dopplerProject string
	dopplerEnv     string
	licenseType    string
	licenseChoice  bool
	licenseText    string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a Bot Box project",
	Long: `Initializes a Bot Box project in the current directory and prompts the user for information about the bot as well as setup other default configurations in a botbox.conf file.
  Will also create the initial project strucutre
  /
  |- README.md
  |- botbox.conf
  |- run.sh
  |- ?LICENSE
  |- ?doppler.yaml
  |- src/
     |- main.py |- cogs/
        |- __init__.py
        |- helloWorld.py
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		botBoxInit()
	},
}

func botBoxInit() {
	fmt.Println(`
    ____        __     ____            
   / __ )____  / /_   / __ )____  _  __
  / __  / __ \\/ __/  / __  / __ \\| |/_/
 / /_/ / /_/ / /_   / /_/ / /_/ />  <  
/_____/\\____/\\__/  /_____/\\____/_/|_|  
  `)

	botBoxConfigForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the name of your bot").
				Prompt("> ").
				Value(&botName),

			huh.NewText().
				Title("Enter a description of your bot").
				Value(&botDescription),

			huh.NewInput().
				Title("Enter the author of your bot").
				Prompt("> ").
				Value(&botAuthor),
			huh.NewInput().
				Title("Enter the command prefix for your bot (default: '!')").
				Prompt("> ").
				Value(&botPrefix),
		),
	)

	botBoxEnvForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How do you want to handle environment variables?").
				Options(
					huh.NewOption("Create a .env file", "env"),
					huh.NewOption("Use Doppler", "doppler"),
				).
				Value(&envChoice),
		),
	)
	botBoxLicenceForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to add a license?").
				Affirmative("yes").
				Negative("no").
				Value(&licenseChoice),
		),
	)
	envForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the bot token").
				Prompt("> ").
				Validate(func(s string) error {
					if len(s) < 10 {
						return fmt.Errorf("token is too short")
					}
					return nil
				}).
				Value(&botToken),
			huh.NewInput().
				Title("Enter the bot guild ID").
				Prompt("> ").
				Value(&botGuild),
		),
	)
	dopplerForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the Doppler project name").
				Prompt("> ").
				Value(&dopplerProject),
			huh.NewInput().
				Title("Enter the Doppler environment name").
				Prompt("> ").
				Value(&dopplerEnv),
		),
	)
	licenceForm := huh.NewForm(
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
	botBoxEnvForm.Run()

	if envChoice == "env" {
		envForm.Run()
	} else if envChoice == "doppler" {
		dopplerForm.Run()
	}

	botBoxLicenceForm.Run()

	if licenseChoice {
		licenceForm.Run()
		var err error
		licenseText, err = fetchLicense(licenseType)
		if err != nil {
			fmt.Println("Error fetching license:", err)
			return
		}
	}

	err := spinner.New().
		Title("Creating project structure...").
		Action(createProjectStructure).
		Run()

	if err != nil {
		fmt.Println("Error creating project structure:", err)
		return
	}
}

func createProjectStructure() {
	directories := []string{
		"src",
		"src/cogs",
	}

	if confOpt, err := createFileOption("botbox.conf"); err == nil && confOpt {
		confFile, err := os.Create("botbox.conf")
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
      "name": "helloWorld",
      "commands": [
          "hello"
      ]
    ]
  }
}`, botName, botPrefix, botAuthor, botDescription)
	} else if err == nil && !confOpt {
		fmt.Println("Not overriding botbox.conf file.")
	} else {
		fmt.Printf("Error creating botbox.conf file: %v\n", err)
		return
	}

	if readmeOpt, err := createFileOption("README.md"); err == nil && readmeOpt {
		readmeFile, err := os.Create("README.md")
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
### License
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
This project is licensed under the %s License - see the [LICENSE](LICENSE) file for details.

## Contributors

- %s
`, botName, botDescription, botAuthor, licenseType, licenseType, botAuthor)
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

	if licenseOpt, err := createFileOption("LICENSE"); err == nil && licenseOpt {
		licenseFile, err := os.Create("LICENSE")
		if err != nil {
			fmt.Printf("Error creating LICENSE file: %v\n", err)
			return
		}
		defer licenseFile.Close()
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

	if envChoice == "doppler" {
		if dopplerOpt, err := createFileOption("doppler.yaml"); err == nil && dopplerOpt {
			dopplerFile, err := os.Create("doppler.yaml")
			if err != nil {
				fmt.Printf("Error creating doppler.yaml file: %v\n", err)
				return
			}
			defer dopplerFile.Close()
			_, err = fmt.Fprintf(dopplerFile, `setup:
  - project: %s
    config: %s
`, dopplerProject, dopplerEnv)
		} else if err == nil && !dopplerOpt {
			fmt.Println("Not overriding doppler.yaml file.")
		} else {
			fmt.Printf("Error creating doppler.yaml file: %v\n", err)
			return
		}
	} else if envChoice == "env" {
		if envOpt, err := createFileOption(".env"); err == nil && envOpt {
			envFile, err := os.Create(".env")
			if err != nil {
				fmt.Printf("Error creating .env file: %v\n", err)
				return
			}
			defer envFile.Close()
			_, err = fmt.Fprintf(envFile, `DISCORD_TOKEN=%s
DISCORD_GUILD=%s
`, botToken, botGuild)
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

	if runOpt, err := createFileOption("run.sh"); err == nil && runOpt {
		runFile, err := os.Create("run.sh")
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
		err = os.Chmod("run.sh", 0755)
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
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	mainFile, err := os.Create("src/main.py")
	if err != nil {
		fmt.Printf("Error creating main.py file: %v\n", err)
		return
	}
	defer mainFile.Close()
	_, err = fmt.Fprint(mainFile, `"""
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
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
        print(f"Synced slash commands for {{self.user}}")

    async def on_command_error(self, ctx, error):
        await ctx.reply(error, ephemeral = True)

load_dotenv()
bot = Bot()
TOKEN = str(os.getenv('DISCORD_TOKEN'))
GUILD = os.getenv('DISCORD_GUILD')

@bot.event
async def on_ready():
    print(f'{{bot.user}} has connected to Discord!')
    print(f'Connected to guild: {{GUILD}}')
    await bot.add_cog(helloWorld.HelloWorld(bot))
    await bot.syncing()

def main():
    print(f"{{bot.name}} is starting up...")
    bot.run(TOKEN)

if __name__ == '__main__':
    main()

"""
Copyright © 2025 Austin "Choice404" Choi

https://github.com/choice404/botbox

Bot Box
The tempalte entry point for the bot.

This code is licensed under the MIT License.
https://github.com/choice404/botbox/license
"""
`)
	if err != nil {
		fmt.Printf("Error writing to main.py file: %v\n", err)
		return
	}
	err = os.Chmod("src/main.py", 0755)
	if err != nil {
		fmt.Printf("Error setting permissions for main.py file: %v\n", err)
		return
	}

	helloWorldFile, err := os.Create("src/cogs/helloWorld.py")
	if err != nil {
		fmt.Printf("Error creating helloWorld.py file: %v\n", err)
		return
	}
	defer helloWorldFile.Close()
	_, err = fmt.Fprintf(helloWorldFile, `"""
Bot Author {config['bot']['author']}

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
`, botName, botDescription)

	if err != nil {
		fmt.Printf("Error writing to helloWorld.py file: %v\n", err)
		return
	}
	err = os.Chmod("src/cogs/helloWorld.py", 0755)
	if err != nil {
		fmt.Printf("Error setting permissions for helloWorld.py file: %v\n", err)
		return
	}

	fmt.Println("Project structure created successfully!")

	_, err = os.Create("src/cogs/__init__.py")
	if err != nil {
		fmt.Printf("Error creating __init__.py file: %v\n", err)
		return
	}

}

func createFileOption(filename string) (bool, error) {
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

func fetchLicense(licenseKey string) (string, error) {
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

func init() {
	rootCmd.AddCommand(initCmd)
}

/*
Copyright © 2025 2025 Austin "Choice404" Choi

Bot Box
A discord bot template generator to help create discord bots quickly and easily

This code is licensed under the MIT License.
Please see the LICENSE file in the root directory of this project for the full license details.
*/
