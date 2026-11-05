/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Block type names in the current palette, everything else is rejected at validation time
const (
	BlockTriggerSlashCommand = "trigger.slash_command"
	BlockActionSendMessage   = "action.send_message"
	BlockActionSendDM        = "action.send_dm"
	BlockActionReact         = "action.react"
	BlockLogicBranch         = "logic.branch"
)

// validBlockTypes lists every block type a graph may contain
var validBlockTypes = []string{
	BlockTriggerSlashCommand,
	BlockActionSendMessage,
	BlockActionSendDM,
	BlockActionReact,
	BlockLogicBranch,
}

// MaxBlocksPerGraph caps how many blocks a single graph can hold,
// every block becomes a module level function so this keeps generated files readable
const MaxBlocksPerGraph = 50

// MaxGraphsPerCog caps how many graphs a single cog can hold,
// every graph trigger is a slash command and discord allows 100 per bot in total
const MaxGraphsPerCog = 20

// DefaultEdgePort is the port an edge uses when the JSON leaves it out
const DefaultEdgePort = "next"

// Builtin run values every graph can branch on, seeded by BlockRun in the generated helper
var builtinRunValues = []string{"user_id", "user_name", "guild_id", "guild_name", "channel_id", "channel_name"}

// Branch blocks only compare with these operators
var validBranchOperators = []string{"equals", "contains"}

// Block ids become python identifiers, so they follow python identifier rules
var blockIDRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// Graph names become python constants and function names after the usual dash conversion
var graphNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

// isTriggerBlock reports whether a block starts a run
func isTriggerBlock(block BlockInfo) bool {
	return strings.HasPrefix(block.Type, "trigger.")
}

// GraphTrigger returns the trigger block of a graph
func GraphTrigger(graph BlockGraph) (BlockInfo, bool) {
	for _, block := range graph.Blocks {
		if isTriggerBlock(block) {
			return block, true
		}
	}
	return BlockInfo{}, false
}

// TriggerArgs reads the ArgInfo list stored as a JSON array string under config["args"]
func TriggerArgs(block BlockInfo) ([]ArgInfo, error) {
	raw := block.Config["args"]
	if raw == "" {
		return nil, nil
	}
	args, err := JSONToArgInfoSlice(raw)
	if err != nil {
		return nil, fmt.Errorf("block '%s': failed to parse trigger args: %w", block.ID, err)
	}
	return args, nil
}

// validatePythonSafe rejects characters that would break the python string the value lands inside
func validatePythonSafe(value, label string) error {
	if strings.ContainsAny(value, "\"\\\n") {
		return fmt.Errorf("%s cannot contain double quotes, backslashes, or newlines", label)
	}
	return nil
}

// validateBraces rejects placeholder syntax format_map cannot parse, unbalanced or nested braces
func validateBraces(value, label string) error {
	depth := 0
	for _, r := range value {
		switch r {
		case '{':
			depth++
			if depth > 1 {
				return fmt.Errorf("%s cannot nest braces", label)
			}
		case '}':
			depth--
			if depth < 0 {
				return fmt.Errorf("%s has a closing brace with no opening brace", label)
			}
		}
	}
	if depth != 0 {
		return fmt.Errorf("%s has an opening brace with no closing brace", label)
	}
	return nil
}

// validateContent runs every check a message body needs before it lands inside run.format
func validateContent(value, label string) error {
	if err := validatePythonSafe(value, label); err != nil {
		return err
	}
	return validateBraces(value, label)
}

// validateBoolKey accepts an absent key or the strings "true" and "false"
func validateBoolKey(block BlockInfo, key string) error {
	if v := block.Config[key]; v != "" && v != "true" && v != "false" {
		return fmt.Errorf("block '%s': %s must be \"true\" or \"false\"", block.ID, key)
	}
	return nil
}

// blockPorts returns the output ports a block type may use
func blockPorts(blockType string) []string {
	if blockType == BlockLogicBranch {
		return []string{"match", "else"}
	}
	return []string{"next"}
}

