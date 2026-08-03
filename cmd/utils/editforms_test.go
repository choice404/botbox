/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package utils

import (
	"testing"
)

// The edit flow indexes mirror the named consts inside EditFormWrapperGenerator
const (
	editIdxSelectCog = iota
	editIdxAction
	editIdxCmdInfo
	editIdxArgStart
	editIdxArgInfo
	editIdxFieldStart
	editIdxFieldInfo
	editIdxAccept
	editIdxMultiPage
	editIdxPageInfo
	editIdxPageFieldStart
	editIdxPageFieldInfo
	editIdxBranchStart
	editIdxBranchInfo
	editIdxPageNext
	editIdxModInfo
	editIdxRedefine
	editIdxRedefineResponses
	editIdxResponseStart
	editIdxResponseInfo
	editIdxPickCommand
	editIdxRemoveCommand
)

// newEditModelValues builds the model value bus the edit flow expects
func newEditModelValues() Values {
	values := map[string]*string{
		"cogName":         new(string),
		"cogEnv":          new(string),
		"backup":          new(string),
		"editingOriginal": new(string),
		"currentCommand":  new(string),
		"currentPage":     new(string),
		"pages":           new(string),
		"slashCommands":   new(string),
		"prefixCommands":  new(string),
	}
	emptySlash := "[]"
	emptyPrefix := "[]"
	emptyPages := "[]"
	values["slashCommands"] = &emptySlash
	values["prefixCommands"] = &emptyPrefix
	values["pages"] = &emptyPages
	return Values{Map: values, Name: "ModelValues"}
}

// editGreetCommand builds a slash command fixture with one argument and one response
func editGreetCommand() CommandInfo {
	return CommandInfo{
		Name:        "greet",
		Type:        "slash",
		Scope:       "guild",
		Description: "Greets someone",
		Args:        []ArgInfo{{Name: "target", Type: "str", Description: "Who to greet"}},
		ReturnType:  "str",
		Responses:   []ResponseInfo{{Type: "message", Content: "hi"}},
	}
}

func TestEditFormWrapperGeneratorFormCount(t *testing.T) {
	forms := EditFormWrapperGenerator()
	if len(forms) != editIdxRemoveCommand+1 {
		t.Fatalf("expected %d forms, got %d", editIdxRemoveCommand+1, len(forms))
	}
}

func TestEditActionRouting(t *testing.T) {
	tests := []struct {
		name   string
		action string
		want   int
	}{
		{"add enters the add command flow", "add", editIdxCmdInfo},
		{"edit picks a command", "edit", editIdxPickCommand},
		{"remove picks a command to remove", "remove", editIdxRemoveCommand},
		{"apply ends the flow", "apply", -2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forms := EditFormWrapperGenerator()
			setFormValue(forms, editIdxAction, "editAction", tt.action)
			got := forms[editIdxAction].BranchCallback(forms[editIdxAction].Values, forms)
			if got != tt.want {
				t.Errorf("action %q routed to %d, want %d", tt.action, got, tt.want)
			}
		})
	}
}

func TestEditActionAddResetsCommandState(t *testing.T) {
	forms := EditFormWrapperGenerator()
	modelValues := newEditModelValues()
	setModelValue(modelValues, "currentCommand", `{"Name":"stale"}`)
	setModelValue(modelValues, "editingOriginal", `{"Name":"stale"}`)
	setModelValue(modelValues, "pages", `[{"Name":"stale"}]`)
	setFormValue(forms, editIdxCmdInfo, "cmdName", "stale")
	setFormValue(forms, editIdxArgInfo, "args", "stale")
	setFormValue(forms, editIdxResponseInfo, "responses", "stale")
	setFormValue(forms, editIdxAction, "editAction", "add")

	forms[editIdxAction].Callback(forms[editIdxAction].Values, modelValues, forms)

	if *modelValues.Map["currentCommand"] != "" {
		t.Errorf("current command = %q, want empty", *modelValues.Map["currentCommand"])
	}
	if *modelValues.Map["editingOriginal"] != "" {
		t.Errorf("editing original = %q, want empty", *modelValues.Map["editingOriginal"])
	}
	if *modelValues.Map["pages"] != "[]" {
		t.Errorf("pages = %q, want empty collection", *modelValues.Map["pages"])
	}
	for index, key := range map[int]string{
		editIdxCmdInfo:      "cmdName",
		editIdxArgInfo:      "args",
		editIdxResponseInfo: "responses",
	} {
		if *forms[index].Values.Map[key] != "" {
			t.Errorf("form %d key %q should reset, got %q", index, key, *forms[index].Values.Map[key])
		}
	}
}

