"""
Copyright © 2025 Austin Choi
See end of file for extended copyright information
"""

import json
import os

def main() -> int:
    try:
        updateMain()
    except FileNotFoundError:
        ("config.json not found, please use the initialize.py script first.")
        return 1
    return 0

def updateMain() -> None:
    if not os.path.exists("cogs"):
        os.makedirs("cogs")
    if not os.path.exists("config.json"):
        raise FileNotFoundError
    with open("config.json", "r") as f:
        config = json.load(f)

    with open("main.py", "w") as f:
        f.write(f"""\"""
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
\"""

import discord
from discord.ext import commands, tasks
from dotenv import load_dotenv
import os
""")
        cog_files = ""
        for cog in config["cogs"]:
            cog_files += f"{cog['file']}, "
        f.write(f"from cogs import {cog_files[:-2]}\n")
        f.write( f"""import json

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
""")
        for cog in config["cogs"]:
            f.write(f"    await bot.add_cog({cog['file']}.{cog['name']}(bot))\n")
        f.write(f"""    await bot.syncing()
def main():
    print(f"{{bot.name}} is starting up...")
    bot.run(TOKEN)

if __name__ == '__main__':
    main()

\"""
Copyright © 2025 Austin "Choice404" Choi

https://github.com/choice404/botbox

Bot Box
The tempalte entry point for the bot.

This code is licensed under the MIT License.
https://github.com/choice404/botbox/license
""\"""")

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