// ValidateBlockGraph checks a graph against the palette rules, the cog it will attach to,
// and the slash command names already taken by the rest of the project,
// replaceName is the graph being swapped out so its own name does not count as a duplicate
func ValidateBlockGraph(graph BlockGraph, config Config, cogIndex int, replaceName string) error {
	if cogIndex < 0 || cogIndex >= len(config.Cogs) {
		return fmt.Errorf("cog index %d is out of range", cogIndex)
	}
	cog := config.Cogs[cogIndex]
	owners := SlashCommandOwners(config)

	if graph.Name == "" {
		return fmt.Errorf("graph name cannot be empty")
	}
	if !graphNameRegex.MatchString(graph.Name) {
		return fmt.Errorf("graph name '%s' must start with a letter and contain only letters, digits, dashes, and underscores", graph.Name)
	}

	// A graph name has to be unique inside its cog, unless it is the one being replaced
	for _, existing := range cog.BlockGraphs {
		if existing.Name == graph.Name && graph.Name != replaceName {
			return fmt.Errorf("graph '%s' already exists on cog '%s'", graph.Name, cog.Name)
		}
	}

	if len(graph.Blocks) == 0 {
		return fmt.Errorf("graph '%s' needs at least one block", graph.Name)
	}
	if len(graph.Blocks) > MaxBlocksPerGraph {
		return fmt.Errorf("graph '%s' can have at most %d blocks", graph.Name, MaxBlocksPerGraph)
	}

	// Collect the block ids first so edges can be checked against them
	blocksByID := make(map[string]BlockInfo)
	triggerCount := 0
	triggerID := ""
	for _, block := range graph.Blocks {
		if block.ID == "" {
			return fmt.Errorf("graph '%s': block ids cannot be empty", graph.Name)
		}
		if !blockIDRegex.MatchString(block.ID) {
			return fmt.Errorf("graph '%s': block id '%s' must be a valid python identifier", graph.Name, block.ID)
		}
		if _, exists := blocksByID[block.ID]; exists {
			return fmt.Errorf("graph '%s': block id '%s' is used more than once", graph.Name, block.ID)
		}
		blocksByID[block.ID] = block

		if !contains(validBlockTypes, block.Type) {
			return fmt.Errorf("graph '%s': block '%s' has unknown type '%s', must be one of %s",
				graph.Name, block.ID, block.Type, strings.Join(validBlockTypes, ", "))
		}
		if isTriggerBlock(block) {
			triggerCount++
			triggerID = block.ID
		}
	}
	if triggerCount != 1 {
		return fmt.Errorf("graph '%s' needs exactly one trigger block, found %d", graph.Name, triggerCount)
	}

	// Every edge leaves a real block through a port its type owns and lands on a real block
	edgeTargets := make(map[string][]string)
	seenPorts := make(map[string]bool)
	for _, edge := range graph.Edges {
		from, exists := blocksByID[edge.From]
		if !exists {
			return fmt.Errorf("graph '%s': edge from unknown block id '%s'", graph.Name, edge.From)
		}
		if _, exists := blocksByID[edge.To]; !exists {
			return fmt.Errorf("graph '%s': edge to unknown block id '%s'", graph.Name, edge.To)
		}
		if edge.To == triggerID {
			return fmt.Errorf("graph '%s': edges cannot point at the trigger block '%s'", graph.Name, triggerID)
		}
		if !contains(blockPorts(from.Type), edge.Port) {
			return fmt.Errorf("graph '%s': block '%s' of type %s has no port '%s', valid ports are %s",
				graph.Name, edge.From, from.Type, edge.Port, strings.Join(blockPorts(from.Type), ", "))
		}
		portKey := edge.From + "/" + edge.Port
		if seenPorts[portKey] {
			return fmt.Errorf("graph '%s': block '%s' has more than one edge on port '%s'", graph.Name, edge.From, edge.Port)
		}
		seenPorts[portKey] = true
		edgeTargets[edge.From] = append(edgeTargets[edge.From], edge.To)
	}

	// Every block has to be reachable from the trigger and the walk must never loop
	if err := checkGraphWalk(graph, triggerID, edgeTargets); err != nil {
		return err
	}

	// Each block type checks its own required config keys
	for _, block := range graph.Blocks {
		if err := validateBlockConfig(block, graph, cog, owners); err != nil {
			return fmt.Errorf("graph '%s': %w", graph.Name, err)
		}
	}

	return nil
}

