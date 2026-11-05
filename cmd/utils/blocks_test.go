/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// testGraph builds a valid slash trigger graph with a branch and two leaf actions
func testGraph() BlockGraph {
	return BlockGraph{
		Name: "welcome",
		Blocks: []BlockInfo{
			{ID: "t1", Type: BlockTriggerSlashCommand, Config: map[string]string{
				"name":        "welcome",
				"description": "Welcomes a member",
				"scope":       "guild",
				"args":        `[{"Name":"mood","Type":"str","Description":"The mood of the welcome"}]`,
			}},
			{ID: "b1", Type: BlockLogicBranch, Config: map[string]string{
				"source": "mood", "operator": "equals", "value": "happy",
			}},
			{ID: "b2", Type: BlockActionSendMessage, Config: map[string]string{
				"content": "Welcome to {guild}, {user}!", "ephemeral": "false",
			}},
			{ID: "b3", Type: BlockActionReact, Config: map[string]string{
				"emoji": "👋",
			}},
		},
		Edges: []EdgeInfo{
			{From: "t1", Port: "next", To: "b1"},
			{From: "b1", Port: "match", To: "b2"},
			{From: "b1", Port: "else", To: "b3"},
		},
	}
}

// mutateGraph copies the base graph and hands the copy to the mutation
func mutateGraph(mutate func(*BlockGraph)) BlockGraph {
	graph := testGraph()

	// Blocks and configs are copied so the mutation cannot touch other cases
	blocks := make([]BlockInfo, len(graph.Blocks))
	for i, block := range graph.Blocks {
		config := make(map[string]string, len(block.Config))
		for key, value := range block.Config {
			config[key] = value
		}
		blocks[i] = BlockInfo{ID: block.ID, Type: block.Type, Config: config}
	}
	graph.Blocks = blocks
	graph.Edges = append([]EdgeInfo{}, graph.Edges...)

	mutate(&graph)
	return graph
}

