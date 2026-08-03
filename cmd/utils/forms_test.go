/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package utils

import (
	"strings"
	"testing"
)

// The add flow indexes mirror the named consts inside AddFormWrapperGenerator
const (
	testIdxFileName = iota
	testIdxCmdStart
	testIdxCmdInfo
	testIdxArgStart
	testIdxArgInfo
	testIdxFieldStart
	testIdxFieldInfo
	testIdxAccept
	testIdxMultiPage
	testIdxPageInfo
	testIdxPageFieldStart
	testIdxPageFieldInfo
	testIdxBranchStart
	testIdxBranchInfo
	testIdxPageNext
	testIdxResponseStart
	testIdxResponseInfo
)

// setFormValue plants a value on a wrapper as if the form had collected it
func setFormValue(forms []FormWrapper, index int, key string, value string) {
	v := value
	forms[index].Values.Map[key] = &v
}

// newAddModelValues builds the model value bus the add flow expects
func newAddModelValues() Values {
	values := map[string]*string{
		"filename":       new(string),
		"currentCommand": new(string),
		"currentPage":    new(string),
		"pages":          new(string),
		"slashCommands":  new(string),
		"prefixCommands": new(string),
	}
	emptySlash := "[]"
	emptyPrefix := "[]"
	emptyPages := "[]"
	values["slashCommands"] = &emptySlash
	values["prefixCommands"] = &emptyPrefix
	values["pages"] = &emptyPages
	return Values{Map: values, Name: "ModelValues"}
}

// setModelValue plants a value on the model value bus
func setModelValue(modelValues Values, key string, value string) {
	v := value
	modelValues.Map[key] = &v
}

// makeFields builds a field slice with unique names for accumulator tests
func makeFields(count int) []FieldInfo {
	fields := make([]FieldInfo, count)
	for i := range fields {
		fields[i] = FieldInfo{Name: "field" + string(rune('a'+i)), Label: "Label", Style: "short"}
	}
	return fields
}

func TestAddFormWrapperGeneratorFormCount(t *testing.T) {
	forms := AddFormWrapperGenerator()
	if len(forms) != testIdxResponseInfo+1 {
		t.Fatalf("expected %d forms, got %d", testIdxResponseInfo+1, len(forms))
	}
}

func TestCmdInfoBranchRouting(t *testing.T) {
	tests := []struct {
		name    string
		cmdType string
		want    int
	}{
		{"modal goes to multi page confirm", "modal", testIdxMultiPage},
		{"slash goes to the arg loop", "slash", -1},
		{"prefix goes to the arg loop", "prefix", -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forms := AddFormWrapperGenerator()
			setFormValue(forms, testIdxCmdInfo, "cmdType", tt.cmdType)
			got := forms[testIdxCmdInfo].BranchCallback(forms[testIdxCmdInfo].Values, forms)
			if got != tt.want {
				t.Errorf("cmdType %q routed to %d, want %d", tt.cmdType, got, tt.want)
			}
		})
	}
}

func TestMultiPageConfirmRouting(t *testing.T) {
	tests := []struct {
		name    string
		confirm string
		want    int
	}{
		{"yes enters the page loop", "yes", testIdxPageInfo},
		{"no keeps the single page field loop", "no", testIdxFieldStart},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forms := AddFormWrapperGenerator()
			setFormValue(forms, testIdxMultiPage, "multiPageConfirm", tt.confirm)
			got := forms[testIdxMultiPage].BranchCallback(forms[testIdxMultiPage].Values, forms)
			if got != tt.want {
				t.Errorf("confirm %q routed to %d, want %d", tt.confirm, got, tt.want)
			}
		})
	}
}

func TestMultiPageCallbackResetsPageCollections(t *testing.T) {
	forms := AddFormWrapperGenerator()
	modelValues := newAddModelValues()
	setModelValue(modelValues, "pages", `[{"Name":"stale"}]`)
	setFormValue(forms, testIdxPageNext, "pages", `[{"Name":"stale"}]`)

	forms[testIdxMultiPage].Callback(forms[testIdxMultiPage].Values, modelValues, forms)

	if *modelValues.Map["pages"] != "[]" {
		t.Errorf("model pages = %q, want empty collection", *modelValues.Map["pages"])
	}
	if *forms[testIdxPageNext].Values.Map["pages"] != "[]" {
		t.Errorf("page next mirror = %q, want empty collection", *forms[testIdxPageNext].Values.Map["pages"])
	}
}

