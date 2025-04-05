/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a Bot Box project",
	Long: `Initializes a Bot Box project in the current directory and prompts the user for information about the bot as well as setup other default configurations in a botbot.conf file.
  Will also create the initial project strucutre
  /
  |- LICENSE
  |- README.md
  |- botbot.conf
  |- requirements.txt
  |- src/
     |- main.py |- cogs/
        |- __init__.py
        |- helloWorld.py
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
	},
}

func createProjectStructure() {
	directories := []string{
		"src",
		"src/cogs",
	}

	fmt.Println(`
    ____        __     ____            
   / __ )____  / /_   / __ )____  _  __
  / __  / __ \\/ __/  / __  / __ \\| |/_/
 / /_/ / /_/ / /_   / /_/ / /_/ />  <  
/_____/\\____/\\__/  /_____/\\____/_/|_|  
  `)

	if _, err := os.Stat("botbot.conf"); err == nil {
		fmt.Println("botbot.conf file already exists. Do you want to overwrite it? (y/N)")
		var response string
		fmt.Scanln(&response)
		if response != "y" {
			fmt.Println("Aborting...")
			return
		}
	}
	var botName string
	var botDescription string
	var botAuthor string
	var botPrefix string

	fmt.Println("Enter the name of your bot:")
	fmt.Scanln(&botName)
	fmt.Println("Enter a description of your bot:")
	fmt.Scanln(&botDescription)
	fmt.Println("Enter the author of your bot:")
	fmt.Scanln(&botAuthor)
	fmt.Println("Enter the command prefix for your bot (default: '!'):")
	fmt.Scanln(&botPrefix)

	if botPrefix == "" {
		botPrefix = "!"
	}

	var envChoice bool
	fmt.Println("Do you want to create a .env file? (y/N)")
	fmt.Scanln(&envChoice)

	if envChoice {
		createEnvFile()
	}

	for _, dir := range directories {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	files := []string{
		"LICENSE",
		"README.md",
		"botbot.conf",
		"requirements.txt",
	}

	for _, file := range files {
		f, err := os.Create(file)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", file, err)
			return
		}
		f.Close()
	}
}

func createEnvFile() {
	var botToken string
	var botGuild string
	fmt.Println("Enter the bot token:")
	fmt.Scanln(&botToken)

	fmt.Println("Enter the bot guild ID:")
	fmt.Scanln(&botGuild)

	envFile, err := os.Create(".env")
	if err != nil {
		fmt.Printf("Error creating .env file: %v\n", err)
		return
	}
	defer envFile.Close()
	envFile.WriteString("BOT_TOKEN=" + botToken + "\n")
	envFile.WriteString("BOT_GUILD=" + botGuild + "\n")
}

func createDopplerFile() {
	var projectName string
	fmt.Print("Enter your doppler project name: ")
	fmt.Scanln(&projectName)

	var envName string
	fmt.Print("Enter your doppler env name: ")
	fmt.Scanln(&envName)

	dopplerFile, err := os.Create("doppler.yaml")
	if err != nil {
		fmt.Printf("Error creating doppler.yaml file: %v\n", err)
		return
	}
	defer dopplerFile.Close()

	dopplerFile.WriteString("setup:")
	dopplerFile.WriteString(" - project: " + projectName + "\n")
	dopplerFile.WriteString("   config: " + envName + "\n")
}

func init() {
	rootCmd.AddCommand(initCmd)
}

/*
Copyright © 2025 2025 Austin "Choice404" Choi

Bot Box
A discord bot template generator to help create discord bots quickly and easily

This code is licensed under the MIT License.
Please see the LICENSE file in the root directory of this project for the full license details.
*/