// TestValidateBlockGraph pins every validation rule with a table of broken graphs
func TestValidateBlockGraph(t *testing.T) {
	tests := []struct {
		name    string
		graph   BlockGraph
		cog     CogConfig
		wantErr string
	}{
		{
			name:  "valid branch graph passes",
			graph: testGraph(),
		},
		{
			name: "valid linear graph passes",
			graph: BlockGraph{
				Name: "ping",
				Blocks: []BlockInfo{
					{ID: "t1", Type: BlockTriggerSlashCommand, Config: map[string]string{
						"name": "ping", "description": "Pings", "scope": "global",
					}},
					{ID: "b1", Type: BlockActionSendMessage, Config: map[string]string{"content": "pong"}},
					{ID: "b2", Type: BlockActionSendDM, Config: map[string]string{"content": "pong again"}},
				},
				Edges: []EdgeInfo{
					{From: "t1", Port: "next", To: "b1"},
					{From: "b1", Port: "next", To: "b2"},
				},
			},
		},
		{
			name:    "empty name fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Name = "" }),
			wantErr: "name cannot be empty",
		},
		{
			name:    "name with spaces fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Name = "wel come" }),
			wantErr: "must start with a letter",
		},
		{
			name:    "duplicate graph name on cog fails",
			graph:   testGraph(),
			cog:     CogConfig{Name: "Greeter", BlockGraphs: []BlockGraph{{Name: "welcome"}}},
			wantErr: "already exists",
		},
		{
			name:    "no blocks fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Blocks = nil; g.Edges = nil }),
			wantErr: "at least one block",
		},
		{
			name: "too many blocks fails",
			graph: mutateGraph(func(g *BlockGraph) {
				for i := 0; i < MaxBlocksPerGraph; i++ {
					g.Blocks = append(g.Blocks, BlockInfo{ID: "x" + string(rune('a'+i)), Type: BlockActionReact, Config: map[string]string{"emoji": "x"}})
				}
			}),
			wantErr: "at most 50 blocks",
		},
		{
			name:    "duplicate block id fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Blocks[2].ID = "b1" }),
			wantErr: "used more than once",
		},
		{
			name: "invalid block id fails",
			graph: mutateGraph(func(g *BlockGraph) {
				g.Blocks[1].ID = "1-bad"
				g.Edges[0].To = "1-bad"
				g.Edges[1].From = "1-bad"
				g.Edges[2].From = "1-bad"
			}),
			wantErr: "python identifier",
		},
		{
			name:    "unknown type fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Blocks[2].Type = "action.explode" }),
			wantErr: "unknown type",
		},
		{
			name:    "no trigger fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Blocks = g.Blocks[1:]; g.Edges = g.Edges[1:] }),
			wantErr: "exactly one trigger",
		},
		{
			name: "two triggers fail",
			graph: mutateGraph(func(g *BlockGraph) {
				g.Blocks = append(g.Blocks, BlockInfo{ID: "t2", Type: BlockTriggerSlashCommand, Config: map[string]string{
					"name": "other", "description": "d", "scope": "global",
				}})
			}),
			wantErr: "exactly one trigger",
		},
		{
			name:    "edge from unknown id fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Edges[1].From = "ghost" }),
			wantErr: "edge from unknown block id",
		},
		{
			name:    "edge to unknown id fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Edges[2].To = "ghost" }),
			wantErr: "edge to unknown block id",
		},
		{
			name:    "edge into trigger fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Edges[2].To = "t1" }),
			wantErr: "cannot point at the trigger",
		},
		{
			name:    "wrong port on action fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Edges = append(g.Edges, EdgeInfo{From: "b2", Port: "match", To: "b3"}) }),
			wantErr: "has no port 'match'",
		},
		{
			name:    "wrong port on branch fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Edges[1].Port = "next" }),
			wantErr: "has no port 'next'",
		},
		{
			name:    "duplicate port edge fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Edges = append(g.Edges, EdgeInfo{From: "b1", Port: "match", To: "b3"}) }),
			wantErr: "more than one edge on port",
		},
		{
			name:    "unreachable block fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Edges = g.Edges[:2] }),
			wantErr: "not reachable",
		},
		{
			name:    "cycle fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Edges = append(g.Edges, EdgeInfo{From: "b2", Port: "next", To: "b1"}) }),
			wantErr: "cycle detected",
		},
		{
			name:    "branch missing source fails",
			graph:   mutateGraph(func(g *BlockGraph) { delete(g.Blocks[1].Config, "source") }),
			wantErr: "needs a source",
		},
		{
			name:    "branch bad operator fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Blocks[1].Config["operator"] = "regex" }),
			wantErr: "operator must be one of",
		},
		{
			name:    "branch missing value fails",
			graph:   mutateGraph(func(g *BlockGraph) { delete(g.Blocks[1].Config, "value") }),
			wantErr: "needs a value",
		},
		{
			name:    "send_message missing content fails",
			graph:   mutateGraph(func(g *BlockGraph) { delete(g.Blocks[2].Config, "content") }),
			wantErr: "needs a content",
		},
		{
			name:    "send_message bad ephemeral fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Blocks[2].Config["ephemeral"] = "yes" }),
			wantErr: "ephemeral must be",
		},
		{
			name:    "content with double quote fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Blocks[2].Config["content"] = `say "hi"` }),
			wantErr: "cannot contain double quotes",
		},
		{
			name:    "react missing emoji fails",
			graph:   mutateGraph(func(g *BlockGraph) { delete(g.Blocks[3].Config, "emoji") }),
			wantErr: "needs an emoji",
		},
		{
			name:    "trigger missing name fails",
			graph:   mutateGraph(func(g *BlockGraph) { delete(g.Blocks[0].Config, "name") }),
			wantErr: "name cannot be empty",
		},
		{
			name:    "trigger missing description fails",
			graph:   mutateGraph(func(g *BlockGraph) { delete(g.Blocks[0].Config, "description") }),
			wantErr: "needs a description",
		},
		{
			name:    "trigger bad scope fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Blocks[0].Config["scope"] = "everywhere" }),
			wantErr: "guild or global",
		},
		{
			name:    "trigger bad args json fails",
			graph:   mutateGraph(func(g *BlockGraph) { g.Blocks[0].Config["args"] = "not json" }),
			wantErr: "failed to parse trigger args",
		},
		{
			name: "trigger arg with bad type fails",
			graph: mutateGraph(func(g *BlockGraph) {
				g.Blocks[0].Config["args"] = `[{"Name":"mood","Type":"complex","Description":"d"}]`
			}),
			wantErr: "argument type must be",
		},
		{
			name:  "trigger name colliding with slash command fails",
			graph: testGraph(),
			cog: CogConfig{Name: "Greeter", SlashCommands: []CommandInfo{
				{Name: "welcome", Type: "slash", Scope: "guild", Description: "d", ReturnType: "None"},
			}},
			wantErr: "collides with an existing command",
		},
		{
			name:  "trigger name colliding with another graph trigger fails",
			graph: testGraph(),
			cog: CogConfig{Name: "Greeter", BlockGraphs: []BlockGraph{{
				Name: "other",
				Blocks: []BlockInfo{{ID: "t1", Type: BlockTriggerSlashCommand, Config: map[string]string{
					"name": "welcome", "description": "d", "scope": "guild",
				}}},
			}}},
			wantErr: "collides with the trigger of graph",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBlockGraph(tt.graph, Config{Cogs: []CogConfig{tt.cog}}, 0, "")
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected an error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want it to contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

// TestBlockGraphsEqual pins the counterpart matched graph comparison
func TestBlockGraphsEqual(t *testing.T) {
	base := testGraph()

	if !blockGraphsEqual([]BlockGraph{base}, []BlockGraph{testGraph()}) {
		t.Error("identical graphs should compare equal")
	}
	if blockGraphsEqual([]BlockGraph{base}, nil) {
		t.Error("different lengths should compare unequal")
	}
	if blockGraphsEqual([]BlockGraph{base}, []BlockGraph{mutateGraph(func(g *BlockGraph) { g.Blocks[2].Config["content"] = "changed" })}) {
		t.Error("changed config value should compare unequal")
	}
	if blockGraphsEqual([]BlockGraph{base}, []BlockGraph{mutateGraph(func(g *BlockGraph) { g.Edges[1].To = "b3"; g.Edges[2].To = "b2" })}) {
		t.Error("changed edges should compare unequal")
	}
	if !blockGraphsEqual(nil, []BlockGraph{}) {
		t.Error("nil and empty should compare equal")
	}
}

// renderTestCog renders a cog holding a plain command plus the test graph and returns its lines
func renderTestCog(t *testing.T, graphs []BlockGraph) string {
	t.Helper()
	content, err := RenderTemplate("cog.py.tmpl", CogTemplateData{
		Author:         "Austin Choi",
		BotName:        "TestBot",
		BotDescription: "A discord bot used by the parser tests",
		ClassName:      "Greeter",
		Filename:       "greeter",
		SlashCommands: []CommandInfo{
			{Name: "hello", Scope: "guild", Type: "slash", Description: "Says hello", ReturnType: "str",
				Args: []ArgInfo{{Name: "target", Type: "str", Description: "Who to greet"}}},
		},
		BlockGraphs: graphs,
	})
	if err != nil {
		t.Fatalf("failed to render cog template: %v", err)
	}
	return content
}

// TestBlockGraphRoundTrip renders a graph into a cog file and parses it back unchanged
func TestBlockGraphRoundTrip(t *testing.T) {
	graphs := []BlockGraph{testGraph()}
	content := renderTestCog(t, graphs)

	path := filepath.Join(t.TempDir(), "greeter.py")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write rendered cog: %v", err)
	}

	parsed, err := parseCogFile(path, "greeter")
	if err != nil {
		t.Fatalf("failed to parse rendered cog: %v", err)
	}

	// The blob is the single source of truth, the parsed graphs must deep equal the input
	if !reflect.DeepEqual(parsed.BlockGraphs, graphs) {
		t.Errorf("round trip drifted:\n got %+v\nwant %+v", parsed.BlockGraphs, graphs)
	}
}

// TestBlockGraphDedup checks that a graph owned slash command never lands in slash_commands
func TestBlockGraphDedup(t *testing.T) {
	content := renderTestCog(t, []BlockGraph{testGraph()})

	path := filepath.Join(t.TempDir(), "greeter.py")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write rendered cog: %v", err)
	}

	parsed, err := parseCogFile(path, "greeter")
	if err != nil {
		t.Fatalf("failed to parse rendered cog: %v", err)
	}

	// The plain command stays, the graph trigger folds into the graph
	if len(parsed.SlashCommands) != 1 || parsed.SlashCommands[0].Name != "hello" {
		t.Errorf("slash commands = %+v, want only hello", parsed.SlashCommands)
	}
	if len(parsed.BlockGraphs) != 1 || parsed.BlockGraphs[0].Name != "welcome" {
		t.Errorf("block graphs = %+v, want only welcome", parsed.BlockGraphs)
	}
}

// TestRenderWithoutGraphsUnchanged pins that an empty graph list adds nothing to the output
func TestRenderWithoutGraphsUnchanged(t *testing.T) {
	withNil := renderTestCog(t, nil)

	if strings.Contains(withNil, "BLOCKS") || strings.Contains(withNil, "BlockRun") {
		t.Error("a cog without graphs should not mention the block runtime")
	}
}

// TestRegenerateCogFileWithGraph regenerates a cog holding a graph and parses the file back
func TestRegenerateCogFileWithGraph(t *testing.T) {
	rootDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(rootDir, "src", "cogs"), 0755); err != nil {
		t.Fatalf("failed to create cogs dir: %v", err)
	}

	cog := CogConfig{
		Name:           "Greeter",
		File:           "greeter",
		Env:            "development",
		SlashCommands:  []CommandInfo{},
		PrefixCommands: []CommandInfo{},
		BlockGraphs:    []BlockGraph{testGraph()},
	}
	config := Config{
		BotInfo: BotConfig{Name: "TestBot", Author: "Austin Choi", Description: "A discord bot used by the parser tests"},
		Cogs:    []CogConfig{cog},
	}

	if err := RegenerateCogFile(rootDir, config, cog, false); err != nil {
		t.Fatalf("failed to regenerate cog file: %v", err)
	}
	assertPythonParses(t, filepath.Join(rootDir, "src", "cogs", "greeter.py"))

	parsed, err := parseCogFile(filepath.Join(rootDir, "src", "cogs", "greeter.py"), "greeter")
	if err != nil {
		t.Fatalf("failed to parse regenerated cog: %v", err)
	}
	if !reflect.DeepEqual(parsed.BlockGraphs, cog.BlockGraphs) {
		t.Errorf("regenerated graphs drifted:\n got %+v\nwant %+v", parsed.BlockGraphs, cog.BlockGraphs)
	}
	if len(parsed.SlashCommands) != 0 {
		t.Errorf("slash commands = %+v, want none, the trigger belongs to the graph", parsed.SlashCommands)
	}
}

