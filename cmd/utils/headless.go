/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package utils

import (
	"fmt"
	"os"
)

// HeadlessMode is set once at startup when the cli runs without the tui
var HeadlessMode bool

/**
 * RunHeadless
 * Runs a model's callbacks directly without starting the bubbletea program
 * The caller fills ModelValues from flags before calling this
 * @param m {Model} - the model built by the same constructors the tui uses
 * @return []error - every error the callbacks produced
 **/
func RunHeadless(m Model) []error {
	// Load the global config the same way Model.Init does
	conf, err := LoadGlobalConfig()
	if err != nil {
		return []error{fmt.Errorf("failed to load global config: %w", err)}
	}
	globalConfig = *conf

	// Make sure the value maps exist before any callback reads them
	if m.ModelValues.Map == nil {
		m.ModelValues.Map = make(map[string]*string)
	}

	allValueMaps := make([]Values, len(m.forms))
	for i, form := range m.forms {
		if form.Values.Map == nil {
			form.Values.Map = make(map[string]*string)
		}
		allValueMaps[i] = form.Values
	}

	// Run the init callback so prefill and validation behave like the tui path
	if m.initCallback != nil {
		m.initCallback(&m, allValueMaps)
	}
	if m.Error != nil {
		return m.Error
	}

	// Run the main callback that does the actual work
	if m.callback != nil {
		return m.callback(&m)
	}
	return nil
}

/**
 * PrintErrors
 * Prints every error to stderr so stdout stays clean for piping
 * @param errs {[]error} - the errors to print
 * @return bool - true if there was at least one error
 **/
func PrintErrors(errs []error) bool {
	printed := false
	for _, err := range errs {
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			printed = true
		}
	}
	return printed
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
