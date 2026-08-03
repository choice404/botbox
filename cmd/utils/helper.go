/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
)

var (
	Version string
)

func SetVersion(version string) {
	Version = version
}

func Banner() {
	fmt.Println(`
    ____        __     ____            
   / __ )____  / /_   / __ )____  _  __
  / __  / __ \/ __/  / __  / __ \| |/_/
 / /_/ / /_/ / /_   / /_/ / /_/ />  <  
/_____/\____/\__/  /_____/\____/_/|_|  
  `)
}

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

func FetchLicense(licenseKey string) (string, error) {
	if licenseKey == "" || licenseKey == "none" {
		return "", fmt.Errorf("no license key provided or selected 'none'")
	}

	apiURL := fmt.Sprintf("https://api.github.com/licenses/%s", licenseKey)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request to %s: %w", apiURL, err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "bot-box")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch license %s: %w", licenseKey, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to fetch license %s: status %s, body: %s",
			licenseKey, resp.Status, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body for %s: %w", licenseKey, err)
	}

	var licenseResp LicenseResponse
	err = json.Unmarshal(bodyBytes, &licenseResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response for %s: %w", licenseKey, err)
	}

	if licenseResp.Body == "" {
		return "", fmt.Errorf("no license body found in response for %s", licenseKey)
	}

	return licenseResp.Body, nil
}

func CreateFileOption(filename string, force bool) (bool, error) {
	if _, err := os.Stat(filename); err != nil {
		return true, nil
	}
	if force {
		return true, nil
	}
	if HeadlessMode {
		fmt.Fprintf(os.Stderr, "Skipping %s, file already exists (use --force to overwrite)\n", filename)
		return false, nil
	}
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
	if err := overrideForm.Run(); err != nil {
		return false, err
	}
	return override, nil
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

func ValidateFileName(fileName string) error {
	if fileName == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	if fileExists(fileName) {
		return fmt.Errorf("file with name '%s' already exists", fileName)
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
	if err != nil {
		return false
	}
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

// RegenerateCogFile rewrites a cog's .py file from its config definition,
// writing a .py.bak copy of the current file first when backup is true
func RegenerateCogFile(rootDir string, config Config, cog CogConfig, backup bool) error {
	filePath := filepath.Join(rootDir, "src", "cogs", cog.File+".py")

	if backup {
		existing, err := os.ReadFile(filePath)
		if err == nil {
			if err := os.WriteFile(filePath+".bak", existing, 0644); err != nil {
				return fmt.Errorf("failed to write backup file: %w", err)
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("failed to read cog file for backup: %w", err)
		}
	}

	content, err := RenderTemplate("cog.py.tmpl", CogTemplateData{
		Author:         config.BotInfo.Author,
		BotName:        config.BotInfo.Name,
		BotDescription: config.BotInfo.Description,
		ClassName:      cog.Name,
		Filename:       cog.File,
		SlashCommands:  cog.SlashCommands,
		PrefixCommands: cog.PrefixCommands,
	})
	if err != nil {
		return fmt.Errorf("failed to render cog template: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write cog file: %w", err)
	}

	return nil
}

func commandExists(commandName string, commandList []CommandInfo) bool {
	for _, cmd := range commandList {
		if cmd.Name == commandName {
			return true
		}
	}
	return false
}

func argExists(argName string, args []ArgInfo) bool {
	for _, arg := range args {
		if arg.Name == argName {
			return true
		}
	}
	return false
}

// DetectEnvProvider infers the env provider from the files in the project root,
// doppler.yaml wins over .env and projects with neither default to env
func DetectEnvProvider(rootDir string) string {
	if _, err := os.Stat(filepath.Join(rootDir, "doppler.yaml")); err == nil {
		return "doppler"
	}
	if _, err := os.Stat(filepath.Join(rootDir, ".env")); err == nil {
		return "env"
	}
	return "env"
}

// ResolveEnvProvider returns the configured provider, falling back to file detection
// for configs written before the env_provider key existed
func ResolveEnvProvider(config Config, rootDir string) string {
	if config.BotInfo.EnvProvider != "" {
		return config.BotInfo.EnvProvider
	}
	return DetectEnvProvider(rootDir)
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
	case "bot.help_style":
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("bot.help_style must be a string")
		}
		if err := ValidateHelpStyle(str); err != nil {
			return err
		}
		config.BotInfo.HelpStyle = str
	case "bot.env_provider":
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("bot.env_provider must be a string")
		}
		// The provider shares the env or doppler choice made at project creation
		if err := ValidateEnvChoice(str); err != nil {
			return err
		}
		config.BotInfo.EnvProvider = str
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
	case "bot.help_style":
		// Projects created before this key existed report the default
		return NormalizeHelpStyle(config.BotInfo.HelpStyle), nil
	case "bot.env_provider":
		// Projects created before this key existed report the detected provider
		rootDir, err := FindBotConf()
		if err != nil {
			return nil, fmt.Errorf("not in a botbox project: %w", err)
		}
		return ResolveEnvProvider(config, rootDir), nil
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