// TestEnsureBlocksHelper writes the helper once and leaves an existing file alone
func TestEnsureBlocksHelper(t *testing.T) {
	rootDir := t.TempDir()

	if err := EnsureBlocksHelper(rootDir); err != nil {
		t.Fatalf("failed to write blocks helper: %v", err)
	}
	path := filepath.Join(rootDir, "src", "utils", "blocks.py")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("blocks helper missing: %v", err)
	}
	if !strings.Contains(string(data), "class BlockRun") {
		t.Error("blocks helper should define BlockRun")
	}

	// A second call must not rewrite a file the user may have edited
	if err := os.WriteFile(path, []byte("edited"), 0644); err != nil {
		t.Fatalf("failed to overwrite helper: %v", err)
	}
	if err := EnsureBlocksHelper(rootDir); err != nil {
		t.Fatalf("second ensure failed: %v", err)
	}
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("blocks helper missing after second ensure: %v", err)
	}
	if string(data) != "edited" {
		t.Error("second ensure should leave the existing file alone")
	}
}

// TestSyncKeepsGraphsWhenBlobCorrupt breaks a blob by hand and checks the config survives a sync
func TestSyncKeepsGraphsWhenBlobCorrupt(t *testing.T) {
	cog := CogConfig{
		Name: "Greeter", File: "greeter", Env: "development",
		SlashCommands: []CommandInfo{}, PrefixCommands: []CommandInfo{},
		BlockGraphs: []BlockGraph{testGraph()},
	}

	config := Config{
		BotInfo: BotConfig{Name: "TestBot", Author: "Austin Choi", Description: "A discord bot used by the parser tests"},
		Cogs:    []CogConfig{cog},
	}

	rootDir := newTestProject(t, config)
	if err := RegenerateCogFile(rootDir, config, cog, false); err != nil {
		t.Fatalf("failed to regenerate cog file: %v", err)
	}

	// Drop a quote from the blob's first line so the JSON no longer parses
	cogPath := filepath.Join(rootDir, "src", "cogs", "greeter.py")
	content, err := os.ReadFile(cogPath)
	if err != nil {
		t.Fatalf("failed to read cog: %v", err)
	}
	broken := strings.Replace(string(content), `"Name": "welcome",`, `"Name: "welcome"`, 1)
	if err := os.WriteFile(cogPath, []byte(broken), 0644); err != nil {
		t.Fatalf("failed to write broken cog: %v", err)
	}

	parsed, err := parseCogFile(cogPath, "greeter")
	if err != nil {
		t.Fatalf("parse failed outright: %v", err)
	}
	if len(parsed.BlockGraphErrors) != 1 {
		t.Fatalf("parse errors = %d, want 1", len(parsed.BlockGraphErrors))
	}

	// Sync has to warn about the blob and leave the config graphs alone
	result, err := SyncCogsWithConfig()
	if err != nil {
		t.Fatalf("sync failed: %v", err)
	}
	if len(result.Errors) != 1 || !strings.Contains(result.Errors[0], "not valid JSON") {
		t.Errorf("sync errors = %v, want one blob warning", result.Errors)
	}

	var saved Config
	readJSONFile(t, filepath.Join(rootDir, "botbox.conf"), &saved)
	if len(saved.Cogs[0].BlockGraphs) != 1 {
		t.Errorf("config graphs = %d after corrupt sync, want 1", len(saved.Cogs[0].BlockGraphs))
	}
}

