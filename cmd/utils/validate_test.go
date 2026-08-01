/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package utils

import "testing"

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
