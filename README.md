# Bot Box

Make the discord bot of your dreams.

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

**Bot Box** is built using GO, cobra-cli, and huh, offering an intuitive cli tool to quickly build Discord bot projects. It includes a cog-based architecture, `.env` management, and built-in utilities for automating bot configuration and extension development.

---

## Features

- Slash command support via `discord.ext.commands`
- Automated cog generation
- Project initialization with `.env` and `botbox.conf`
- Auto-regeneration of `main.py` to match active cogs
- Easily extendable and modular design

---

## Installation

### macOS
```
brew install choice404/tap/botbox
```

### Windows/Linux (macOS alternative)
This method will require that the user installs go

It is recommended to use the [golang documentation](https://go.dev/doc/install) for installation

Or alternatively use the following commands for
**Linux/macOS**
```
curl -sS https://webi.sh/golang | sh; \
source ~/.config/envman/PATH.env
```
**Windows**
```
curl.exe https://webi.ms/golang | powershell
```

After this do the following command
```
go install github.com/choice404/botbox/v2@latest
```


## Usage

### Create a new Bot Box project
```sh
botbox create
```

This command will prompt the user to provide information about the project then create the project with initial files.

### Initialize a Bot Box project in the current directory
```sh
botbox init
```

This command will initialize a Bot Box project in the current directory.

### Add a new cog to the current Bot Box project
```sh
botbox add <name>
```

You'll be prompted to define:

- Cog name
- Command names and descriptions
- Argument names, types, and descriptions
- Return types

The cog will be saved in `cogs/` and automatically added to `config.json`. It also triggers a rewrite of `main.py` using `updateMain.py`.

### Display current Bot Box project configuration
```sh
botbox config
```

This command will display the configuration of the current Bot Box project.

---

## Configuration

### `.env`

Stores environment-specific values. Created during creation, initialization or manually:

```
DISCORD_TOKEN=your_token_here
DISCORD_GUILD=your_guild_id
OTHER_ENV=custom_value
```

### `botbox.conf`

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
      "file": helloWorld,
      "commands": ["hello"]
    }
  ]
}
```

---

## Troubleshooting

- **Missing `botbox.conf`?** Run `botbox initialize` to generate it.
- **Cogs not loading?** Check `botbox.conf` for correct names and verify files exist in `cogs/`.
- **Token errors?** Make sure your `.env` file is present and properly formatted.

---

## TODO
- [x] Refactor?
  - [x] Considering refactoring this into a golang cli tool or some other low level compiled language.
  - [ ] Maybe a npm cli tool for this?
    - [ ] If so then should use discord.js library instead of discord.py
- [ ] Expand botbox.conf to include the following
  - [ ] More details about each command provided to the `botbox add` command
  - [ ] Expected bot response
- [ ] New commands
  - [x] Remove command
  - [ ] Edit command
- [ ] Dynamic bot building
  - [ ] Create "blocks" that will be dynamically connected together for command functionality


---

## Version History

- 1.0.0 - Initial version which includes basic features such as generate basic boilerplate code for cogs and the main file.
- 2.0.0 - A major refactor of the project in golang. Scrapped python for this...
- 2.0.1 - Github releases using goreleaser.
- 2.0.2 - Brew release through taps and updated the imports to use github in the project.
- 2.0.3 - Updated the CLI form so custom prefixes are single non-alphanumeric characters only.
- 2.0.4 - Updated README.md to include instructions on how to install botbox cli and updated imports with v2.
- 2.0.5 - Updated version number.
- 2.1.0 - Fixed some more bugs like the issues with setting prefixes, license creation, and finally added flags to the config.
- 2.2.0 - Changed configuration to be through botbox.conf and modified add command.
- 2.2.1 - Updated version number.
- 2.2.2 - Updated version number.
- 2.2.3 - Updated version number.
- 2.2.4 - Updated README.md generation and token input is hidden.
- 2.2.5 - Updated version number.
- 2.2.6 - Fixed add command bug. Added MIT license to each file. Updated project template generation.
- 2.3.0 - Added remove function and updated botbox.conf schema. Added prefix command generation support.

## License

This project is licensed under the [MIT License](LICENSE), © 2025 [Austin \"Choice404\" Choi](https://github.com/choice404).

---

## Contributors

- **[Austin Choi](https://github.com/choice404)** — Original author and maintainer
