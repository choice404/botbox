/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/choice404/botbox/v2/cmd/utils"
)

// blockTestGraphJSON is the demo graph, a slash trigger into a branch with two leaf actions
const blockTestGraphJSON = `{
  "Name": "welcome",
  "Blocks": [
    {"ID": "t1", "Type": "trigger.slash_command",
     "Config": {"name": "welcome", "description": "Welcomes a member", "scope": "guild",
                "args": "[{\"Name\":\"mood\",\"Type\":\"str\",\"Description\":\"The mood of the welcome\"}]"}},
    {"ID": "b1", "Type": "logic.branch",
     "Config": {"source": "mood", "operator": "equals", "value": "happy"}},
    {"ID": "b2", "Type": "action.send_message",
     "Config": {"content": "Welcome to {guild}, {user}!", "ephemeral": "false"}},
    {"ID": "b3", "Type": "action.react",
     "Config": {"emoji": "wave"}}
  ],
  "Edges": [
    {"From": "t1", "Port": "next", "To": "b1"},
    {"From": "b1", "Port": "match", "To": "b2"},
    {"From": "b1", "Port": "else", "To": "b3"}
  ]
}`

// setupBlockProject builds a temp project holding one cog through the add pipeline
func setupBlockProject(t *testing.T) string {
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
	addCogName = "greeter"
	t.Cleanup(func() { addCogName = "" })
	model := utils.AddModel(addCallback, addInitCallback)
	slashJSON, err := utils.CmdInfoSliceToJSON([]utils.CommandInfo{editTestHello()})
	if err != nil {
		t.Fatalf("failed to marshal slash commands: %v", err)
	}
	prefixJSON, err := utils.CmdInfoSliceToJSON([]utils.CommandInfo{})
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

// loadBlockCog loads the config and returns the greeter cog
func loadBlockCog(t *testing.T) utils.CogConfig {
	t.Helper()
	config, err := utils.LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	for _, cog := range config.Cogs {
		if cog.File == "greeter" {
			return cog
		}
	}
	t.Fatal("greeter cog missing from config")
	return utils.CogConfig{}
}

// assertSyncNoChanges runs sync twice and fails on any reported change
func assertSyncNoChanges(t *testing.T) {
	t.Helper()
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

func TestBlockAddHeadless(t *testing.T) {
	project := setupBlockProject(t)

	if err := runBlockAdd("greeter", blockTestGraphJSON, false, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The generated file holds the blob, the woven functions, and the trigger command
	data, err := os.ReadFile(filepath.Join(project, "src", "cogs", "greeter.py"))
	if err != nil {
		t.Fatalf("failed to read cog file: %v", err)
	}
	content := string(data)
	for _, want := range []string{
		"WELCOME_BLOCKS = json.loads(r'''",
		"async def _welcome_b1(cog, run):  # logic.branch",
		"WELCOME_STEPS = {",
		`@app_commands.command(name="welcome", description="Welcomes a member")`,
		"from utils.blocks import BlockRun",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("cog file missing %q", want)
		}
	}

	// The helper lands in projects created before the block system existed
	if _, err := os.Stat(filepath.Join(project, "src", "utils", "blocks.py")); err != nil {
		t.Errorf("blocks helper missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(project, "src", "cogs", "greeter.py.bak")); err != nil {
		t.Errorf("backup file missing: %v", err)
	}

	cog := loadBlockCog(t)
	if len(cog.BlockGraphs) != 1 || cog.BlockGraphs[0].Name != "welcome" {
		t.Errorf("block graphs = %+v, want only welcome", cog.BlockGraphs)
	}
	// The trigger owns the command, it must not appear in slash_commands
	if len(cog.SlashCommands) != 1 || cog.SlashCommands[0].Name != "hello" {
		t.Errorf("slash commands = %+v, want only hello", cog.SlashCommands)
	}

	assertSyncNoChanges(t)
}

func TestBlockAddErrors(t *testing.T) {
	setupBlockProject(t)

	tests := []struct {
		name  string
		cog   string
		graph string
	}{
		{"unknown cog fails", "ghost", blockTestGraphJSON},
		{"missing graph fails", "greeter", ""},
		{"bad json fails", "greeter", "not json"},
		{"invalid graph fails", "greeter", `{"Name":"broken","Blocks":[],"Edges":[]}`},
		{"trigger colliding with command fails", "greeter", `{"Name":"clash","Blocks":[{"ID":"t1","Type":"trigger.slash_command","Config":{"name":"hello","description":"d","scope":"guild"}}],"Edges":[]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runBlockAdd(tt.cog, tt.graph, true, false); err == nil {
				t.Error("expected an error, got nil")
			}
		})
	}
}

func TestBlockAddDuplicateGraphName(t *testing.T) {
	setupBlockProject(t)

	if err := runBlockAdd("greeter", blockTestGraphJSON, true, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := runBlockAdd("greeter", blockTestGraphJSON, true, false); err == nil {
		t.Error("expected an error adding the same graph name twice")
	}
}

func TestBlockListHeadless(t *testing.T) {
	setupBlockProject(t)

	if err := runBlockAdd("greeter", blockTestGraphJSON, true, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Listing every cog and one cog both succeed, an unknown cog fails
	if err := runBlockList("", false); err != nil {
		t.Errorf("list all failed: %v", err)
	}
	if err := runBlockList("greeter", false); err != nil {
		t.Errorf("list greeter failed: %v", err)
	}
	if err := runBlockList("ghost", false); err == nil {
		t.Error("expected an error listing an unknown cog")
	}
}

func TestBlockRemoveHeadless(t *testing.T) {
	project := setupBlockProject(t)

	if err := runBlockAdd("greeter", blockTestGraphJSON, true, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	beforeRemove, err := os.ReadFile(filepath.Join(project, "src", "cogs", "greeter.py"))
	if err != nil {
		t.Fatalf("failed to read cog file: %v", err)
	}

	if err := runBlockRemove("greeter", "welcome", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(project, "src", "cogs", "greeter.py"))
	if err != nil {
		t.Fatalf("failed to read cog file: %v", err)
	}
	if strings.Contains(string(data), "WELCOME_BLOCKS") {
		t.Error("regenerated cog file should not mention the removed graph")
	}
	if string(data) == string(beforeRemove) {
		t.Error("cog file should change after the graph is removed")
	}

	cog := loadBlockCog(t)
	if len(cog.BlockGraphs) != 0 {
		t.Errorf("block graphs = %+v, want none", cog.BlockGraphs)
	}

	// A removed graph leaves no block_graphs key behind, old configs stay byte stable
	conf, err := os.ReadFile(filepath.Join(project, "botbox.conf"))
	if err != nil {
		t.Fatalf("failed to read botbox.conf: %v", err)
	}
	if strings.Contains(string(conf), "block_graphs") {
		t.Error("botbox.conf should not hold an empty block_graphs key")
	}

	assertSyncNoChanges(t)
}

func TestBlockRemoveErrors(t *testing.T) {
	setupBlockProject(t)

	if err := runBlockRemove("ghost", "welcome", true); err == nil {
		t.Error("expected an error removing from an unknown cog")
	}
	if err := runBlockRemove("greeter", "ghost", true); err == nil {
		t.Error("expected an error removing an unknown graph")
	}
}

func TestBlockGraphCountLimit(t *testing.T) {
	setupBlockProject(t)

	// Filling the cog to the limit succeeds, one more fails
	for i := 0; i < utils.MaxGraphsPerCog; i++ {
		graph := strings.ReplaceAll(blockTestGraphJSON, "welcome", "welcome"+string(rune('a'+i)))
		if err := runBlockAdd("greeter", graph, true, false); err != nil {
			t.Fatalf("add %d failed: %v", i, err)
		}
	}
	extra := strings.ReplaceAll(blockTestGraphJSON, "welcome", "welcomez")
	if err := runBlockAdd("greeter", extra, true, false); err == nil {
		t.Error("expected an error past the graph limit")
	}
}

// captureStdout runs fn and returns everything it printed on stdout
func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	saved := os.Stdout
	os.Stdout = writer
	runErr := fn()
	os.Stdout = saved
	writer.Close()

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, reader); err != nil {
		t.Fatalf("failed to read captured stdout: %v", err)
	}
	return buffer.String(), runErr
}

// TestBlockAddReplace swaps a graph in place and keeps the config at one graph
func TestBlockAddReplace(t *testing.T) {
	project := setupBlockProject(t)

	if err := runBlockAdd("greeter", blockTestGraphJSON, true, false); err != nil {
		t.Fatalf("first add failed: %v", err)
	}

	// Replace flips the message content, the graph name stays the same
	changed := strings.Replace(blockTestGraphJSON, "Welcome to {guild}, {user}!", "Hey there {user}", 1)
	if err := runBlockAdd("greeter", changed, true, true); err != nil {
		t.Fatalf("replace failed: %v", err)
	}

	cog := loadBlockCog(t)
	if len(cog.BlockGraphs) != 1 {
		t.Fatalf("block graphs = %d after replace, want 1", len(cog.BlockGraphs))
	}
	if got := cog.BlockGraphs[0].Blocks[2].Config["content"]; got != "Hey there {user}" {
		t.Errorf("content = %q, want the replaced content", got)
	}

	data, err := os.ReadFile(filepath.Join(project, "src", "cogs", "greeter.py"))
	if err != nil {
		t.Fatalf("failed to read cog file: %v", err)
	}
	if !strings.Contains(string(data), "Hey there {user}") {
		t.Error("cog file still holds the old content")
	}

	// Replace without a matching graph and a plain add of a duplicate both fail
	if err := runBlockAdd("greeter", strings.ReplaceAll(blockTestGraphJSON, "welcome", "other"), true, true); err == nil {
		t.Error("expected an error replacing a graph that does not exist")
	}
	if err := runBlockAdd("greeter", blockTestGraphJSON, true, false); err == nil {
		t.Error("expected a duplicate error without --replace")
	}

	assertSyncNoChanges(t)
}

// TestBlockAddDefaults fills scope, ephemeral, and the edge port when the JSON leaves them out
func TestBlockAddDefaults(t *testing.T) {
	setupBlockProject(t)

	bare := `{
  "Name": "ping",
  "Blocks": [
    {"ID": "t1", "Type": "trigger.slash_command", "Config": {"name": "ping", "description": "Pings"}},
    {"ID": "a1", "Type": "action.send_message", "Config": {"content": "pong"}}
  ],
  "Edges": [{"From": "t1", "To": "a1"}]
}`
	if err := runBlockAdd("greeter", bare, true, false); err != nil {
		t.Fatalf("add with bare keys failed: %v", err)
	}

	cog := loadBlockCog(t)
	graph := cog.BlockGraphs[0]
	if got := graph.Blocks[0].Config["scope"]; got != "global" {
		t.Errorf("scope = %q, want global", got)
	}
	if got := graph.Blocks[1].Config["ephemeral"]; got != "false" {
		t.Errorf("ephemeral = %q, want false", got)
	}
	if got := graph.Edges[0].Port; got != utils.DefaultEdgePort {
		t.Errorf("port = %q, want %q", got, utils.DefaultEdgePort)
	}

	assertSyncNoChanges(t)
}

// TestBlockShowHeadless prints the stored graph as JSON that block add accepts again
func TestBlockShowHeadless(t *testing.T) {
	setupBlockProject(t)

	if err := runBlockAdd("greeter", blockTestGraphJSON, true, false); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	output, err := captureStdout(t, func() error { return runBlockShow("greeter", "welcome") })
	if err != nil {
		t.Fatalf("show failed: %v", err)
	}

	var graph utils.BlockGraph
	if err := json.Unmarshal([]byte(output), &graph); err != nil {
		t.Fatalf("show output is not a graph: %v\n%s", err, output)
	}
	if graph.Name != "welcome" || len(graph.Blocks) != 4 {
		t.Errorf("shown graph = %+v, want welcome with 4 blocks", graph)
	}

	if err := runBlockShow("greeter", "ghost"); err == nil {
		t.Error("expected an error showing an unknown graph")
	}
	if err := runBlockShow("ghost", "welcome"); err == nil {
		t.Error("expected an error showing an unknown cog")
	}
}

// TestBlockValidateHeadless checks graphs without touching the project
func TestBlockValidateHeadless(t *testing.T) {
	project := setupBlockProject(t)

	before, err := os.ReadFile(filepath.Join(project, "botbox.conf"))
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	if _, err := captureStdout(t, func() error { return runBlockValidate("", blockTestGraphJSON) }); err != nil {
		t.Errorf("validate against an empty cog failed: %v", err)
	}
	if _, err := captureStdout(t, func() error { return runBlockValidate("greeter", blockTestGraphJSON) }); err != nil {
		t.Errorf("validate against greeter failed: %v", err)
	}

	// The hello command already lives on greeter, so a graph triggering /hello fails there only
	collides := strings.ReplaceAll(blockTestGraphJSON, `"name": "welcome"`, `"name": "hello"`)
	if _, err := captureStdout(t, func() error { return runBlockValidate("greeter", collides) }); err == nil {
		t.Error("expected a collision error validating against greeter")
	}
	if _, err := captureStdout(t, func() error { return runBlockValidate("", collides) }); err != nil {
		t.Errorf("validate against an empty cog should ignore project commands: %v", err)
	}

	if err := runBlockValidate("", `{"Name": "broken"}`); err == nil {
		t.Error("expected an error validating a graph with no blocks")
	}
	if err := runBlockValidate("ghost", blockTestGraphJSON); err == nil {
		t.Error("expected an error validating against an unknown cog")
	}

	after, err := os.ReadFile(filepath.Join(project, "botbox.conf"))
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if !bytes.Equal(before, after) {
		t.Error("validate changed botbox.conf")
	}
}

// TestBlockTypesHeadless prints the palette in both formats
func TestBlockTypesHeadless(t *testing.T) {
	output, err := captureStdout(t, func() error { return runBlockTypes(true) })
	if err != nil {
		t.Fatalf("types --json failed: %v", err)
	}

	var types []utils.BlockTypeInfo
	if err := json.Unmarshal([]byte(output), &types); err != nil {
		t.Fatalf("types output is not JSON: %v\n%s", err, output)
	}
	if len(types) != 5 {
		t.Errorf("types = %d, want 5", len(types))
	}

	plain, err := captureStdout(t, func() error { return runBlockTypes(false) })
	if err != nil {
		t.Fatalf("types failed: %v", err)
	}
	for _, want := range []string{"trigger.slash_command", "logic.branch", "builtin values: user_id"} {
		if !strings.Contains(plain, want) {
			t.Errorf("types output missing %q", want)
		}
	}
}

// TestBlockListJSON prints the rows as objects with stable keys
func TestBlockListJSON(t *testing.T) {
	setupBlockProject(t)

	if err := runBlockAdd("greeter", blockTestGraphJSON, true, false); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	output, err := captureStdout(t, func() error { return runBlockList("", true) })
	if err != nil {
		t.Fatalf("list --json failed: %v", err)
	}

	var rows []blockListRow
	if err := json.Unmarshal([]byte(output), &rows); err != nil {
		t.Fatalf("list output is not JSON: %v\n%s", err, output)
	}
	if len(rows) != 1 || rows[0].Cog != "greeter" || rows[0].Graph != "welcome" || rows[0].Trigger != "/welcome" {
		t.Errorf("rows = %+v, want one greeter/welcome row", rows)
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
