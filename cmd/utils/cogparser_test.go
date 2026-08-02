/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Resolved before any test changes the working directory
var fixtureCogsDir = mustFixtureCogsDir()

func mustFixtureCogsDir() string {
	dir, err := filepath.Abs(filepath.Join("testdata", "cogs"))
	if err != nil {
		panic(err)
	}
	return dir
}

// TestParseCogFile pins the parsed shape of every fixture cog
func TestParseCogFile(t *testing.T) {
	tests := []struct {
		name string
		file string
		want ParsedCogInfo
	}{
		{
			name: "guild slash command and prefix command",
			file: "validCog",
			want: ParsedCogInfo{
				FileName:    "validCog",
				CogName:     "ValidCog",
				Author:      "Austin Choi",
				ProjectName: "TestBot",
				Description: "A discord bot used by the parser tests",
				SlashCommands: []CommandInfo{
					{
						Name:        "greet",
						Scope:       "guild",
						Type:        "slash",
						Description: "Greets a member",
						ReturnType:  "str",
						Args: []ArgInfo{
							{Name: "member", Type: "discord.Member", Description: "The member to greet"},
							{Name: "times", Type: "int", Description: "How many times to greet"},
						},
					},
				},
				PrefixCommands: []CommandInfo{
					{
						Name:        "ping",
						Scope:       "global",
						Type:        "prefix",
						Description: `Replies with pong when the user types "/ping"`,
						ReturnType:  "None",
						Args: []ArgInfo{
							{Name: "count", Type: "int"},
						},
					},
				},
			},
		},
		{
			name: "global slash command with single line describe",
			file: "globalCog",
			want: ParsedCogInfo{
				FileName:    "globalCog",
				CogName:     "GlobalCog",
				Author:      "Austin Choi",
				ProjectName: "TestBot",
				Description: "A discord bot used by the parser tests",
				SlashCommands: []CommandInfo{
					{
						Name:        "echo",
						Scope:       "global",
						Type:        "slash",
						Description: "Echoes a message",
						ReturnType:  "str",
						Args: []ArgInfo{
							{Name: "message", Type: "str", Description: "The message to echo"},
						},
					},
				},
			},
		},
		{
			name: "header with nothing but an author",
			file: "partialHeader",
			want: ParsedCogInfo{
				FileName: "partialHeader",
				CogName:  "PartialHeader",
				Author:   "Austin Choi",
			},
		},
		{
			name: "header with an empty project name",
			file: "emptyName",
			want: ParsedCogInfo{
				FileName:    "emptyName",
				CogName:     "EmptyName",
				Author:      "Austin Choi",
				Description: "A discord bot used by the parser tests",
			},
		},
		{
			name: "modal command with two fields",
			file: "modalCog",
			want: ParsedCogInfo{
				FileName:    "modalCog",
				CogName:     "ModalCog",
				Author:      "Austin Choi",
				ProjectName: "TestBot",
				Description: "A discord bot used by the parser tests",
				SlashCommands: []CommandInfo{
					{
						Name:        "feedback",
						Scope:       "guild",
						Type:        "modal",
						Description: "Collects feedback from a member",
						ReturnType:  "None",
						Fields: []FieldInfo{
							{Name: "summary", Label: "Feedback summary", Style: "short", Required: true, Placeholder: "Short summary"},
							{Name: "details", Label: "Feedback details", Style: "paragraph", Required: false},
						},
					},
				},
			},
		},
		{
			name: "decorators the command regex cannot match",
			file: "oddShape",
			want: ParsedCogInfo{
				FileName:    "oddShape",
				CogName:     "OddShape",
				Author:      "Austin Choi",
				ProjectName: "TestBot",
				Description: "A discord bot used by the parser tests",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(fixtureCogsDir, tt.file+".py")
			got, err := parseCogFile(path, tt.file)
			if err != nil {
				t.Fatalf("parseCogFile(%s) returned error: %v", path, err)
			}

			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("parseCogFile(%s)\ngot:  %+v\nwant: %+v", path, *got, tt.want)
			}
		})
	}
}

