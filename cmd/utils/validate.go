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
	validCommandTypes  = []string{"slash", "prefix", "modal"}
	validCommandScopes = []string{"guild", "global"}
	validReturnTypes   = []string{"str", "int", "float", "bool", "None"}
	validArgTypes      = []string{"str", "int", "float", "bool", "discord.Member", "discord.Role"}
	validFieldStyles   = []string{"short", "paragraph"}
	validLicenses      = []string{"mit", "apache-2.0", "gpl-3.0", "bsd-3-clause", "unlicense", "no-license"}
	validHelpStyles    = []string{"compact", "detailed"}
)

// DefaultHelpStyle is used when a project predates the help_style key or leaves it unset
const DefaultHelpStyle = "compact"

// Discord allows at most five text inputs on a single modal page
const MaxModalFields = 5

// MaxFlowPages caps how many pages a multi page modal flow can chain together
const MaxFlowPages = 10

// MaxCommandResponses caps how many expected responses a command can declare
const MaxCommandResponses = 3

// Valid response types, only plain messages exist today
var validResponseTypes = []string{"message"}

// Discord caps text input labels at 45 characters
const maxFieldLabelLength = 45

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

func ValidateHelpStyle(s string) error {
	if s == "" {
		return fmt.Errorf("Please select a help style")
	}
	if !contains(validHelpStyles, s) {
		return fmt.Errorf("help style must be one of %s", strings.Join(validHelpStyles, ", "))
	}
	return nil
}

// NormalizeHelpStyle turns an unset help style into the default so readers never see ""
func NormalizeHelpStyle(s string) string {
	if s == "" {
		return DefaultHelpStyle
	}
	return s
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
		return fmt.Errorf("command type must be one of %s", strings.Join(validCommandTypes, ", "))
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

func fieldExists(fieldName string, fields []FieldInfo) bool {
	for _, field := range fields {
		if field.Name == fieldName {
			return true
		}
	}
	return false
}

func ValidateFieldName(s string, existing []FieldInfo) error {
	if s == "" {
		return fmt.Errorf("field name cannot be empty")
	}
	if strings.Contains(s, " ") {
		return fmt.Errorf("field name cannot contain spaces")
	}
	if strings.Contains(s, "-") {
		return fmt.Errorf("field name cannot contain dashes")
	}
	if fieldExists(s, existing) {
		return fmt.Errorf("field name already exists")
	}
	return nil
}

func ValidateFieldLabel(s string) error {
	if s == "" {
		return fmt.Errorf("field label cannot be empty")
	}
	if len([]rune(s)) > maxFieldLabelLength {
		return fmt.Errorf("field label must be at most %d characters", maxFieldLabelLength)
	}
	return nil
}

func ValidateFieldStyle(s string) error {
	if s == "" {
		return fmt.Errorf("field style cannot be empty")
	}
	if !contains(validFieldStyles, s) {
		return fmt.Errorf("field style must be one of %s", strings.Join(validFieldStyles, ", "))
	}
	return nil
}

// validateFields runs the shared field checks used by single page modals and flow pages
func validateFields(fields []FieldInfo) error {
	if len(fields) == 0 {
		return fmt.Errorf("modal commands need at least one field")
	}
	if len(fields) > MaxModalFields {
		return fmt.Errorf("modal commands can have at most %d fields", MaxModalFields)
	}
	for i, field := range fields {
		if err := ValidateFieldName(field.Name, fields[:i]); err != nil {
			return fmt.Errorf("field '%s': %w", field.Name, err)
		}
		if err := ValidateFieldLabel(field.Label); err != nil {
			return fmt.Errorf("field '%s': %w", field.Name, err)
		}
		if err := ValidateFieldStyle(field.Style); err != nil {
			return fmt.Errorf("field '%s': %w", field.Name, err)
		}
	}
	return nil
}

// ValidatePages checks a multi page modal flow, the first page in the slice starts the flow
func ValidatePages(pages []PageInfo) error {
	if len(pages) == 0 {
		return fmt.Errorf("multi page modal commands need at least one page")
	}
	if len(pages) > MaxFlowPages {
		return fmt.Errorf("multi page modal commands can have at most %d pages", MaxFlowPages)
	}

	// Collect every page name first so branches can point forward or backward
	pageNames := make(map[string]bool)
	for _, page := range pages {
		if page.Name == "" {
			return fmt.Errorf("page names cannot be empty")
		}
		if strings.Contains(page.Name, " ") {
			return fmt.Errorf("page '%s': page names cannot contain spaces", page.Name)
		}
		if pageNames[page.Name] {
			return fmt.Errorf("page name '%s' is used more than once", page.Name)
		}
		pageNames[page.Name] = true
	}

	for _, page := range pages {
		if err := validateFields(page.Fields); err != nil {
			return fmt.Errorf("page '%s': %w", page.Name, err)
		}

		// Branches can only test fields the user filled in on this page
		for _, branch := range page.Branches {
			if !fieldExists(branch.Field, page.Fields) {
				return fmt.Errorf("page '%s': branch field '%s' is not a field on this page", page.Name, branch.Field)
			}
			if branch.Goto == "" {
				return fmt.Errorf("page '%s': branch goto cannot be empty", page.Name)
			}
			if !pageNames[branch.Goto] {
				return fmt.Errorf("page '%s': branch goto '%s' is not a page name", page.Name, branch.Goto)
			}
		}

		// An empty next ends the flow, anything else has to be a real page
		if page.Next != "" && !pageNames[page.Next] {
			return fmt.Errorf("page '%s': next page '%s' is not a page name", page.Name, page.Next)
		}
	}

	return nil
}

// ValidateResponses checks the expected responses a command declares
func ValidateResponses(responses []ResponseInfo) error {
	if len(responses) > MaxCommandResponses {
		return fmt.Errorf("commands can have at most %d responses", MaxCommandResponses)
	}
	for i, response := range responses {
		if !contains(validResponseTypes, response.Type) {
			return fmt.Errorf("response %d: type must be one of %s", i+1, strings.Join(validResponseTypes, ", "))
		}
		if response.Content == "" {
			return fmt.Errorf("response %d: content cannot be empty", i+1)
		}
		// The content lands inside a python string literal so these characters would break the generated file
		if strings.ContainsAny(response.Content, "\"\\\n") {
			return fmt.Errorf("response %d: content cannot contain double quotes, backslashes, or newlines", i+1)
		}
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
	if err := ValidateResponses(command.Responses); err != nil {
		return err
	}
	if command.Type == "modal" {
		if len(command.Args) > 0 {
			return fmt.Errorf("modal commands cannot have arguments")
		}
		if len(command.Fields) > 0 && len(command.Pages) > 0 {
			return fmt.Errorf("modal commands cannot have both fields and pages")
		}
		if len(command.Pages) > 0 {
			return ValidatePages(command.Pages)
		}
		return validateFields(command.Fields)
	}
	if len(command.Fields) > 0 {
		return fmt.Errorf("only modal commands can have fields")
	}
	if len(command.Pages) > 0 {
		return fmt.Errorf("only modal commands can have pages")
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