// TestValidateBlockGraphCrossCogTrigger rejects a trigger whose name another cog already registers
func TestValidateBlockGraphCrossCogTrigger(t *testing.T) {
	config := Config{Cogs: []CogConfig{
		{Name: "Util", File: "util", SlashCommands: []CommandInfo{{Name: "welcome", Type: "slash", Scope: "global", Description: "Says hi", ReturnType: "None"}}},
		{Name: "Greeter", File: "greeter"},
	}}

	err := ValidateBlockGraph(testGraph(), config, 1, "")
	if err == nil {
		t.Fatal("expected a cross cog collision error")
	}
	if !strings.Contains(err.Error(), "cog 'util'") {
		t.Errorf("error %q does not name the owning cog", err)
	}

	// The same graph on the util cog itself is caught by the same cog check with its own message
	err = ValidateBlockGraph(testGraph(), config, 0, "")
	if err == nil || !strings.Contains(err.Error(), "existing command on cog") {
		t.Errorf("same cog collision error = %v", err)
	}
}

// TestValidateBlockGraphReplace lets a graph keep its own name and trigger when replacing itself
func TestValidateBlockGraphReplace(t *testing.T) {
	config := Config{Cogs: []CogConfig{{Name: "Greeter", File: "greeter", BlockGraphs: []BlockGraph{testGraph()}}}}

	if err := ValidateBlockGraph(testGraph(), config, 0, ""); err == nil {
		t.Error("expected a duplicate graph error without a replace name")
	}
	if err := ValidateBlockGraph(testGraph(), config, 0, "welcome"); err != nil {
		t.Errorf("replace of the same graph failed: %v", err)
	}
}

