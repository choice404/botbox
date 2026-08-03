/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/choice404/botbox/v2/cmd/utils"
)

// writeTestGlobalConfig plants a minimal global config under the test home directory
func writeTestGlobalConfig(t *testing.T, home string) {
	t.Helper()
	dir := filepath.Join(home, ".config", "botbox")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create global config dir: %v", err)
	}
	config := `{"cli":{"version":"2.11.0","check_updates":false,"auto_update":false},"user":{"default_user":"Tester","github_username":""},"display":{"scroll_enabled":true,"color_scheme":"default"},"defaults":{"command_prefix":"!","python_version":"3.12","auto_git_init":false},"dev":{"editor":"vim"}}`
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(config), 0644); err != nil {
		t.Fatalf("failed to write global config: %v", err)
	}
}

// editTestHello builds a slash command fixture with an argument
func editTestHello() utils.CommandInfo {
	return utils.CommandInfo{
		Name:        "hello",
		Scope:       "guild",
		Type:        "slash",
		Description: "Says hello",
		Args:        []utils.ArgInfo{{Name: "target", Type: "str", Description: "Who to greet"}},
		ReturnType:  "str",
	}
}

// editTestWizard builds a single page modal command fixture
func editTestWizard() utils.CommandInfo {
	return utils.CommandInfo{
		Name:        "wizard",
		Scope:       "guild",
		Type:        "modal",
		Description: "Opens a wizard",
		Fields:      []utils.FieldInfo{{Name: "answer", Label: "Your answer", Style: "short", Required: true}},
		ReturnType:  "None",
	}
}

// editTestWave builds a prefix command fixture
func editTestWave() utils.CommandInfo {
	return utils.CommandInfo{
		Name:        "wave",
		Scope:       "guild",
		Type:        "prefix",
		Description: "Waves back",
		ReturnType:  "None",
	}
}

// setupEditProject builds a temp project holding one cog with three commands
func setupEditProject(t *testing.T) string {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	writeTestGlobalConfig(t, home)

	project := t.TempDir()
	conf := `{"botbox":{"version":"2.11.0"},"bot":{"name":"TestBot","command_prefix":"!","author":"Tester","description":"A test bot","help_style":"compact","env_provider":"env"},"cogs":[]}`
	if err := os.WriteFile(filepath.Join(project, "botbox.conf"), []byte(conf), 0644); err != nil {
		t.Fatalf("failed to write botbox.conf: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(project, "src", "cogs"), 0755); err != nil {
		t.Fatalf("failed to create cogs dir: %v", err)
	}
	t.Chdir(project)

	// The starting cog goes through the same path the add command uses
	addCogName = "greetings"
	t.Cleanup(func() { addCogName = "" })
	model := utils.AddModel(addCallback, addInitCallback)
	slashJSON, err := utils.CmdInfoSliceToJSON([]utils.CommandInfo{editTestHello(), editTestWizard()})
	if err != nil {
		t.Fatalf("failed to marshal slash commands: %v", err)
	}
	prefixJSON, err := utils.CmdInfoSliceToJSON([]utils.CommandInfo{editTestWave()})
	if err != nil {
		t.Fatalf("failed to marshal prefix commands: %v", err)
	}
	model.ModelValues.Map["slashCommands"] = &slashJSON
	model.ModelValues.Map["prefixCommands"] = &prefixJSON
	if errs := utils.RunHeadless(model); len(errs) > 0 {
		t.Fatalf("failed to add starting cog: %v", errs)
	}

	return project
}

// runEditForTest runs the headless edit pipeline without touching os.Exit
func runEditForTest(t *testing.T, cogName string, opts editOptions) error {
	t.Helper()
	editCogName = cogName
	t.Cleanup(func() {
		editCogName = ""
		editWrittenPath = ""
		editBackupPath = ""
	})

	model, err := buildEditHeadlessModel(opts)
	if err != nil {
		return err
	}
	if errs := utils.RunHeadless(*model); len(errs) > 0 {
		t.Fatalf("edit callbacks failed: %v", errs)
	}
	return nil
}

// readCogFile reads the generated cog file back for content assertions
func readCogFile(t *testing.T, project string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(project, "src", "cogs", "greetings.py"))
	if err != nil {
		t.Fatalf("failed to read cog file: %v", err)
	}
	return string(data)
}

// loadEditedCog loads the config and returns the greetings cog
func loadEditedCog(t *testing.T) utils.CogConfig {
	t.Helper()
	config, err := utils.LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	for _, cog := range config.Cogs {
		if cog.File == "greetings" {
			return cog
		}
	}
	t.Fatal("greetings cog missing from config")
	return utils.CogConfig{}
}

