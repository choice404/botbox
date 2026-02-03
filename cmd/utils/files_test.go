/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"testing"
)

func TestValidateFileName(t *testing.T) {
	// Note: validateFileName calls fileExists which depends on FindBotConf
	// We test the validation logic that doesn't require filesystem access
	// by checking for invalid characters and empty names

	tests := []struct {
		name     string
		fileName string
		wantErr  bool
	}{
		{
			name:     "empty filename",
			fileName: "",
			wantErr:  true,
		},
		{
			name:     "filename with space",
			fileName: "my file",
			wantErr:  true,
		},
		{
			name:     "filename with dot",
			fileName: "my.file",
			wantErr:  true,
		},
		{
			name:     "filename with forward slash",
			fileName: "my/file",
			wantErr:  true,
		},
		{
			name:     "filename with backslash",
			fileName: "my\\file",
			wantErr:  true,
		},
		{
			name:     "filename with dash",
			fileName: "my-file",
			wantErr:  true,
		},
		{
			name:     "filename with colon",
			fileName: "my:file",
			wantErr:  true,
		},
		{
			name:     "filename with asterisk",
			fileName: "my*file",
			wantErr:  true,
		},
		{
			name:     "filename with question mark",
			fileName: "my?file",
			wantErr:  true,
		},
		{
			name:     "filename with quote",
			fileName: `my"file`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFileName(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFileName(%q) error = %v, wantErr %v", tt.fileName, err, tt.wantErr)
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
