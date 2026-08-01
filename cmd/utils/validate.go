/*
Copyright © 2025 Austin "Choice404" Choi
See end of file for extended copyright information
*/

package utils

import (
	"fmt"
	"strings"
	"unicode"
)

// Valid option sets shared by the forms and the headless flag parsing
var (
	validCommandTypes = []string{"slash", "prefix"}
	validCommandScopes = []string{"guild", "global"}
	validReturnTypes = []string{"str", "int", "float", "bool", "None"}
	validArgTypes = []string{"str", "int", "float", "bool", "discord.Member", "discord.Role"}
	validLicenses = []string{"mit", "apache-2.0", "gpl-3.0", "bsd-3-clause", "unlicense", "no-license"}
)

func contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func ValidateBotName(s string) error {
	if s == "" {
		return fmt.Errorf("Bot name cannot be empty")
	}
	if len(s) > 20 {
		return fmt.Errorf("Bot name is too long")
	}
	r := []rune(s)[0]
	if !unicode.IsLetter(r) {
		return fmt.Errorf("Bot name must start with a letter")
	}
	if strings.ContainsRune(s, ' ') || strings.ContainsRune(s, '\t') {
		return fmt.Errorf("Bot name cannot contain whitespace")
	}
	if strings.ContainsAny(s, "!@#$%^&*()_+={}[]|\\:;\"'<>,.?/~`") {
		return fmt.Errorf("Bot name cannot contain special characters")
	}
	return nil
}

func ValidateBotDescription(s string) error {
	if s == "" {
		return fmt.Errorf("Description cannot be empty")
	}
	return nil
}

func ValidateBotAuthor(s string) error {
	if s == "" {
		return fmt.Errorf("Author name cannot be empty")
	}
	return nil
}

func ValidateBotPrefix(s string) error {
	if len(s) > 1 {
		return fmt.Errorf("Command prefix must be a single character")
	}
	if s == "" {
		return nil
	}
	r := []rune(s)[0]
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return fmt.Errorf("Command prefix can not be an alphanumeric character")
	}
	return nil
}

func ValidateEnvChoice(s string) error {
	if s != "env" && s != "doppler" {
		return fmt.Errorf("Please select either 'env' or 'doppler'")
	}
	return nil
}

func ValidateToken(envChoice string, s string) error {
	if envChoice == "env" {
		if s == "" {
			return fmt.Errorf("Token cannot be empty")
		}
		if len(s) < 9 {
			return fmt.Errorf("Token is too short")
		}
	}
	return nil
}

func ValidateLicense(s string) error {
	if s == "" {
		return fmt.Errorf("Please select a license type")
	}
	if !contains(validLicenses, s) {
		return fmt.Errorf("license must be one of %s", strings.Join(validLicenses, ", "))
	}
	return nil
}

func ValidateCommandName(s string, existing []CommandInfo) error {
	if s == "" {
		return fmt.Errorf("command name cannot be empty")
	}
	if strings.Contains(s, " ") {
		return fmt.Errorf("command name cannot contain spaces")
	}
	if commandExists(s, existing) {
		return fmt.Errorf("command name already exists")
	}
	return nil
}

func ValidateCommandType(s string) error {
	if s == "" {
		return fmt.Errorf("command type cannot be empty")
	}
	if !contains(validCommandTypes, s) {
		return fmt.Errorf("command type must be either slash or prefix")
	}
	return nil
}

func ValidateCommandScope(s string) error {
	if s == "" {
		return fmt.Errorf("command scope cannot be empty")
	}
	if !contains(validCommandScopes, s) {
		return fmt.Errorf("command scope must be either guild or global")
	}
	return nil
}

func ValidateCommandDescription(s string) error {
	if s == "" {
		return fmt.Errorf("description cannot be empty")
	}
	return nil
}

func ValidateReturnType(s string) error {
	if s == "" {
		return fmt.Errorf("return type cannot be empty")
	}
	if !contains(validReturnTypes, s) {
		return fmt.Errorf("return type must be one of %s", strings.Join(validReturnTypes, ", "))
	}
	return nil
}

func ValidateArgName(s string, existing []ArgInfo) error {
	if s == "" {
		return fmt.Errorf("Argument name cannot be empty")
	}
	if strings.Contains(s, " ") {
		return fmt.Errorf("Argument name cannot contain spaces")
	}
	if strings.Contains(s, "-") {
		return fmt.Errorf("Argument name cannot contain dashes")
	}
	if argExists(s, existing) {
		return fmt.Errorf("Argument name already exists")
	}
	return nil
}

func ValidateArgDescription(s string) error {
	if s == "" {
		return fmt.Errorf("argument description cannot be empty")
	}
	return nil
}

func ValidateArgType(s string) error {
	if s == "" {
		return fmt.Errorf("argument type cannot be empty")
	}
	if !contains(validArgTypes, s) {
		return fmt.Errorf("argument type must be one of %s", strings.Join(validArgTypes, ", "))
	}
	return nil
}

// ValidateCommand checks a full command definition against the already accepted commands
func ValidateCommand(command CommandInfo, existing []CommandInfo) error {
	if err := ValidateCommandName(command.Name, existing); err != nil {
		return err
	}
	if err := ValidateCommandType(command.Type); err != nil {
		return err
	}
	if err := ValidateCommandScope(command.Scope); err != nil {
		return err
	}
	if err := ValidateCommandDescription(command.Description); err != nil {
		return err
	}
	if err := ValidateReturnType(command.ReturnType); err != nil {
		return err
	}
	for i, arg := range command.Args {
		if err := ValidateArgName(arg.Name, command.Args[:i]); err != nil {
			return fmt.Errorf("argument '%s': %w", arg.Name, err)
		}
		if err := ValidateArgDescription(arg.Description); err != nil {
			return fmt.Errorf("argument '%s': %w", arg.Name, err)
		}
		if err := ValidateArgType(arg.Type); err != nil {
			return fmt.Errorf("argument '%s': %w", arg.Name, err)
		}
	}
	return nil
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
