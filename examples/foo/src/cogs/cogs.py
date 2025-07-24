"""
Bot Author: foobar

foo
bar
"""

import discord
from discord import app_commands
from discord.ext import commands
from dotenv import load_dotenv
import json
import os

GUILD_ID = int(os.getenv("GUILD_ID", 0))
GUILD = discord.Object(id=GUILD_ID) if GUILD_ID else None

class CogManagement(commands.Cog, name="Cog Management"):
    def __init__(self, bot):
        self.bot = bot

    @app_commands.command(name="reload-cog", description="Reloads a cog by name")
    @app_commands.describe(
        cog_name="The name of the cog to reload (without .py cog)"
    )
    @app_commands.guilds(GUILD) if GUILD else app_commands.default_permissions
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

        self.bot.syncing()

    @app_commands.command(name="reload-all-cogs", description="Reloads all cogs")
    @app_commands.guilds(GUILD) if GUILD else app_commands.default_permissions
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

        self.bot.syncing()

    @app_commands.command(name="list-cogs", description="Lists all available cogs")
    @app_commands.guilds(GUILD) if GUILD else app_commands.default_permissions
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
    @app_commands.guilds(GUILD) if GUILD else app_commands.default_permissions
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
        self.bot.syncing()

    @app_commands.command(name="load-cog", description="Loads a cog by name")
    @app_commands.describe(
        cog_name="The name of the cog to load (without .py cog)"
    )
    @app_commands.guilds(GUILD) if GUILD else app_commands.default_permissions
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
        self.bot.syncing()

async def setup(bot):
    await bot.add_cog(CogManagement(bot))

"""
File generated by BotBox - https://github.com/choice404/botbox
"""