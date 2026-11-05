/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/choice404/botbox/v2/cmd/utils"
	"github.com/spf13/cobra"
)

var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "Manage block graphs on the cogs of the current Bot Box project",
	Long: `Manage block graphs, small directed graphs of trigger, action, and logic
blocks that generate multi step behaviors inside a cog.

A graph is a JSON document with a Name, a Blocks list, and an Edges list. It
holds exactly one trigger block and every other block hangs off it through
named ports. The graph JSON is embedded in the generated cog file and read
back by config sync, so the file and botbox.conf never drift.

Block graphs are managed headlessly through this command group for now, the
TUI builder comes later.`,
}

var blockAddCmd = &cobra.Command{
	Use:   "add <cog>",
	Short: "Add a block graph to a cog, or replace one with --replace",
	Long: `Add a block graph to an existing cog and regenerate its .py file.

The graph is given through --graph as inline JSON, @path/to/file.json, or -
for stdin. Optional keys are filled in before validation, scope defaults to
global, ephemeral to false, and an edge with no Port uses next. With
--replace the graph swaps out the existing graph of the same name in place.

Regeneration rewrites the cog file from its definition, so custom code in
command bodies is not preserved. A .py.bak copy of the previous file is
written first unless --no-backup is given.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rawGraph, _ := cmd.Flags().GetString("graph")
		noBackup, _ := cmd.Flags().GetBool("no-backup")
		replace, _ := cmd.Flags().GetBool("replace")
		if err := runBlockAdd(args[0], rawGraph, noBackup, replace); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	},
}

var blockListCmd = &cobra.Command{
	Use:   "list [cog]",
	Short: "List the block graphs of one cog or every cog",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cogName := ""
		if len(args) > 0 {
			cogName = args[0]
		}
		asJSON, _ := cmd.Flags().GetBool("json")
		if err := runBlockList(cogName, asJSON); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	},
}

var blockShowCmd = &cobra.Command{
	Use:   "show <cog> <graph-name>",
	Short: "Print one block graph as JSON, the same document block add accepts",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runBlockShow(args[0], args[1]); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	},
}

var blockValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check a block graph without writing anything",
	Long: `Normalize and validate a graph given through --graph without touching the
project. With --cog the graph is checked against that cog's commands and the
rest of the project, otherwise it is checked against an empty cog so the
command also works outside a project.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rawGraph, _ := cmd.Flags().GetString("graph")
		cogName, _ := cmd.Flags().GetString("cog")
		if err := runBlockValidate(cogName, rawGraph); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	},
}

var blockTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List every block type with the config keys and ports it uses",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		asJSON, _ := cmd.Flags().GetBool("json")
		if err := runBlockTypes(asJSON); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	},
}

