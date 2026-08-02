"""
Bot Author: Austin Choi


A discord bot used by the parser tests
"""

import discord
from discord.ext import commands

class EmptyName(commands.Cog, name="EmptyName"):
    def __init__(self, bot) -> None:
        self.bot = bot

async def setup(bot):
    await bot.add_cog(EmptyName(bot))
