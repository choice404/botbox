/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreateProject(rootDir string, values Values) error {
	directories := []string{
		"src",
		"src/cogs",
	}

	for _, dir := range directories {
		fullPath := filepath.Join(rootDir, dir)
		err := os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %w", fullPath, err)
		}
	}

	if confOpt, err := CreateFileOption(filepath.Join(rootDir, "botbox.conf")); err == nil && confOpt {
		confFile, err := os.Create(filepath.Join(rootDir, "botbox.conf"))
		if err != nil {
			return fmt.Errorf("Error creating botbox.conf file: %v\n", err)
		}
		defer confFile.Close()

		var confContent strings.Builder

		fmt.Fprintf(&confContent, `{
  "botbox": {
    "version": "%s"
  },
  "bot": {
    "name": "%s",
    "command_prefix": "%s",
    "author": "%s",
    "description": "%s"
  },`, Version, *values.Map["botName"], *values.Map["botPrefix"], *values.Map["botAuthor"], *values.Map["botDescription"])
		confContent.WriteString(`
  "cogs": [
    {
      "name": "HelloWorld",
      "file": "helloWorld",
      "env": "development",
      "slash_commands": [
        {
          "Name": "hello",
          "Scope": "guild",
          "Type": "slash",
          "Description": "Bot responds with world",
          "Args": null,
          "ReturnType": "None"
        }
      ],
      "prefix_commands": []
    },
    {
      "name": "CogManagement",
      "file": "cogs",
      "env": "production",
      "slash_commands": [
        {
          "Name": "reload-cog",
          "Scope": "guild",
          "Type": "slash",
          "Description": "Reloads a cog by name",
          "Args": [
          {
            "Name": "cog_name",
            "Type": "str",
            "Description": "The name of the cog to reload (without .py cog)"
          }
          ],
          "ReturnType": "None"
        },
        {
          "Name": "reload-all-cogs",
          "Scope": "guild",
          "Type": "slash",
          "Description": "Reloads all cogs",
          "Args": null,
          "ReturnType": "None"
        },
        {
          "Name": "list-cogs",
          "Scope": "guild",
          "Type": "slash",
          "Description": "Lists all available cogs",
          "Args": null,
          "ReturnType": "None"
        },
        {
          "Name": "unload-cog",
          "Scope": "guild",
          "Type": "slash",
          "Description": "Unloads a cog by name",
          "Args": [
          {
            "Name": "cog_name",
            "Type": "str",
            "Description": "The name of the cog to unload (without .py cog)"
          }
          ],
          "ReturnType": "None"
        },
        {
          "Name": "load-cog",
          "Scope": "guild",
          "Type": "slash",
          "Description": "Loads a cog by name",
          "Args": [
          {
            "Name": "cog_name",
            "Type": "str",
            "Description": "The name of the cog to load (without .py cog)"
          }
          ],
          "ReturnType": "None"
        }
      ],
      "prefix_commands": []
    }
  ]
}`)

		_, err = confFile.WriteString(confContent.String())

	} else if err == nil && !confOpt {
		fmt.Println("Not overriding botbox.conf file.")
	} else {
		return fmt.Errorf("error creating botbox.conf file: %w", err)
	}

	if readmeOpt, err := CreateFileOption(filepath.Join(rootDir, "README.md")); err == nil && readmeOpt {
		readmeFile, err := os.Create(filepath.Join(rootDir, "README.md"))
		if err != nil {
			return fmt.Errorf("Error creating README.md file: %v\n", err)
		}
		defer readmeFile.Close()

		var readmeContent strings.Builder

		fmt.Fprintf(&readmeContent, `# %s

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
1. Install the required dependencies
`+"```bash"+`
# if necessary generate and activate a virtual environment
python3 -m venv venv
source venv/bin/activate

# install the dependencies using pip and the provided requirements.txt file
python3 -m pip install -r requirements.txt
`+"```"+`

2. Run the bot
`+"```bash"+`
# If the command does not have the requirements to run, you can run it with the command:
chmod +x run.sh

# Run the bot using the provided run.sh script
./run.sh
`+"```"+`

## License
`, *values.Map["botName"], *values.Map["botDescription"], *values.Map["botAuthor"])

		if *values.Map["licenseType"] != "no-license" && *values.Map["licenseType"] != "" {
			fmt.Fprintf(&readmeContent, `This project is licensed under the %s License - see the [LICENSE](LICENSE) file for details.
    `, *values.Map["licenseType"])
		} else {
			readmeContent.WriteString(`All rights reserved.`)
		}
		fmt.Fprintf(&readmeContent, `
    ## Contributors

- %s`, *values.Map["botAuthor"])

		readmeContent.WriteString(`
Bot generated using BotBox - https://github.com/choice404/botbox`)

		_, err = readmeFile.WriteString(readmeContent.String())
		if err != nil {
			return fmt.Errorf("error writing to README.md file: %w", err)
		}

	} else if err == nil && !readmeOpt {
		fmt.Println("Not overriding README.md file.")
	} else {
		return fmt.Errorf("error creating README.md file: %w", err)
	}

	if *values.Map["licenseType"] != "no-license" && *values.Map["licenseType"] != "" {
		if licenseOpt, err := CreateFileOption(filepath.Join(rootDir, "LICENSE")); err == nil && licenseOpt {
			licenseFile, err := os.Create(filepath.Join(rootDir, "LICENSE"))
			if err != nil {
				return fmt.Errorf("Error creating LICENSE file: %v\n", err)
			}
			defer licenseFile.Close()
			LicenseText, err := FetchLicense(*values.Map["licenseType"])

			if err != nil {
				return fmt.Errorf("Error fetching license %s: %v\n", *values.Map["licenseType"], err)
			}

			_, err = fmt.Fprint(licenseFile, LicenseText)
			if err != nil {
				return fmt.Errorf("Error writing to LICENSE file: %v\n", err)
			}
		} else if err == nil && !licenseOpt {
			fmt.Println("Not overriding LICENSE file.")
		} else {
			return fmt.Errorf("Error creating LICENSE file: %v\n", err)
		}
	}

	if *values.Map["envChoice"] == "doppler" {
		if dopplerOpt, err := CreateFileOption(filepath.Join(rootDir, "doppler.yaml")); err == nil && dopplerOpt {
			dopplerFile, err := os.Create(filepath.Join(rootDir, "doppler.yaml"))
			if err != nil {
				return fmt.Errorf("Error creating doppler.yaml file: %v\n", err)
			}
			defer dopplerFile.Close()
			_, err = fmt.Fprintf(dopplerFile, `setup:
  - project: %s
    config: %s
`, *values.Map["botTokenDopplerProject"], *values.Map["botGuildDopplerEnv"])
		} else if err == nil && !dopplerOpt {
			fmt.Println("Not overriding doppler.yaml file.")
		} else {
			return fmt.Errorf("Error creating doppler.yaml file: %v\n", err)
		}
	} else if *values.Map["envChoice"] == "env" {
		if envOpt, err := CreateFileOption(filepath.Join(rootDir, ".env")); err == nil && envOpt {
			envFile, err := os.Create(filepath.Join(rootDir, ".env"))
			if err != nil {
				return fmt.Errorf("Error creating .env file: %v\n", err)
			}
			defer envFile.Close()
			_, err = fmt.Fprintf(envFile, `DISCORD_TOKEN=%s
DISCORD_GUILD=%s
ENVIRONMENTS=production,development
`, *values.Map["botTokenDopplerProject"], *values.Map["botGuildDopplerEnv"])
			if err != nil {
				return fmt.Errorf("Error writing to .env file: %v\n", err)
			}
		} else if err == nil && !envOpt {
			fmt.Println("Not overriding .env file.")
		} else {
			return fmt.Errorf("Error creating .env file: %v\n", err)
		}
	} else if *values.Map["envChoice"] == "none" {
		fmt.Println("No environment file will be created.")
	} else {
		return fmt.Errorf("Invalid environment choice: %s", *values.Map["envChoice"])
	}

	if runOpt, err := CreateFileOption(filepath.Join(rootDir, "run.sh")); err == nil && runOpt {
		runFile, err := os.Create(filepath.Join(rootDir, "run.sh"))
		if err != nil {
			return fmt.Errorf("Error creating run.sh file: %v\n", err)
		}
		defer runFile.Close()

		var runScriptContent strings.Builder

		runScriptContent.WriteString(`#!/bin/bash
# This script will run the bot
# Make sure to give it execute permissions with chmod +x run.sh
# and run it with ./run.sh
# If you are using botbox to run your bot, you can run it with the command: botbox run

`)

		if *values.Map["envChoice"] == "doppler" {
			runScriptContent.WriteString("doppler run -- \\\n")
		}
		runScriptContent.WriteString(`python3 src/main.py

# Script generated by BotBox - https://github.com/choice404/botbox`)

		_, err = fmt.Fprint(runFile, runScriptContent.String())

		if err != nil {
			return fmt.Errorf("Error writing to run.sh file: %v\n", err)
		}
		err = os.Chmod(filepath.Join(rootDir, "run.sh"), 0755)
		if err != nil {
			return fmt.Errorf("Error setting permissions for run.sh file: %v\n", err)
		}
	} else if err == nil && !runOpt {
		fmt.Println("Not overriding run.sh file.")
	} else {
		return fmt.Errorf("Error creating run.sh file: %v\n", err)
	}

	if mainOpt, err := CreateFileOption(filepath.Join(rootDir, "src", "main.py")); err == nil && mainOpt {
		mainFile, err := os.Create(filepath.Join(rootDir, "src", "main.py"))
		if err != nil {
			return fmt.Errorf("Error creating main.py file: %v\n", err)
		}
		defer mainFile.Close()
		_, err = fmt.Fprintf(mainFile, `"""
Bot Author: %s

%s
%s
"""

import discord
from discord.ext import commands
from dotenv import load_dotenv
import os
import json

class Bot(commands.Bot):
    def __init__(self):
        with open('botbox.conf') as f:
            config = json.load(f)
        self.name =  config['bot']['name']
        self.environments = os.getenv('ENVIRONMENTS', 'production,development').split(',')
        intents = discord.Intents.all()
        intents.message_content = True
        super().__init__(command_prefix = config['bot']['command_prefix'], intents=intents, help_command = None)
        self.guild = discord.Object(id=os.getenv("DISCORD_GUILD", ""))
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
    print(f"{bot.name} is starting up...")

    with open('botbox.conf', 'r') as f:
        config = json.load(f)

    for cog_config in config['cogs']:
        if 'file' not in cog_config:
            print("❌ Cog configuration is missing 'file' key.")
            continue
        if 'name' not in cog_config:
            print("❌ Cog configuration is missing 'name' key.")
            continue
        if 'env' not in cog_config:
            print(f"❌ Cog configuration is missing 'env' key.")
            continue
        if cog_config['env'] not in bot.environments:
            print(f"❌ Skipping cog {cog_config['name']}: Not in current environments -  {bot.environments}")
            continue
        cog_file = cog_config['file']
        try:
            await bot.load_extension(f'cogs.{cog_file}')
            print(f"✅ Loaded cog: {cog_file}")
        except Exception as e:
            print(f"❌ Failed to load cog {cog_file}: {e}")

    await bot.syncing()
    print("Bot is ready!")

def main():
    print(f"{bot.name} is starting up...")
    bot.run(TOKEN)

if __name__ == '__main__':
    main()

"""
File generated by BotBox - https://github.com/choice404/botbox
"""`, *values.Map["botAuthor"], *values.Map["botName"], *values.Map["botDescription"])
		if err != nil {
			return fmt.Errorf("Error writing to main.py file: %v\n", err)
		}
		err = os.Chmod(filepath.Join(rootDir, "src", "main.py"), 0755)
		if err != nil {
			return fmt.Errorf("Error setting permissions for main.py file: %v\n", err)
		}
	}

	if helloWorldOpt, err := CreateFileOption(filepath.Join(rootDir, "src", "cogs", "helloWorld.py")); err == nil && helloWorldOpt {
		helloWorldFile, err := os.Create(filepath.Join(rootDir, "src", "cogs", "helloWorld.py"))
		if err != nil {
			return fmt.Errorf("Error creating helloWorld.py file: %v\n", err)
		}
		defer helloWorldFile.Close()
		_, err = fmt.Fprintf(helloWorldFile, `"""
Bot Author: %s

%s
%s

This is an example file. Delete using the command "botbox remove"
"""

import discord
from discord import app_commands
from discord.ext import commands
from dotenv import load_dotenv
import os

load_dotenv()

GUILD_ID = os.getenv('DISCORD_GUILD', 0)
GUILD = discord.Object(id=GUILD_ID)

class HelloWorld(commands.Cog):
    def __init__(self, bot) -> None:
        self.bot = bot

    @app_commands.command(name="hello", description="Bot responds with world")
    @app_commands.guilds(GUILD)
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

async def setup(bot):
    await bot.add_cog(HelloWorld(bot))

"""
File generated by BotBox - https://github.com/choice404/botbox
"""`, *values.Map["botAuthor"], *values.Map["botName"], *values.Map["botDescription"])

		if err != nil {
			return fmt.Errorf("Error writing to helloWorld.py file: %v\n", err)
		}
		err = os.Chmod(filepath.Join(rootDir, "src", "cogs", "helloWorld.py"), 0755)
		if err != nil {
			return fmt.Errorf("Error setting permissions for helloWorld.py file: %v\n", err)
		}
	}

	if cogsOpt, err := CreateFileOption(filepath.Join(rootDir, "src", "cogs", "cogs.py")); err == nil && cogsOpt {
		cogsFile, err := os.Create(filepath.Join(rootDir, "src", "cogs", "cogs.py"))
		if err != nil {
			return fmt.Errorf("Error creating cogs.py file: %v\n", err)
		}
		defer cogsFile.Close()
		_, err = fmt.Fprintf(cogsFile, `"""
Bot Author: %s

%s
%s
"""

import discord
from discord import app_commands
from discord.ext import commands
from dotenv import load_dotenv
import json
import os

load_dotenv()

GUILD_ID = int(os.getenv("GUILD_ID", 0))
GUILD = discord.Object(id=GUILD_ID)

class CogManagement(commands.Cog, name="Cog Management"):
    def __init__(self, bot):
        self.bot = bot

    @app_commands.command(name="reload-cog", description="Reloads a cog by name")
    @app_commands.describe(
        cog_name="The name of the cog to reload (without .py cog)"
    )
    @app_commands.guilds(GUILD)
    async def reload_cog(self, interaction: discord.Interaction, cog_name: str) -> None:
        """
        Reloads a cog by name.

            Parameters:
                interaction (discord.Interaction): The interaction context.
                cog_name (str): The name of the cog to reload (without .py cog).

            Returns:
                None
        """

        with open('botbox.conf', 'r') as f:
            config = json.load(f)

        if not any(cog['file'] == cog_name for cog in config['cogs']):
            await interaction.response.send_message(f'{cog_name} is not a valid cog name.', ephemeral=True)
            return

        try:
            await self.bot.reload_extension(f'cogs.{cog_name}')
            await interaction.response.send_message(f'✅ Reloaded {cog_name}', ephemeral=True)
        except Exception as e:
            await interaction.response.send_message(f'❌ Failed to reload {cog_name}: {e}', ephemeral=True)

        await self.bot.syncing()

    @app_commands.command(name="reload-all-cogs", description="Reloads all cogs")
    @app_commands.guilds(GUILD)
    async def reload_all_cogs(self, interaction: discord.Interaction) -> None:
        """
        Reloads all cogs.

            Parameters:
                interaction (discord.Interaction): The interaction context.

            Returns:
                None
        """

        with open('botbox.conf', 'r') as f:
            config = json.load(f)

        failed_cogs = []
        success_count = 0

        for cog_config in config['cogs']:
            cog_file = cog_config['file']
            try:
                await self.bot.reload_extension(f'cogs.{cog_file}')
                success_count += 1
            except Exception as e:
                failed_cogs.append(f"{cog_file}: {e}")

        if failed_cogs:
            await interaction.response.send_message(f'✅ Reloaded {success_count} cogs\n❌ Failed:\n- {"\n- ".join(failed_cogs)}', ephemeral=True)
        else:
            await interaction.response.send_message(f'✅ Successfully reloaded all {success_count} cogs!', ephemeral=True)

        await self.bot.syncing()

    @app_commands.command(name="list-cogs", description="Lists all available cogs")
    @app_commands.guilds(GUILD)
    async def list_cogs(self, interaction: discord.Interaction) -> None:
        """
        Lists all available cogs.

            Parameters:
                interaction (discord.Interaction): The interaction context.

            Returns:
                None
        """
        with open('botbox.conf', 'r') as f:
            config = json.load(f)

        cog_list = [cog['file'] for cog in config['cogs']]
        if cog_list:
            await interaction.response.send_message(f'Available cogs:\n- {"\n- ".join(cog_list)}', ephemeral=True)
        else:
            await interaction.response.send_message('No cogs available.', ephemeral=True)

    @app_commands.command(name="unload-cog", description="Unloads a cog by name")
    @app_commands.describe(
        cog_name="The name of the cog to unload (without .py cog)"
    )
    @app_commands.guilds(GUILD)
    async def unload_cog(self, interaction: discord.Interaction, cog_name: str) -> None:
        """
        Unloads a cog by name.

            Parameters:
                interaction (discord.Interaction): The interaction context.
                cog_name (str): The name of the cog to unload (without .py cog).

            Returns:
                None
        """

        with open('botbox.conf', 'r') as f:
            config = json.load(f)

        if not any(cog['file'] == cog_name for cog in config['cogs']):
            await interaction.response.send_message(f'{cog_name} is not a valid cog name.', ephemeral=True)
            return
        try:
            await self.bot.unload_extension(f'cogs.{cog_name}')
            await interaction.response.send_message(f'✅ Unloaded {cog_name}', ephemeral=True)
        except Exception as e:
            await interaction.response.send_message(f'❌ Failed to unload {cog_name}: {e}', ephemeral=True)
        await self.bot.syncing()

    @app_commands.command(name="load-cog", description="Loads a cog by name")
    @app_commands.describe(
        cog_name="The name of the cog to load (without .py cog)"
    )
    @app_commands.guilds(GUILD)
    async def load_cog(self, interaction: discord.Interaction, cog_name: str) -> None:
        """
        Loads a cog by name.

            Parameters:
                interaction (discord.Interaction): The interaction context.
                cog_name (str): The name of the cog to load (without .py cog).

            Returns:
                None
        """

        with open('botbox.conf', 'r') as f:
            config = json.load(f)

        if not any(cog['file'] == cog_name for cog in config['cogs']):
            await interaction.response.send_message(f'{cog_name} is not a valid cog name.', ephemeral=True)
            return
        try:
            await self.bot.load_extension(f'cogs.{cog_name}')
            await interaction.response.send_message(f'✅ Loaded {cog_name}', ephemeral=True)
        except Exception as e:
            await interaction.response.send_message(f'❌ Failed to load {cog_name}: {e}', ephemeral=True)
        await self.bot.syncing()

async def setup(bot):
    await bot.add_cog(CogManagement(bot))

"""
File generated by BotBox - https://github.com/choice404/botbox
"""`, *values.Map["botAuthor"], *values.Map["botName"], *values.Map["botDescription"])

		if err != nil {
			return fmt.Errorf("Error writing to cogs.py file: %v\n", err)
		}
		err = os.Chmod(filepath.Join(rootDir, "src", "cogs", "cogs.py"), 0755)
		if err != nil {
			return fmt.Errorf("Error setting permissions for cogs.py file: %v\n", err)
		}
	}

	if initOpt, err := CreateFileOption(filepath.Join(rootDir, "src", "cogs", "__init__.py")); err == nil && initOpt {
		_, err := os.Create(filepath.Join(rootDir, "src", "cogs", "__init__.py"))
		if err != nil {
			return fmt.Errorf("error creating __init__.py file: %w", err)
		}
	}

	return nil
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