// TestValidateBraces pins which placeholder shapes the content check accepts
func TestValidateBraces(t *testing.T) {
	cases := []struct {
		value   string
		wantErr bool
	}{
		{"plain text", false},
		{"hi {user}", false},
		{"{user} and {guild}", false},
		{"bad {", true},
		{"bad }", true},
		{"nested {{user}}", true},
	}
	for _, c := range cases {
		err := validateBraces(c.value, "content")
		if (err != nil) != c.wantErr {
			t.Errorf("validateBraces(%q) err = %v, wantErr %v", c.value, err, c.wantErr)
		}
	}

	// The graph validator runs the same check on message content
	graph := mutateGraph(func(g *BlockGraph) { g.Blocks[2].Config["content"] = "oops {" })
	if err := ValidateBlockGraph(graph, Config{Cogs: []CogConfig{{Name: "Greeter", File: "greeter"}}}, 0, ""); err == nil {
		t.Error("expected the graph validator to reject unbalanced braces")
	}
}

// TestValidateDeferKeys accepts boolean strings only for defer and defer_ephemeral
func TestValidateDeferKeys(t *testing.T) {
	config := Config{Cogs: []CogConfig{{Name: "Greeter", File: "greeter"}}}

	good := mutateGraph(func(g *BlockGraph) {
		g.Blocks[0].Config["defer"] = "true"
		g.Blocks[0].Config["defer_ephemeral"] = "false"
	})
	if err := ValidateBlockGraph(good, config, 0, ""); err != nil {
		t.Errorf("valid defer keys rejected: %v", err)
	}

	bad := mutateGraph(func(g *BlockGraph) { g.Blocks[0].Config["defer"] = "yes" })
	if err := ValidateBlockGraph(bad, config, 0, ""); err == nil {
		t.Error("expected an error for defer=yes")
	}
}