var blockRemoveCmd = &cobra.Command{
	Use:   "remove <cog> <graph-name>",
	Short: "Remove a block graph from a cog",
	Long: `Remove a block graph from a cog and regenerate its .py file.

Regeneration rewrites the cog file from its definition, so custom code in
command bodies is not preserved. A .py.bak copy of the previous file is
written first unless --no-backup is given.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		noBackup, _ := cmd.Flags().GetBool("no-backup")
		if err := runBlockRemove(args[0], args[1], noBackup); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	},
}

// blockListRow is one line of block list, the JSON shape is the public one so the names stay stable
type blockListRow struct {
	Cog     string `json:"cog"`
	Graph   string `json:"graph"`
	Blocks  int    `json:"blocks"`
	Trigger string `json:"trigger"`
}

/**
 * findGraph
 * Finds a graph by name on a cog
 * @param cog {utils.CogConfig} - the cog to search
 * @param graphName {string} - the graph name
 * @return int - the index of the graph or -1
 **/
func findGraph(cog utils.CogConfig, graphName string) int {
	for i, graph := range cog.BlockGraphs {
		if graph.Name == graphName {
			return i
		}
	}
	return -1
}

/**
 * loadGraphInput
 * Reads the --graph flag value and fills in the optional keys
 * @param rawGraph {string} - the --graph flag value
 * @return utils.BlockGraph - the normalized graph
 * @return error - the read or parse failure
 **/
func loadGraphInput(rawGraph string) (utils.BlockGraph, error) {
	if rawGraph == "" {
		return utils.BlockGraph{}, fmt.Errorf("a graph is required, set --graph with inline JSON, @path/to/file.json, or - for stdin")
	}

	data, err := readJSONInput(rawGraph, "graph")
	if err != nil {
		return utils.BlockGraph{}, err
	}

	parsedGraph, err := utils.JSONToBlockGraph(string(data))
	if err != nil {
		return utils.BlockGraph{}, fmt.Errorf("error parsing graph JSON: %w", err)
	}

	// Defaults are filled in before validation so the saved blob always carries every key
	return utils.NormalizeBlockGraph(*parsedGraph), nil
}

/**
 * runBlockAdd
 * Validates a graph, attaches it to the cog, regenerates the cog file, and saves the config
 * @param cogName {string} - the cog name or file base to attach the graph to
 * @param rawGraph {string} - the --graph flag value
 * @param noBackup {bool} - true to skip the .py.bak backup
 * @param replace {bool} - true to swap out the graph of the same name instead of rejecting it
 * @return error - the first validation, render, or write failure
 **/
func runBlockAdd(cogName string, rawGraph string, noBackup bool, replace bool) error {
	rootDir, err := utils.FindBotConf()
	if err != nil {
		return fmt.Errorf("not in a botbox project: %w", err)
	}

	config, err := utils.LoadConfig()
	if err != nil {
		return err
	}

	cogIndex, found := findProjectCog(config, cogName)
	if !found {
		return fmt.Errorf("cog '%s' does not exist in the project", cogName)
	}
	cog := config.Cogs[cogIndex]

	graph, err := loadGraphInput(rawGraph)
	if err != nil {
		return err
	}

	// A replace swaps the graph in place so the generated file keeps its order
	replaceIndex := -1
	replaceName := ""
	if replace {
		replaceIndex = findGraph(cog, graph.Name)
		if replaceIndex < 0 {
			return fmt.Errorf("graph '%s' does not exist on cog '%s', drop --replace to add it", graph.Name, cog.Name)
		}
		replaceName = graph.Name
	}

	if replaceIndex < 0 && len(cog.BlockGraphs) >= utils.MaxGraphsPerCog {
		return fmt.Errorf("cog '%s' already has %d graphs, the maximum", cog.Name, utils.MaxGraphsPerCog)
	}

	if err := utils.ValidateBlockGraph(graph, config, cogIndex, replaceName); err != nil {
		return err
	}

	// Older projects predate src/utils/blocks.py, the generated code imports it
	if err := utils.EnsureBlocksHelper(rootDir); err != nil {
		return err
	}

	if replaceIndex >= 0 {
		cog.BlockGraphs[replaceIndex] = graph
	} else {
		cog.BlockGraphs = append(cog.BlockGraphs, graph)
	}

	if err := utils.RegenerateCogFile(rootDir, config, cog, !noBackup); err != nil {
		return fmt.Errorf("error regenerating cog file: %w", err)
	}

	config.Cogs[cogIndex] = cog
	if err := saveBlockConfig(rootDir, config); err != nil {
		return err
	}

	printBlockResult(rootDir, cog, noBackup)
	return nil
}

/**
 * runBlockList
 * Prints one line per graph, filtered to one cog when a name is given
 * @param cogName {string} - the cog name or file base to filter on, empty for every cog
 * @param asJSON {bool} - true to print the rows as a JSON array
 * @return error - the load failure or an unknown cog name
 **/
func runBlockList(cogName string, asJSON bool) error {
	if _, err := utils.FindBotConf(); err != nil {
		return fmt.Errorf("not in a botbox project: %w", err)
	}

	config, err := utils.LoadConfig()
	if err != nil {
		return err
	}

	cogs := config.Cogs
	if cogName != "" {
		cogIndex, found := findProjectCog(config, cogName)
		if !found {
			return fmt.Errorf("cog '%s' does not exist in the project", cogName)
		}
		cogs = config.Cogs[cogIndex : cogIndex+1]
	}

	rows := []blockListRow{}
	for _, cog := range cogs {
		for _, graph := range cog.BlockGraphs {
			// The trigger identifies the graph on discord, so the listing shows it
			trigger := ""
			if block, found := utils.GraphTrigger(graph); found && block.Type == utils.BlockTriggerSlashCommand {
				trigger = "/" + block.Config["name"]
			}
			rows = append(rows, blockListRow{Cog: cog.File, Graph: graph.Name, Blocks: len(graph.Blocks), Trigger: trigger})
		}
	}

	if asJSON {
		return printJSON(rows)
	}

	for _, row := range rows {
		fmt.Printf("%s %s blocks=%d trigger=%s\n", row.Cog, row.Graph, row.Blocks, row.Trigger)
	}

	return nil
}

/**
 * runBlockShow
 * Prints one graph as indented JSON, this doubles as the export path
 * @param cogName {string} - the cog name or file base holding the graph
 * @param graphName {string} - the name of the graph to print
 * @return error - the lookup failure
 **/
func runBlockShow(cogName string, graphName string) error {
	if _, err := utils.FindBotConf(); err != nil {
		return fmt.Errorf("not in a botbox project: %w", err)
	}

	config, err := utils.LoadConfig()
	if err != nil {
		return err
	}

	cogIndex, found := findProjectCog(config, cogName)
	if !found {
		return fmt.Errorf("cog '%s' does not exist in the project", cogName)
	}
	cog := config.Cogs[cogIndex]

	graphIndex := findGraph(cog, graphName)
	if graphIndex < 0 {
		return fmt.Errorf("graph '%s' does not exist on cog '%s'", graphName, cog.Name)
	}

	return printJSON(cog.BlockGraphs[graphIndex])
}

/**
 * runBlockValidate
 * Normalizes and validates a graph without writing anything, for scripts and the builder's live check
 * @param cogName {string} - the cog the graph would attach to, empty validates against a blank cog
 * @param rawGraph {string} - the --graph flag value
 * @return error - the validation failure
 **/
func runBlockValidate(cogName string, rawGraph string) error {
	graph, err := loadGraphInput(rawGraph)
	if err != nil {
		return err
	}

	// Outside a project or without a cog the graph is checked against an empty cog
	config := utils.Config{Cogs: []utils.CogConfig{{Name: "Scratch", File: "scratch"}}}
	cogIndex := 0
	if cogName != "" {
		if _, err := utils.FindBotConf(); err != nil {
			return fmt.Errorf("not in a botbox project: %w", err)
		}
		loaded, err := utils.LoadConfig()
		if err != nil {
			return err
		}
		index, found := findProjectCog(loaded, cogName)
		if !found {
			return fmt.Errorf("cog '%s' does not exist in the project", cogName)
		}
		config = loaded
		cogIndex = index
	}

	if err := utils.ValidateBlockGraph(graph, config, cogIndex, ""); err != nil {
		return err
	}

	fmt.Printf("graph '%s' is valid, %d blocks, %d edges\n", graph.Name, len(graph.Blocks), len(graph.Edges))
	return nil
}

/**
 * runBlockTypes
 * Prints the block palette, one line per type or the JSON list
 * @param asJSON {bool} - true to print the palette as JSON
 * @return error - the marshal failure
 **/
func runBlockTypes(asJSON bool) error {
	types := utils.BlockTypes()
	if asJSON {
		return printJSON(types)
	}

	for _, blockType := range types {
		fmt.Printf("%s required=%s optional=%s ports=%s\n",
			blockType.Type,
			strings.Join(blockType.Required, ","),
			strings.Join(blockType.Optional, ","),
			strings.Join(blockType.Ports, ","))
		fmt.Printf("    %s\n", blockType.Notes)
	}
	fmt.Printf("builtin values: %s\n", strings.Join(utils.BuiltinRunValues(), ", "))

	return nil
}

/**
 * runBlockRemove
 * Removes a graph from a cog, regenerates the cog file, and saves the config
 * @param cogName {string} - the cog name or file base holding the graph
 * @param graphName {string} - the name of the graph to remove
 * @param noBackup {bool} - true to skip the .py.bak backup
 * @return error - the first lookup, render, or write failure
 **/
func runBlockRemove(cogName string, graphName string, noBackup bool) error {
	rootDir, err := utils.FindBotConf()
	if err != nil {
		return fmt.Errorf("not in a botbox project: %w", err)
	}

	config, err := utils.LoadConfig()
	if err != nil {
		return err
	}

	cogIndex, found := findProjectCog(config, cogName)
	if !found {
		return fmt.Errorf("cog '%s' does not exist in the project", cogName)
	}
	cog := config.Cogs[cogIndex]

	graphIndex := findGraph(cog, graphName)
	if graphIndex < 0 {
		return fmt.Errorf("graph '%s' does not exist on cog '%s'", graphName, cog.Name)
	}

	cog.BlockGraphs = slices.Delete(cog.BlockGraphs, graphIndex, graphIndex+1)
	// An empty slice marshals as [] while nil is omitted, nil keeps old configs byte stable
	if len(cog.BlockGraphs) == 0 {
		cog.BlockGraphs = nil
	}

	if err := utils.RegenerateCogFile(rootDir, config, cog, !noBackup); err != nil {
		return fmt.Errorf("error regenerating cog file: %w", err)
	}

	config.Cogs[cogIndex] = cog
	if err := saveBlockConfig(rootDir, config); err != nil {
		return err
	}

	printBlockResult(rootDir, cog, noBackup)
	return nil
}

/**
 * printJSON
 * Prints any value as indented JSON on stdout
 * @param value {any} - the value to print
 * @return error - the marshal failure
 **/
func printJSON(value any) error {
	jsonData, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}

/**
 * saveBlockConfig
 * Writes the updated config back to botbox.conf
 * @param rootDir {string} - the project root holding botbox.conf
 * @param config {utils.Config} - the config to write
 * @return error - the marshal or write failure
 **/
func saveBlockConfig(rootDir string, config utils.Config) error {
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	confPath := filepath.Join(rootDir, "botbox.conf")
	if err := os.WriteFile(confPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write updated botbox.conf: %w", err)
	}
	return nil
}

/**
 * printBlockResult
 * Prints the regeneration warning to stderr and the written cog path to stdout
 * @param rootDir {string} - the project root
 * @param cog {utils.CogConfig} - the cog whose file was regenerated
 * @param noBackup {bool} - true when the backup was skipped
 **/
func printBlockResult(rootDir string, cog utils.CogConfig, noBackup bool) {
	writtenPath := filepath.Join(rootDir, "src", "cogs", cog.File+".py")
	warning := "regeneration rewrites the cog file from its definition, custom code in command bodies is not preserved"
	if !noBackup {
		if _, err := os.Stat(writtenPath + ".bak"); err == nil {
			warning += ", previous version saved to " + writtenPath + ".bak"
		}
	}
	fmt.Fprintln(os.Stderr, "Warning:", warning)
	fmt.Println(writtenPath)
}

func init() {
	rootCmd.AddCommand(blockCmd)
	blockCmd.AddCommand(blockAddCmd)
	blockCmd.AddCommand(blockListCmd)
	blockCmd.AddCommand(blockShowCmd)
	blockCmd.AddCommand(blockValidateCmd)
	blockCmd.AddCommand(blockTypesCmd)
	blockCmd.AddCommand(blockRemoveCmd)
	blockAddCmd.Flags().String("graph", "", "BlockGraph JSON document, accepts inline JSON, @path/to/file.json, or - for stdin")
	blockAddCmd.Flags().Bool("no-backup", false, "Skip the .py.bak backup written before the cog file is regenerated")
	blockAddCmd.Flags().Bool("replace", false, "Replace the existing graph of the same name instead of rejecting it")
	blockListCmd.Flags().Bool("json", false, "Print the graphs as a JSON array")
	blockValidateCmd.Flags().String("graph", "", "BlockGraph JSON document, accepts inline JSON, @path/to/file.json, or - for stdin")
	blockValidateCmd.Flags().String("cog", "", "Cog to validate the graph against, defaults to an empty cog")
	blockTypesCmd.Flags().Bool("json", false, "Print the palette as JSON")
	blockRemoveCmd.Flags().Bool("no-backup", false, "Skip the .py.bak backup written before the cog file is regenerated")
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
