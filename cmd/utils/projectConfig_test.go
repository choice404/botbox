/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"testing"
)

func TestConvertLegacyCommands(t *testing.T) {
	tests := []struct {
		name        string
		commandNames []string
		commandType string
		wantLen     int
	}{
		{
			name:        "empty slice",
			commandNames: []string{},
			commandType: "slash",
			wantLen:     0,
		},
		{
			name:        "single slash command",
			commandNames: []string{"hello"},
			commandType: "slash",
			wantLen:     1,
		},
		{
			name:        "single prefix command",
			commandNames: []string{"ping"},
			commandType: "prefix",
			wantLen:     1,
		},
		{
			name:        "multiple commands",
			commandNames: []string{"cmd1", "cmd2", "cmd3"},
			commandType: "slash",
			wantLen:     3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertLegacyCommands(tt.commandNames, tt.commandType)

			if len(got) != tt.wantLen {
				t.Errorf("convertLegacyCommands() returned %d commands, want %d", len(got), tt.wantLen)
			}

			for i, cmd := range got {
				if cmd.Name != tt.commandNames[i] {
					t.Errorf("Command %d name = %q, want %q", i, cmd.Name, tt.commandNames[i])
				}
				if cmd.Type != tt.commandType {
					t.Errorf("Command %d type = %q, want %q", i, cmd.Type, tt.commandType)
				}
				if cmd.Scope != "guild" {
					t.Errorf("Command %d scope = %q, want %q", i, cmd.Scope, "guild")
				}
				if cmd.ReturnType != "None" {
					t.Errorf("Command %d return type = %q, want %q", i, cmd.ReturnType, "None")
				}
				if cmd.Description == "" {
					t.Errorf("Command %d should have a description", i)
				}
			}
		})
	}
}

func TestConvertLegacyCommandsDescriptionFormat(t *testing.T) {
	slashCmds := convertLegacyCommands([]string{"test"}, "slash")
	if len(slashCmds) > 0 && slashCmds[0].Description != "Legacy slash command" {
		t.Errorf("Slash command description = %q, want %q", slashCmds[0].Description, "Legacy slash command")
	}

	prefixCmds := convertLegacyCommands([]string{"test"}, "prefix")
	if len(prefixCmds) > 0 && prefixCmds[0].Description != "Legacy prefix command" {
		t.Errorf("Prefix command description = %q, want %q", prefixCmds[0].Description, "Legacy prefix command")
	}
}

func TestIsModernCogFormat(t *testing.T) {
	tests := []struct {
		name string
		cog  CogConfig
		want bool
	}{
		{
			name: "empty cog is modern",
			cog: CogConfig{
				Name:           "EmptyCog",
				SlashCommands:  []CommandInfo{},
				PrefixCommands: []CommandInfo{},
			},
			want: true,
		},
		{
			name: "modern format with slash command description",
			cog: CogConfig{
				Name: "ModernCog",
				SlashCommands: []CommandInfo{
					{Name: "hello", Description: "Says hello"},
				},
				PrefixCommands: []CommandInfo{},
			},
			want: true,
		},
		{
			name: "modern format with prefix command description",
			cog: CogConfig{
				Name:          "ModernCog",
				SlashCommands: []CommandInfo{},
				PrefixCommands: []CommandInfo{
					{Name: "ping", Description: "Pong!"},
				},
			},
			want: true,
		},
		{
			name: "legacy format - slash command without description",
			cog: CogConfig{
				Name: "LegacyCog",
				SlashCommands: []CommandInfo{
					{Name: "hello", Description: ""},
				},
				PrefixCommands: []CommandInfo{},
			},
			want: false,
		},
		{
			name: "legacy format - prefix command without description",
			cog: CogConfig{
				Name:          "LegacyCog",
				SlashCommands: []CommandInfo{},
				PrefixCommands: []CommandInfo{
					{Name: "ping", Description: ""},
				},
			},
			want: false,
		},
		{
			name: "checks first slash command only",
			cog: CogConfig{
				Name: "MixedCog",
				SlashCommands: []CommandInfo{
					{Name: "first", Description: "Has desc"},
					{Name: "second", Description: ""},
				},
				PrefixCommands: []CommandInfo{},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isModernCogFormat(tt.cog)
			if got != tt.want {
				t.Errorf("isModernCogFormat() = %v, want %v", got, tt.want)
			}
		})
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