// TestParseAllCogFiles checks that every fixture in the directory is picked up
func TestParseAllCogFiles(t *testing.T) {
	parsed, err := parseAllCogFiles(fixtureCogsDir)
	if err != nil {
		t.Fatalf("parseAllCogFiles returned error: %v", err)
	}

	want := map[string]bool{
		"validCog":      false,
		"globalCog":     false,
		"partialHeader": false,
		"emptyName":     false,
		"oddShape":      false,
		"modalCog":      false,
	}

	if len(parsed) != len(want) {
		t.Fatalf("parsed %d cogs, want %d", len(parsed), len(want))
	}

	for _, cog := range parsed {
		if _, ok := want[cog.FileName]; !ok {
			t.Errorf("unexpected cog parsed: %s", cog.FileName)
			continue
		}
		want[cog.FileName] = true
	}

	for name, seen := range want {
		if !seen {
			t.Errorf("cog %s was not parsed", name)
		}
	}
}

// TestParseHeaderComment covers the header shapes the sync command has to read
func TestParseHeaderComment(t *testing.T) {
	tests := []struct {
		name            string
		lines           []string
		wantAuthor      string
		wantProjectName string
		wantDescription string
	}{
		{
			name: "generated header",
			lines: []string{
				`"""`,
				"Bot Author: Austin Choi",
				"",
				"TestBot",
				"A test description",
				`"""`,
			},
			wantAuthor:      "Austin Choi",
			wantProjectName: "TestBot",
			wantDescription: "A test description",
		},
		{
			name: "empty project name does not swallow the description",
			lines: []string{
				`"""`,
				"Bot Author: Austin Choi",
				"",
				"",
				"A test description",
				`"""`,
			},
			wantAuthor:      "Austin Choi",
			wantDescription: "A test description",
		},
		{
			name: "header without a blank separator line",
			lines: []string{
				`"""`,
				"Bot Author: Austin Choi",
				"TestBot",
				"A test description",
				`"""`,
			},
			wantAuthor:      "Austin Choi",
			wantProjectName: "TestBot",
			wantDescription: "A test description",
		},
		{
			name: "author only header still yields the author",
			lines: []string{
				`"""`,
				"Bot Author: Austin Choi",
				`"""`,
			},
			wantAuthor: "Austin Choi",
		},
		{
			name: "trailing notes do not overwrite the description",
			lines: []string{
				`"""`,
				"Bot Author: Austin Choi",
				"",
				"TestBot",
				"A test description",
				"",
				`This is an example file. Delete using the command "botbox remove"`,
				`"""`,
			},
			wantAuthor:      "Austin Choi",
			wantProjectName: "TestBot",
			wantDescription: "A test description",
		},
		{
			name:  "no docstring at all",
			lines: []string{"import discord"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := &ParsedCogInfo{}
			parseHeaderComment(tt.lines, parsed)

			if parsed.Author != tt.wantAuthor {
				t.Errorf("author = %q, want %q", parsed.Author, tt.wantAuthor)
			}
			if parsed.ProjectName != tt.wantProjectName {
				t.Errorf("project name = %q, want %q", parsed.ProjectName, tt.wantProjectName)
			}
			if parsed.Description != tt.wantDescription {
				t.Errorf("description = %q, want %q", parsed.Description, tt.wantDescription)
			}
		})
	}
}

// TestParseSlashCommandRejectsUnusableDecorators makes sure no nameless command reaches the config
func TestParseSlashCommandRejectsUnusableDecorators(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
	}{
		{
			name: "decorator split across lines",
			lines: []string{
				"    @app_commands.command(",
				`        name="split",`,
				`        description="Split decorator",`,
				"    )",
				"    async def split(self, interaction: discord.Interaction) -> None:",
			},
		},
		{
			name: "decorator without a name",
			lines: []string{
				`    @app_commands.command(description="No name")`,
				"    async def nameless(self, interaction: discord.Interaction) -> None:",
			},
		},
		{
			name: "decorator without a following function",
			lines: []string{
				`    @app_commands.command(name="orphan", description="No function")`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if cmd := parseSlashCommand(tt.lines, 0); cmd != nil {
				t.Errorf("parseSlashCommand returned %+v, want nil", *cmd)
			}
		})
	}
}

