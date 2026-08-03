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

	// Never let an empty parse result wipe cogs that are still recorded in the config
	if len(newCogs) == 0 && len(config.Cogs) > 0 {
		result.RemovedCogs = nil
		result.Errors = append(result.Errors, fmt.Sprintf("no cog files parsed from %s, keeping the %d cog entries already in botbox.conf", cogsDir, len(config.Cogs)))
		return result, nil
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
	docstringLines := collectHeaderDocstring(lines)
	if len(docstringLines) == 0 {
		return
	}

	// Everything after the author line is positional, so find that line first
	bodyStart := 0
	for i, line := range docstringLines {
		if strings.HasPrefix(line, "Bot Author:") {
			parsed.Author = strings.TrimSpace(strings.TrimPrefix(line, "Bot Author:"))
			bodyStart = i + 1
			break
		}
	}

	body := docstringLines[bodyStart:]

	// Generated headers separate the author from the project block with a single blank line
	if len(body) > 0 && body[0] == "" {
		body = body[1:]
	}

	if len(body) > 0 {
		parsed.ProjectName = body[0]
	}

	if len(body) > 1 {
		parsed.Description = body[1]
	}
}

func collectHeaderDocstring(lines []string) []string {
	inDocstring := false
	docstringLines := []string{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, `"""`) {
			if inDocstring {
				break
			}

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
		} else if inDocstring {
			docstringLines = append(docstringLines, trimmed)
		}
	}

	return docstringLines
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

// Line budgets for locating the function a command decorator belongs to
const (
	maxSlashDecoratorLines  = 30
	maxPrefixDecoratorLines = 5
)

func parseSlashCommand(lines []string, startIndex int) *CommandInfo {
	commandRegex := regexp.MustCompile(`@app_commands\.command\s*\(\s*name\s*=\s*["']([^"']+)["']\s*,\s*description\s*=\s*["']([^"']+)["']\s*\)`)
	matches := commandRegex.FindStringSubmatch(lines[startIndex])

	// An unrecognized decorator carries no command identity, so record nothing
	if matches == nil {
		return nil
	}

	cmd := &CommandInfo{
		Type:        "slash",
		Scope:       "global",
		Name:        matches[1],
		Description: matches[2],
	}

	funcIndex, describeIndex := scanDecoratorBlock(lines, startIndex, cmd)

	if funcIndex == -1 {
		return nil
	}

	parseCommandFunction(strings.TrimSpace(lines[funcIndex]), cmd)

	// Argument descriptions are applied after the arguments themselves exist
	if describeIndex != -1 {
		parseDescribeDecorator(lines, describeIndex, cmd)
	}

	parseCommandDocstring(lines, funcIndex, cmd)

	// A command whose body opens a modal is recorded as a modal command
	if modalClass, found := findSendModal(lines, funcIndex); found {
		cmd.Type = "modal"
		cmd.Args = nil
		// A FLOW blob is the single source for a multi page command, only single page modals fall back to the class
		if flow, ok := parseCommandFlow(lines, cmd.Name); ok {
			cmd.Pages = flow.Pages
			cmd.Responses = flow.Responses
		} else if modalClass != "" {
			cmd.Fields = parseModalFields(lines, modalClass)
		}
	} else {
		parseCommandResponse(lines, funcIndex, cmd, slashResponseRegex)
	}

	return cmd
}

// parseCommandFlow reads the FLOW JSON blob generated next to a multi page modal command
func parseCommandFlow(lines []string, commandName string) (*commandFlow, bool) {
	marker := CommandConstName(commandName) + "_FLOW = json.loads(r'''"

	// Find the line that opens the raw triple quoted JSON string
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == marker {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return nil, false
	}

	// Collect every line until the closing quotes so the whole blob can be unmarshaled at once
	var jsonLines []string
	end := -1
	for j := start; j < len(lines); j++ {
		if strings.TrimSpace(lines[j]) == "''')" {
			end = j
			break
		}
		jsonLines = append(jsonLines, lines[j])
	}
	if end == -1 {
		return nil, false
	}

	var flow commandFlow
	if err := json.Unmarshal([]byte(strings.Join(jsonLines, "\n")), &flow); err != nil {
		return nil, false
	}

	return &flow, true
}