// TestNormalizeBlockGraph fills the optional keys and leaves the input alone
func TestNormalizeBlockGraph(t *testing.T) {
	graph := mutateGraph(func(g *BlockGraph) {
		delete(g.Blocks[0].Config, "scope")
		delete(g.Blocks[2].Config, "ephemeral")
		g.Edges[0].Port = ""
	})

	normalized := NormalizeBlockGraph(graph)

	if got := normalized.Blocks[0].Config["scope"]; got != "global" {
		t.Errorf("scope = %q, want global", got)
	}
	if got := normalized.Blocks[2].Config["ephemeral"]; got != "false" {
		t.Errorf("ephemeral = %q, want false", got)
	}
	if got := normalized.Edges[0].Port; got != DefaultEdgePort {
		t.Errorf("port = %q, want %q", got, DefaultEdgePort)
	}
	if _, present := graph.Blocks[0].Config["scope"]; present {
		t.Error("normalize mutated the input graph")
	}
	if graph.Edges[0].Port != "" {
		t.Error("normalize mutated the input edges")
	}
	if err := ValidateBlockGraph(normalized, Config{Cogs: []CogConfig{{Name: "Greeter", File: "greeter"}}}, 0, ""); err != nil {
		t.Errorf("normalized graph failed validation: %v", err)
	}
}

// TestBlockTypesCoverPalette keeps the types listing in step with the validator's palette
func TestBlockTypesCoverPalette(t *testing.T) {
	types := BlockTypes()
	if len(types) != len(validBlockTypes) {
		t.Fatalf("BlockTypes has %d entries, validator knows %d", len(types), len(validBlockTypes))
	}
	for i, blockType := range types {
		if blockType.Type != validBlockTypes[i] {
			t.Errorf("type %d = %s, want %s", i, blockType.Type, validBlockTypes[i])
		}
		if len(blockType.Ports) == 0 {
			t.Errorf("type %s lists no ports", blockType.Type)
		}
	}
}