// TestParsePrefixCommandRejectsUnusableDecorators covers the same guard for prefix commands
func TestParsePrefixCommandRejectsUnusableDecorators(t *testing.T) {
	lines := []string{
		"    @commands.command()",
		"    def not_async(self, ctx: commands.Context) -> None:",
		"        return None",
	}

	if cmd := parsePrefixCommand(lines, 0); cmd != nil {
		t.Errorf("parsePrefixCommand returned %+v, want nil", *cmd)
	}
}

// TestModalCommandTemplateParseRoundTrip renders a modal command and checks the parser reads back the same command
func TestModalCommandTemplateParseRoundTrip(t *testing.T) {
	modal := CommandInfo{
		Name:        "feedback",
		Scope:       "guild",
		Type:        "modal",
		Description: "Collects feedback from a member",
		ReturnType:  "None",
		Fields: []FieldInfo{
			{Name: "summary", Label: "Feedback summary", Style: "short", Required: true, Placeholder: "Short summary"},
			{Name: "details", Label: "Feedback details", Style: "paragraph", Required: false},
		},
	}

	content, err := RenderTemplate("cog.py.tmpl", CogTemplateData{
		Author:         "Austin Choi",
		BotName:        "TestBot",
		BotDescription: "A discord bot used by the parser tests",
		ClassName:      "ModalCog",
		Filename:       "modalCog",
		SlashCommands:  []CommandInfo{modal},
	})
	if err != nil {
		t.Fatalf("RenderTemplate returned error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "modalCog.py")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write rendered cog: %v", err)
	}

	parsed, err := parseCogFile(path, "modalCog")
	if err != nil {
		t.Fatalf("parseCogFile returned error: %v", err)
	}

	if !commandsEqual(parsed.SlashCommands, []CommandInfo{modal}) {
		t.Errorf("round trip changed the command\ngot:  %+v\nwant: %+v", parsed.SlashCommands, modal)
	}
}