func TestEditPickCommandPrefillsInfoForm(t *testing.T) {
	forms := EditFormWrapperGenerator()
	modelValues := newEditModelValues()
	command := editGreetCommand()
	slashJSON, _ := CmdInfoSliceToJSON([]CommandInfo{command})
	setModelValue(modelValues, "slashCommands", slashJSON)
	setFormValue(forms, editIdxPickCommand, "editCmdName", "greet")

	forms[editIdxPickCommand].Callback(forms[editIdxPickCommand].Values, modelValues, forms)

	// The picked command leaves its list so name checks run against the rest
	remaining, _ := JSONToCmdInfoSlice(*modelValues.Map["slashCommands"])
	if len(remaining) != 0 {
		t.Errorf("slash list holds %d commands after pick, want 0", len(remaining))
	}

	current, err := JSONToCmdInfo(*modelValues.Map["currentCommand"])
	if err != nil {
		t.Fatalf("failed to parse current command: %v", err)
	}
	if current.Name != "greet" || len(current.Args) != 1 || len(current.Responses) != 1 {
		t.Errorf("current command = %+v, want the picked command", current)
	}
	if *modelValues.Map["editingOriginal"] == "" {
		t.Error("editing original should hold the picked command")
	}

	// The info form starts prefilled with the picked command's values
	for key, want := range map[string]string{
		"cmdName":        "greet",
		"cmdType":        "slash",
		"cmdScope":       "guild",
		"cmdDescription": "Greets someone",
		"cmdReturnType":  "str",
	} {
		if got := *forms[editIdxModInfo].Values.Map[key]; got != want {
			t.Errorf("mod info %q = %q, want %q", key, got, want)
		}
	}

	if got := forms[editIdxPickCommand].BranchCallback(forms[editIdxPickCommand].Values, forms); got != editIdxModInfo {
		t.Errorf("pick routed to %d, want %d", got, editIdxModInfo)
	}
}

func TestEditPickCommandMissingNameReturnsToAction(t *testing.T) {
	forms := EditFormWrapperGenerator()
	modelValues := newEditModelValues()
	setFormValue(forms, editIdxPickCommand, "editCmdName", "ghost")

	forms[editIdxPickCommand].Callback(forms[editIdxPickCommand].Values, modelValues, forms)

	if got := forms[editIdxPickCommand].BranchCallback(forms[editIdxPickCommand].Values, forms); got != editIdxAction {
		t.Errorf("missing command routed to %d, want %d", got, editIdxAction)
	}
}

func TestEditModInfoPreservesCollectedSets(t *testing.T) {
	forms := EditFormWrapperGenerator()
	modelValues := newEditModelValues()
	command := editGreetCommand()
	commandString, _ := command.ToJSON()
	setModelValue(modelValues, "currentCommand", commandString)
	setFormValue(forms, editIdxModInfo, "cmdName", "salute")
	setFormValue(forms, editIdxModInfo, "cmdType", "slash")
	setFormValue(forms, editIdxModInfo, "cmdScope", "global")
	setFormValue(forms, editIdxModInfo, "cmdDescription", "Salutes someone")
	setFormValue(forms, editIdxModInfo, "cmdReturnType", "None")

	forms[editIdxModInfo].Callback(forms[editIdxModInfo].Values, modelValues, forms)

	current, err := JSONToCmdInfo(*modelValues.Map["currentCommand"])
	if err != nil {
		t.Fatalf("failed to parse current command: %v", err)
	}
	if current.Name != "salute" || current.Scope != "global" || current.Description != "Salutes someone" || current.ReturnType != "None" {
		t.Errorf("current command info = %+v, want the new values", current)
	}
	if len(current.Args) != 1 || len(current.Responses) != 1 {
		t.Errorf("collected sets should survive the info edit, got %+v", current)
	}

	if got := forms[editIdxModInfo].BranchCallback(forms[editIdxModInfo].Values, forms); got != editIdxRedefine {
		t.Errorf("mod info routed to %d, want %d", got, editIdxRedefine)
	}
}

func TestEditRedefineNoKeepsSets(t *testing.T) {
	forms := EditFormWrapperGenerator()
	modelValues := newEditModelValues()
	command := editGreetCommand()
	commandString, _ := command.ToJSON()
	setModelValue(modelValues, "currentCommand", commandString)
	setFormValue(forms, editIdxRedefine, "redefineConfirm", "no")

	forms[editIdxRedefine].Callback(forms[editIdxRedefine].Values, modelValues, forms)

	current, _ := JSONToCmdInfo(*modelValues.Map["currentCommand"])
	if len(current.Args) != 1 {
		t.Errorf("redefine no should keep the args, got %+v", current.Args)
	}
	if got := forms[editIdxRedefine].BranchCallback(forms[editIdxRedefine].Values, forms); got != editIdxRedefineResponses {
		t.Errorf("redefine no routed to %d, want %d", got, editIdxRedefineResponses)
	}
}