func TestPageInfoCallbackBuildsCurrentPage(t *testing.T) {
	forms := AddFormWrapperGenerator()
	modelValues := newAddModelValues()
	setFormValue(forms, testIdxPageInfo, "pageName", "intro")
	setFormValue(forms, testIdxPageInfo, "pageTitle", "Intro Page")
	setFormValue(forms, testIdxPageFieldInfo, "pageFields", "stale")

	forms[testIdxPageInfo].Callback(forms[testIdxPageInfo].Values, modelValues, forms)

	page, err := jsonToPage(*modelValues.Map["currentPage"])
	if err != nil {
		t.Fatalf("failed to parse current page: %v", err)
	}
	if page.Name != "intro" || page.Title != "Intro Page" {
		t.Errorf("current page = %+v, want name intro and title Intro Page", page)
	}
	if len(page.Fields) != 0 || len(page.Branches) != 0 {
		t.Errorf("new page should start with no fields or branches, got %+v", page)
	}
	if *forms[testIdxPageFieldInfo].Values.Map["pageFields"] != "" {
		t.Error("page info callback should reset the page field accumulator")
	}
}

func TestPageFieldStartRouting(t *testing.T) {
	forms := AddFormWrapperGenerator()
	setFormValue(forms, testIdxPageFieldStart, "pageFieldStartConfirm", "yes")
	if got := forms[testIdxPageFieldStart].BranchCallback(forms[testIdxPageFieldStart].Values, forms); got != -1 {
		t.Errorf("yes routed to %d, want -1", got)
	}
	setFormValue(forms, testIdxPageFieldStart, "pageFieldStartConfirm", "no")
	if got := forms[testIdxPageFieldStart].BranchCallback(forms[testIdxPageFieldStart].Values, forms); got != testIdxBranchStart {
		t.Errorf("no routed to %d, want %d", got, testIdxBranchStart)
	}
}

func TestPageFieldInfoCallbackAppendsFieldToCurrentPage(t *testing.T) {
	forms := AddFormWrapperGenerator()
	modelValues := newAddModelValues()
	pageString, _ := pageToJSON(PageInfo{Name: "intro", Title: "Intro", Fields: []FieldInfo{}})
	setModelValue(modelValues, "currentPage", pageString)
	setFormValue(forms, testIdxPageFieldInfo, "fieldName", "answer")
	setFormValue(forms, testIdxPageFieldInfo, "fieldLabel", "Your answer")
	setFormValue(forms, testIdxPageFieldInfo, "fieldStyle", "short")
	setFormValue(forms, testIdxPageFieldInfo, "fieldRequired", "yes")
	setFormValue(forms, testIdxPageFieldInfo, "fieldPlaceholder", "type here")

	forms[testIdxPageFieldInfo].Callback(forms[testIdxPageFieldInfo].Values, modelValues, forms)

	page, err := jsonToPage(*modelValues.Map["currentPage"])
	if err != nil {
		t.Fatalf("failed to parse current page: %v", err)
	}
	if len(page.Fields) != 1 {
		t.Fatalf("expected 1 field on the page, got %d", len(page.Fields))
	}
	field := page.Fields[0]
	if field.Name != "answer" || field.Label != "Your answer" || field.Style != "short" || !field.Required || field.Placeholder != "type here" {
		t.Errorf("field = %+v, want the collected values", field)
	}

	accumulated, _ := JSONToFieldInfoSlice(*forms[testIdxPageFieldInfo].Values.Map["pageFields"])
	if len(accumulated) != 1 {
		t.Errorf("page field accumulator holds %d fields, want 1", len(accumulated))
	}
}

func TestPageFieldInfoBranchStopsAtMaxFields(t *testing.T) {
	forms := AddFormWrapperGenerator()
	fullFields, _ := FieldInfoSliceToJSON(makeFields(MaxModalFields))
	setFormValue(forms, testIdxPageFieldInfo, "pageFields", fullFields)
	if got := forms[testIdxPageFieldInfo].BranchCallback(forms[testIdxPageFieldInfo].Values, forms); got != testIdxBranchStart {
		t.Errorf("full page routed to %d, want %d", got, testIdxBranchStart)
	}

	partialFields, _ := FieldInfoSliceToJSON(makeFields(2))
	setFormValue(forms, testIdxPageFieldInfo, "pageFields", partialFields)
	if got := forms[testIdxPageFieldInfo].BranchCallback(forms[testIdxPageFieldInfo].Values, forms); got != testIdxPageFieldStart {
		t.Errorf("partial page routed to %d, want %d", got, testIdxPageFieldStart)
	}
}

