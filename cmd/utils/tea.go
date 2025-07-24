/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"fmt"
	"maps"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const maxWidth = 180

var (
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
	KeyText,
	ValueText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(cyan).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cyan).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.KeyText = lg.NewStyle().
		Foreground(indigo).
		Bold(true)
	s.ValueText = lg.NewStyle().
		Foreground(emerald)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type state int

const (
	statusNormal state = iota
	stateDone
)

type FormField struct {
	Key   string
	Value any
	Field huh.Field
}

type FormWrapper struct {
	Name               string
	Form               func(map[string]*string, map[string]*string) *huh.Form
	Values             map[string]*string
	Callback           func(formValues map[string]*string, modelValues map[string]*string, allForms []FormWrapper)
	BranchCallback     func(map[string]*string) int
	ShowStatus         bool
	FormGroup          string
	BranchValueHandler func(targetFormIndex int, targetValues map[string]*string)
}

func (fw *FormWrapper) GetValues() map[string]string {
	result := make(map[string]string)
	for key, ptr := range fw.Values {
		if ptr != nil {
			result[key] = *ptr
		}
	}
	return result
}

func (fw *FormWrapper) ExecuteCallback(ModelValues map[string]*string, allForms []FormWrapper) {
	if fw.Callback != nil {
		fw.Callback(fw.Values, ModelValues, allForms)
	}
}

func (fw *FormWrapper) ExecuteBranchCallback() int {
	if fw.BranchCallback != nil {
		return fw.BranchCallback(fw.Values)
	}
	return -1
}

func (m *Model) Init() tea.Cmd {
	originalFormsCount := len(m.forms)

	m.forms = append(m.forms, FormWrapper{
		Name:       "Complete",
		Form:       finalCompleteFormGenerator,
		Values:     m.modelValues,
		ShowStatus: true,
		FormGroup:  "final",
	})
	allValueMaps := make([]map[string]*string, len(m.forms))
	for i, form := range m.forms {
		allValueMaps[i] = form.Values
	}

	if m.initCallback != nil {
		m.initCallback(m.modelValues, allValueMaps)
	}

	if originalFormsCount == 0 {
		m.currentFormPtr = len(m.forms) - 1
		m.currentForm = m.forms[m.currentFormPtr].Form(m.forms[m.currentFormPtr].Values, m.modelValues)
		if m.currentForm == nil {
			return tea.Quit
		}
		m.currentForm.WithShowHelp(false)
		m.state = stateDone
		return m.currentForm.Init()
	}

	m.currentForm = m.forms[m.currentFormPtr].Form(m.forms[m.currentFormPtr].Values, m.modelValues)
	if m.currentForm == nil {
		return tea.Quit
	}
	m.currentForm.WithShowHelp(false)
	return m.currentForm.Init()
}

