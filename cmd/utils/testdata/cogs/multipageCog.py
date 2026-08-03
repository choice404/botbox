"""
Bot Author: Austin Choi

TestBot
A discord bot used by the parser tests
"""

import discord
from discord import app_commands
from discord.ext import commands
from dotenv import load_dotenv
import os

try:
    from utils.logger import get_logger
    logger = get_logger(__name__)
except ImportError:
    import logging
    logger = logging.getLogger(__name__)

load_dotenv()

GUILD_ID = int(os.getenv("DISCORD_GUILD", 0))
GUILD = discord.Object(id=GUILD_ID)

import json

SURVEY_FLOW = json.loads(r'''
{
    "Pages": [
        {
            "Name": "start",
            "Title": "Survey start",
            "Fields": [
                {
                    "Name": "track",
                    "Label": "Track",
                    "Style": "short",
                    "Required": true,
                    "Placeholder": "backend or frontend"
                }
            ],
            "Branches": [
                {
                    "Field": "track",
                    "Equals": "backend",
                    "Goto": "backend"
                }
            ],
            "Next": "wrap"
        },
        {
            "Name": "backend",
            "Title": "Backend questions",
            "Fields": [
                {
                    "Name": "language",
                    "Label": "Favorite language",
                    "Style": "short",
                    "Required": true,
                    "Placeholder": ""
                }
            ],
            "Branches": null,
            "Next": "wrap"
        },
        {
            "Name": "wrap",
            "Title": "Wrap up",
            "Fields": [
                {
                    "Name": "comments",
                    "Label": "Comments",
                    "Style": "paragraph",
                    "Required": false,
                    "Placeholder": ""
                }
            ],
            "Branches": null,
            "Next": ""
        }
    ],
    "Responses": [
        {
            "Type": "message",
            "Content": "Thanks, track {track} recorded",
            "Ephemeral": true
        }
    ]
}
''')

SURVEY_PAGES = {page["Name"]: page for page in SURVEY_FLOW["Pages"]}

class SafeDict(dict):
    def __missing__(self, key):
        return "{" + key + "}"

class SurveyStartModal(discord.ui.Modal, title="Survey start"):
    track = discord.ui.TextInput(label="Track", style=discord.TextStyle.short, required=True, placeholder="backend or frontend")

    def __init__(self, cog):
        super().__init__()
        self.cog = cog

    async def on_submit(self, interaction: discord.Interaction):
        session = self.cog.survey_sessions.setdefault(interaction.user.id, {})
        session["track"] = self.track.value
        await survey_advance(self.cog, interaction, "start", session)

class SurveyBackendModal(discord.ui.Modal, title="Backend questions"):
    language = discord.ui.TextInput(label="Favorite language", style=discord.TextStyle.short, required=True)

    def __init__(self, cog):
        super().__init__()
        self.cog = cog

    async def on_submit(self, interaction: discord.Interaction):
        session = self.cog.survey_sessions.setdefault(interaction.user.id, {})
        session["language"] = self.language.value
        await survey_advance(self.cog, interaction, "backend", session)

class SurveyWrapModal(discord.ui.Modal, title="Wrap up"):
    comments = discord.ui.TextInput(label="Comments", style=discord.TextStyle.paragraph, required=False)

    def __init__(self, cog):
        super().__init__()
        self.cog = cog

    async def on_submit(self, interaction: discord.Interaction):
        session = self.cog.survey_sessions.setdefault(interaction.user.id, {})
        session["comments"] = self.comments.value
        await survey_advance(self.cog, interaction, "wrap", session)

class SurveyContinueView(discord.ui.View):
    def __init__(self, cog, next_page, user_id):
        super().__init__(timeout=120)
        self.cog = cog
        self.next_page = next_page
        self.user_id = user_id
        self.message = None

    @discord.ui.button(label="Continue", style=discord.ButtonStyle.primary)
    async def continue_page(self, interaction: discord.Interaction, button: discord.ui.Button):
        await interaction.response.send_modal(SURVEY_MODALS[self.next_page](self.cog))

    async def on_timeout(self):
        for child in self.children:
            child.disabled = True
        if self.message is not None:
            await self.message.edit(view=self)
        self.cog.survey_sessions.pop(self.user_id, None)

async def survey_advance(cog, interaction, page_name, session):
    page = SURVEY_PAGES[page_name]
    next_page = page.get("Next") or ""
    for branch in page.get("Branches") or []:
        if session.get(branch["Field"]) == branch["Equals"]:
            next_page = branch["Goto"]
            break
    if not next_page:
        await survey_finish(cog, interaction, session)
        return
    view = SurveyContinueView(cog, next_page, interaction.user.id)
    await interaction.response.send_message(f"Continue to {SURVEY_PAGES[next_page]['Title']}", view=view, ephemeral=True)
    view.message = await interaction.original_response()

async def survey_finish(cog, interaction, session):
    responses = SURVEY_FLOW.get("Responses") or []
    if responses:
        content = responses[0]["Content"].format_map(SafeDict(session))
        ephemeral = bool(responses[0].get("Ephemeral"))
    else:
        content = "survey submitted: " + " ".join(f"{key}={value}" for key, value in session.items())
        ephemeral = True
    await interaction.response.send_message(content, ephemeral=ephemeral)
    cog.survey_sessions.pop(interaction.user.id, None)

SURVEY_MODALS = {
    "start": SurveyStartModal,
    "backend": SurveyBackendModal,
    "wrap": SurveyWrapModal,
}

class MultipageCog(commands.Cog, name="MultipageCog"):
    def __init__(self, bot) -> None:
        self.bot = bot
        self.survey_sessions = {}
        logger.info("MultipageCog cog loaded")

    @app_commands.command(name="survey", description="Runs a short survey")
    @app_commands.guilds(GUILD)
    async def survey(self, interaction: discord.Interaction) -> None:
        """
        Runs a short survey when the user types "/survey"

            Returns:
                    None
        """

        self.survey_sessions[interaction.user.id] = {}
        await interaction.response.send_modal(SurveyStartModal(self))


async def setup(bot):
    await bot.add_cog(MultipageCog(bot))

"""
File generated by BotBox - https://github.com/choice404/botbox
"""