func TestBranchStartRouting(t *testing.T) {
	forms := AddFormWrapperGenerator()
	fields, _ := FieldInfoSliceToJSON(makeFields(1))
	setFormValue(forms, testIdxPageFieldInfo, "pageFields", fields)

	setFormValue(forms, testIdxBranchStart, "branchStartConfirm", "yes")
	if got := forms[testIdxBranchStart].BranchCallback(forms[testIdxBranchStart].Values, forms); got != -1 {
		t.Errorf("yes with fields routed to %d, want -1", got)
	}

	setFormValue(forms, testIdxBranchStart, "branchStartConfirm", "no")
	if got := forms[testIdxBranchStart].BranchCallback(forms[testIdxBranchStart].Values, forms); got != testIdxPageNext {
		t.Errorf("no routed to %d, want %d", got, testIdxPageNext)
	}

	// A page without fields has nothing to branch on, so yes skips ahead
	setFormValue(forms, testIdxPageFieldInfo, "pageFields", "[]")
	setFormValue(forms, testIdxBranchStart, "branchStartConfirm", "yes")
	if got := forms[testIdxBranchStart].BranchCallback(forms[testIdxBranchStart].Values, forms); got != testIdxPageNext {
		t.Errorf("yes without fields routed to %d, want %d", got, testIdxPageNext)
	}
}

func TestBranchInfoCallbackAppendsBranchToCurrentPage(t *testing.T) {
	forms := AddFormWrapperGenerator()
	modelValues := newAddModelValues()
	pageString, _ := pageToJSON(PageInfo{Name: "intro", Title: "Intro", Fields: makeFields(1)})
	setModelValue(modelValues, "currentPage", pageString)
	setFormValue(forms, testIdxBranchInfo, "branchField", "fielda")
	setFormValue(forms, testIdxBranchInfo, "branchEquals", "yes")
	setFormValue(forms, testIdxBranchInfo, "branchGoto", "details")

	forms[testIdxBranchInfo].Callback(forms[testIdxBranchInfo].Values, modelValues, forms)

	page, err := jsonToPage(*modelValues.Map["currentPage"])
	if err != nil {
		t.Fatalf("failed to parse current page: %v", err)
	}
	if len(page.Branches) != 1 {
		t.Fatalf("expected 1 branch on the page, got %d", len(page.Branches))
	}
	branch := page.Branches[0]
	if branch.Field != "fielda" || branch.Equals != "yes" || branch.Goto != "details" {
		t.Errorf("branch = %+v, want the collected values", branch)
	}

	// The branch loop always returns to the branch confirm
	if got := forms[testIdxBranchInfo].BranchCallback(forms[testIdxBranchInfo].Values, forms); got != testIdxBranchStart {
		t.Errorf("branch info routed to %d, want %d", got, testIdxBranchStart)
	}
}

