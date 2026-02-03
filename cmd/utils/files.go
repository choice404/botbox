/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
)

func FindBotConf() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	originalDir := currentDir

	for {
		confDir := filepath.Join(currentDir)

		_, err := os.Stat(filepath.Join(confDir, "botbox.conf"))
		if err == nil {
			confPath, err := filepath.Abs(filepath.Join(confDir, "botbox.conf"))
			if err != nil {
				return "", fmt.Errorf("failed to get absolute path of %s: %w", confPath, err)
			}

			return confDir, nil
		}

		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("error checking file %s: %w", confDir, err)
		}

		parentDir := filepath.Dir(currentDir)

		if parentDir == currentDir {
			break
		}

		currentDir = parentDir
	}

	return "", fmt.Errorf("Not a botbox project: %s", originalDir)
}

func CreateFileOption(filename string) (bool, error) {
	var override bool
	formTitle := fmt.Sprintf("The file %s already exists. Do you want to override it?", filename)
	overrideForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(formTitle).
				Affirmative("yes").
				Negative("no").
				Value(&override),
		),
	)
	if _, err := os.Stat(filename); err == nil {
		overrideForm.Run()
		if override {
			return true, nil
		} else {
			return false, nil
		}
	}
	return true, nil
}

func LoadConfig() (Config, error) {
	var cfg Config

	confDir, err := FindBotConf()
	if err != nil {
		return cfg, fmt.Errorf("failed to find config directory: %w", err)
	}

	confPath := filepath.Join(confDir, "botbox.conf")

	jsonData, err := os.ReadFile(confPath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file %s: %w", confPath, err)
	}

	err = json.Unmarshal(jsonData, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to parse config JSON from %s: %w", confPath, err)
	}

	return cfg, nil
}

func validateFileName(fileName string) error {
	if fileExists(fileName) {
		return fmt.Errorf("file with name '%s' already exists", fileName)
	}
	if fileName == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	if strings.Contains(fileName, " ") {
		return fmt.Errorf("filename cannot contain spaces")
	}
	if strings.Contains(fileName, ".") || strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		return fmt.Errorf("filename cannot contain '.' or '/' or '\\'")
	}
	if strings.Contains(fileName, "-") || strings.Contains(fileName, ":") || strings.Contains(fileName, "*") || strings.Contains(fileName, "?") || strings.Contains(fileName, "\"") {
		return fmt.Errorf("filename cannot contain '-', ':', '*', '?', or '\"'")
	}
	return nil
}

func fileExists(fileName string) bool {
	rootDir, err := FindBotConf()
	filePath := filepath.Join(rootDir, "src", "cogs", fileName+".py")
	_, err = os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func SetLocalConfigValue(key string, value any) error {
	rootDir, err := FindBotConf()
	if err != nil {
		return fmt.Errorf("not in a botbox project: %w", err)
	}

	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch key {
	case "bot.name":
		if str, ok := value.(string); ok {
			config.BotInfo.Name = str
		} else {
			return fmt.Errorf("bot.name must be a string")
		}
	case "bot.description":
		if str, ok := value.(string); ok {
			config.BotInfo.Description = str
		} else {
			return fmt.Errorf("bot.description must be a string")
		}
	case "bot.command_prefix":
		if str, ok := value.(string); ok {
			config.BotInfo.CommandPrefix = str
		} else {
			return fmt.Errorf("bot.command_prefix must be a string")
		}
	case "bot.author":
		if str, ok := value.(string); ok {
			config.BotInfo.Author = str
		} else {
			return fmt.Errorf("bot.author must be a string")
		}
	default:
		return fmt.Errorf("invalid local config key: %s", key)
	}

	return saveConfig(rootDir, config)
}

func GetLocalConfigValue(key string) (any, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	switch key {
	case "bot.name":
		return config.BotInfo.Name, nil
	case "bot.description":
		return config.BotInfo.Description, nil
	case "bot.command_prefix":
		return config.BotInfo.CommandPrefix, nil
	case "bot.author":
		return config.BotInfo.Author, nil
	default:
		return nil, fmt.Errorf("invalid local config key: %s", key)
	}
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
