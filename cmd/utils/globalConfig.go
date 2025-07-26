/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".config", "botbox", "config.json"), nil
}

func GlobalConfigExists() (bool, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error checking config file: %w", err)
	}

	return true, nil
}

func CreateGlobalConfig() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	configDir := filepath.Dir(configPath)

	if err := createConfigDirectories(configDir); err != nil {
		return fmt.Errorf("failed to create config directories: %w", err)
	}

	defaultConfig := createDefaultConfig()

	viper.SetConfigType("json")

	setGlobalConfigViperDefaults(defaultConfig)

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Created config file: %s\n", configPath)
	return nil
}

func LoadGlobalConfig() (*GlobalConfig, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	configDir := filepath.Dir(configPath)

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config GlobalConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func createConfigDirectories(configDir string) error {
	configParent := filepath.Dir(configDir)
	if _, err := os.Stat(configParent); os.IsNotExist(err) {
		if err := os.MkdirAll(configParent, 0755); err != nil {
			return fmt.Errorf("failed to create .config directory: %w", err)
		}
	}

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create botbox directory: %w", err)
		}
	}

	return nil
}

func createDefaultConfig() GlobalConfig {
	configVersion := strings.TrimPrefix(Version, "v")

	defaultCLI := GlobalCLI{
		Version:      configVersion,
		CheckUpdates: true,
		AutoUpdate:   false,
	}

	defaultUser := GlobalUser{
		DefaultUser:    "",
		GithubUsername: "",
	}

	defaultDisplay := GlobalDisplay{
		ScrollEnabled: true,
		ColorScheme:   "default",
	}

	defaultDefaults := GlobalDefaults{
		CommandPrefix: "!",
		PythonVersion: "3.11",
		AutoGitInit:   true,
	}

	defaultDev := GlobalDev{
		Editor: "code",
	}

	return GlobalConfig{
		CLI:      defaultCLI,
		User:     defaultUser,
		Display:  defaultDisplay,
		Defaults: defaultDefaults,
		Dev:      defaultDev,
	}
}

func setGlobalConfigViperDefaults(config GlobalConfig) {
	viper.SetDefault("cli.version", config.CLI.Version)
	viper.SetDefault("cli.check_updates", config.CLI.CheckUpdates)
	viper.SetDefault("cli.auto_update", config.CLI.AutoUpdate)

	viper.SetDefault("user.default_user", config.User.DefaultUser)
	viper.SetDefault("user.github_username", config.User.GithubUsername)

	viper.SetDefault("display.scroll_enabled", config.Display.ScrollEnabled)
	viper.SetDefault("display.color_scheme", config.Display.ColorScheme)

	viper.SetDefault("defaults.command_prefix", config.Defaults.CommandPrefix)
	viper.SetDefault("defaults.python_version", config.Defaults.PythonVersion)
	viper.SetDefault("defaults.auto_git_init", config.Defaults.AutoGitInit)

	viper.SetDefault("dev.editor", config.Dev.Editor)
}

func GetGlobalConfigValue(key string) any {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil
	}

	configDir := filepath.Dir(configPath)

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		return nil
	}

	return viper.Get(key)
}