// checkGraphWalk walks the edges from the trigger, failing on cycles and unreachable blocks
func checkGraphWalk(graph BlockGraph, triggerID string, edgeTargets map[string][]string) error {
	// States are 0 unvisited, 1 on the current path, 2 done
	state := make(map[string]int)

	var walk func(id string) error
	walk = func(id string) error {
		if state[id] == 1 {
			return fmt.Errorf("graph '%s': cycle detected through block '%s'", graph.Name, id)
		}
		if state[id] == 2 {
			return nil
		}
		state[id] = 1
		for _, next := range edgeTargets[id] {
			if err := walk(next); err != nil {
				return err
			}
		}
		state[id] = 2
		return nil
	}

	if err := walk(triggerID); err != nil {
		return err
	}

	for _, block := range graph.Blocks {
		if state[block.ID] != 2 {
			return fmt.Errorf("graph '%s': block '%s' is not reachable from the trigger", graph.Name, block.ID)
		}
	}
	return nil
}

// validateBlockConfig checks the required config keys of one block against its type
func validateBlockConfig(block BlockInfo, graph BlockGraph, cog CogConfig, owners map[string]string) error {
	switch block.Type {
	case BlockTriggerSlashCommand:
		return validateSlashTriggerConfig(block, graph, cog, owners)
	case BlockActionSendMessage:
		if block.Config["content"] == "" {
			return fmt.Errorf("block '%s': send_message needs a content config key", block.ID)
		}
		if err := validateContent(block.Config["content"], "content"); err != nil {
			return fmt.Errorf("block '%s': %w", block.ID, err)
		}
		// Ephemeral is optional but has to be a boolean string when present
		if err := validateBoolKey(block, "ephemeral"); err != nil {
			return err
		}
	case BlockActionSendDM:
		if block.Config["content"] == "" {
			return fmt.Errorf("block '%s': send_dm needs a content config key", block.ID)
		}
		if err := validateContent(block.Config["content"], "content"); err != nil {
			return fmt.Errorf("block '%s': %w", block.ID, err)
		}
	case BlockActionReact:
		if block.Config["emoji"] == "" {
			return fmt.Errorf("block '%s': react needs an emoji config key", block.ID)
		}
		if err := validatePythonSafe(block.Config["emoji"], "emoji"); err != nil {
			return fmt.Errorf("block '%s': %w", block.ID, err)
		}
	case BlockLogicBranch:
		if block.Config["source"] == "" {
			return fmt.Errorf("block '%s': branch needs a source config key", block.ID)
		}
		if !contains(validBranchOperators, block.Config["operator"]) {
			return fmt.Errorf("block '%s': branch operator must be one of %s", block.ID, strings.Join(validBranchOperators, ", "))
		}
		if block.Config["value"] == "" {
			return fmt.Errorf("block '%s': branch needs a value config key", block.ID)
		}
		if err := validatePythonSafe(block.Config["source"], "source"); err != nil {
			return fmt.Errorf("block '%s': %w", block.ID, err)
		}
		if err := validatePythonSafe(block.Config["value"], "value"); err != nil {
			return fmt.Errorf("block '%s': %w", block.ID, err)
		}
	}
	return nil
}