func TestEditHeadlessRemoveCommand(t *testing.T) {
	project := setupEditProject(t)

	if err := runEditForTest(t, "greetings", editOptions{removeCommands: []string{"hello"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := readCogFile(t, project)
	if strings.Contains(content, "hello") {
		t.Error("regenerated cog file should not mention the removed command")
	}
	cog := loadEditedCog(t)
	if len(cog.SlashCommands) != 1 || cog.SlashCommands[0].Name != "wizard" {
		t.Errorf("slash commands = %+v, want only wizard", cog.SlashCommands)
	}
	if len(cog.PrefixCommands) != 1 || cog.PrefixCommands[0].Name != "wave" {
		t.Errorf("prefix commands = %+v, want only wave", cog.PrefixCommands)
	}
	if _, err := os.Stat(filepath.Join(project, "src", "cogs", "greetings.py.bak")); err != nil {
		t.Errorf("backup file missing: %v", err)
	}
}

func TestEditHeadlessAddCommands(t *testing.T) {
	project := setupEditProject(t)

	addJSON := `[{"Name":"bye","Scope":"guild","Type":"slash","Description":"Says goodbye","Args":[],"ReturnType":"str"}]`
	if err := runEditForTest(t, "greetings", editOptions{addCommands: addJSON}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := readCogFile(t, project)
	if !strings.Contains(content, "bye") {
		t.Error("regenerated cog file should contain the added command")
	}
	cog := loadEditedCog(t)
	if len(cog.SlashCommands) != 3 {
		t.Errorf("slash command count = %d, want 3", len(cog.SlashCommands))
	}
	if len(cog.PrefixCommands) != 1 {
		t.Errorf("prefix command count = %d, want 1", len(cog.PrefixCommands))
	}
}

func TestEditHeadlessReplaceCommands(t *testing.T) {
	project := setupEditProject(t)

	replaceJSON := `[{"Name":"solo","Scope":"guild","Type":"slash","Description":"The only command","Args":[],"ReturnType":"None"}]`
	if err := runEditForTest(t, "greetings", editOptions{replaceCommands: replaceJSON, replaceSet: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := readCogFile(t, project)
	if !strings.Contains(content, "solo") {
		t.Error("regenerated cog file should contain the replacement command")
	}
	for _, gone := range []string{"hello", "wizard", "wave"} {
		if strings.Contains(content, gone) {
			t.Errorf("regenerated cog file should not mention replaced command %q", gone)
		}
	}
	cog := loadEditedCog(t)
	if len(cog.SlashCommands) != 1 || cog.SlashCommands[0].Name != "solo" {
		t.Errorf("slash commands = %+v, want only solo", cog.SlashCommands)
	}
	if len(cog.PrefixCommands) != 0 {
		t.Errorf("prefix commands = %+v, want none", cog.PrefixCommands)
	}
}

func TestEditHeadlessSetEnv(t *testing.T) {
	setupEditProject(t)

	if err := runEditForTest(t, "greetings", editOptions{env: "production"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cog := loadEditedCog(t)
	if cog.Env != "production" {
		t.Errorf("cog env = %q, want production", cog.Env)
	}
}

func TestEditHeadlessNoBackup(t *testing.T) {
	project := setupEditProject(t)

	if err := runEditForTest(t, "greetings", editOptions{removeCommands: []string{"wave"}, noBackup: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(project, "src", "cogs", "greetings.py.bak")); !os.IsNotExist(err) {
		t.Errorf("backup file should not exist, stat error = %v", err)
	}
}

func TestEditHeadlessErrors(t *testing.T) {
	setupEditProject(t)

	tests := []struct {
		name string
		cog  string
		opts editOptions
	}{
		{"unknown remove name fails", "greetings", editOptions{removeCommands: []string{"nope"}}},
		{"unknown cog fails", "ghost", editOptions{removeCommands: []string{"hello"}}},
		{"duplicate added name fails", "greetings", editOptions{addCommands: `[{"Name":"wave","Scope":"guild","Type":"slash","Description":"d","Args":[],"ReturnType":"None"}]`}},
		{"bad env fails", "greetings", editOptions{env: "staging"}},
		{"bad json fails", "greetings", editOptions{addCommands: "not json"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editCogName = tt.cog
			t.Cleanup(func() { editCogName = "" })
			if _, err := buildEditHeadlessModel(tt.opts); err == nil {
				t.Error("expected an error, got nil")
			}
		})
	}
}

func TestEditHeadlessMissingCogName(t *testing.T) {
	setupEditProject(t)
	editCogName = ""
	if _, err := buildEditHeadlessModel(editOptions{removeCommands: []string{"hello"}}); err == nil {
		t.Error("expected an error when no cog name is given")
	}
}

func TestEditThenSyncReportsNoChanges(t *testing.T) {
	setupEditProject(t)

	addJSON := `[{"Name":"bye","Scope":"guild","Type":"slash","Description":"Says goodbye","Args":[],"ReturnType":"str"}]`
	if err := runEditForTest(t, "greetings", editOptions{removeCommands: []string{"hello"}, addCommands: addJSON}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The regenerated file and config must round trip cleanly through sync, twice
	for round := 1; round <= 2; round++ {
		result, err := utils.SyncCogsWithConfig()
		if err != nil {
			t.Fatalf("sync round %d failed: %v", round, err)
		}
		if len(result.Errors) > 0 {
			t.Fatalf("sync round %d reported errors: %v", round, result.Errors)
		}
		if len(result.AddedCogs) > 0 || len(result.UpdatedCogs) > 0 || len(result.RemovedCogs) > 0 {
			t.Errorf("sync round %d found changes: added=%v updated=%v removed=%v",
				round, result.AddedCogs, result.UpdatedCogs, result.RemovedCogs)
		}
	}
}

/*
Copyright © 2025 Austin "Choice404" Choi

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
