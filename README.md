<div align="center">
  <h1>🤖 Bot Box 🤖</h1>
  <p>Forget the boring parts. Make building bots fun!</p>
  <p>A powerful CLI tool to scaffold, configure, and manage your Discord bot projects quickly and efficiently.</p>
</div>

---

## 📖 Table of Contents

- [About](#-about)
- [Features](#-features)
- [Installation](#%EF%B8%8F-installation)
- [Usage](#-usage)
- [Configuration](#%EF%B8%8F-configuration)
- [Troubleshooting](#-troubleshooting)
- [Roadmap](#%EF%B8%8F-roadmap-todo)
- [Version History](#-version-history)
- [License](#%EF%B8%8F-license)
- [Contributors](#-contributors)

---

## ✨ About

**Bot Box** is your ultimate companion for creating Discord bots with ease. Forget the tedious boilerplate and dive straight into building the unique functionalities that make your bot stand out.

Built with **Go**, `cobra-cli`, and `huh`, Bot Box offers an intuitive command-line interface to quickly generate Discord bot projects. It champions a **cog-based architecture** for modularity, simplifies `.env` management, and provides built-in utilities to automate bot configuration and extension development.

---

## 🚀 Features

-   **Slash Command Support**: Seamless integration via `discord.ext.commands`.
-   **Automated Cog Generation**: Generate new cogs with predefined commands and arguments effortlessly.
-   **Project Initialization**: Quick setup with `.env` and `botbox.conf` files.
-   **Dynamic cog maintenance**: Used `botbox.conf` to dynamically load cogs and provide an interface for users to load, reload, and unload cogs as needed.
-   **Modular Design**: Easily extendable and maintainable structure.

---

## 🛠️ Installation

Bot Box requires **Go** to be installed on your system.

### Prerequisites: Go Installation

It is highly recommended to follow the official [Golang documentation](https://go.dev/doc/install) for the most up-to-date installation instructions.

Alternatively, you can use `webi` for a quick installation:

**Linux/macOS:**

```sh
curl -sS https://webi.sh/golang | sh; \
source ~/.config/envman/PATH.env
```

**Windows:**

```powershell
curl.exe https://webi.ms/golang | powershell
```

### Install Bot Box CLI

Once Go is installed, use the following command to install Bot Box:

```sh
go install github.com/choice404/botbox/v2@latest
```

---

## 💡 Usage

Bot Box provides several commands to help you manage your bot project.

### Create a new Bot Box project

```sh
botbox create
```

This command will prompt you to provide project details (like bot name, prefix, etc.) and then generate a new project with initial files.

### Initialize a Bot Box project in the current directory

```sh
botbox init
```

Use this command to set up a new Bot Box project in your current working directory.

### Add a new cog to the current Bot Box project

```sh
botbox add [name]
```

It is highly recommended to use `botbox add` to generate new cogs within your projects. While you *can* manually add cogs, this will incur overhead as you'll need to manually update `botbox.conf` and potentially `main.py`.

You'll be prompted to define:

-   Cog name (if not provided in the command)
-   Command names and descriptions
-   Argument names, types, and descriptions
-   Return types

The new cog will be saved in the `cogs/` directory and automatically registered in `botbox.conf`.

### Remove a cog from the current Bot Box project

```sh
botbox remove [name]
```

You'll be prompted to select a cog to remove if `[name]` is not provided. This command will also update `botbox.conf` accordingly.

### Display current Bot Box project configuration

```sh
botbox config [flags]
```

This command displays the configuration details of your current Bot Box project as defined in `botbox.conf`.

---

## ⚙️ Configuration

### `botbox.conf`

This is the central configuration file for your Bot Box project, crucial for dynamically loading, reloading, and unloading cogs via `/src/main.py` and `/src/cogs/cogs.py`.

The Bot Box CLI tool automatically keeps `botbox.conf` synchronized with your project. If you choose to add or remove cogs manually (without the CLI tool), you *must* manually update `botbox.conf` or manage cog loading within `/src/main.py` and `/src/cogs/cogs.py`.

Example `botbox.conf` structure:

```json
{
  "bot": {
    "name": "botbox",
    "command_prefix": "!",
    "author": "Austin \"Choice404\" Choi",
    "description": "The bot description"
  },
  "cogs": [
    {
      "name": "HelloWorld",
      "file": "helloWorld",
      "slash_commands": [
        "hello"
      ],
      "prefix_commands": []
    }
  ]
}
```

---

## 🐛 Troubleshooting

-   **Missing `botbox.conf`?** Run `botbox init` in your project directory to generate it.
-   **Cogs not loading?**
    -   Verify cog names in `botbox.conf` match the actual file names in the `cogs/` directory.
    -   Ensure the cog files exist.
-   **Token errors?** Make sure your `.env` file is present in the project root and contains `DISCORD_BOT_TOKEN=YOUR_TOKEN_HERE`.

---

## 🛣️ Roadmap (TODO)

-   [ ] Expand `botbox.conf` to include:
    -   [ ] More details about each command provided via `botbox add`.
    -   [ ] Expected bot responses.
-   [ ] New commands:
    -   [ ] `botbox edit` command for modifying existing cogs/commands.
-   [ ] Advanced dynamic bot building:
    -   [ ] Create cogs containing "blocks" that can be dynamically added and connected for complex functionality.

---

## 📜 Version History

-   **2.3.0**: Added `remove` command; updated `botbox.conf` schema; improved prefix command generation support.
-   **2.2.6**: Fixed `add` command bug; added MIT license to generated files; updated project template.
-   **2.2.4**: Improved `README.md` generation; token input is now hidden.
-   **2.2.0**: Configuration transitioned to `botbox.conf`; `add` command modified.
-   **2.1.0**: Fixed issues with setting prefixes, license creation; added flags to `config` command.
-   **2.0.4**: Updated `README.md` with installation instructions and v2 imports.
-   **2.0.3**: Enforced single non-alphanumeric character custom prefixes in CLI form.
-   **2.0.2**: Brew release via taps; updated imports to use GitHub paths.
-   **2.0.1**: GitHub releases implemented using `goreleaser`.
-   **2.0.0**: Major refactor of the project in Go (Python scrapped for CLI core).
-   **1.0.0**: Initial version with basic boilerplate generation for cogs and main file.

*(Minor patch versions between major/minor releases are omitted for brevity)*

---

## ⚖️ License

This project is licensed under the [MIT License](LICENSE), © 2025 [Austin \"Choice404\" Choi](https://github.com/choice404).

---

## 🤝 Contributors

-   **[Austin Choi](https://github.com/choice404)** — Original author and maintainer