// validateSlashTriggerConfig checks a slash trigger against the command name rules and the cog's commands
func validateSlashTriggerConfig(block BlockInfo, graph BlockGraph, cog CogConfig, owners map[string]string) error {
	name := block.Config["name"]
	if err := ValidateCommandName(name, nil); err != nil {
		return fmt.Errorf("block '%s': trigger %w", block.ID, err)
	}

	// A graph trigger registers a real slash command, so it cannot collide with the cog's commands
	if commandExists(name, cog.SlashCommands) || commandExists(name, cog.PrefixCommands) {
		return fmt.Errorf("block '%s': trigger name '%s' collides with an existing command on cog '%s'", block.ID, name, cog.Name)
	}

	// Another cog owning the name would make discord.py raise CommandAlreadyRegistered at startup,
	// the same cog is skipped because the checks above already cover it
	if owner, taken := owners[name]; taken && owner != cog.File {
		return fmt.Errorf("block '%s': trigger name '%s' is already registered by cog '%s'", block.ID, name, owner)
	}

	// Nor with the trigger of another graph on the same cog
	for _, other := range cog.BlockGraphs {
		if other.Name == graph.Name {
			continue
		}
		if trigger, found := GraphTrigger(other); found && trigger.Config["name"] == name {
			return fmt.Errorf("block '%s': trigger name '%s' collides with the trigger of graph '%s'", block.ID, name, other.Name)
		}
	}

	if block.Config["description"] == "" {
		return fmt.Errorf("block '%s': trigger needs a description config key", block.ID)
	}
	if err := validatePythonSafe(block.Config["description"], "description"); err != nil {
		return fmt.Errorf("block '%s': %w", block.ID, err)
	}
	if err := ValidateCommandScope(block.Config["scope"]); err != nil {
		return fmt.Errorf("block '%s': trigger %w", block.ID, err)
	}

	// Defer and its ephemeral flag are optional but have to be boolean strings when present
	if err := validateBoolKey(block, "defer"); err != nil {
		return err
	}
	if err := validateBoolKey(block, "defer_ephemeral"); err != nil {
		return err
	}

	// Trigger args ride config["args"] as a JSON array string and validate like command args
	args, err := TriggerArgs(block)
	if err != nil {
		return err
	}
	for i, arg := range args {
		if err := ValidateArgName(arg.Name, args[:i]); err != nil {
			return fmt.Errorf("block '%s': argument '%s': %w", block.ID, arg.Name, err)
		}
		if err := ValidateArgDescription(arg.Description); err != nil {
			return fmt.Errorf("block '%s': argument '%s': %w", block.ID, arg.Name, err)
		}
		if err := ValidateArgType(arg.Type); err != nil {
			return fmt.Errorf("block '%s': argument '%s': %w", block.ID, arg.Name, err)
		}
	}

	return nil
}

// NormalizeBlockGraph returns a copy of the graph with the optional keys filled in with their defaults,
// the input is never mutated so callers can keep the original around
func NormalizeBlockGraph(graph BlockGraph) BlockGraph {
	normalized := BlockGraph{Name: graph.Name}

	normalized.Blocks = make([]BlockInfo, 0, len(graph.Blocks))
	for _, block := range graph.Blocks {
		// Copy the config so the defaults land on the new graph only
		config := make(map[string]string, len(block.Config)+2)
		maps.Copy(config, block.Config)

		if block.Type == BlockTriggerSlashCommand && config["scope"] == "" {
			config["scope"] = "global"
		}
		if block.Type == BlockActionSendMessage && config["ephemeral"] == "" {
			config["ephemeral"] = "false"
		}

		normalized.Blocks = append(normalized.Blocks, BlockInfo{ID: block.ID, Type: block.Type, Config: config})
	}

	normalized.Edges = make([]EdgeInfo, 0, len(graph.Edges))
	for _, edge := range graph.Edges {
		// Edges are values, so changing the port here only changes the copy
		if edge.Port == "" {
			edge.Port = DefaultEdgePort
		}
		normalized.Edges = append(normalized.Edges, edge)
	}

	return normalized
}

// BlockTypeInfo describes one block type for the types command and later for the builder forms
type BlockTypeInfo struct {
	Type     string   `json:"type"`
	Required []string `json:"required"`
	Optional []string `json:"optional"`
	Ports    []string `json:"ports"`
	Notes    string   `json:"notes"`
}