func TestPageNextCallbackCollectsPage(t *testing.T) {
	forms := AddFormWrapperGenerator()
	modelValues := newAddModelValues()
	command := CommandInfo{Name: "wizard", Type: "modal", Scope: "guild", Description: "d", ReturnType: "None"}
	commandString, _ := command.ToJSON()
	setModelValue(modelValues, "currentCommand", commandString)
	pageString, _ := pageToJSON(PageInfo{Name: "intro", Title: "Intro", Fields: makeFields(1)})
	setModelValue(modelValues, "currentPage", pageString)
	setFormValue(forms, testIdxPageNext, "pageNext", "details")
	setFormValue(forms, testIdxPageInfo, "pageName", "intro")
	setFormValue(forms, testIdxPageInfo, "pageTitle", "Intro")

	forms[testIdxPageNext].Callback(forms[testIdxPageNext].Values, modelValues, forms)

	pages, err := JSONToPageInfoSlice(*modelValues.Map["pages"])
	if err != nil {
		t.Fatalf("failed to parse collected pages: %v", err)
	}
	if len(pages) != 1 {
		t.Fatalf("expected 1 collected page, got %d", len(pages))
	}
	if pages[0].Name != "intro" || pages[0].Next != "details" {
		t.Errorf("collected page = %+v, want intro with next details", pages[0])
	}

	// The pages also ride on the command for the accept summary and save path
	currentCommand, err := JSONToCmdInfo(*modelValues.Map["currentCommand"])
	if err != nil {
		t.Fatalf("failed to parse current command: %v", err)
	}
	if len(currentCommand.Pages) != 1 || currentCommand.Pages[0].Name != "intro" {
		t.Errorf("command pages = %+v, want the collected page", currentCommand.Pages)
	}

	// The branch callback reads the mirrored pages off the form values
	mirrored, _ := JSONToPageInfoSlice(*forms[testIdxPageNext].Values.Map["pages"])
	if len(mirrored) != 1 {
		t.Errorf("mirrored pages hold %d entries, want 1", len(mirrored))
	}

	// The next page starts with clean page forms
	if *forms[testIdxPageInfo].Values.Map["pageName"] != "" || *forms[testIdxPageInfo].Values.Map["pageTitle"] != "" {
		t.Error("page info values should reset after a page is collected")
	}
	if *forms[testIdxPageNext].Values.Map["pageNext"] != "" {
		t.Error("page next value should reset after a page is collected")
	}
}

func TestPageNextBranchRouting(t *testing.T) {
	onePage, _ := PageInfoSliceToJSON([]PageInfo{{Name: "intro"}})
	fullPages := make([]PageInfo, MaxFlowPages)
	for i := range fullPages {
		fullPages[i] = PageInfo{Name: "page" + string(rune('a'+i))}
	}
	fullPagesString, _ := PageInfoSliceToJSON(fullPages)

	tests := []struct {
		name    string
		pages   string
		confirm string
		want    int
	}{
		{"another page loops back", onePage, "yes", testIdxPageInfo},
		{"done moves to responses", onePage, "no", testIdxResponseStart},
		{"full flow moves to responses even on yes", fullPagesString, "yes", testIdxResponseStart},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forms := AddFormWrapperGenerator()
			setFormValue(forms, testIdxPageNext, "pages", tt.pages)
			setFormValue(forms, testIdxPageNext, "pageAnotherConfirm", tt.confirm)
			got := forms[testIdxPageNext].BranchCallback(forms[testIdxPageNext].Values, forms)
			if got != tt.want {
				t.Errorf("routed to %d, want %d", got, tt.want)
			}
		})
	}
}

func TestArgAndFieldLoopsRouteToResponseConfirm(t *testing.T) {
	forms := AddFormWrapperGenerator()

	setFormValue(forms, testIdxArgStart, "argStartConfirm", "no")
	if got := forms[testIdxArgStart].BranchCallback(forms[testIdxArgStart].Values, forms); got != testIdxResponseStart {
		t.Errorf("arg start no routed to %d, want %d", got, testIdxResponseStart)
	}

	setFormValue(forms, testIdxFieldStart, "fieldStartConfirm", "no")
	if got := forms[testIdxFieldStart].BranchCallback(forms[testIdxFieldStart].Values, forms); got != testIdxResponseStart {
		t.Errorf("field start no routed to %d, want %d", got, testIdxResponseStart)
	}

	fullFields, _ := FieldInfoSliceToJSON(makeFields(MaxModalFields))
	setFormValue(forms, testIdxFieldInfo, "fields", fullFields)
	if got := forms[testIdxFieldInfo].BranchCallback(forms[testIdxFieldInfo].Values, forms); got != testIdxResponseStart {
		t.Errorf("full field loop routed to %d, want %d", got, testIdxResponseStart)
	}
}

func TestResponseStartRouting(t *testing.T) {
	forms := AddFormWrapperGenerator()
	setFormValue(forms, testIdxResponseStart, "responseStartConfirm", "yes")
	if got := forms[testIdxResponseStart].BranchCallback(forms[testIdxResponseStart].Values, forms); got != -1 {
		t.Errorf("yes routed to %d, want -1", got)
	}
	setFormValue(forms, testIdxResponseStart, "responseStartConfirm", "no")
	if got := forms[testIdxResponseStart].BranchCallback(forms[testIdxResponseStart].Values, forms); got != testIdxAccept {
		t.Errorf("no routed to %d, want %d", got, testIdxAccept)
	}
}

