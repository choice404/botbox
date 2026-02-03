/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"testing"
)

func TestCommandExists(t *testing.T) {
	tests := []struct {
		name        string
		commandName string
		commandList []CommandInfo
		want        bool
	}{
		{
			name:        "found in list",
			commandName: "hello",
			commandList: []CommandInfo{
				{Name: "hello", Type: "slash"},
				{Name: "ping", Type: "prefix"},
			},
			want: true,
		},
		{
			name:        "not found in list",
			commandName: "goodbye",
			commandList: []CommandInfo{
				{Name: "hello", Type: "slash"},
				{Name: "ping", Type: "prefix"},
			},
			want: false,
		},
		{
			name:        "empty list",
			commandName: "hello",
			commandList: []CommandInfo{},
			want:        false,
		},
		{
			name:        "single item found",
			commandName: "test",
			commandList: []CommandInfo{
				{Name: "test", Type: "slash"},
			},
			want: true,
		},
		{
			name:        "single item not found",
			commandName: "other",
			commandList: []CommandInfo{
				{Name: "test", Type: "slash"},
			},
			want: false,
		},
		{
			name:        "case sensitive not found",
			commandName: "Hello",
			commandList: []CommandInfo{
				{Name: "hello", Type: "slash"},
			},
			want: false,
		},
		{
			name:        "empty command name",
			commandName: "",
			commandList: []CommandInfo{
				{Name: "hello", Type: "slash"},
			},
			want: false,
		},
		{
			name:        "empty command name in list",
			commandName: "",
			commandList: []CommandInfo{
				{Name: "", Type: "slash"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := commandExists(tt.commandName, tt.commandList)
			if got != tt.want {
				t.Errorf("commandExists(%q, ...) = %v, want %v", tt.commandName, got, tt.want)
			}
		})
	}
}

func TestArgExists(t *testing.T) {
	tests := []struct {
		name    string
		argName string
		args    []ArgInfo
		want    bool
	}{
		{
			name:    "found in list",
			argName: "user",
			args: []ArgInfo{
				{Name: "user", Type: "str", Description: "The user"},
				{Name: "count", Type: "int", Description: "The count"},
			},
			want: true,
		},
		{
			name:    "not found in list",
			argName: "missing",
			args: []ArgInfo{
				{Name: "user", Type: "str"},
				{Name: "count", Type: "int"},
			},
			want: false,
		},
		{
			name:    "empty list",
			argName: "user",
			args:    []ArgInfo{},
			want:    false,
		},
		{
			name:    "single item found",
			argName: "arg1",
			args: []ArgInfo{
				{Name: "arg1", Type: "str"},
			},
			want: true,
		},
		{
			name:    "single item not found",
			argName: "other",
			args: []ArgInfo{
				{Name: "arg1", Type: "str"},
			},
			want: false,
		},
		{
			name:    "case sensitive not found",
			argName: "User",
			args: []ArgInfo{
				{Name: "user", Type: "str"},
			},
			want: false,
		},
		{
			name:    "empty arg name",
			argName: "",
			args: []ArgInfo{
				{Name: "user", Type: "str"},
			},
			want: false,
		},
		{
			name:    "empty arg name in list",
			argName: "",
			args: []ArgInfo{
				{Name: "", Type: "str"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := argExists(tt.argName, tt.args)
			if got != tt.want {
				t.Errorf("argExists(%q, ...) = %v, want %v", tt.argName, got, tt.want)
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