type Model struct {
	title           string
	state           state
	lg              *lipgloss.Renderer
	styles          *Styles
	currentFormPtr  int
	currentForm     *huh.Form
	forms           []FormWrapper
	width           int
	initCallback    func(map[string]*string, []map[string]*string)
	callback        func(values map[string]string)
	displayCallback func() string
	modelValues     map[string]*string
	displayKeys     []string
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc":
			return m, tea.Quit
		case "q":
			if m.state == stateDone {
				return m, tea.Quit
			}
		case "enter":
			if m.state == stateDone {
				if m.forms[m.currentFormPtr].ShowStatus || m.currentFormPtr >= len(m.forms)-1 {
					m.executeCallback()
					return m, tea.Quit
				}
				return m.handleFormCompletion()
			}
		}
	}

	var cmds []tea.Cmd

	if m.state != stateDone {
		form, cmd := m.currentForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.currentForm = f
			cmds = append(cmds, cmd)
		}

		if m.currentForm.State == huh.StateCompleted {
			m.forms[m.currentFormPtr].ExecuteCallback(m.modelValues, m.forms)
			m.state = stateDone

			if !m.forms[m.currentFormPtr].ShowStatus {
				return m.handleFormCompletion()
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) handleFormCompletion() (tea.Model, tea.Cmd) {
	if m.forms[m.currentFormPtr].BranchCallback != nil {
		branchIndex := m.forms[m.currentFormPtr].ExecuteBranchCallback()

		if branchIndex == -1 {
			if m.currentFormPtr < len(m.forms)-1 {
				m.currentFormPtr++
			} else {
				return m.handleFinalCompletion()
			}
		} else if branchIndex == -2 {
			return m.handleFinalCompletion()
		} else if branchIndex >= 0 && branchIndex < len(m.forms) {
			if m.forms[m.currentFormPtr].BranchValueHandler != nil {
				m.forms[m.currentFormPtr].BranchValueHandler(branchIndex, m.forms[branchIndex].Values)
			}
			m.currentFormPtr = branchIndex
		}
	} else {
		if m.currentFormPtr < len(m.forms)-1 {
			m.currentFormPtr++
		} else {
			return m.handleFinalCompletion()
		}
	}

	if m.currentFormPtr < len(m.forms) {
		newForm := m.forms[m.currentFormPtr].Form(m.forms[m.currentFormPtr].Values, m.modelValues)

		if newForm.State == huh.StateCompleted {
			m.currentForm = newForm
			m.forms[m.currentFormPtr].ExecuteCallback(m.modelValues, m.forms)
			m.state = stateDone

			if m.forms[m.currentFormPtr].ShowStatus || m.currentFormPtr >= len(m.forms)-1 {
				return m, nil
			}

			return m.handleFormCompletion()
		}

		m.currentForm = newForm
		m.currentForm.WithShowHelp(false)
		initCmd := m.currentForm.Init()
		m.state = statusNormal

		updatedForm, updateCmd := m.currentForm.Update(tea.WindowSizeMsg{Width: m.width, Height: 24})
		if f, ok := updatedForm.(*huh.Form); ok {
			m.currentForm = f
		}

		return m, tea.Batch(initCmd, updateCmd)
	}

	return m, tea.Batch()
}

func (m *Model) handleFinalCompletion() (tea.Model, tea.Cmd) {
	m.createFinalStatusForm()
	m.state = stateDone
	return m, nil
}

func (m *Model) createFinalStatusForm() {
	finalForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(" ").
				Description(" "),
		),
	)
	finalForm.State = huh.StateCompleted
	m.currentForm = finalForm

	m.forms[m.currentFormPtr].ShowStatus = true
}

