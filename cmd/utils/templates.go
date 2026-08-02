/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed all:templates
var templateFS embed.FS

// CogTemplateData holds the values rendered into cog.py.tmpl
type CogTemplateData struct {
	Author         string
	BotName        string
	BotDescription string
	ClassName      string
	Filename       string
	SlashCommands  []CommandInfo
	PrefixCommands []CommandInfo
}

// templateFuncs holds the helpers available inside all templates
var templateFuncs = template.FuncMap{
	"returnValue": GetReturnValue,
	"argString":   BuildArgString,
	"underscore":  underscoreName,
}

// RenderTemplate renders the named embedded template with the given data
// Templates use << >> delimiters so Python brace syntax stays untouched
func RenderTemplate(name string, data any) (string, error) {
	tmpl, err := template.New(name).Delims("<<", ">>").Funcs(templateFuncs).ParseFS(templateFS, "templates/"+name)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	var content strings.Builder
	if err := tmpl.Execute(&content, data); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", name, err)
	}

	return content.String(), nil
}

// GetReturnValue maps a Python return type to its zero value literal
func GetReturnValue(returnType string) string {
	switch returnType {
	case "str":
		return `""`
	case "int":
		return "0"
	case "float":
		return "0.0"
	case "bool":
		return "False"
	default:
		return "None"
	}
}

// BuildArgString joins command args into a Python parameter list
func BuildArgString(args []ArgInfo) string {
	if len(args) == 0 {
		return ""
	}
	var argBuilder strings.Builder
	for i, arg := range args {
		fmt.Fprintf(&argBuilder, "%s: %s", arg.Name, arg.Type)
		if i < len(args)-1 {
			argBuilder.WriteString(", ")
		}
	}
	return argBuilder.String()
}

// underscoreName converts dashes in a command name to underscores
func underscoreName(name string) string {
	return strings.ReplaceAll(name, "-", "_")
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