func TestEditRedefineYesClearsSetsAndRoutesByType(t *testing.T) {
	tests := []struct {
		name    string
		cmdType string
		want    int
	}{
		{"slash redefine enters the arg loop", "slash", editIdxArgStart},
		{"modal redefine enters the multi page confirm", "modal", editIdxMultiPage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forms := EditFormWrapperGenerator()
			modelValues := newEditModelValues()
			command := editGreetCommand()
			command.Type = tt.cmdType
			commandString, _ := command.ToJSON()
			setModelValue(modelValues, "currentCommand", commandString)
			setModelValue(modelValues, "pages", `[{"Name":"stale"}]`)
			setFormValue(forms, editIdxArgInfo, "args", "stale")
			setFormValue(forms, editIdxRedefine, "redefineConfirm", "yes")

			forms[editIdxRedefine].Callback(forms[editIdxRedefine].Values, modelValues, forms)

			current, _ := JSONToCmdInfo(*modelValues.Map["currentCommand"])
			if len(current.Args) != 0 || len(current.Fields) != 0 || len(current.Pages) != 0 {
				t.Errorf("redefine yes should clear the sets, got %+v", current)
			}
			if *modelValues.Map["pages"] != "[]" {
				t.Errorf("pages = %q, want empty collection", *modelValues.Map["pages"])
			}
			if *forms[editIdxArgInfo].Values.Map["args"] != "" {
				t.Error("redefine yes should reset the arg accumulator")
			}

			if got := forms[editIdxRedefine].BranchCallback(forms[editIdxRedefine].Values, forms); got != tt.want {
				t.Errorf("redefine yes routed to %d, want %d", got, tt.want)
			}
		})
	}
}