// Reply call shapes the generator writes into slash and prefix command bodies
var (
	slashResponseRegex  = regexp.MustCompile(`await interaction\.response\.send_message\(f?"((?:[^"\\]|\\.)*)"\s*,\s*ephemeral\s*=\s*(True|False)\s*\)`)
	prefixResponseRegex = regexp.MustCompile(`await ctx\.send\(f?"((?:[^"\\]|\\.)*)"\s*,\s*ephemeral\s*=\s*(True|False)\s*\)`)
)

// parseCommandResponse reads the generated reply call in a command body into the expected responses
// Only the generated shape counts, the first statement after the docstring must be a try block whose first line is the reply
func parseCommandResponse(lines []string, funcIndex int, cmd *CommandInfo, replyRegex *regexp.Regexp) {
	inDocstring := false
	sawTry := false

	for j := funcIndex + 1; j < len(lines) && j < funcIndex+maxCommandBodyLines; j++ {
		line := strings.TrimSpace(lines[j])

		// Track the docstring so its text is never mistaken for code
		if strings.HasPrefix(line, `"""`) {
			if !inDocstring {
				inDocstring = strings.Count(line, `"""`) == 1
			} else {
				inDocstring = false
			}
			continue
		}
		if inDocstring || line == "" {
			continue
		}

		// The first real statement has to open the try block or the body is not the generated shape
		if !sawTry {
			if line == "try:" {
				sawTry = true
				continue
			}
			return
		}

		matches := replyRegex.FindStringSubmatch(line)
		if matches == nil {
			return
		}

		content := matches[1]
		ephemeral := matches[2] == "True"

		// The default generated reply echoes the command name, that exact shape means no expected responses
		if content == cmd.Name && ephemeral {
			return
		}

		cmd.Responses = []ResponseInfo{{Type: "message", Content: content, Ephemeral: ephemeral}}
		return
	}
}

// Line budget for scanning a command body for a send_modal call
const maxCommandBodyLines = 30

// findSendModal reports whether the function body calls send_modal and which modal class it opens
func findSendModal(lines []string, funcIndex int) (string, bool) {
	sendModalRegex := regexp.MustCompile(`send_modal\(\s*(\w+)\s*\(`)

	for j := funcIndex + 1; j < len(lines) && j < funcIndex+maxCommandBodyLines; j++ {
		line := strings.TrimSpace(lines[j])

		// Stop before crossing into the next decorator, function, or class
		if strings.HasPrefix(line, "@") || strings.HasPrefix(line, "async def ") ||
			strings.HasPrefix(line, "def ") || strings.HasPrefix(line, "class ") {
			break
		}

		if !strings.Contains(line, "send_modal(") {
			continue
		}

		if matches := sendModalRegex.FindStringSubmatch(line); matches != nil {
			return matches[1], true
		}

		// A send_modal call whose class cannot be read still marks the command as modal
		return "", true
	}

	return "", false
}