func TestResponseInfoCallbackAppendsResponse(t *testing.T) {
	forms := AddFormWrapperGenerator()
	modelValues := newAddModelValues()
	command := CommandInfo{Name: "greet", Type: "slash", Scope: "guild", Description: "d", ReturnType: "str"}
	commandString, _ := command.ToJSON()
	setModelValue(modelValues, "currentCommand", commandString)
	setFormValue(forms, testIdxResponseInfo, "responseContent", "Hello there")
	setFormValue(forms, testIdxResponseInfo, "responseEphemeral", "yes")

	forms[testIdxResponseInfo].Callback(forms[testIdxResponseInfo].Values, modelValues, forms)

	currentCommand, err := JSONToCmdInfo(*modelValues.Map["currentCommand"])
	if err != nil {
		t.Fatalf("failed to parse current command: %v", err)
	}
	if len(currentCommand.Responses) != 1 {
		t.Fatalf("expected 1 response on the command, got %d", len(currentCommand.Responses))
	}
	response := currentCommand.Responses[0]
	if response.Type != "message" || response.Content != "Hello there" || !response.Ephemeral {
		t.Errorf("response = %+v, want the collected values", response)
	}

	accumulated, _ := JSONToResponseInfoSlice(*forms[testIdxResponseInfo].Values.Map["responses"])
	if len(accumulated) != 1 {
		t.Errorf("response accumulator holds %d entries, want 1", len(accumulated))
	}
}

func TestResponseInfoBranchStopsAtMaxResponses(t *testing.T) {
	forms := AddFormWrapperGenerator()
	fullResponses := make([]ResponseInfo, MaxCommandResponses)
	for i := range fullResponses {
		fullResponses[i] = ResponseInfo{Type: "message", Content: "c"}
	}
	fullString, _ := ResponseInfoSliceToJSON(fullResponses)
	setFormValue(forms, testIdxResponseInfo, "responses", fullString)
	if got := forms[testIdxResponseInfo].BranchCallback(forms[testIdxResponseInfo].Values, forms); got != testIdxAccept {
		t.Errorf("full responses routed to %d, want %d", got, testIdxAccept)
	}

	oneString, _ := ResponseInfoSliceToJSON([]ResponseInfo{{Type: "message", Content: "c"}})
	setFormValue(forms, testIdxResponseInfo, "responses", oneString)
	if got := forms[testIdxResponseInfo].BranchCallback(forms[testIdxResponseInfo].Values, forms); got != testIdxResponseStart {
		t.Errorf("partial responses routed to %d, want %d", got, testIdxResponseStart)
	}
}

func TestCmdStartCallbackResetsPageAndResponseState(t *testing.T) {
	forms := AddFormWrapperGenerator()
	modelValues := newAddModelValues()
	setModelValue(modelValues, "pages", `[{"Name":"stale"}]`)
	setModelValue(modelValues, "currentPage", `{"Name":"stale"}`)
	setFormValue(forms, testIdxPageInfo, "pageName", "stale")
	setFormValue(forms, testIdxPageInfo, "pageTitle", "stale")
	setFormValue(forms, testIdxPageFieldInfo, "pageFields", "stale")
	setFormValue(forms, testIdxPageNext, "pages", "stale")
	setFormValue(forms, testIdxResponseInfo, "responses", "stale")

	forms[testIdxCmdStart].Callback(forms[testIdxCmdStart].Values, modelValues, forms)

	if *modelValues.Map["pages"] != "[]" {
		t.Errorf("model pages = %q, want empty collection", *modelValues.Map["pages"])
	}
	if *modelValues.Map["currentPage"] != "" {
		t.Errorf("current page = %q, want empty", *modelValues.Map["currentPage"])
	}
	for index, key := range map[int]string{
		testIdxPageInfo:      "pageName",
		testIdxPageFieldInfo: "pageFields",
		testIdxPageNext:      "pages",
		testIdxResponseInfo:  "responses",
	} {
		if *forms[index].Values.Map[key] != "" {
			t.Errorf("form %d key %q should reset, got %q", index, key, *forms[index].Values.Map[key])
		}
	}
}