func (m Model) View() string {
	s := m.styles

	if m.currentForm == nil {
		return s.Base.Render("Loading...")
	}

	switch m.currentForm.State {
	case huh.StateCompleted:
		if m.forms[m.currentFormPtr].ShowStatus || m.currentFormPtr >= len(m.forms)-1 {
			var message string
			if m.displayCallback == nil {
				var values map[string]string
				if m.displayKeys == nil || len(m.displayKeys) == 0 {
					values = m.getValues()
				} else {
					values = m.GetValuesByKeys(m.displayKeys)
				}
				var b strings.Builder

				for key, value := range values {
					key := s.HeaderText.Render(key)
					value := s.Highlight.Render(value)
					if key == "" && value == "" {
						continue
					}

					if key != "" {
						b.WriteString(fmt.Sprintf("%s: %s\n", key, value))
					} else {
						b.WriteString(value + "\n")
					}
				}

				header := m.appBoundaryView(m.title)
				statusBox := s.Status.Margin(0, 0).Padding(1, 2).Width(80).Render(b.String())
				message = header + "\n" + statusBox
			} else {
				displayContent := m.displayCallback()
				header := m.appBoundaryView(m.title)
				if displayContent == "" {
					displayContent = s.Highlight.Render("No values to display.")
				}

				statusBox := s.Status.Margin(0, 0).Padding(1, 2).Width(80).Render(displayContent)
				message = header + "\n" + statusBox
			}

			confirmPrompt := s.Highlight.Render("Press Enter to submit, or Esc/Q to cancel.")
			return s.Base.Render(message + "\n\n" + confirmPrompt + "\n\n")
		}
		header := m.appBoundaryView(m.title)
		message := s.Highlight.Render("Press Enter to continue.")
		return s.Base.Render(header + "\n" + message + "\n\n")
	default:
		v := strings.TrimSuffix(m.currentForm.View(), "\n\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		errors := m.currentForm.Errors()
		header := m.appBoundaryView(m.title)

		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinHorizontal(lipgloss.Left, form)

		footer := m.appBoundaryView(m.currentForm.Help().ShortHelpView(m.currentForm.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		return s.Base.Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m Model) errorView() string {
	var s string
	for _, err := range m.currentForm.Errors() {
		s += err.Error()
	}
	return s
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("="),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func (m Model) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("="),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func (m Model) GetValuesByKeys(keys []string) map[string]string {
	values := m.GetAllValuesFlat()
	for _, key := range keys {
		if value, exists := m.modelValues[key]; exists && value != nil {
			values[key] = *value
		}
	}
	return values
}

func (m Model) getValues() map[string]string {
	return m.forms[m.currentFormPtr].GetValues()
}

func CupSleeve(cup tea.Model) {
	_, err := tea.NewProgram(cup).Run()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func (m *Model) executeCallback() {
	if m.callback != nil {
		allValues := m.GetAllValuesFlat()
		m.callback(allValues)
	}
}

func (m *Model) SetTitle(title string) {
	m.title = title
}

func (m *Model) SetForms(forms []FormWrapper) {
	m.forms = forms
}

func (m *Model) GetAllValues() map[string]string {
	allValues := make(map[string]string)

	for i, form := range m.forms {
		formValues := form.GetValues()

		var prefix string
		if form.Name != "" {
			prefix = strings.ToLower(strings.ReplaceAll(form.Name, " ", "_"))
		} else {
			prefix = fmt.Sprintf("form_%d", i)
		}

		for key, value := range formValues {
			uniqueKey := fmt.Sprintf("%s_%s", prefix, key)
			allValues[uniqueKey] = value
		}
	}

	return allValues
}

func (m *Model) GetAllValuesFlat() map[string]string {
	allValues := make(map[string]string)

	for key, ptr := range m.modelValues {
		if ptr != nil && *ptr != "" {
			allValues[key] = *ptr
		}
	}

	for _, form := range m.forms {
		formValues := form.GetValues()
		maps.Copy(allValues, formValues)
	}

	return allValues
}

func (m *Model) GetValuesByFormName(formName string) map[string]string {
	for _, form := range m.forms {
		if form.Name == formName {
			return form.GetValues()
		}
	}
	return make(map[string]string)
}

func (m *Model) GetValuesByFormIndex(index int) map[string]string {
	if index >= 0 && index < len(m.forms) {
		return m.forms[index].GetValues()
	}
	return make(map[string]string)
}

func (m *Model) GetValue(formName, key string) (string, bool) {
	formValues := m.GetValuesByFormName(formName)
	value, exists := formValues[key]
	return value, exists
}

func (m *Model) GetValueFlat(key string) (string, bool) {
	for _, form := range m.forms {
		formValues := form.GetValues()
		if value, exists := formValues[key]; exists {
			return value, true
		}
	}
	return "", false
}

func ResetFormValues(values map[string]*string) {
	for key := range values {
		values[key] = new(string)
	}
}

func PreserveFormValues(values map[string]*string) {}

func SetSpecificValues(values map[string]*string, updates map[string]string) {
	for key, value := range updates {
		if _, exists := values[key]; exists {
			newValue := value
			values[key] = &newValue
		}
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
