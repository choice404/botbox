/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"testing"
)

func TestCreateDefaultConfig(t *testing.T) {
	// Set a known version for testing
	oldVersion := Version
	Version = "2.5.3"
	defer func() { Version = oldVersion }()

	config := createDefaultConfig()

	t.Run("CLI defaults", func(t *testing.T) {
		if config.CLI.Version != "2.5.3" {
			t.Errorf("CLI.Version = %q, want %q", config.CLI.Version, "2.5.3")
		}
		if config.CLI.CheckUpdates != true {
			t.Errorf("CLI.CheckUpdates = %v, want %v", config.CLI.CheckUpdates, true)
		}
		if config.CLI.AutoUpdate != false {
			t.Errorf("CLI.AutoUpdate = %v, want %v", config.CLI.AutoUpdate, false)
		}
	})

	t.Run("User defaults", func(t *testing.T) {
		if config.User.DefaultUser != "" {
			t.Errorf("User.DefaultUser = %q, want empty string", config.User.DefaultUser)
		}
		if config.User.GithubUsername != "" {
			t.Errorf("User.GithubUsername = %q, want empty string", config.User.GithubUsername)
		}
	})

	t.Run("Display defaults", func(t *testing.T) {
		if config.Display.ScrollEnabled != true {
			t.Errorf("Display.ScrollEnabled = %v, want %v", config.Display.ScrollEnabled, true)
		}
		if config.Display.ColorScheme != "default" {
			t.Errorf("Display.ColorScheme = %q, want %q", config.Display.ColorScheme, "default")
		}
	})

	t.Run("Defaults defaults", func(t *testing.T) {
		if config.Defaults.CommandPrefix != "!" {
			t.Errorf("Defaults.CommandPrefix = %q, want %q", config.Defaults.CommandPrefix, "!")
		}
		if config.Defaults.PythonVersion != "3.11" {
			t.Errorf("Defaults.PythonVersion = %q, want %q", config.Defaults.PythonVersion, "3.11")
		}
		if config.Defaults.AutoGitInit != true {
			t.Errorf("Defaults.AutoGitInit = %v, want %v", config.Defaults.AutoGitInit, true)
		}
	})

	t.Run("Dev defaults", func(t *testing.T) {
		if config.Dev.Editor != "code" {
			t.Errorf("Dev.Editor = %q, want %q", config.Dev.Editor, "code")
		}
	})
}

func TestCreateDefaultConfigVersionPrefix(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		wantVersion string
	}{
		{
			name:        "version without prefix",
			version:     "2.5.3",
			wantVersion: "2.5.3",
		},
		{
			name:        "version with v prefix",
			version:     "v2.5.3",
			wantVersion: "2.5.3",
		},
		{
			name:        "empty version",
			version:     "",
			wantVersion: "",
		},
		{
			name:        "only v",
			version:     "v",
			wantVersion: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldVersion := Version
			Version = tt.version
			defer func() { Version = oldVersion }()

			config := createDefaultConfig()
			if config.CLI.Version != tt.wantVersion {
				t.Errorf("createDefaultConfig().CLI.Version = %q, want %q", config.CLI.Version, tt.wantVersion)
			}
		})
	}
}

func TestGlobalConfigStructure(t *testing.T) {
	// Test that GlobalConfig has the expected structure
	config := GlobalConfig{
		CLI: GlobalCLI{
			Version:      "1.0.0",
			CheckUpdates: true,
			AutoUpdate:   false,
		},
		User: GlobalUser{
			DefaultUser:    "testuser",
			GithubUsername: "testgh",
		},
		Display: GlobalDisplay{
			ScrollEnabled: true,
			ColorScheme:   "dark",
		},
		Defaults: GlobalDefaults{
			CommandPrefix: "!",
			PythonVersion: "3.10",
			AutoGitInit:   true,
		},
		Dev: GlobalDev{
			Editor: "vim",
		},
	}

	if config.CLI.Version != "1.0.0" {
		t.Errorf("Expected CLI.Version to be '1.0.0', got %q", config.CLI.Version)
	}
	if config.User.DefaultUser != "testuser" {
		t.Errorf("Expected User.DefaultUser to be 'testuser', got %q", config.User.DefaultUser)
	}
	if config.Display.ColorScheme != "dark" {
		t.Errorf("Expected Display.ColorScheme to be 'dark', got %q", config.Display.ColorScheme)
	}
	if config.Defaults.PythonVersion != "3.10" {
		t.Errorf("Expected Defaults.PythonVersion to be '3.10', got %q", config.Defaults.PythonVersion)
	}
	if config.Dev.Editor != "vim" {
		t.Errorf("Expected Dev.Editor to be 'vim', got %q", config.Dev.Editor)
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
