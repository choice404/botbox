/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new Bot Box project",
	Long: `Creates a directory containing a new Bot Box project.
  The project directory will contain the following file structure:

  projectName/
  |- README.md
  |- botbox.conf
  |- run.sh
  |- LICENSE (optional)
  |- doppler.yaml (optional)
  |- src/
     |- main.py
     |- cogs/
        |- __init__.py
        |- helloWorld.py
  `,
	Run: func(cmd *cobra.Command, args []string) {
		var createNewProject bool
		botBoxExistsForm := huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("The current directory is in an existing Bot Box project."),
				huh.NewConfirm().
					Title("Would you like to create a new Bot Box project?").
					Affirmative("yes").
					Negative("no").
					Value(&createNewProject),
			),
		)

		_, err := FindBotConf()
		if err == nil {
			botBoxExistsForm.Run()
			if !createNewProject {
				return
			}
		}
		BotBoxCreate(createProjectWrapper)
	},
}

func createProjectWrapper() {
	rootDir := botName
	if _, err := os.Stat(rootDir); err == nil && !os.IsNotExist(err) {
		fmt.Printf("Directory %s already exists. Please choose a different name.\n", rootDir)
		return
	} else if os.IsNotExist(err) {
		err = os.MkdirAll(rootDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", rootDir, err)
			return
		}
	} else {
		fmt.Printf("Error checking directory %s: %v\n", rootDir, err)
		return
	}

	CreateProject(rootDir)
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
