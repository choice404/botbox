/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCogConfigSliceToJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   []CogConfig
		wantErr bool
	}{
		{
			name:    "empty slice",
			input:   []CogConfig{},
			wantErr: false,
		},
		{
			name: "single cog",
			input: []CogConfig{
				{
					Name: "TestCog",
					Env:  "development",
					File: "testCog",
					SlashCommands: []CommandInfo{
						{Name: "test", Type: "slash", Scope: "guild"},
					},
					PrefixCommands: []CommandInfo{},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple cogs",
			input: []CogConfig{
				{Name: "Cog1", Env: "production", File: "cog1"},
				{Name: "Cog2", Env: "development", File: "cog2"},
				{Name: "Cog3", Env: "production", File: "cog3"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CogConfigSliceToJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CogConfigSliceToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("CogConfigSliceToJSON() returned empty string for valid input")
			}
			// Verify it's valid JSON
			if !tt.wantErr {
				var result []CogConfig
				if err := json.Unmarshal([]byte(got), &result); err != nil {
					t.Errorf("CogConfigSliceToJSON() produced invalid JSON: %v", err)
				}
			}
		})
	}
}

func TestJSONToCogConfigSlice(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []CogConfig
		wantErr bool
	}{
		{
			name:    "valid JSON empty array",
			input:   "[]",
			want:    []CogConfig{},
			wantErr: false,
		},
		{
			name:  "valid JSON single cog",
			input: `[{"name":"TestCog","env":"development","file":"testCog","slash_commands":[],"prefix_commands":[]}]`,
			want: []CogConfig{
				{Name: "TestCog", Env: "development", File: "testCog", SlashCommands: []CommandInfo{}, PrefixCommands: []CommandInfo{}},
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{invalid json}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "null JSON",
			input:   "null",
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONToCogConfigSlice(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONToCogConfigSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONToCogConfigSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCogConfigRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []CogConfig
	}{
		{
			name:  "empty slice",
			input: []CogConfig{},
		},
		{
			name: "single cog with commands",
			input: []CogConfig{
				{
					Name: "TestCog",
					Env:  "development",
					File: "testCog",
					SlashCommands: []CommandInfo{
						{
							Name:        "hello",
							Type:        "slash",
							Scope:       "guild",
							Description: "Says hello",
							Args: []ArgInfo{
								{Name: "user", Type: "str", Description: "User to greet"},
							},
							ReturnType: "None",
						},
					},
					PrefixCommands: []CommandInfo{
						{Name: "ping", Type: "prefix", Description: "Pong!"},
					},
				},
			},
		},
		{
			name: "multiple cogs",
			input: []CogConfig{
				{Name: "Cog1", Env: "production", File: "cog1", SlashCommands: []CommandInfo{}, PrefixCommands: []CommandInfo{}},
				{Name: "Cog2", Env: "development", File: "cog2", SlashCommands: []CommandInfo{}, PrefixCommands: []CommandInfo{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize
			jsonStr, err := CogConfigSliceToJSON(tt.input)
			if err != nil {
				t.Fatalf("CogConfigSliceToJSON() error = %v", err)
			}

			// Deserialize
			got, err := JSONToCogConfigSlice(jsonStr)
			if err != nil {
				t.Fatalf("JSONToCogConfigSlice() error = %v", err)
			}

			// Compare
			if !reflect.DeepEqual(got, tt.input) {
				t.Errorf("Round trip failed: got %v, want %v", got, tt.input)
			}
		})
	}
}

func TestJSONToCogConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *CogConfig
		wantErr bool
	}{
		{
			name:  "valid cog config",
			input: `{"name":"TestCog","env":"development","file":"testCog","slash_commands":[],"prefix_commands":[]}`,
			want: &CogConfig{
				Name:           "TestCog",
				Env:            "development",
				File:           "testCog",
				SlashCommands:  []CommandInfo{},
				PrefixCommands: []CommandInfo{},
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{not valid}`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONToCogConfig(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONToCogConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONToCogConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCogConfig_ToJSON(t *testing.T) {
	tests := []struct {
		name    string
		cog     CogConfig
		wantErr bool
	}{
		{
			name: "basic cog",
			cog: CogConfig{
				Name:           "TestCog",
				Env:            "development",
				File:           "testCog",
				SlashCommands:  []CommandInfo{},
				PrefixCommands: []CommandInfo{},
			},
			wantErr: false,
		},
		{
			name: "cog with commands",
			cog: CogConfig{
				Name: "CmdCog",
				Env:  "production",
				File: "cmdCog",
				SlashCommands: []CommandInfo{
					{Name: "cmd1", Type: "slash"},
				},
				PrefixCommands: []CommandInfo{
					{Name: "cmd2", Type: "prefix"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cog.ToJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("CogConfig.ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify it's valid JSON by parsing it back
				var result CogConfig
				if err := json.Unmarshal([]byte(got), &result); err != nil {
					t.Errorf("CogConfig.ToJSON() produced invalid JSON: %v", err)
				}
			}
		})
	}
}

func TestCommandInfoJSON(t *testing.T) {
	t.Run("CmdInfoSliceToJSON", func(t *testing.T) {
		input := []CommandInfo{
			{Name: "cmd1", Type: "slash", Scope: "guild", Description: "Command 1"},
			{Name: "cmd2", Type: "prefix", Description: "Command 2"},
		}
		got, err := CmdInfoSliceToJSON(input)
		if err != nil {
			t.Fatalf("CmdInfoSliceToJSON() error = %v", err)
		}
		// Verify valid JSON
		var result []CommandInfo
		if err := json.Unmarshal([]byte(got), &result); err != nil {
			t.Errorf("CmdInfoSliceToJSON() produced invalid JSON: %v", err)
		}
	})

	t.Run("CmdInfoSliceToJSON empty", func(t *testing.T) {
		input := []CommandInfo{}
		got, err := CmdInfoSliceToJSON(input)
		if err != nil {
			t.Fatalf("CmdInfoSliceToJSON() error = %v", err)
		}
		if got != "[]" {
			t.Errorf("CmdInfoSliceToJSON() = %v, want []", got)
		}
	})

	t.Run("JSONToCmdInfoSlice", func(t *testing.T) {
		input := `[{"Name":"test","Type":"slash","Scope":"guild","Description":"Test cmd"}]`
		got, err := JSONToCmdInfoSlice(input)
		if err != nil {
			t.Fatalf("JSONToCmdInfoSlice() error = %v", err)
		}
		if len(got) != 1 || got[0].Name != "test" {
			t.Errorf("JSONToCmdInfoSlice() = %v, want slice with test command", got)
		}
	})

	t.Run("JSONToCmdInfoSlice invalid", func(t *testing.T) {
		_, err := JSONToCmdInfoSlice(`{invalid}`)
		if err == nil {
			t.Error("JSONToCmdInfoSlice() expected error for invalid JSON")
		}
	})

	t.Run("JSONToCmdInfo", func(t *testing.T) {
		input := `{"Name":"test","Type":"slash","Scope":"guild","Description":"Test cmd"}`
		got, err := JSONToCmdInfo(input)
		if err != nil {
			t.Fatalf("JSONToCmdInfo() error = %v", err)
		}
		if got.Name != "test" || got.Type != "slash" {
			t.Errorf("JSONToCmdInfo() = %v, want command with name=test, type=slash", got)
		}
	})

	t.Run("JSONToCmdInfo invalid", func(t *testing.T) {
		_, err := JSONToCmdInfo(`{invalid}`)
		if err == nil {
			t.Error("JSONToCmdInfo() expected error for invalid JSON")
		}
	})

	t.Run("CommandInfo.ToJSON", func(t *testing.T) {
		cmd := CommandInfo{
			Name:        "test",
			Type:        "slash",
			Scope:       "guild",
			Description: "Test command",
			Args:        []ArgInfo{{Name: "arg1", Type: "str"}},
			ReturnType:  "None",
		}
		got, err := cmd.ToJSON()
		if err != nil {
			t.Fatalf("CommandInfo.ToJSON() error = %v", err)
		}
		// Verify round trip
		parsed, err := JSONToCmdInfo(got)
		if err != nil {
			t.Fatalf("Failed to parse ToJSON output: %v", err)
		}
		if parsed.Name != cmd.Name || parsed.Type != cmd.Type {
			t.Errorf("Round trip failed: got %v, want %v", parsed, cmd)
		}
	})
}

func TestArgInfoJSON(t *testing.T) {
	t.Run("ArgInfoSliceToJSON", func(t *testing.T) {
		input := []ArgInfo{
			{Name: "arg1", Type: "str", Description: "First argument"},
			{Name: "arg2", Type: "int", Description: "Second argument"},
		}
		got, err := ArgInfoSliceToJSON(input)
		if err != nil {
			t.Fatalf("ArgInfoSliceToJSON() error = %v", err)
		}
		// Verify valid JSON
		var result []ArgInfo
		if err := json.Unmarshal([]byte(got), &result); err != nil {
			t.Errorf("ArgInfoSliceToJSON() produced invalid JSON: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("ArgInfoSliceToJSON() result has %d items, want 2", len(result))
		}
	})

	t.Run("ArgInfoSliceToJSON empty", func(t *testing.T) {
		input := []ArgInfo{}
		got, err := ArgInfoSliceToJSON(input)
		if err != nil {
			t.Fatalf("ArgInfoSliceToJSON() error = %v", err)
		}
		if got != "[]" {
			t.Errorf("ArgInfoSliceToJSON() = %v, want []", got)
		}
	})

	t.Run("JSONToArgInfoSlice", func(t *testing.T) {
		input := `[{"Name":"user","Type":"str","Description":"User name"}]`
		got, err := JSONToArgInfoSlice(input)
		if err != nil {
			t.Fatalf("JSONToArgInfoSlice() error = %v", err)
		}
		if len(got) != 1 || got[0].Name != "user" {
			t.Errorf("JSONToArgInfoSlice() = %v, want slice with user arg", got)
		}
	})

	t.Run("JSONToArgInfoSlice invalid", func(t *testing.T) {
		_, err := JSONToArgInfoSlice(`{invalid}`)
		if err == nil {
			t.Error("JSONToArgInfoSlice() expected error for invalid JSON")
		}
	})

	t.Run("ArgInfo round trip", func(t *testing.T) {
		input := []ArgInfo{
			{Name: "arg1", Type: "str", Description: "Desc 1"},
			{Name: "arg2", Type: "int", Description: "Desc 2"},
		}
		jsonStr, err := ArgInfoSliceToJSON(input)
		if err != nil {
			t.Fatalf("ArgInfoSliceToJSON() error = %v", err)
		}
		got, err := JSONToArgInfoSlice(jsonStr)
		if err != nil {
			t.Fatalf("JSONToArgInfoSlice() error = %v", err)
		}
		if !reflect.DeepEqual(got, input) {
			t.Errorf("Round trip failed: got %v, want %v", got, input)
		}
	})
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
