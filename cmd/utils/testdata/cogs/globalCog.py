"""
Bot Author: Austin Choi

TestBot
A discord bot used by the parser tests
"""

import discord
from discord import app_commands
from discord.ext import commands

class GlobalCog(commands.Cog, name="GlobalCog"):
    def __init__(self, bot) -> None:
        self.bot = bot

    @app_commands.command(name="echo", description="Echoes a message")
    @app_commands.describe(message="The message to echo")
    async def echo(self, interaction: discord.Interaction, message: str) -> str:
        """
        Echoes a message when the user types "/echo"

            Parameters:
                    message (str): The message to echo

            Returns:
                    str
        """

        return message

async def setup(bot):
    await bot.add_cog(GlobalCog(bot))