// parseModalFields reads the TextInput assignments out of the named modal class body
func parseModalFields(lines []string, modalClass string) []FieldInfo {
	classRegex := regexp.MustCompile(`^class\s+` + regexp.QuoteMeta(modalClass) + `\s*\(.*discord\.ui\.Modal`)
	textInputRegex := regexp.MustCompile(`^(\w+)\s*=\s*discord\.ui\.TextInput\((.*)\)\s*$`)
	labelRegex := regexp.MustCompile(`label\s*=\s*["']([^"']*)["']`)
	styleRegex := regexp.MustCompile(`style\s*=\s*discord\.TextStyle\.(\w+)`)
	requiredRegex := regexp.MustCompile(`required\s*=\s*(True|False)`)
	placeholderRegex := regexp.MustCompile(`placeholder\s*=\s*["']([^"']*)["']`)

	classIndex := -1
	for i, line := range lines {
		if classRegex.MatchString(strings.TrimSpace(line)) {
			classIndex = i
			break
		}
	}

	if classIndex == -1 {
		return nil
	}

	var fields []FieldInfo
	for j := classIndex + 1; j < len(lines); j++ {
		trimmed := strings.TrimSpace(lines[j])

		// A non-indented statement ends the class body
		if trimmed != "" && !strings.HasPrefix(lines[j], " ") && !strings.HasPrefix(lines[j], "\t") {
			break
		}

		matches := textInputRegex.FindStringSubmatch(trimmed)
		if matches == nil {
			continue
		}

		field := FieldInfo{
			Name: matches[1],
			// TextInput defaults to a required short style input when the arguments are absent
			Style:    "short",
			Required: true,
		}

		callArgs := matches[2]
		if labelMatch := labelRegex.FindStringSubmatch(callArgs); labelMatch != nil {
			field.Label = labelMatch[1]
		}
		if styleMatch := styleRegex.FindStringSubmatch(callArgs); styleMatch != nil {
			field.Style = styleMatch[1]
		}
		if requiredMatch := requiredRegex.FindStringSubmatch(callArgs); requiredMatch != nil {
			field.Required = requiredMatch[1] == "True"
		}
		if placeholderMatch := placeholderRegex.FindStringSubmatch(callArgs); placeholderMatch != nil {
			field.Placeholder = placeholderMatch[1]
		}

		fields = append(fields, field)
	}

	return fields
}

// scanDecoratorBlock walks the decorators between startIndex and the function they decorate
func scanDecoratorBlock(lines []string, startIndex int, cmd *CommandInfo) (funcIndex, describeIndex int) {
	funcIndex = -1
	describeIndex = -1

	for j := startIndex + 1; j < len(lines) && j < startIndex+maxSlashDecoratorLines; j++ {
		line := strings.TrimSpace(lines[j])

		if strings.HasPrefix(line, "async def ") {
			funcIndex = j
			break
		}

		if strings.Contains(line, "@app_commands.guilds(GUILD)") {
			cmd.Scope = "guild"
		}

		if describeIndex == -1 && strings.Contains(line, "@app_commands.describe") {
			describeIndex = j
		}
	}

	return funcIndex, describeIndex
}

func parsePrefixCommand(lines []string, startIndex int) *CommandInfo {
	cmd := &CommandInfo{
		Type:  "prefix",
		Scope: "global"}

	funcIndex := -1
	for j := startIndex + 1; j < len(lines) && j < startIndex+maxPrefixDecoratorLines; j++ {
		line := strings.TrimSpace(lines[j])
		if strings.HasPrefix(line, "async def ") {
			funcIndex = j
			break
		}
	}

	// Without a function definition there is no command name to record
	if funcIndex == -1 {
		return nil
	}

	funcLine := strings.TrimSpace(lines[funcIndex])
	parseCommandFunction(funcLine, cmd)

	funcRegex := regexp.MustCompile(`async def (\w+)\s*\(`)
	if matches := funcRegex.FindStringSubmatch(funcLine); matches != nil {
		cmd.Name = matches[1]
	}

	if cmd.Name == "" {
		return nil
	}

	parseCommandDocstring(lines, funcIndex, cmd)

	// The generator appends this phrase to the docstring, stripping it keeps descriptions round trip stable
	generatedSuffix := fmt.Sprintf(" when the user types \"/%s\"", cmd.Name)
	cmd.Description = strings.TrimSuffix(cmd.Description, generatedSuffix)

	parseDocstringArgDescriptions(lines, funcIndex, cmd)

	parseCommandResponse(lines, funcIndex, cmd, prefixResponseRegex)

	return cmd
}

