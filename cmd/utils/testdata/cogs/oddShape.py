"""
Bot Author: Austin Choi

TestBot
A discord bot used by the parser tests
"""

import discord
from discord import app_commands
from discord.ext import commands

GUILD = discord.Object(id=0)

class OddShape(commands.Cog, name="OddShape"):
    def __init__(self, bot) -> None:
        self.bot = bot

    @app_commands.command(
        name="split",
        description="Decorator split across several lines",
    )
    @app_commands.guilds(GUILD)
    async def split(self, interaction: discord.Interaction) -> None:
        """
        Decorator split across several lines
        """

        return None

    @app_commands.command(description="Decorator without a name")
    async def nameless(self, interaction: discord.Interaction) -> None:
        """
        Decorator without a name
        """

        return None

    @commands.command()
    def not_async(self, ctx: commands.Context) -> None:
        """
        Not a coroutine so there is no command to record
        """

        return None
