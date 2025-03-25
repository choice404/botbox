"""
Copyright © 2025 Austin Choi
See end of file for extended copyright information
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


"""
Copyright © 2025 Austin "Choice404" Choi

https://github.com/choice404/botbox

Bot Box
A discord bot template generator to help create discord bots quickly and easily

This code is licensed under the MIT License.
Please see the LICENSE file in the root directory of this project for the full license details.
"""