// parseDocstringArgDescriptions fills empty arg descriptions from the docstring Parameters block,
// prefix commands have no describe decorator so the docstring is their only source
func parseDocstringArgDescriptions(lines []string, funcIndex int, cmd *CommandInfo) {
	if funcIndex < 0 || len(cmd.Args) == 0 {
		return
	}

	argLineRegex := regexp.MustCompile(`^(\w+)\s*\(([^)]*)\):\s*(.+)$`)

	inDocstring := false
	for j := funcIndex + 1; j < len(lines); j++ {
		line := strings.TrimSpace(lines[j])
		if strings.HasPrefix(line, `"""`) {
			if inDocstring {
				return
			}
			inDocstring = true
			continue
		}
		if !inDocstring {
			continue
		}

		matches := argLineRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		for i := range cmd.Args {
			if cmd.Args[i].Name == matches[1] && cmd.Args[i].Description == "" {
				cmd.Args[i].Description = strings.TrimSpace(matches[3])
				break
			}
		}
	}
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
	fullDecorator := strings.TrimSpace(lines[startIndex])

	// Only a decorator left open on its first line continues onto the following lines
	if !strings.HasSuffix(fullDecorator, ")") {
		for j := startIndex + 1; j < len(lines); j++ {
			nextLine := strings.TrimSpace(lines[j])
			if strings.HasSuffix(nextLine, ")") {
				fullDecorator += " " + nextLine
				break
			} else if nextLine != "" {
				fullDecorator += " " + nextLine
			}
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

// parseCommandDocstring fills in the description from the function docstring when nothing else supplied one
func parseCommandDocstring(lines []string, funcIndex int, cmd *CommandInfo) {
	if funcIndex < 0 || cmd.Description != "" {
		return
	}

	for j := funcIndex + 1; j < len(lines) && j < funcIndex+10; j++ {
		line := strings.TrimSpace(lines[j])
		if !strings.HasPrefix(line, `"""`) {
			continue
		}

		content := strings.TrimSuffix(strings.TrimPrefix(line, `"""`), `"""`)
		content = strings.TrimSpace(content)

		// Generated docstrings open on their own line and carry the summary on the next one
		if content == "" && j+1 < len(lines) {
			content = strings.TrimSpace(lines[j+1])
		}

		cmd.Description = content
		break
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

	// Each entry needs its own counterpart so repeated names cannot hide a difference
	used := make([]bool, len(b))

	for _, cmdA := range a {
		matched := false
		for i := range b {
			if used[i] || !commandEqual(cmdA, b[i]) {
				continue
			}
			used[i] = true
			matched = true
			break
		}

		if !matched {
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

	if len(a.Fields) != len(b.Fields) {
		return false
	}

	for i, fieldA := range a.Fields {
		if fieldA != b.Fields[i] {
			return false
		}
	}

	if len(a.Pages) != len(b.Pages) {
		return false
	}

	for i := range a.Pages {
		if !pageEqual(a.Pages[i], b.Pages[i]) {
			return false
		}
	}

	if len(a.Responses) != len(b.Responses) {
		return false
	}

	for i, responseA := range a.Responses {
		if responseA != b.Responses[i] {
			return false
		}
	}

	return true
}

// pageEqual compares two flow pages including their fields and branch rules
func pageEqual(a, b PageInfo) bool {
	if a.Name != b.Name || a.Title != b.Title || a.Next != b.Next {
		return false
	}

	if len(a.Fields) != len(b.Fields) {
		return false
	}

	for i, fieldA := range a.Fields {
		if fieldA != b.Fields[i] {
			return false
		}
	}

	if len(a.Branches) != len(b.Branches) {
		return false
	}

	for i, branchA := range a.Branches {
		if branchA != b.Branches[i] {
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

	// Write to a sibling temp file first so a failed write cannot destroy the existing config
	tempPath := configPath + ".tmp"
	if err := os.WriteFile(tempPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write temp config file: %w", err)
	}

	if err := os.Rename(tempPath, configPath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to replace config file: %w", err)
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
		fmt.Fprintf(os.Stderr, "📝 Created global config with version %s\n", currentVersion)
		return nil
	}

	configVersion := GetGlobalConfigValue("cli.version")

	configVersionStr, ok := configVersion.(string)
	if !ok || configVersionStr != currentVersion {
		if err := SetGlobalConfigValue("cli.version", currentVersion); err != nil {
			return fmt.Errorf("failed to sync version in global config: %w", err)
		}
		fmt.Fprintf(os.Stderr, "📝 Synced global config version to %s\n", currentVersion)
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