func SetGlobalConfigValue(key string, value any) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	configDir := filepath.Dir(configPath)

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := CreateGlobalConfig(); err != nil {
				return fmt.Errorf("failed to create global config: %w", err)
			}
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("failed to read newly created config: %w", err)
			}
		} else {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	viper.Set(key, value)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func SyncCogsWithConfig() (*SyncResult, error) {
	result := &SyncResult{}

	rootDir, err := FindBotConf()
	if err != nil {
		return nil, fmt.Errorf("failed to find botbox project: %w", err)
	}

	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	cogsDir := filepath.Join(rootDir, "src", "cogs")
	parsedCogs, err := parseAllCogFiles(cogsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cog files: %w", err)
	}

	existingCogs := make(map[string]*CogConfig)
	for i := range config.Cogs {
		existingCogs[config.Cogs[i].File] = &config.Cogs[i]
	}

	var newCogs []CogConfig
	for _, parsed := range parsedCogs {
		if existing, exists := existingCogs[parsed.FileName]; exists {
			updated := updateCogConfig(existing, parsed)
			if updated {
				result.UpdatedCogs = append(result.UpdatedCogs, parsed.FileName)
			}
			newCogs = append(newCogs, *existing)
			delete(existingCogs, parsed.FileName)
		} else {
			newCog := createCogConfigFromParsed(parsed)
			newCogs = append(newCogs, newCog)
			result.AddedCogs = append(result.AddedCogs, parsed.FileName)
		}

		if headerIssue := checkHeaderIssues(parsed, &config.BotInfo); headerIssue != "" {
			result.HeaderIssues = append(result.HeaderIssues, fmt.Sprintf("%s: %s", parsed.FileName, headerIssue))
		}
	}

	for fileName := range existingCogs {
		result.RemovedCogs = append(result.RemovedCogs, fileName)
	}

	config.Cogs = newCogs

	if err := saveConfig(rootDir, config); err != nil {
		return nil, fmt.Errorf("failed to save updated config: %w", err)
	}

	return result, nil
}

func parseAllCogFiles(cogsDir string) ([]ParsedCogInfo, error) {
	var parsedCogs []ParsedCogInfo

	files, err := os.ReadDir(cogsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read cogs directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".py") && file.Name() != "__init__.py" {
			filePath := filepath.Join(cogsDir, file.Name())
			fileName := strings.TrimSuffix(file.Name(), ".py")

			parsed, err := parseCogFile(filePath, fileName)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s: %w", file.Name(), err)
			}

			parsedCogs = append(parsedCogs, *parsed)
		}
	}

	return parsedCogs, nil
}

func parseCogFile(filePath, fileName string) (*ParsedCogInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	parsed := &ParsedCogInfo{
		FileName: fileName,
	}

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	parseHeaderComment(lines, parsed)

	parseCogClassName(lines, parsed)

	parseCommands(lines, parsed)

	return parsed, nil
}

func parseHeaderComment(lines []string, parsed *ParsedCogInfo) {
	inDocstring := false
	docstringLines := []string{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, `"""`) {
			if inDocstring {
				break
			} else {
				inDocstring = true
				if strings.Count(trimmed, `"""`) == 2 {
					content := strings.Trim(trimmed, `"`)
					docstringLines = append(docstringLines, content)
					break
				}
				content := strings.TrimPrefix(trimmed, `"""`)
				if content != "" {
					docstringLines = append(docstringLines, content)
				}
			}
		} else if inDocstring {
			docstringLines = append(docstringLines, trimmed)
		}
	}

	if len(docstringLines) >= 3 {
		for i, line := range docstringLines {
			if strings.HasPrefix(line, "Bot Author:") {
				parsed.Author = strings.TrimSpace(strings.TrimPrefix(line, "Bot Author:"))
			} else if i == 1 || (i > 0 && parsed.ProjectName == "" && !strings.HasPrefix(line, "Bot Author:")) {
				if parsed.ProjectName == "" && line != "" {
					parsed.ProjectName = line
				}
			} else if i == 2 || (i > 1 && parsed.Description == "" && line != "") {
				if parsed.Description == "" && line != "" {
					parsed.Description = line
				}
			}
		}
	}
}

func parseCogClassName(lines []string, parsed *ParsedCogInfo) {
	classRegex := regexp.MustCompile(`^class\s+(\w+)\s*\(.*commands\.Cog.*\):`)

	for _, line := range lines {
		if matches := classRegex.FindStringSubmatch(strings.TrimSpace(line)); matches != nil {
			parsed.CogName = matches[1]
			break
		}
	}

	if parsed.CogName == "" {
		parsed.CogName = parsed.FileName
	}
}

func parseCommands(lines []string, parsed *ParsedCogInfo) {
	for i := range lines {
		line := strings.TrimSpace(lines[i])

		if strings.Contains(line, "@app_commands.command") {
			cmd := parseSlashCommand(lines, i)
			if cmd != nil {
				parsed.SlashCommands = append(parsed.SlashCommands, *cmd)
			}
		}

		if strings.Contains(line, "@commands.command") {
			cmd := parsePrefixCommand(lines, i)
			if cmd != nil {
				parsed.PrefixCommands = append(parsed.PrefixCommands, *cmd)
			}
		}
	}
}

