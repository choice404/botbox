/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidateBotName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "MyBot", false},
		{"empty", "", true},
		{"too long", "ThisBotNameIsWayTooLongOk", true},
		{"starts with digit", "1bot", true},
		{"contains space", "my bot", true},
		{"contains special char", "my_bot", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBotName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBotName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateBotPrefix(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"bang", "!", false},
		{"question", "?", false},
		{"empty allowed", "", false},
		{"letter", "a", true},
		{"digit", "1", true},
		{"multiple chars", "!!", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBotPrefix(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBotPrefix(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEnvChoice(t *testing.T) {
	if err := ValidateEnvChoice("env"); err != nil {
		t.Errorf("env should be valid, got %v", err)
	}
	if err := ValidateEnvChoice("doppler"); err != nil {
		t.Errorf("doppler should be valid, got %v", err)
	}
	if err := ValidateEnvChoice("other"); err == nil {
		t.Error("other should be invalid")
	}
	if err := ValidateEnvChoice(""); err == nil {
		t.Error("empty should be invalid")
	}
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name      string
		envChoice string
		token     string
		wantErr   bool
	}{
		{"valid env token", "env", "1234567890abc", false},
		{"empty env token", "env", "", true},
		{"short env token", "env", "12345", true},
		{"doppler skips check", "doppler", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateToken(tt.envChoice, tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken(%q, %q) error = %v, wantErr %v", tt.envChoice, tt.token, err, tt.wantErr)
			}
		})
	}
}

func TestValidateLicense(t *testing.T) {
	for _, valid := range []string{"mit", "apache-2.0", "gpl-3.0", "bsd-3-clause", "unlicense", "no-license"} {
		if err := ValidateLicense(valid); err != nil {
			t.Errorf("%q should be valid, got %v", valid, err)
		}
	}
	if err := ValidateLicense(""); err == nil {
		t.Error("empty should be invalid")
	}
	if err := ValidateLicense("wtfpl"); err == nil {
		t.Error("unknown license should be invalid")
	}
}

func TestValidateHelpStyle(t *testing.T) {
	for _, valid := range []string{"compact", "detailed"} {
		if err := ValidateHelpStyle(valid); err != nil {
			t.Errorf("%q should be valid, got %v", valid, err)
		}
	}
	if err := ValidateHelpStyle(""); err == nil {
		t.Error("empty should be invalid")
	}
	if err := ValidateHelpStyle("verbose"); err == nil {
		t.Error("unknown help style should be invalid")
	}
	if err := ValidateHelpStyle("Compact"); err == nil {
		t.Error("help style should be case sensitive")
	}
}

func TestNormalizeHelpStyle(t *testing.T) {
	if got := NormalizeHelpStyle(""); got != DefaultHelpStyle {
		t.Errorf("empty help style should normalize to %q, got %q", DefaultHelpStyle, got)
	}
	if got := NormalizeHelpStyle("detailed"); got != "detailed" {
		t.Errorf("set help style should be kept, got %q", got)
	}
}

func TestValidateCommand(t *testing.T) {
	valid := CommandInfo{
		Name:        "greet",
		Scope:       "guild",
		Type:        "slash",
		Description: "Greets a user",
		Args: []ArgInfo{
			{Name: "user", Type: "discord.Member", Description: "who to greet"},
		},
		ReturnType: "None",
	}

	if err := ValidateCommand(valid, nil); err != nil {
		t.Errorf("valid command should pass, got %v", err)
	}

	dupe := valid
	if err := ValidateCommand(dupe, []CommandInfo{valid}); err == nil {
		t.Error("duplicate command name should fail")
	}

	badType := valid
	badType.Name = "other"
	badType.Type = "hybrid"
	if err := ValidateCommand(badType, nil); err == nil {
		t.Error("unknown command type should fail")
	}

	badArg := valid
	badArg.Name = "another"
	badArg.Args = []ArgInfo{{Name: "bad-arg", Type: "str", Description: "has a dash"}}
	if err := ValidateCommand(badArg, nil); err == nil {
		t.Error("arg name with dash should fail")
	}

	dupeArgs := valid
	dupeArgs.Name = "yetanother"
	dupeArgs.Args = []ArgInfo{
		{Name: "user", Type: "str", Description: "first"},
		{Name: "user", Type: "str", Description: "second"},
	}
	if err := ValidateCommand(dupeArgs, nil); err == nil {
		t.Error("duplicate arg names should fail")
	}
}

func TestValidateFieldName(t *testing.T) {
	existing := []FieldInfo{{Name: "summary", Label: "Summary", Style: "short"}}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "details", false},
		{"empty", "", true},
		{"contains space", "field name", true},
		{"contains dash", "field-name", true},
		{"duplicate", "summary", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFieldName(tt.input, existing)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFieldName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateFieldLabel(t *testing.T) {
	if err := ValidateFieldLabel("Feedback summary"); err != nil {
		t.Errorf("valid label should pass, got %v", err)
	}
	if err := ValidateFieldLabel(""); err == nil {
		t.Error("empty label should fail")
	}
	if err := ValidateFieldLabel(strings.Repeat("a", 46)); err == nil {
		t.Error("label over 45 characters should fail")
	}
	if err := ValidateFieldLabel(strings.Repeat("a", 45)); err != nil {
		t.Errorf("label of exactly 45 characters should pass, got %v", err)
	}
}

func TestValidateFieldStyle(t *testing.T) {
	for _, valid := range []string{"short", "paragraph"} {
		if err := ValidateFieldStyle(valid); err != nil {
			t.Errorf("%q should be valid, got %v", valid, err)
		}
	}
	if err := ValidateFieldStyle(""); err == nil {
		t.Error("empty style should be invalid")
	}
	if err := ValidateFieldStyle("long"); err == nil {
		t.Error("unknown style should be invalid")
	}
}

func TestValidateCommandModal(t *testing.T) {
	valid := CommandInfo{
		Name:        "feedback",
		Scope:       "guild",
		Type:        "modal",
		Description: "Collects feedback",
		ReturnType:  "None",
		Fields: []FieldInfo{
			{Name: "summary", Label: "Summary", Style: "short", Required: true},
			{Name: "details", Label: "Details", Style: "paragraph"},
		},
	}

	if err := ValidateCommand(valid, nil); err != nil {
		t.Errorf("valid modal command should pass, got %v", err)
	}

	noFields := valid
	noFields.Fields = nil
	if err := ValidateCommand(noFields, nil); err == nil {
		t.Error("modal command without fields should fail")
	}

	tooMany := valid
	tooMany.Fields = nil
	for i := 0; i < MaxModalFields+1; i++ {
		tooMany.Fields = append(tooMany.Fields, FieldInfo{
			Name:  fmt.Sprintf("field%d", i),
			Label: fmt.Sprintf("Field %d", i),
			Style: "short",
		})
	}
	if err := ValidateCommand(tooMany, nil); err == nil {
		t.Error("modal command with six fields should fail")
	}

	withArgs := valid
	withArgs.Args = []ArgInfo{{Name: "user", Type: "str", Description: "who"}}
	if err := ValidateCommand(withArgs, nil); err == nil {
		t.Error("modal command with arguments should fail")
	}

	dupeFields := valid
	dupeFields.Fields = []FieldInfo{
		{Name: "summary", Label: "First", Style: "short"},
		{Name: "summary", Label: "Second", Style: "short"},
	}
	if err := ValidateCommand(dupeFields, nil); err == nil {
		t.Error("duplicate field names should fail")
	}

	slashWithFields := valid
	slashWithFields.Type = "slash"
	if err := ValidateCommand(slashWithFields, nil); err == nil {
		t.Error("non modal command with fields should fail")
	}
}

func TestValidateCommandPages(t *testing.T) {
	valid := CommandInfo{
		Name:        "survey",
		Scope:       "guild",
		Type:        "modal",
		Description: "Runs a survey",
		ReturnType:  "None",
		Pages: []PageInfo{
			{
				Name:     "start",
				Title:    "Start",
				Fields:   []FieldInfo{{Name: "track", Label: "Track", Style: "short", Required: true}},
				Branches: []BranchRule{{Field: "track", Equals: "backend", Goto: "wrap"}},
				Next:     "wrap",
			},
			{
				Name:   "wrap",
				Title:  "Wrap up",
				Fields: []FieldInfo{{Name: "comments", Label: "Comments", Style: "paragraph"}},
			},
		},
	}

	if err := ValidateCommand(valid, nil); err != nil {
		t.Errorf("valid multi page command should pass, got %v", err)
	}

	bothShapes := valid
	bothShapes.Fields = []FieldInfo{{Name: "summary", Label: "Summary", Style: "short"}}
	if err := ValidateCommand(bothShapes, nil); err == nil {
		t.Error("modal command with both fields and pages should fail")
	}

	neitherShape := valid
	neitherShape.Pages = nil
	if err := ValidateCommand(neitherShape, nil); err == nil {
		t.Error("modal command without fields or pages should fail")
	}

	slashWithPages := valid
	slashWithPages.Type = "slash"
	if err := ValidateCommand(slashWithPages, nil); err == nil {
		t.Error("non modal command with pages should fail")
	}

	tooManyPages := valid
	tooManyPages.Pages = nil
	for i := 0; i < MaxFlowPages+1; i++ {
		tooManyPages.Pages = append(tooManyPages.Pages, PageInfo{
			Name:   fmt.Sprintf("page%d", i),
			Title:  fmt.Sprintf("Page %d", i),
			Fields: []FieldInfo{{Name: "field", Label: "Field", Style: "short"}},
		})
	}
	if err := ValidateCommand(tooManyPages, nil); err == nil {
		t.Error("eleven pages should fail")
	}

	dupePageNames := valid
	dupePageNames.Pages = []PageInfo{
		{Name: "start", Title: "First", Fields: []FieldInfo{{Name: "a", Label: "A", Style: "short"}}},
		{Name: "start", Title: "Second", Fields: []FieldInfo{{Name: "b", Label: "B", Style: "short"}}},
	}
	if err := ValidateCommand(dupePageNames, nil); err == nil {
		t.Error("duplicate page names should fail")
	}

	emptyPageName := valid
	emptyPageName.Pages = []PageInfo{
		{Name: "", Title: "First", Fields: []FieldInfo{{Name: "a", Label: "A", Style: "short"}}},
	}
	if err := ValidateCommand(emptyPageName, nil); err == nil {
		t.Error("empty page name should fail")
	}

	sixFieldPage := valid
	sixFieldPage.Pages = []PageInfo{{Name: "start", Title: "Start"}}
	for i := 0; i < MaxModalFields+1; i++ {
		sixFieldPage.Pages[0].Fields = append(sixFieldPage.Pages[0].Fields, FieldInfo{
			Name:  fmt.Sprintf("field%d", i),
			Label: fmt.Sprintf("Field %d", i),
			Style: "short",
		})
	}
	if err := ValidateCommand(sixFieldPage, nil); err == nil {
		t.Error("page with six fields should fail")
	}

	fieldlessPage := valid
	fieldlessPage.Pages = []PageInfo{{Name: "start", Title: "Start"}}
	if err := ValidateCommand(fieldlessPage, nil); err == nil {
		t.Error("page without fields should fail")
	}

	badBranchField := valid
	badBranchField.Pages = []PageInfo{
		{
			Name:     "start",
			Title:    "Start",
			Fields:   []FieldInfo{{Name: "track", Label: "Track", Style: "short"}},
			Branches: []BranchRule{{Field: "missing", Equals: "x", Goto: "start"}},
		},
	}
	if err := ValidateCommand(badBranchField, nil); err == nil {
		t.Error("branch field that is not on the page should fail")
	}

	badGoto := valid
	badGoto.Pages = []PageInfo{
		{
			Name:     "start",
			Title:    "Start",
			Fields:   []FieldInfo{{Name: "track", Label: "Track", Style: "short"}},
			Branches: []BranchRule{{Field: "track", Equals: "x", Goto: "nowhere"}},
		},
	}
	if err := ValidateCommand(badGoto, nil); err == nil {
		t.Error("branch goto pointing at a missing page should fail")
	}

	badNext := valid
	badNext.Pages = []PageInfo{
		{
			Name:   "start",
			Title:  "Start",
			Fields: []FieldInfo{{Name: "track", Label: "Track", Style: "short"}},
			Next:   "nowhere",
		},
	}
	if err := ValidateCommand(badNext, nil); err == nil {
		t.Error("next pointing at a missing page should fail")
	}
}

func TestValidateResponses(t *testing.T) {
	valid := []ResponseInfo{{Type: "message", Content: "done", Ephemeral: true}}
	if err := ValidateResponses(valid); err != nil {
		t.Errorf("valid responses should pass, got %v", err)
	}

	if err := ValidateResponses(nil); err != nil {
		t.Errorf("no responses should pass, got %v", err)
	}

	badType := []ResponseInfo{{Type: "embed", Content: "done"}}
	if err := ValidateResponses(badType); err == nil {
		t.Error("unknown response type should fail")
	}

	emptyContent := []ResponseInfo{{Type: "message", Content: ""}}
	if err := ValidateResponses(emptyContent); err == nil {
		t.Error("empty content should fail")
	}

	quotedContent := []ResponseInfo{{Type: "message", Content: `say "hi"`}}
	if err := ValidateResponses(quotedContent); err == nil {
		t.Error("content with double quotes should fail")
	}

	var tooMany []ResponseInfo
	for i := 0; i < MaxCommandResponses+1; i++ {
		tooMany = append(tooMany, ResponseInfo{Type: "message", Content: fmt.Sprintf("reply %d", i)})
	}
	if err := ValidateResponses(tooMany); err == nil {
		t.Error("four responses should fail")
	}
}

func TestValidateCommandResponses(t *testing.T) {
	valid := CommandInfo{
		Name:        "greet",
		Scope:       "global",
		Type:        "slash",
		Description: "Greets a member",
		ReturnType:  "None",
		Responses:   []ResponseInfo{{Type: "message", Content: "Hello {member}", Ephemeral: false}},
	}

	if err := ValidateCommand(valid, nil); err != nil {
		t.Errorf("slash command with responses should pass, got %v", err)
	}

	badResponse := valid
	badResponse.Responses = []ResponseInfo{{Type: "message", Content: ""}}
	if err := ValidateCommand(badResponse, nil); err == nil {
		t.Error("command with an empty response content should fail")
	}
}

func TestValidateCommandScopeAndReturnType(t *testing.T) {
	if err := ValidateCommandScope("guild"); err != nil {
		t.Errorf("guild should be valid, got %v", err)
	}
	if err := ValidateCommandScope("everywhere"); err == nil {
		t.Error("unknown scope should be invalid")
	}
	if err := ValidateReturnType("None"); err != nil {
		t.Errorf("None should be valid, got %v", err)
	}
	if err := ValidateReturnType("list"); err == nil {
		t.Error("unknown return type should be invalid")
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