// fullPaletteGraph uses every block type, both branch ports, a DM, and a deferred trigger
func fullPaletteGraph() BlockGraph {
	return BlockGraph{
		Name: "full-tour",
		Blocks: []BlockInfo{
			{ID: "t1", Type: BlockTriggerSlashCommand, Config: map[string]string{
				"name": "tour", "description": "Runs every block", "scope": "global",
				"defer": "true", "defer_ephemeral": "true",
			}},
			{ID: "b1", Type: BlockLogicBranch, Config: map[string]string{
				"source": "user_id", "operator": "contains", "value": "1",
			}},
			{ID: "m1", Type: BlockActionSendMessage, Config: map[string]string{
				"content": "Hi {user_name} from {guild_name}", "ephemeral": "true",
			}},
			{ID: "d1", Type: BlockActionSendDM, Config: map[string]string{
				"content": "Psst {user}",
			}},
			{ID: "r1", Type: BlockActionReact, Config: map[string]string{
				"emoji": "✅",
			}},
		},
		Edges: []EdgeInfo{
			{From: "t1", Port: "next", To: "b1"},
			{From: "b1", Port: "match", To: "m1"},
			{From: "b1", Port: "else", To: "d1"},
			{From: "m1", Port: "next", To: "r1"},
		},
	}
}

// TestRenderFullPalette renders every block type and checks the python parses and calls the helper
func TestRenderFullPalette(t *testing.T) {
	cog := CogConfig{
		Name: "Tour", File: "tour", Env: "development",
		SlashCommands: []CommandInfo{}, PrefixCommands: []CommandInfo{},
		BlockGraphs: []BlockGraph{fullPaletteGraph()},
	}
	config := Config{
		BotInfo: BotConfig{Name: "TestBot", Author: "Austin Choi", Description: "A discord bot used by the parser tests"},
		Cogs:    []CogConfig{cog},
	}
	if err := ValidateBlockGraph(fullPaletteGraph(), Config{Cogs: []CogConfig{{Name: "Tour", File: "tour"}}}, 0, ""); err != nil {
		t.Fatalf("full palette graph is invalid: %v", err)
	}

	rootDir := newTestProject(t, config)
	if err := RegenerateCogFile(rootDir, config, cog, false); err != nil {
		t.Fatalf("failed to regenerate cog file: %v", err)
	}
	cogPath := filepath.Join(rootDir, "src", "cogs", "tour.py")
	assertPythonParses(t, cogPath)

	data, err := os.ReadFile(cogPath)
	if err != nil {
		t.Fatalf("failed to read cog: %v", err)
	}
	content := string(data)
	for _, want := range []string{
		`await run.send_dm(run.format("Psst {user}"))`,
		"await interaction.response.defer(ephemeral=True)",
		`if "1" in run.get("user_id"):`,
	} {
		if !strings.Contains(content, want) {
			t.Errorf("rendered cog missing %q", want)
		}
	}
	if strings.Contains(content, "run.user.send(") {
		t.Error("rendered cog still reaches into run.user.send")
	}

	// Defer has to happen before the run is built or the first send races the timeout
	if strings.Index(content, "response.defer(") > strings.Index(content, "BlockRun.from_interaction(") {
		t.Error("defer call comes after the run is built")
	}

	// The blob round trips through the parser untouched
	parsed, err := parseCogFile(cogPath, "tour")
	if err != nil {
		t.Fatalf("failed to parse rendered cog: %v", err)
	}
	if !reflect.DeepEqual(parsed.BlockGraphs, cog.BlockGraphs) {
		t.Errorf("rendered graphs drifted:\n got %+v\nwant %+v", parsed.BlockGraphs, cog.BlockGraphs)
	}
}

// TestBlocksHelperParses runs the embedded helper through python's parser
func TestBlocksHelperParses(t *testing.T) {
	rootDir := t.TempDir()
	if err := EnsureBlocksHelper(rootDir); err != nil {
		t.Fatalf("failed to write helper: %v", err)
	}
	assertPythonParses(t, filepath.Join(rootDir, "src", "utils", "blocks.py"))
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
