"""
Copyright © 2025 Austin Choi
See end of file for extended copyright information
"""

import json
import os
from updateMain import *

def main() -> int:
    if not os.path.exists("cogs"):
        os.makedirs("cogs")
    try:
        print("""
   ________          _           ____        __ 
  / ____/ /_  ____  (_)_______  / __ )____  / /_
 / /   / __ \\/ __ \\/ / ___/ _ \\/ __  / __ \\/ __/
/ /___/ / / / /_/ / / /__/  __/ /_/ / /_/ / /_  
\\____/_/ /_/\\____/_/\\___/\\___/_____/\\____/\\__/  
""")
        createCog()
    except FileNotFoundError:
        ("config.json not found, please use the projectInit.py script first.")
        return 1

    updateMain()
    return 0

def createCog() -> None:
    if not os.path.exists("config.json"):
        raise FileNotFoundError
    with open("config.json", "r") as f:
        config = json.load(f)
    commands = []
    name = input("Enter the cog name: ")
    command = ""
    while True:
        command = input("Enter the command name (enter '!' to stop): ")
        if command == "!":
            break
        description: str = input("Enter the command description: ")
        return_type = input("Enter the command return type (str, int, float, bool, None): ")
        args = []
        while True:
            argument = input("Enter the argument name (enter '!' to stop): ")
            if argument == "!":
                break
            arg_type = input("Enter the argument type (str, int, float, bool): ")
            arg_description= input("Enter the argument description: ")
            args.append({"name": argument, "type": arg_type, "description": arg_description})
            print()
        commands.append({"name": command, "description": description, "args": args, "return_type": return_type})
        print("\n")

    with open(f"cogs/{name}.py", "w") as f:
        f.write(
f"""\"""
Bot Author {config['bot']['author']}

{config['bot']['name']}
{config['bot']['description']}
\"""

import discord
from discord import app_commands
from discord.ext import commands
from dotenv import load_dotenv

load_dotenv()

class {name[0].upper()}{name[1:]}(commands.Cog):
    def __init__(self, bot) -> None:
        self.bot = bot
        print("{name} cog loaded")
""")

        for command in commands:
            full_arg = ""
            for arg in command["args"]:
                full_arg += f"{arg['name']}: {arg['type']}, "
            f.write( f"""
    @app_commands.command(name="{command['name']}", description="{command['description']}")
    async def {command['name']}(self, interaction: discord.Interaction, {full_arg[:-2]}) -> {command['return_type']}:
        \"""
        {command['description']} when the user types "/{command['name']}"
            
            Parameters:
""")
            for arg in command["args"]:
                f.write( f"""                    {arg['name']} ({arg['type']}): {arg['description']}
""")
            f.write(f"""
            Returns:
                    {command['return_type']}
        \"""

        try:
            await interaction.response.send_message(f"{command['name']}", ephemeral=True)
        except Exception as e:
            print(f"Error: {{e}}")
            await interaction.response.send_message(f"Error: {{e}}", ephemeral=True)
""")
            ret_type = ""
            if command["return_type"] == "str":
                ret_type = "\"\""
            elif command["return_type"] == "int":
                ret_type = 0
            elif command["return_type"] == "float":
                ret_type = 0.0
            elif command["return_type"] == "bool":
                ret_type = False
            f.write(f"""
        return {ret_type}
""")
                
    cog = {
        "name": name[0].upper() + name[1:],
        "file": name[0].lower() + name[1:],
        "commands": []
    }
    for command in commands:
        cog["commands"].append(command["name"])

    config["cogs"].append(cog)

    with open("config.json", "w") as f:
        f.write(json.dumps(config, indent=4))


if __name__ == '__main__':
    main()

"""
Copyright © 2025 Austin "Choice404" Choi

https://github.com/choice404/botbox

Bot Box
A discord bot template generator to help create discord bots quickly and easily

This code is licensed under the MIT License.
Please see the LICENSE file in the root directory of this project for the full license details.
"""
