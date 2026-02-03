/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		// Equal versions
		{name: "equal simple", v1: "1.0.0", v2: "1.0.0", want: 0},
		{name: "equal two parts", v1: "1.0", v2: "1.0", want: 0},
		{name: "equal single part", v1: "1", v2: "1", want: 0},
		{name: "equal with trailing zero", v1: "1.0.0", v2: "1.0", want: 0},

		// Less than
		{name: "less than major", v1: "1.0.0", v2: "2.0.0", want: -1},
		{name: "less than minor", v1: "1.1.0", v2: "1.2.0", want: -1},
		{name: "less than patch", v1: "1.0.1", v2: "1.0.2", want: -1},
		{name: "less than different lengths", v1: "1.0", v2: "1.0.1", want: -1},

		// Greater than
		{name: "greater than major", v1: "2.0.0", v2: "1.0.0", want: 1},
		{name: "greater than minor", v1: "1.2.0", v2: "1.1.0", want: 1},
		{name: "greater than patch", v1: "1.0.2", v2: "1.0.1", want: 1},
		{name: "greater than different lengths", v1: "1.0.1", v2: "1.0", want: 1},

		// Multi-digit versions
		{name: "multi-digit equal", v1: "10.20.30", v2: "10.20.30", want: 0},
		{name: "multi-digit less", v1: "9.99.99", v2: "10.0.0", want: -1},
		{name: "multi-digit greater", v1: "10.0.0", v2: "9.99.99", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompareVersions(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("CompareVersions(%q, %q) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

func TestCompareVersionsWithPrefix(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		{name: "both with v prefix equal", v1: "v1.0.0", v2: "v1.0.0", want: 0},
		{name: "first with v prefix", v1: "v1.0.0", v2: "1.0.0", want: 0},
		{name: "second with v prefix", v1: "1.0.0", v2: "v1.0.0", want: 0},
		{name: "v prefix less than", v1: "v1.0.0", v2: "v2.0.0", want: -1},
		{name: "v prefix greater than", v1: "v2.0.0", v2: "v1.0.0", want: 1},
		{name: "mixed v prefix comparison", v1: "v1.5.0", v2: "1.4.0", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompareVersions(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("CompareVersions(%q, %q) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

func TestCompareVersionsEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		// Missing parts (padded with 0)
		{name: "short vs long equal", v1: "1", v2: "1.0.0", want: 0},
		{name: "short vs long less", v1: "1", v2: "1.0.1", want: -1},
		{name: "long vs short greater", v1: "1.0.1", v2: "1", want: 1},

		// Empty strings (treated as 0)
		{name: "empty vs empty", v1: "", v2: "", want: 0},
		{name: "empty vs zero", v1: "", v2: "0", want: 0},
		{name: "empty vs version", v1: "", v2: "1.0.0", want: -1},
		{name: "version vs empty", v1: "1.0.0", v2: "", want: 1},

		// Non-numeric parts (Sscanf will parse as 0)
		{name: "pre-release suffix ignored", v1: "1.0.0-alpha", v2: "1.0.0-beta", want: 0},
		{name: "numeric vs pre-release", v1: "1.0.0", v2: "1.0.0-alpha", want: 0},

		// Extra parts
		{name: "four parts equal", v1: "1.0.0.0", v2: "1.0.0.0", want: 0},
		{name: "three vs four parts", v1: "1.0.0", v2: "1.0.0.0", want: 0},
		{name: "three vs four parts different", v1: "1.0.0", v2: "1.0.0.1", want: -1},

		// Only v prefix
		{name: "only v vs empty", v1: "v", v2: "", want: 0},
		{name: "v vs v", v1: "v", v2: "v", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompareVersions(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("CompareVersions(%q, %q) = %d, want %d", tt.v1, tt.v2, got, tt.want)
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
