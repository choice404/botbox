/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	Primary   lipgloss.AdaptiveColor
	Secondary lipgloss.AdaptiveColor
	Ternary   lipgloss.AdaptiveColor

	red      = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo   = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green    = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
	cyan     = lipgloss.AdaptiveColor{Light: "#00FFFF", Dark: "#0066aa"}
	white    = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}
	blue     = lipgloss.AdaptiveColor{Light: "#0077FF", Dark: "#0000FF"}
	navy     = lipgloss.AdaptiveColor{Light: "#000080", Dark: "#000080"}
	sapphire = lipgloss.AdaptiveColor{Light: "#0F52BA", Dark: "#0F52BA"}
	emerald  = lipgloss.AdaptiveColor{Light: "#50C878", Dark: "#50C878"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	ErrorText,
	KeyText,
	ValueText,
	Help lipgloss.Style
	FooterText lipgloss.Style
}

func SetColorScheme(scheme string) {
	switch scheme {
	case "default":
		Primary = cyan
		Secondary = emerald
		Ternary = indigo
	case "dark":
		Primary = white
		Secondary = sapphire
		Ternary = emerald
	case "ubuntu":
		Primary = lipgloss.AdaptiveColor{Light: "#E95420", Dark: "#E95420"}
		Secondary = lipgloss.AdaptiveColor{Light: "#D3A625", Dark: "#D3A625"}
		Ternary = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	default:
		Primary = lipgloss.AdaptiveColor{Light: "#00FFFF", Dark: "#0066aa"}
		Secondary = lipgloss.AdaptiveColor{Light: "#50C878", Dark: "#50C878"}
		Ternary = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	}
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	SetColorScheme(globalConfig.Display.ColorScheme)

	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(Primary).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Primary).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(Secondary).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.ErrorText = lg.NewStyle().
		Foreground(red).
		Bold(true)
	s.KeyText = lg.NewStyle().
		Foreground(Ternary).
		Bold(true)
	s.ValueText = lg.NewStyle().
		Foreground(Secondary)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	s.FooterText = lg.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true).
		Padding(0, 1, 0, 2)
	return &s
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
