/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func UpgradeConfig() (*UpgradeResult, error) {
	rootDir, err := FindBotConf()
	if err != nil {
		return nil, fmt.Errorf("not in a botbox project: %w", err)
	}

	configPath := filepath.Join(rootDir, "botbox.conf")

	jsonData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var legacyConfig LegacyConfig
	if err := json.Unmarshal(jsonData, &legacyConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	var modernConfig Config
	if err := json.Unmarshal(jsonData, &modernConfig); err == nil {
		if len(modernConfig.Cogs) > 0 && isModernCogFormat(modernConfig.Cogs[0]) {
			return &UpgradeResult{
				AlreadyUpgraded: true,
				Message:         "Configuration is already in the latest format",
			}, nil
		}
	}

	fmt.Println("🔄 Upgrading botbox.conf to latest schema...")

	upgradedConfig := Config{
		BotBox: BotBoxConfig{
			Version: strings.TrimPrefix(Version, "v"),
		},
		BotInfo: legacyConfig.BotInfo,
		Cogs:    []CogConfig{},
	}

	if legacyConfig.BotBox != nil {
		upgradedConfig.BotBox = *legacyConfig.BotBox
		upgradedConfig.BotBox.Version = strings.TrimPrefix(Version, "v")
	}

	result := &UpgradeResult{
		UpgradedCogs: []string{},
		Errors:       []string{},
	}

	cogsDir := filepath.Join(rootDir, "src", "cogs")

	for _, legacyCog := range legacyConfig.Cogs {
		fmt.Printf("📝 Upgrading cog: %s\n", legacyCog.Name)

		upgradedCog := CogConfig{
			Name:           legacyCog.Name,
			Env:            legacyCog.Env,
			File:           legacyCog.File,
			SlashCommands:  []CommandInfo{},
			PrefixCommands: []CommandInfo{},
		}

		cogFilePath := filepath.Join(cogsDir, legacyCog.File+".py")
		if _, err := os.Stat(cogFilePath); err == nil {
			if parsedCog, err := parseCogFile(cogFilePath, legacyCog.File); err == nil {
				upgradedCog.SlashCommands = parsedCog.SlashCommands
				upgradedCog.PrefixCommands = parsedCog.PrefixCommands
			} else {
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to parse %s: %v", legacyCog.File, err))
				upgradedCog.SlashCommands = convertLegacyCommands(legacyCog.SlashCommands, "slash")
				upgradedCog.PrefixCommands = convertLegacyCommands(legacyCog.PrefixCommands, "prefix")
			}
		} else {
			result.Errors = append(result.Errors, fmt.Sprintf("Cog file %s.py not found, using legacy command names", legacyCog.File))
			upgradedCog.SlashCommands = convertLegacyCommands(legacyCog.SlashCommands, "slash")
			upgradedCog.PrefixCommands = convertLegacyCommands(legacyCog.PrefixCommands, "prefix")
		}

		upgradedConfig.Cogs = append(upgradedConfig.Cogs, upgradedCog)
		result.UpgradedCogs = append(result.UpgradedCogs, legacyCog.Name)
	}

	backupPath := configPath + ".backup"
	if err := os.WriteFile(backupPath, jsonData, 0644); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to create backup: %v", err))
	} else {
		result.BackupCreated = true
		result.BackupPath = backupPath
	}

	if err := saveConfig(rootDir, upgradedConfig); err != nil {
		return nil, fmt.Errorf("failed to save upgraded config: %w", err)
	}

	result.Success = true
	result.Message = "Configuration successfully upgraded to latest schema"

	return result, nil
}

func convertLegacyCommands(commandNames []string, commandType string) []CommandInfo {
	var commands []CommandInfo

	for _, name := range commandNames {
		cmd := CommandInfo{
			Name:        name,
			Type:        commandType,
			Scope:       "guild",
			Description: fmt.Sprintf("Legacy %s command", commandType),
			Args:        nil,
			ReturnType:  "None",
		}
		commands = append(commands, cmd)
	}

	return commands
}

func isModernCogFormat(cog CogConfig) bool {
	if len(cog.SlashCommands) > 0 {
		return cog.SlashCommands[0].Description != ""
	}
	if len(cog.PrefixCommands) > 0 {
		return cog.PrefixCommands[0].Description != ""
	}
	return true
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