func TestCommandEqual(t *testing.T) {
	base := CommandInfo{
		Name:        "greet",
		Scope:       "guild",
		Type:        "slash",
		Description: "Greets a member",
		ReturnType:  "str",
		Args:        []ArgInfo{{Name: "member", Type: "discord.Member", Description: "The member to greet"}},
	}

	withScope := base
	withScope.Scope = "global"

	withDescription := base
	withDescription.Description = "Greets everyone"

	withReturn := base
	withReturn.ReturnType = "None"

	withArgType := base
	withArgType.Args = []ArgInfo{{Name: "member", Type: "str", Description: "The member to greet"}}

	withArgDescription := base
	withArgDescription.Args = []ArgInfo{{Name: "member", Type: "discord.Member", Description: "Someone"}}

	withExtraArg := base
	withExtraArg.Args = append(append([]ArgInfo{}, base.Args...), ArgInfo{Name: "times", Type: "int"})

	modal := CommandInfo{
		Name:        "feedback",
		Scope:       "guild",
		Type:        "modal",
		Description: "Collects feedback",
		ReturnType:  "None",
		Fields:      []FieldInfo{{Name: "summary", Label: "Summary", Style: "short", Required: true}},
	}

	withFieldLabel := modal
	withFieldLabel.Fields = []FieldInfo{{Name: "summary", Label: "Other", Style: "short", Required: true}}

	withFieldRequired := modal
	withFieldRequired.Fields = []FieldInfo{{Name: "summary", Label: "Summary", Style: "short", Required: false}}

	withExtraField := modal
	withExtraField.Fields = append(append([]FieldInfo{}, modal.Fields...), FieldInfo{Name: "details", Label: "Details", Style: "paragraph"})

	tests := []struct {
		name string
		a    CommandInfo
		b    CommandInfo
		want bool
	}{
		{"identical", base, base, true},
		{"different scope", base, withScope, false},
		{"different description", base, withDescription, false},
		{"different return type", base, withReturn, false},
		{"different arg type", base, withArgType, false},
		{"different arg description", base, withArgDescription, false},
		{"extra arg", base, withExtraArg, false},
		{"identical modal", modal, modal, true},
		{"different field label", modal, withFieldLabel, false},
		{"different field required", modal, withFieldRequired, false},
		{"extra field", modal, withExtraField, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commandEqual(tt.a, tt.b); got != tt.want {
				t.Errorf("commandEqual = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandsEqual(t *testing.T) {
	greet := CommandInfo{Name: "greet", Scope: "guild", Type: "slash", Description: "Greets", ReturnType: "None"}
	echo := CommandInfo{Name: "echo", Scope: "global", Type: "slash", Description: "Echoes", ReturnType: "str"}
	greetGlobal := greet
	greetGlobal.Scope = "global"

	tests := []struct {
		name string
		a    []CommandInfo
		b    []CommandInfo
		want bool
	}{
		{"both empty", nil, nil, true},
		{"same order", []CommandInfo{greet, echo}, []CommandInfo{greet, echo}, true},
		{"different order", []CommandInfo{greet, echo}, []CommandInfo{echo, greet}, true},
		{"different length", []CommandInfo{greet}, []CommandInfo{greet, echo}, false},
		{"one command changed", []CommandInfo{greet, echo}, []CommandInfo{greetGlobal, echo}, false},
		{"duplicate name hides a different command", []CommandInfo{greet, greet}, []CommandInfo{greet, echo}, false},
		{"duplicate name against a changed duplicate", []CommandInfo{greet, greet}, []CommandInfo{greet, greetGlobal}, false},
		{"matching duplicates", []CommandInfo{greet, greet}, []CommandInfo{greet, greet}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commandsEqual(tt.a, tt.b); got != tt.want {
				t.Errorf("commandsEqual = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSaveConfigRoundTrip checks the config survives a write and leaves no temp file behind
func TestSaveConfigRoundTrip(t *testing.T) {
	rootDir := t.TempDir()

	config := Config{
		BotBox:  BotBoxConfig{Version: "2.5.4"},
		BotInfo: BotConfig{Name: "TestBot", CommandPrefix: "!", Author: "Austin Choi", Description: "A test bot"},
		Cogs: []CogConfig{
			{
				Name:          "ValidCog",
				Env:           "development",
				File:          "validCog",
				SlashCommands: []CommandInfo{{Name: "greet", Scope: "guild", Type: "slash", Description: "Greets", ReturnType: "None"}},
			},
		},
	}

	if err := saveConfig(rootDir, config); err != nil {
		t.Fatalf("saveConfig returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(rootDir, "botbox.conf.tmp")); !os.IsNotExist(err) {
		t.Errorf("temp config file was left behind: %v", err)
	}

	var reloaded Config
	readJSONFile(t, filepath.Join(rootDir, "botbox.conf"), &reloaded)

	if !reflect.DeepEqual(reloaded, config) {
		t.Errorf("reloaded config\ngot:  %+v\nwant: %+v", reloaded, config)
	}
}

// TestSaveConfigOverwritesExisting checks a second write replaces the previous contents
func TestSaveConfigOverwritesExisting(t *testing.T) {
	rootDir := t.TempDir()

	first := Config{BotInfo: BotConfig{Name: "First"}}
	if err := saveConfig(rootDir, first); err != nil {
		t.Fatalf("first saveConfig returned error: %v", err)
	}

	second := Config{BotInfo: BotConfig{Name: "Second"}}
	if err := saveConfig(rootDir, second); err != nil {
		t.Fatalf("second saveConfig returned error: %v", err)
	}

	var reloaded Config
	readJSONFile(t, filepath.Join(rootDir, "botbox.conf"), &reloaded)

	if reloaded.BotInfo.Name != "Second" {
		t.Errorf("bot name = %q, want %q", reloaded.BotInfo.Name, "Second")
	}
}

// TestSyncCogsWithConfigKeepsCogsWhenNoneParsed guards against an empty cogs directory wiping the config
func TestSyncCogsWithConfigKeepsCogsWhenNoneParsed(t *testing.T) {
	rootDir := newTestProject(t, Config{
		BotBox:  BotBoxConfig{Version: "2.5.4"},
		BotInfo: BotConfig{Name: "TestBot", Author: "Austin Choi", Description: "A discord bot used by the parser tests"},
		Cogs: []CogConfig{
			{Name: "ValidCog", Env: "development", File: "validCog"},
		},
	})

	result, err := SyncCogsWithConfig()
	if err != nil {
		t.Fatalf("SyncCogsWithConfig returned error: %v", err)
	}

	if len(result.Errors) == 0 {
		t.Errorf("expected a sync error explaining that nothing was parsed")
	}

	if len(result.RemovedCogs) != 0 {
		t.Errorf("removed cogs = %v, want none", result.RemovedCogs)
	}

	var reloaded Config
	readJSONFile(t, filepath.Join(rootDir, "botbox.conf"), &reloaded)

	if len(reloaded.Cogs) != 1 || reloaded.Cogs[0].File != "validCog" {
		t.Errorf("cogs = %+v, want the original single entry", reloaded.Cogs)
	}
}

// TestSyncCogsWithConfigAddsParsedCogs covers the normal sync path end to end
func TestSyncCogsWithConfigAddsParsedCogs(t *testing.T) {
	rootDir := newTestProject(t, Config{
		BotBox:  BotBoxConfig{Version: "2.5.4"},
		BotInfo: BotConfig{Name: "TestBot", Author: "Austin Choi", Description: "A discord bot used by the parser tests"},
	})

	copyFixtureCog(t, "validCog", rootDir)

	result, err := SyncCogsWithConfig()
	if err != nil {
		t.Fatalf("SyncCogsWithConfig returned error: %v", err)
	}

	if len(result.Errors) != 0 {
		t.Errorf("sync errors = %v, want none", result.Errors)
	}

	if len(result.HeaderIssues) != 0 {
		t.Errorf("header issues = %v, want none", result.HeaderIssues)
	}

	if !reflect.DeepEqual(result.AddedCogs, []string{"validCog"}) {
		t.Errorf("added cogs = %v, want [validCog]", result.AddedCogs)
	}

	var reloaded Config
	readJSONFile(t, filepath.Join(rootDir, "botbox.conf"), &reloaded)

	if len(reloaded.Cogs) != 1 {
		t.Fatalf("cogs = %+v, want a single entry", reloaded.Cogs)
	}

	if reloaded.Cogs[0].Name != "ValidCog" || reloaded.Cogs[0].File != "validCog" {
		t.Errorf("cog entry = %+v, want ValidCog/validCog", reloaded.Cogs[0])
	}

	if len(reloaded.Cogs[0].SlashCommands) != 1 || len(reloaded.Cogs[0].PrefixCommands) != 1 {
		t.Errorf("cog commands = %+v, want one slash and one prefix command", reloaded.Cogs[0])
	}
}

// newTestProject writes a botbox project into a temp dir and makes it the working directory
func newTestProject(t *testing.T, config Config) string {
	t.Helper()

	rootDir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(rootDir, "src", "cogs"), 0755); err != nil {
		t.Fatalf("failed to create cogs directory: %v", err)
	}

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(filepath.Join(rootDir, "botbox.conf"), jsonData, 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	t.Chdir(rootDir)

	return rootDir
}

// readJSONFile decodes a JSON file written by the code under test
func readJSONFile(t *testing.T, path string, target any) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}

	if err := json.Unmarshal(content, target); err != nil {
		t.Fatalf("failed to parse %s: %v", path, err)
	}
}

// copyFixtureCog places one of the parser fixtures into a test project
func copyFixtureCog(t *testing.T, name, rootDir string) {
	t.Helper()

	source := filepath.Join(fixtureCogsDir, name+".py")

	content, err := os.ReadFile(source)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", name, err)
	}

	if err := os.WriteFile(filepath.Join(rootDir, "src", "cogs", name+".py"), content, 0644); err != nil {
		t.Fatalf("failed to write fixture %s: %v", name, err)
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