func TestEditRedefineResponsesSkipCondition(t *testing.T) {
	command := editGreetCommand()
	commandString, _ := command.ToJSON()
	noResponses := command
	noResponses.Responses = nil
	noResponsesString, _ := noResponses.ToJSON()

	tests := []struct {
		name     string
		original string
		current  string
		want     bool
	}{
		{"add action skips", "", commandString, true},
		{"edited command with responses shows", commandString, commandString, false},
		{"edited command without responses skips", commandString, noResponsesString, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forms := EditFormWrapperGenerator()
			modelValues := newEditModelValues()
			setModelValue(modelValues, "editingOriginal", tt.original)
			setModelValue(modelValues, "currentCommand", tt.current)
			got := forms[editIdxRedefineResponses].SkipCondition(modelValues, forms, editIdxRedefineResponses)
			if got != tt.want {
				t.Errorf("skip = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditRedefineResponsesRoutingAndClear(t *testing.T) {
	forms := EditFormWrapperGenerator()
	modelValues := newEditModelValues()
	command := editGreetCommand()
	commandString, _ := command.ToJSON()
	setModelValue(modelValues, "currentCommand", commandString)
	setFormValue(forms, editIdxResponseInfo, "responses", "stale")

	setFormValue(forms, editIdxRedefineResponses, "redefineResponsesConfirm", "yes")
	forms[editIdxRedefineResponses].Callback(forms[editIdxRedefineResponses].Values, modelValues, forms)
	current, _ := JSONToCmdInfo(*modelValues.Map["currentCommand"])
	if len(current.Responses) != 0 {
		t.Errorf("redefine responses yes should clear the responses, got %+v", current.Responses)
	}
	if *forms[editIdxResponseInfo].Values.Map["responses"] != "" {
		t.Error("redefine responses yes should reset the response accumulator")
	}
	if got := forms[editIdxRedefineResponses].BranchCallback(forms[editIdxRedefineResponses].Values, forms); got != -1 {
		t.Errorf("yes routed to %d, want -1", got)
	}

	setFormValue(forms, editIdxRedefineResponses, "redefineResponsesConfirm", "no")
	if got := forms[editIdxRedefineResponses].BranchCallback(forms[editIdxRedefineResponses].Values, forms); got != editIdxAccept {
		t.Errorf("no routed to %d, want %d", got, editIdxAccept)
	}
}

func TestEditAcceptRestoresOriginalOnNo(t *testing.T) {
	forms := EditFormWrapperGenerator()
	modelValues := newEditModelValues()
	original := editGreetCommand()
	originalString, _ := original.ToJSON()
	edited := original
	edited.Name = "salute"
	editedString, _ := edited.ToJSON()
	setModelValue(modelValues, "editingOriginal", originalString)
	setModelValue(modelValues, "currentCommand", editedString)
	setFormValue(forms, editIdxAccept, "cmdAcceptConfirm", "no")

	forms[editIdxAccept].Callback(forms[editIdxAccept].Values, modelValues, forms)

	slashCommands, _ := JSONToCmdInfoSlice(*modelValues.Map["slashCommands"])
	if len(slashCommands) != 1 || slashCommands[0].Name != "greet" {
		t.Errorf("slash list = %+v, want the restored original", slashCommands)
	}
	if *modelValues.Map["editingOriginal"] != "" {
		t.Error("editing original should clear after the accept decision")
	}
	if got := forms[editIdxAccept].BranchCallback(forms[editIdxAccept].Values, forms); got != editIdxAction {
		t.Errorf("accept routed to %d, want %d", got, editIdxAction)
	}
}

func TestEditAcceptYesAppendsEditedCommand(t *testing.T) {
	forms := EditFormWrapperGenerator()
	modelValues := newEditModelValues()
	original := editGreetCommand()
	originalString, _ := original.ToJSON()
	edited := original
	edited.Name = "salute"
	editedString, _ := edited.ToJSON()
	setModelValue(modelValues, "editingOriginal", originalString)
	setModelValue(modelValues, "currentCommand", editedString)
	setFormValue(forms, editIdxAccept, "cmdAcceptConfirm", "yes")

	forms[editIdxAccept].Callback(forms[editIdxAccept].Values, modelValues, forms)

	slashCommands, _ := JSONToCmdInfoSlice(*modelValues.Map["slashCommands"])
	if len(slashCommands) != 1 || slashCommands[0].Name != "salute" {
		t.Errorf("slash list = %+v, want the edited command", slashCommands)
	}
	if *modelValues.Map["editingOriginal"] != "" {
		t.Error("editing original should clear after the accept decision")
	}
}

func TestEditRemoveCommandCallback(t *testing.T) {
	forms := EditFormWrapperGenerator()
	modelValues := newEditModelValues()
	prefixJSON, _ := CmdInfoSliceToJSON([]CommandInfo{{Name: "wave", Type: "prefix", Scope: "guild", Description: "d", ReturnType: "None"}})
	setModelValue(modelValues, "prefixCommands", prefixJSON)

	// A declined confirm leaves the command in place
	setFormValue(forms, editIdxRemoveCommand, "removeCmdName", "wave")
	setFormValue(forms, editIdxRemoveCommand, "removeConfirm", "no")
	forms[editIdxRemoveCommand].Callback(forms[editIdxRemoveCommand].Values, modelValues, forms)
	prefixCommands, _ := JSONToCmdInfoSlice(*modelValues.Map["prefixCommands"])
	if len(prefixCommands) != 1 {
		t.Errorf("declined remove left %d commands, want 1", len(prefixCommands))
	}

	setFormValue(forms, editIdxRemoveCommand, "removeConfirm", "yes")
	forms[editIdxRemoveCommand].Callback(forms[editIdxRemoveCommand].Values, modelValues, forms)
	prefixCommands, _ = JSONToCmdInfoSlice(*modelValues.Map["prefixCommands"])
	if len(prefixCommands) != 0 {
		t.Errorf("confirmed remove left %d commands, want 0", len(prefixCommands))
	}

	if got := forms[editIdxRemoveCommand].BranchCallback(forms[editIdxRemoveCommand].Values, forms); got != editIdxAction {
		t.Errorf("remove routed to %d, want %d", got, editIdxAction)
	}
}

func TestEditLoopExitsRouteToRedefineResponses(t *testing.T) {
	forms := EditFormWrapperGenerator()

	setFormValue(forms, editIdxArgStart, "argStartConfirm", "no")
	if got := forms[editIdxArgStart].BranchCallback(forms[editIdxArgStart].Values, forms); got != editIdxRedefineResponses {
		t.Errorf("arg start no routed to %d, want %d", got, editIdxRedefineResponses)
	}

	setFormValue(forms, editIdxFieldStart, "fieldStartConfirm", "no")
	if got := forms[editIdxFieldStart].BranchCallback(forms[editIdxFieldStart].Values, forms); got != editIdxRedefineResponses {
		t.Errorf("field start no routed to %d, want %d", got, editIdxRedefineResponses)
	}

	onePage, _ := PageInfoSliceToJSON([]PageInfo{{Name: "intro"}})
	setFormValue(forms, editIdxPageNext, "pages", onePage)
	setFormValue(forms, editIdxPageNext, "pageAnotherConfirm", "no")
	if got := forms[editIdxPageNext].BranchCallback(forms[editIdxPageNext].Values, forms); got != editIdxRedefineResponses {
		t.Errorf("page next done routed to %d, want %d", got, editIdxRedefineResponses)
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