func TestAcceptCallbackAppendsMultiPageModal(t *testing.T) {
	forms := AddFormWrapperGenerator()
	modelValues := newAddModelValues()
	command := CommandInfo{
		Name:        "wizard",
		Type:        "modal",
		Scope:       "guild",
		Description: "d",
		ReturnType:  "None",
		Pages: []PageInfo{
			{Name: "intro", Title: "Intro", Fields: makeFields(1), Next: ""},
		},
		Responses: []ResponseInfo{{Type: "message", Content: "done"}},
	}
	commandString, _ := command.ToJSON()
	setModelValue(modelValues, "currentCommand", commandString)
	setFormValue(forms, testIdxAccept, "cmdAcceptConfirm", "yes")

	forms[testIdxAccept].Callback(forms[testIdxAccept].Values, modelValues, forms)

	slashCommands, err := JSONToCmdInfoSlice(*modelValues.Map["slashCommands"])
	if err != nil {
		t.Fatalf("failed to parse slash commands: %v", err)
	}
	if len(slashCommands) != 1 {
		t.Fatalf("expected 1 slash command, got %d", len(slashCommands))
	}
	saved := slashCommands[0]
	if len(saved.Pages) != 1 || saved.Pages[0].Name != "intro" {
		t.Errorf("saved pages = %+v, want the collected page", saved.Pages)
	}
	if len(saved.Responses) != 1 || saved.Responses[0].Content != "done" {
		t.Errorf("saved responses = %+v, want the collected response", saved.Responses)
	}
}

func TestValidateAcceptedCommand(t *testing.T) {
	validModal := CommandInfo{
		Name:        "wizard",
		Type:        "modal",
		Scope:       "guild",
		Description: "d",
		ReturnType:  "None",
		Pages: []PageInfo{
			{Name: "intro", Title: "Intro", Fields: makeFields(1), Next: ""},
		},
	}
	brokenGoto := validModal
	brokenGoto.Pages = []PageInfo{
		{Name: "intro", Title: "Intro", Fields: makeFields(1), Branches: []BranchRule{{Field: "fielda", Equals: "x", Goto: "missing"}}},
	}
	validSlash := CommandInfo{Name: "greet", Type: "slash", Scope: "guild", Description: "d", ReturnType: "str"}

	tests := []struct {
		name     string
		command  CommandInfo
		existing string
		wantErr  bool
	}{
		{"valid multi page modal passes", validModal, "[]", false},
		{"branch goto to a missing page fails", brokenGoto, "[]", true},
		{"valid slash command passes", validSlash, "[]", false},
		{"duplicate command name fails", validSlash, `[{"Name":"greet","Type":"slash"}]`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelValues := newAddModelValues()
			commandString, _ := tt.command.ToJSON()
			setModelValue(modelValues, "currentCommand", commandString)
			setModelValue(modelValues, "slashCommands", tt.existing)
			err := validateAcceptedCommand(modelValues)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAcceptedCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// A missing current command is always an error
	modelValues := newAddModelValues()
	if err := validateAcceptedCommand(modelValues); err == nil {
		t.Error("empty current command should fail validation")
	}
}

func TestValidatePageName(t *testing.T) {
	existing := []PageInfo{{Name: "intro"}}
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "details", false},
		{"empty", "", true},
		{"contains space", "my page", true},
		{"duplicate", "intro", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePageName(tt.input, existing)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePageName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateResponseContent(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid content", "Hello there", false},
		{"empty", "", true},
		{"double quote", `say "hi"`, true},
		{"backslash", `a\b`, true},
		{"newline", "a\nb", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateResponseContent(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateResponseContent(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestBuildCommandSummaryShowsPagesAndResponses(t *testing.T) {
	command := CommandInfo{
		Name:        "wizard",
		Type:        "modal",
		Scope:       "guild",
		Description: "d",
		ReturnType:  "None",
		Pages: []PageInfo{
			{Name: "intro", Title: "Intro", Fields: makeFields(2), Branches: []BranchRule{{Field: "fielda", Equals: "x", Goto: "details"}}, Next: "details"},
			{Name: "details", Title: "Details", Fields: makeFields(1)},
		},
		Responses: []ResponseInfo{{Type: "message", Content: "done", Ephemeral: true}},
	}

	summary := buildCommandSummary(command)

	for _, want := range []string{
		"intro (2 fields, 1 branches, next: details)",
		"details (1 fields, 0 branches, next: end)",
		"message: done, ephemeral",
	} {
		if !strings.Contains(summary, want) {
			t.Errorf("summary missing %q:\n%s", want, summary)
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