func parseSlashCommand(lines []string, startIndex int) *CommandInfo {
	cmd := &CommandInfo{
		Type: "slash",
	}

	commandRegex := regexp.MustCompile(`@app_commands\.command\s*\(\s*name\s*=\s*["']([^"']+)["']\s*,\s*description\s*=\s*["']([^"']+)["']\s*\)`)
	if matches := commandRegex.FindStringSubmatch(lines[startIndex]); matches != nil {
		cmd.Name = matches[1]
		cmd.Description = matches[2]
	}

	cmd.Scope = "global"
	for j := startIndex + 1; j < len(lines) && j < startIndex+5; j++ {
		if strings.Contains(lines[j], "@app_commands.guilds(GUILD)") {
			cmd.Scope = "guild"
			break
		}
	}

	for j := startIndex + 1; j < len(lines) && j < startIndex+10; j++ {
		line := strings.TrimSpace(lines[j])
		if strings.HasPrefix(line, "async def ") {
			parseCommandFunction(line, cmd)
			break
		}
	}

	for j := startIndex + 1; j < len(lines) && j < startIndex+10; j++ {
		line := strings.TrimSpace(lines[j])
		if strings.Contains(line, "@app_commands.describe") {
			parseDescribeDecorator(lines, j, cmd)
			break
		}
	}

	return cmd
}

func parsePrefixCommand(lines []string, startIndex int) *CommandInfo {
	cmd := &CommandInfo{
		Type:  "prefix",
		Scope: "global"}

	for j := startIndex + 1; j < len(lines) && j < startIndex+5; j++ {
		line := strings.TrimSpace(lines[j])
		if strings.HasPrefix(line, "async def ") {
			parseCommandFunction(line, cmd)
			if cmd.Name == "" {
				funcRegex := regexp.MustCompile(`async def (\w+)\s*\(`)
				if matches := funcRegex.FindStringSubmatch(line); matches != nil {
					cmd.Name = matches[1]
				}
			}
			break
		}
	}

	parseCommandDocstring(lines, startIndex, cmd)

	return cmd
}

func parseCommandFunction(line string, cmd *CommandInfo) {
	returnRegex := regexp.MustCompile(`->\s*([^:]+):`)
	if matches := returnRegex.FindStringSubmatch(line); matches != nil {
		cmd.ReturnType = strings.TrimSpace(matches[1])
	} else {
		cmd.ReturnType = "None"
	}

	paramRegex := regexp.MustCompile(`\(([^)]+)\)`)
	if matches := paramRegex.FindStringSubmatch(line); matches != nil {
		for param := range strings.SplitSeq(matches[1], ",") {
			param = strings.TrimSpace(param)
			if param == "self" || strings.HasPrefix(param, "interaction:") || strings.HasPrefix(param, "ctx:") {
				continue
			}

			parts := strings.Split(param, ":")
			if len(parts) >= 2 {
				name := strings.TrimSpace(parts[0])
				typeAndDefault := strings.TrimSpace(parts[1])

				typePart := strings.Split(typeAndDefault, "=")[0]
				typePart = strings.TrimSpace(typePart)

				cmd.Args = append(cmd.Args, ArgInfo{
					Name: name,
					Type: typePart,
				})
			}
		}
	}
}

func parseDescribeDecorator(lines []string, startIndex int, cmd *CommandInfo) {
	line := lines[startIndex]

	fullDecorator := line
	for j := startIndex + 1; j < len(lines); j++ {
		nextLine := strings.TrimSpace(lines[j])
		if strings.HasSuffix(nextLine, ")") {
			fullDecorator += " " + nextLine
			break
		} else if nextLine != "" {
			fullDecorator += " " + nextLine
		}
	}

	argRegex := regexp.MustCompile(`(\w+)\s*=\s*["']([^"']+)["']`)
	matches := argRegex.FindAllStringSubmatch(fullDecorator, -1)

	for _, match := range matches {
		argName := match[1]
		description := match[2]

		for i := range cmd.Args {
			if cmd.Args[i].Name == argName {
				cmd.Args[i].Description = description
				break
			}
		}
	}
}

