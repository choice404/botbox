/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCommandsInputInline(t *testing.T) {
	commands, err := parseCommandsInput(`[{"Name":"greet","Scope":"guild","Type":"slash","Description":"Greets","Args":[],"ReturnType":"None"}]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(commands))
	}
	if commands[0].Name != "greet" || commands[0].Type != "slash" {
		t.Errorf("parsed command has wrong fields: %+v", commands[0])
	}
}

func TestParseCommandsInputEmpty(t *testing.T) {
	commands, err := parseCommandsInput("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if commands != nil {
		t.Errorf("expected nil commands for empty input, got %v", commands)
	}
}

func TestParseCommandsInputBadJSON(t *testing.T) {
	if _, err := parseCommandsInput("not json"); err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseCommandsInputFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "commands.json")
	content := `[{"Name":"wave","Scope":"guild","Type":"prefix","Description":"Waves","Args":[],"ReturnType":"None"}]`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	commands, err := parseCommandsInput("@" + path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commands) != 1 || commands[0].Name != "wave" {
		t.Errorf("parsed file command wrong: %+v", commands)
	}
}

func TestParseCommandsInputMissingFile(t *testing.T) {
	if _, err := parseCommandsInput("@/nonexistent/commands.json"); err == nil {
		t.Error("expected error for missing file")
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
