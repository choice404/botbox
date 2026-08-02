/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Manage Docker files for the current Bot Box project",
	Long: `Manage Docker support for an existing Bot Box project.

Use the init subcommand to generate a Dockerfile, docker-compose.yml, and
.dockerignore tuned to the project's environment provider (env file or Doppler).`,
}

var dockerInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate Docker files for the current project",
	Long: `Generate Docker files in the root of the current Bot Box project.

This command creates:
  - Dockerfile running the bot as a non root user from /app
  - docker-compose.yml with restart policy, log volume, and env wiring
  - .dockerignore keeping secrets and local clutter out of the image

The environment wiring follows bot.env_provider from botbox.conf. Projects
created before that key existed fall back to detecting doppler.yaml or .env
in the project root. Existing files are skipped unless --force is given.`,
	Run: func(cmd *cobra.Command, args []string) {
		runDockerInit(cmd)
	},
}

/**
 * runDockerInit
 * Generates the docker files for the project the command runs inside
 * @param cmd {*cobra.Command} - the command holding the flags
 * @return ...
 **/
func runDockerInit(cmd *cobra.Command) {
	// This command is flag driven, so existing files are skipped instead of prompting
	utils.HeadlessMode = true

	rootDir, err := utils.FindBotConf()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Current directory is not in a botbox project.")
		os.Exit(1)
	}

	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	force, _ := cmd.Flags().GetBool("force")

	// An unset flag falls back to the global default, then the built in default
	pythonVersion, _ := cmd.Flags().GetString("python")
	if pythonVersion == "" && GlobalConfig != nil {
		pythonVersion = GlobalConfig.Defaults.PythonVersion
	}
	if pythonVersion == "" {
		pythonVersion = utils.DefaultPythonVersion
	}

	envProvider := utils.ResolveEnvProvider(config, rootDir)

	written, err := utils.GenerateDockerFiles(rootDir, pythonVersion, envProvider, force)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	for _, path := range written {
		fmt.Println(path)
	}
}

func init() {
	rootCmd.AddCommand(dockerCmd)
	dockerCmd.AddCommand(dockerInitCmd)

	dockerInitCmd.Flags().String("python", "", "Python base image version (defaults to defaults.python_version, then 3.11)")
	dockerInitCmd.Flags().Bool("force", false, "Overwrite existing files without prompting")
}

/*
Copyright © 2025 Austin "Choice404" Choi

https://github.com/choice404/botbox

Bot Box
A discord bot template generator to help create discord bots quickly and easily

This code is licensed under the MIT License.

MIT License

Copyright (c) 2025 Austin Choi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
