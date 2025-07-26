/*
Copyright © 2025 Austin "Choice404" Choi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project management commands",
	Long:  `Commands for managing Bot Box project configurations and maintenance.`,
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade botbox.conf to the latest schema",
	Long: `Upgrade your botbox.conf file to the latest schema format.

This command will:
1. Read your existing botbox.conf file
2. Parse your cog files to extract detailed command information
3. Create a backup of your original config
4. Upgrade the configuration to the latest schema format

The upgrade process will attempt to preserve all existing configuration
while adding missing fields and converting legacy command arrays to
the new CommandInfo format.

Examples:
  botbox project upgrade                # Upgrade current project config
`,
	RunE: runConfigUpgrade,
}

func runConfigUpgrade(cmd *cobra.Command, args []string) error {
	_, err := utils.FindBotConf()
	if err != nil {
		return fmt.Errorf("not in a botbox project: %w", err)
	}

	fmt.Println("🔍 Analyzing botbox.conf for upgrade...")

	result, err := utils.UpgradeConfig()
	if err != nil {
		return fmt.Errorf("upgrade failed: %w", err)
	}

	if result.AlreadyUpgraded {
		fmt.Printf("✅ %s\n", result.Message)
		return nil
	}

	if result.Success {
		fmt.Printf("✅ %s\n", result.Message)

		if result.BackupCreated {
			fmt.Printf("💾 Backup created: %s\n", result.BackupPath)
		}

		if len(result.UpgradedCogs) > 0 {
			fmt.Printf("📦 Upgraded cogs: %s\n", formatList(result.UpgradedCogs))
		}

		if len(result.Errors) > 0 {
			fmt.Println("\n⚠️  Warnings during upgrade:")
			for _, errMsg := range result.Errors {
				fmt.Printf("  • %s\n", errMsg)
			}
		}

		fmt.Println("\n🎉 Upgrade completed successfully!")
		fmt.Println("💡 Run 'botbox config sync' to ensure everything is synchronized.")

		return nil
	}

	return fmt.Errorf("upgrade completed with errors")
}

func formatList(items []string) string {
	if len(items) == 0 {
		return "none"
	}
	if len(items) == 1 {
		return items[0]
	}
	if len(items) == 2 {
		return items[0] + " and " + items[1]
	}

	result := ""
	for i, item := range items {
		if i == len(items)-1 {
			result += "and " + item
		} else if i == len(items)-2 {
			result += item + " "
		} else {
			result += item + ", "
		}
	}
	return result
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(upgradeCmd)
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
