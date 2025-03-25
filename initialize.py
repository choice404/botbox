"""
Copyright © 2025 Austin Choi
See end of file for extended copyright information
"""

import os
import json

def main():
    print("""
   ________          _           ____        __ 
  / ____/ /_  ____  (_)_______  / __ )____  / /_
 / /   / __ \\/ __ \\/ / ___/ _ \\/ __  / __ \\/ __/
/ /___/ / / / /_/ / / /__/  __/ /_/ / /_/ / /_  
\\____/_/ /_/\\____/_/\\___/\\___/_____/\\____/\\__/  
                                                
    """)

    if not os.path.exists("config.json"):
        setBotConfig()

    if not os.path.exists(".env"):
        setEnv()

    install = input("Do you want to install the dependencies? (y/N): ").lower or "n"

    if install == "y":
        os.system("pip install -r requirements.txt")
        print("Dependencies installed.")

def setEnv() -> None:
    token: str = input("Enter your discord bot token: ")
    guild: str = input("Enter your discord bot guild (default: 'None'): ") or 'None'

    environs = []
    while True:
        key: str = input("Enter the environment variable key (enter '!' to stop): ")
        if key == "!":
            break
        value: str = input("Enter the environment variable value: ")
        environs.append((key.capitalize(), value))

    with open(".env", "w") as envFile:
        envFile.write(f"DISCORD_TOKEN={token}\n")
        envFile.write(f"DISCORD_GUILD={guild}\n")
        for key, value in environs:
            envFile.write(f"{key}={value}\n")


def setBotConfig() -> None:
    name: str = input("Enter the name of your discord bot: ")

    command_prefix: str = input("Enter the command prefix for your bot (default: '!'): ") or '!'

    author: str = input("Enter the author of the bot: ")

    description: str = input("Enter the description of the bot: ") or "A generic bot built using Bot Box (https://github.com/choice404/botbox"

    bot = {
        "bot": {
            "name": name,
            "command_prefix": command_prefix,
            "author": author,
            "description": description,
        },
        "cogs": [
            {
                "name": "HelloWorld",
                "commands": [
                    "hello"
                ]
            }
        ]
    }

    json_bot = json.dumps(bot, indent=4)

    with open("config.json", "w") as f:
        f.write(json_bot)


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