// BlockTypes returns the palette in a stable order with the config keys each type reads
func BlockTypes() []BlockTypeInfo {
	return []BlockTypeInfo{
		{
			Type:     BlockTriggerSlashCommand,
			Required: []string{"name", "description"},
			Optional: []string{"scope", "args", "defer", "defer_ephemeral"},
			Ports:    blockPorts(BlockTriggerSlashCommand),
			Notes:    "scope defaults to global, args is a JSON array of command arguments, defer acknowledges the interaction before the graph runs",
		},
		{
			Type:     BlockActionSendMessage,
			Required: []string{"content"},
			Optional: []string{"ephemeral"},
			Ports:    blockPorts(BlockActionSendMessage),
			Notes:    "content accepts {user}, {guild}, argument names, and the builtin values as placeholders",
		},
		{
			Type:     BlockActionSendDM,
			Required: []string{"content"},
			Optional: []string{},
			Ports:    blockPorts(BlockActionSendDM),
			Notes:    "content accepts the same placeholders as send_message, closed DMs are logged and skipped",
		},
		{
			Type:     BlockActionReact,
			Required: []string{"emoji"},
			Optional: []string{},
			Ports:    blockPorts(BlockActionReact),
			Notes:    "reacts to the triggering message, or to the graph's own reply for slash commands",
		},
		{
			Type:     BlockLogicBranch,
			Required: []string{"source", "operator", "value"},
			Optional: []string{},
			Ports:    blockPorts(BlockLogicBranch),
			Notes:    "source is an argument name or one of " + strings.Join(builtinRunValues, ", ") + ", operator is one of " + strings.Join(validBranchOperators, ", "),
		},
	}
}

// BuiltinRunValues returns the names every graph can read without declaring them
func BuiltinRunValues() []string {
	return append([]string{}, builtinRunValues...)
}

// blockGraphsEqual reports whether two graph lists hold the same graphs, counterpart matched like commandsEqual
func blockGraphsEqual(a, b []BlockGraph) bool {
	if len(a) != len(b) {
		return false
	}

	// Each entry needs its own counterpart so repeated names cannot hide a difference
	used := make([]bool, len(b))

	for _, graphA := range a {
		matched := false
		for i := range b {
			if used[i] || !blockGraphEqual(graphA, b[i]) {
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

// blockGraphEqual compares two graphs, order sensitive within blocks and edges
func blockGraphEqual(a, b BlockGraph) bool {
	if a.Name != b.Name {
		return false
	}

	if len(a.Blocks) != len(b.Blocks) {
		return false
	}
	for i, blockA := range a.Blocks {
		blockB := b.Blocks[i]
		if blockA.ID != blockB.ID || blockA.Type != blockB.Type {
			return false
		}
		if !maps.Equal(blockA.Config, blockB.Config) {
			return false
		}
	}

	if len(a.Edges) != len(b.Edges) {
		return false
	}
	for i, edgeA := range a.Edges {
		if edgeA != b.Edges[i] {
			return false
		}
	}

	return true
}

// graphOwnedCommandNames collects the slash command names registered by graph triggers,
// the parser folds these into the graph instead of recording them as plain commands
func graphOwnedCommandNames(graphs []BlockGraph) map[string]bool {
	owned := make(map[string]bool)
	for _, graph := range graphs {
		if trigger, found := GraphTrigger(graph); found && trigger.Type == BlockTriggerSlashCommand {
			owned[trigger.Config["name"]] = true
		}
	}
	return owned
}

// EnsureBlocksHelper writes src/utils/blocks.py for projects created before the block system existed
func EnsureBlocksHelper(rootDir string) error {
	path := filepath.Join(rootDir, "src", "utils", "blocks.py")
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check blocks helper: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create utils directory: %w", err)
	}

	if err := renderToFile(path, "blocks.py.tmpl", nil); err != nil {
		return fmt.Errorf("failed to write blocks helper: %w", err)
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
