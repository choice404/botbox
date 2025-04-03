# Bot Box

Make the discord bot of your dreams.

![](./readme_assets/botbox_showcase.gif)

## Table of Contents

- [About](#about)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Creating a Cog](#creating-a-cog)
- [Dependencies](#dependencies)
- [Development Scripts](#development-scripts)
- [Troubleshooting](#troubleshooting)
- [License](#license)
- [Contributors](#contributors)

---

## About

A discord bot template generator to help create discord bots quickly and easily. Forget about the boilerplate and focus on what really matters, what your bot will do.

**Bot Box** is built using Python and `discord.py`, offering a boilerplate codebase for quickly building Discord bot projects. It includes a cog-based architecture, `.env` management, and built-in utilities for automating bot configuration and extension development.

---

## Features

- Slash command support via `discord.ext.commands`
- Automated cog generation with `createCog.py`
- Project initialization with `.env` and `config.json` via `initialize.py`
- Auto-regeneration of `main.py` to match active cogs
- Easily extendable and modular design

---

## Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/your-username/botbox.git
   cd botbox
   ```

2. **Initialize the project**

   Run the following script to set up your `.env` and `config.json` files:

   ```bash
   python3 initialize.py
   ```

3. **Install dependencies**

   If not done during initialization, install manually with:

   ```bash
   pip install -r requirements.txt
   ```

---

## Usage

Once configured, you can start the bot using:

```bash
python3 main.py
```

On startup, the bot connects to Discord and syncs all slash commands. It also dynamically loads all cogs listed in `config.json`.

---

## Configuration

### `.env`

Stores environment-specific values. Created during setup or manually:

```
DISCORD_TOKEN=your_token_here
DISCORD_GUILD=your_guild_id
OTHER_ENV=custom_value
```

### `config.json`

Defines the bot name, command prefix, and active cogs:

```json
{
  "bot": {
    "name": "botbox",
    "command_prefix": "!"
    "author": "Austin \"Choice404\" Choi"
    "description": "The bot description"
  },
  "cogs": [
    {
      "name": "HelloWorld",
      "commands": ["hello"]
    }
  ]
}
```

---

## Creating a Cog

To create a new command set (cog):

```bash
python3 createCog.py
```

You'll be prompted to define:

- Cog name
- Command names and descriptions
- Argument names, types, and descriptions
- Return types

The cog will be saved in `cogs/` and automatically added to `config.json`. It also triggers a rewrite of `main.py` using `updateMain.py`.

---

## Dependencies

See [`requirements.txt`](requirements.txt):

This file was generated using pipreqs to help make installing modules simple. If you add new modules and dependencies, use pipreqs to update this file in case of distribution and collaboration.

- `discord.py==2.3.2`
- `python-dotenv==1.0.1`

---

## Development Scripts

- `initialize.py` – Initializes your bot setup with `.env` and `config.json`
- `createCog.py` – Interactive script for creating new cog modules
- `updateMain.py` – Regenerates `main.py` based on `config.json`

---

## Troubleshooting

- **Missing `config.json`?** Run `initialize.py` to generate it.
- **Cogs not loading?** Check `config.json` for correct names and verify files exist in `cogs/`.
- **Token errors?** Make sure your `.env` file is present and properly formatted.

---

## TODO
- [ ] Refactor?
  - [ ] Considering refactoring this into a golang cli tool or some other low level compiled language.
  - [ ] Maybe a npm cli tool for this?
    - [ ] If so then should use discord.js library instead of discord.py
- [ ] Expand config.json to include the following
  - [ ] More details about each command provided to the createCog script
  - [ ] Expected bot response
- [ ] Linting script inputs for any potential errors before creating any files
- [ ] Should also run pipreqs when creating new files to update pip dependencies


---

## Version History

- 1.0.0 - Initial version which includes basic features such as generate basic boilerplate code for cogs and the main file

## License

This project is licensed under the [MIT License](LICENSE), © 2025 [Austin \"Choice404\" Choi](https://github.com/choice404).

---

## Contributors

- **[Austin Choi](https://github.com/choice404)** — Original author and maintainer