func parseCommandDocstring(lines []string, startIndex int, cmd *CommandInfo) {
	funcIndex := -1
	for j := startIndex + 1; j < len(lines) && j < startIndex+5; j++ {
		if strings.HasPrefix(strings.TrimSpace(lines[j]), "async def ") {
			funcIndex = j
			break
		}
	}

	if funcIndex == -1 {
		return
	}

	for j := funcIndex + 1; j < len(lines) && j < funcIndex+10; j++ {
		line := strings.TrimSpace(lines[j])
		if strings.HasPrefix(line, `"""`) {
			content := strings.TrimPrefix(line, `"""`)
			content = strings.TrimSuffix(content, `"""`)
			content = strings.TrimSpace(content)
			if content != "" {
				cmd.Description = content
			}
			break
		}
	}
}

func updateCogConfig(existing *CogConfig, parsed ParsedCogInfo) bool {
	updated := false

	if existing.Name != parsed.CogName {
		existing.Name = parsed.CogName
		updated = true
	}

	if !commandsEqual(existing.SlashCommands, parsed.SlashCommands) {
		existing.SlashCommands = parsed.SlashCommands
		updated = true
	}

	if !commandsEqual(existing.PrefixCommands, parsed.PrefixCommands) {
		existing.PrefixCommands = parsed.PrefixCommands
		updated = true
	}

	return updated
}

func createCogConfigFromParsed(parsed ParsedCogInfo) CogConfig {
	return CogConfig{
		Name: parsed.CogName,
		File: parsed.FileName,
		Env:  "development", SlashCommands: parsed.SlashCommands,
		PrefixCommands: parsed.PrefixCommands,
	}
}

func checkHeaderIssues(parsed ParsedCogInfo, botInfo *BotConfig) string {
	var issues []string

	if parsed.Author != botInfo.Author {
		issues = append(issues, fmt.Sprintf("author mismatch (file: %s, config: %s)", parsed.Author, botInfo.Author))
	}

	if parsed.ProjectName != botInfo.Name {
		issues = append(issues, fmt.Sprintf("project name mismatch (file: %s, config: %s)", parsed.ProjectName, botInfo.Name))
	}

	if parsed.Description != botInfo.Description {
		issues = append(issues, fmt.Sprintf("description mismatch (file: %s, config: %s)", parsed.Description, botInfo.Description))
	}

	if len(issues) > 0 {
		return strings.Join(issues, "; ")
	}

	return ""
}

func commandsEqual(a, b []CommandInfo) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := make(map[string]CommandInfo)
	bMap := make(map[string]CommandInfo)

	for _, cmd := range a {
		aMap[cmd.Name] = cmd
	}

	for _, cmd := range b {
		bMap[cmd.Name] = cmd
	}

	for name, cmdA := range aMap {
		cmdB, exists := bMap[name]
		if !exists {
			return false
		}

		if !commandEqual(cmdA, cmdB) {
			return false
		}
	}

	return true
}

func commandEqual(a, b CommandInfo) bool {
	if a.Name != b.Name || a.Type != b.Type || a.Scope != b.Scope ||
		a.Description != b.Description || a.ReturnType != b.ReturnType {
		return false
	}

	if len(a.Args) != len(b.Args) {
		return false
	}

	for i, argA := range a.Args {
		argB := b.Args[i]
		if argA.Name != argB.Name || argA.Type != argB.Type || argA.Description != argB.Description {
			return false
		}
	}

	return true
}

func saveConfig(rootDir string, config Config) error {
	configPath := filepath.Join(rootDir, "botbox.conf")

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func SyncGlobalConfigVersion() error {
	if Version == "" {
		return fmt.Errorf("version not set")
	}

	exists, err := GlobalConfigExists()
	if err != nil {
		return fmt.Errorf("failed to check global config: %w", err)
	}

	currentVersion := strings.TrimPrefix(Version, "v")

	if currentVersion == "" {
		return fmt.Errorf("invalid version after trimming: %s", Version)
	}

	if !exists {
		if err := CreateGlobalConfig(); err != nil {
			return fmt.Errorf("failed to create global config: %w", err)
		}
		fmt.Printf("📝 Created global config with version %s\n", currentVersion)
		return nil
	}

	configVersion := GetGlobalConfigValue("cli.version")

	if configVersion == nil || configVersion.(string) != currentVersion {
		if err := SetGlobalConfigValue("cli.version", currentVersion); err != nil {
			return fmt.Errorf("failed to sync version in global config: %w", err)
		}
		fmt.Printf("📝 Synced global config version to %s\n", currentVersion)
	}

	return nil
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
