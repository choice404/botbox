/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the bot",
	Long:  `Run the bot`,
	Run: func(cmd *cobra.Command, args []string) {
		rootDir, err := utils.FindBotConf()
		if err != nil {
			fmt.Println("Error finding root directory:", err)
			return
		}

		runCmd := exec.Command("bash", filepath.Join(rootDir, "run.sh"))

		output, err := runCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error running the bot:", err)
			fmt.Println("Output:", string(output))
			return
		}

		fmt.Println("Bot is running...")
		fmt.Println("Output:", string(output))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
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
