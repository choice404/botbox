/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"fmt"
	"maps"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	globalConfig GlobalConfig
)

const maxWidth = 180

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

type Values struct {
	Map  map[string]*string
	Name string
}

type FormWrapper struct {
	Name               string
	Form               func(Values, Values) *huh.Form
	Values             Values
	SkipCondition      func(ModelValues Values, allForms []FormWrapper, currentIndex int) bool
	SkipCallback       func(ModelValues Values, allForms []FormWrapper, currentIndex int)
	Callback           func(formValues Values, ModelValues Values, allForms []FormWrapper)
	BranchCallback     func(Values, []FormWrapper) int
	ShowStatus         bool
	FormGroup          string
	BranchValueHandler func(targetFormIndex int, targetValues Values)
}

func (fw *FormWrapper) GetValues() Values {
	resultMap := make(map[string]string)
	for key, ptr := range fw.Values.Map {
		if ptr != nil {
			resultMap[key] = *ptr
		}
	}
	result := Values{
		Map:  make(map[string]*string),
		Name: fmt.Sprintf("%s_%s_%s", fw.FormGroup, fw.Name, fw.Values.Name),
	}
	return result
}

func (fw *FormWrapper) ExecuteCallback(ModelValues Values, allForms []FormWrapper) {
	if fw.Callback != nil {
		fw.Callback(fw.Values, ModelValues, allForms)
	}
}

func (fw *FormWrapper) ExecuteBranchCallback(allForms []FormWrapper) int {
	if fw.BranchCallback != nil {
		return fw.BranchCallback(fw.Values, allForms)
	}
	return -1
}

func (fw *FormWrapper) ShouldSkip(ModelValues Values, allForms []FormWrapper, currentIndex int) bool {
	if fw.SkipCondition != nil {
		return fw.SkipCondition(ModelValues, allForms, currentIndex)
	}
	return false
}

func (fw *FormWrapper) ExecuteSkipCallback(ModelValues Values, allForms []FormWrapper, currentIndex int) {
	if fw.SkipCallback != nil {
		fw.SkipCallback(ModelValues, allForms, currentIndex)
	}
}

func (m *Model) Init() tea.Cmd {
	conf, err := LoadGlobalConfig()
	if err != nil {
		errors := []error{fmt.Errorf("failed to load global config: %w", err)}
		m.HandleError(errors)
	}
	globalConfig = *conf

	originalFormsCount := len(m.forms)

	m.forms = append(m.forms, FormWrapper{
		Name:       "Complete",
		Form:       finalCompleteFormGenerator,
		Values:     m.ModelValues,
		ShowStatus: true,
		FormGroup:  "final",
	})

	allValueMaps := make([]Values, len(m.forms))
	for i, form := range m.forms {
		if form.Values.Map == nil {
			form.Values.Map = make(map[string]*string)
		}
		allValueMaps[i] = form.Values
	}

	if m.ModelValues.Map == nil {
		m.ModelValues.Map = make(map[string]*string)
	}

	if m.initCallback != nil {
		m.initCallback(m, allValueMaps)
	}

	if m.Error != nil {
		m.jumpToFinalStatus()
		m.currentForm.WithShowHelp(false)
		return m.currentForm.Init()
	}

	if originalFormsCount == 0 {
		m.currentFormPtr = len(m.forms) - 1
		m.currentForm = m.forms[m.currentFormPtr].Form(m.forms[m.currentFormPtr].Values, m.ModelValues)
		if m.currentForm == nil {
			return tea.Quit
		}
		m.currentForm.WithShowHelp(false)
		m.state = stateDone
		return m.currentForm.Init()
	}

	if !m.setCurrentFormToValidForm(0) {
		return m.handleDirectCompletion()
	}

	m.currentForm = m.forms[m.currentFormPtr].Form(m.forms[m.currentFormPtr].Values, m.ModelValues)
	if m.currentForm == nil {
		return tea.Quit
	}

	m.currentForm.WithShowHelp(false)
	return m.currentForm.Init()
}

type Model struct {
	viewportOffset    int
	viewportHeight    int
	maxViewportHeight int
	title             string
	state             state
	lg                *lipgloss.Renderer
	styles            *Styles
	currentFormPtr    int
	currentForm       *huh.Form
	forms             []FormWrapper
	width             int
	initCallback      func(*Model, []Values)
	callback          func(*Model) []error
	displayCallback   func() string
	ModelValues       Values
	displayKeys       []string
	Error             []error
}

func (m *Model) resetViewport() {
	m.viewportOffset = 0
}

func (m *Model) getMaxViewportOffset() int {
	if m.state != stateDone {
		return 0
	}

	var content string

	if m.Error != nil {
		var b strings.Builder
		b.WriteString("Errors:\n\n")
		for _, err := range m.Error {
			if b.Len() > 0 {
				b.WriteString("\n")
			}
			b.WriteString(err.Error())
		}
		content = b.String()
	} else if m.displayCallback == nil {
		var b strings.Builder
		var values Values
		if m.displayKeys == nil || len(m.displayKeys) == 0 {
			values = m.getValues()
		} else {
			values = m.GetValuesByKeys(m.displayKeys)
		}

		for key, value := range values.Map {
			if key == "" && *value == "" {
				continue
			}
			if key != "" {
				b.WriteString(fmt.Sprintf("%s: %s\n", key, *value))
			} else {
				b.WriteString(*value + "\n")
			}
		}
		content = b.String()
	} else {
		content = m.displayCallback()
		if content == "" {
			content = "No values to display."
		}
	}

	lines := strings.Split(content, "\n")
	availableContentHeight := max(1, m.viewportHeight-6)

	if len(lines) <= availableContentHeight {
		return 0
	}
	return len(lines) - availableContentHeight
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = int(math.Min(float64(msg.Width), float64(maxWidth))) - m.styles.Base.GetHorizontalFrameSize()
		m.viewportHeight = msg.Height - 8
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt

		case "up":
			if m.state == stateDone && m.viewportOffset > 0 {
				m.viewportOffset--
				return m, nil
			}
		case "down":
			if m.state == stateDone {
				maxOffset := m.getMaxViewportOffset()
				if m.viewportOffset < maxOffset {
					m.viewportOffset++
					return m, nil
				}
			}

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

	if m.Error != nil && m.state != stateDone {
		m.jumpToFinalStatus()
		return m, nil
	}

	var cmds []tea.Cmd

	if m.state != stateDone {
		form, cmd := m.currentForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.currentForm = f
			cmds = append(cmds, cmd)
		}

		if m.currentForm.State == huh.StateCompleted {
			m.forms[m.currentFormPtr].ExecuteCallback(m.ModelValues, m.forms)
			m.state = stateDone

			if !m.forms[m.currentFormPtr].ShowStatus {
				return m.handleFormCompletion()
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) handleFormCompletion() (tea.Model, tea.Cmd) {
	if m.Error != nil {
		m.jumpToFinalStatus()
		return m, nil
	}
	if m.forms[m.currentFormPtr].BranchCallback != nil {
		branchIndex := m.forms[m.currentFormPtr].ExecuteBranchCallback(m.forms)

		if branchIndex == -1 {
			if !m.moveToNextValidForm() {
				return m.handleFinalCompletion()
			}
		} else if branchIndex == -2 {
			return m.handleFinalCompletion()
		} else if branchIndex >= 0 && branchIndex < len(m.forms) {
			if m.forms[m.currentFormPtr].BranchValueHandler != nil {
				m.forms[m.currentFormPtr].BranchValueHandler(branchIndex, m.forms[branchIndex].Values)
			}

			if m.forms[branchIndex].ShouldSkip(m.ModelValues, m.forms, branchIndex) {
				m.forms[branchIndex].ExecuteSkipCallback(m.ModelValues, m.forms, branchIndex)
				if !m.setCurrentFormToValidForm(branchIndex + 1) {
					return m.handleFinalCompletion()
				}
			} else {
				m.currentFormPtr = branchIndex
			}
		}
	} else {
		if !m.moveToNextValidForm() {
			return m.handleFinalCompletion()
		}
	}

	if m.currentFormPtr < len(m.forms) {
		m.resetViewport()

		newForm := m.forms[m.currentFormPtr].Form(m.forms[m.currentFormPtr].Values, m.ModelValues)

		if newForm.State == huh.StateCompleted {
			m.currentForm = newForm
			m.forms[m.currentFormPtr].ExecuteCallback(m.ModelValues, m.forms)
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
		versionMessage := s.Highlight.Render(fmt.Sprintf("BotBox Version: %s", Version))
		if m.forms[m.currentFormPtr].ShowStatus || m.currentFormPtr >= len(m.forms)-1 {
			var content string
			var header string

			if m.Error != nil {
				var b strings.Builder
				b.WriteString("Errors:\n\n")
				for _, err := range m.Error {
					if b.Len() > 0 {
						b.WriteString("\n")
					}
					b.WriteString(err.Error())
				}
				content = b.String()
				header = m.appBoundaryView(m.title)
			} else if m.displayCallback == nil {
				var b strings.Builder
				var values Values
				if m.displayKeys == nil || len(m.displayKeys) == 0 {
					values = m.getValues()
				} else {
					values = m.GetValuesByKeys(m.displayKeys)
				}

				for key, value := range values.Map {
					key := s.HeaderText.Render(key)
					value := s.Highlight.Render(*value)
					if key == "" && value == "" {
						continue
					}

					if key != "" {
						b.WriteString(fmt.Sprintf("%s: %s\n", key, value))
					} else {
						b.WriteString(value + "\n")
					}
				}
				content = b.String()
				header = m.appBoundaryView(m.title)
			} else {
				content = m.displayCallback()
				header = m.appBoundaryView(m.title)
				if content == "" {
					content = s.Highlight.Render("No values to display.")
				}
			}

			lines := strings.Split(content, "\n")
			scrollableContent := content

			availableContentHeight := max(1, m.viewportHeight-6)

			if len(lines) > availableContentHeight {
				if m.viewportOffset > 0 {
					startLine := min(m.viewportOffset, len(lines)-1)
					endLine := min(startLine+availableContentHeight, len(lines))
					if startLine < len(lines) && endLine > startLine {
						scrollableContent = strings.Join(lines[startLine:endLine], "\n")
					}
				} else {
					scrollableContent = strings.Join(lines[:availableContentHeight], "\n")
				}
			}

			statusBox := s.Status.Margin(0, 0).Padding(1, 2).Width(80).Render(scrollableContent)

			scrollIndicator := ""
			maxOffset := m.getMaxViewportOffset()
			if maxOffset > 0 {
				totalLines := len(lines)
				currentEndLine := min(m.viewportOffset+availableContentHeight, totalLines)
				scrollIndicator = s.Help.Render(fmt.Sprintf("↑/↓ to scroll (line %d-%d of %d)",
					m.viewportOffset+1,
					currentEndLine,
					totalLines))
			}

			confirmPrompt := s.Highlight.Render("Press Enter to submit, or Esc/Q to cancel.")

			message := header + "\n" + statusBox
			if scrollIndicator != "" {
				message += "\n" + scrollIndicator
			}
			message += "\n\n" + confirmPrompt

			result := s.Base.Render(message + "\n")
			result += "\n" + versionMessage + "\n\n"

			return result
		}

		header := m.appBoundaryView(m.title)
		message := s.Highlight.Render("Press Enter to continue.")
		return s.Base.Render(header + "\n" + message + "\n\n" + versionMessage + "\n\n")
	default:
		v := strings.TrimSuffix(m.currentForm.View(), "\n\n")

		lines := strings.Split(v, "\n")
		if m.viewportOffset > 0 && len(lines) > m.viewportHeight {
			startLine := m.viewportOffset
			endLine := min(m.viewportOffset+m.viewportHeight, len(lines))
			if startLine < len(lines) {
				v = strings.Join(lines[startLine:endLine], "\n")
			}
		} else if len(lines) > m.viewportHeight {
			v = strings.Join(lines[:m.viewportHeight], "\n")
		}

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

		// NOTE: Uncomment to enable scrolling functionality
		// scrollIndicator := ""
		// maxOffset := m.getMaxViewportOffset()
		// if maxOffset > 0 {
		// 	scrollIndicator = s.Help.Render(fmt.Sprintf("↑/↓ to scroll (line %d-%d of %d)",
		// 		m.viewportOffset+1,
		// 		int(math.Min(float64(m.viewportOffset+m.viewportHeight), float64(len(strings.Split(strings.TrimSuffix(m.currentForm.View(), "\n\n"), "\n"))))),
		// 		len(strings.Split(strings.TrimSuffix(m.currentForm.View(), "\n\n"), "\n"))))
		// }

		return s.Base.Render(header + "\n" + body + "\n\n" +
			footer + "\n" +
			// NOTE: Uncomment to enable scrolling functionality
			// scrollIndicator +
			"\n\n" + s.Highlight.Render("BotBox Version: "+Version) + "\n\n")
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

func (m Model) GetValuesByKeys(keys []string) Values {
	values := m.GetAllValuesFlat()
	for _, key := range keys {
		if value, exists := m.ModelValues.Map[key]; exists && value != nil {
			values.Map[key] = value
		}
	}
	return values
}

func (m *Model) getValues() Values {
	return m.forms[m.currentFormPtr].GetValues()
}

func CupSleeve(cup Model) {
	_, err := tea.NewProgram(&cup, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func (m *Model) executeCallback() {
	if m.callback != nil {
		m.HandleError(m.callback(m))
	}
}

func (m *Model) SetTitle(title string) {
	m.title = title
}

func (m *Model) SetForms(forms []FormWrapper) {
	m.forms = forms
}

func (m *Model) GetAllValues() Values {
	allValues := make(map[string]*string)

	for i, form := range m.forms {
		formValues := form.GetValues()

		var prefix string
		if form.Name != "" {
			prefix = strings.ToLower(strings.ReplaceAll(form.Name, " ", "_"))
		} else {
			prefix = fmt.Sprintf("form_%d", i)
		}

		for key, value := range formValues.Map {
			uniqueKey := fmt.Sprintf("%s_%s", prefix, key)
			allValues[uniqueKey] = value
		}
	}
	result := Values{
		Map:  allValues,
		Name: fmt.Sprintf("%s_%s", m.title, "all_values"),
	}
	return result
}

func (m *Model) FlattenModelValuesWithPrefix() {
	for key, value := range m.ModelValues.Map {
		if value != nil {
			m.ModelValues.Map[key] = new(string)
			*m.ModelValues.Map[key] = *value
		} else {
			m.ModelValues.Map[key] = new(string)
		}
	}

	for i, form := range m.forms {
		prefix := strings.ToLower(strings.ReplaceAll(form.Name, " ", "_"))
		if prefix == "" {
			prefix = fmt.Sprintf("form_%d", i)
		}

		for key, value := range form.Values.Map {
			if value != nil && *value != "" {
				newKey := fmt.Sprintf("%s_%s", prefix, key)
				m.ModelValues.Map[newKey] = value
			}
		}
	}
}

func (m *Model) GetAllValuesFlat() Values {
	allValues := make(map[string]*string)
	for key, ptr := range m.ModelValues.Map {
		if ptr != nil && *ptr != "" {
			allValues[key] = ptr
		}
	}

	for _, form := range m.forms {
		formValues := form.GetValues()
		maps.Copy(allValues, formValues.Map)
	}

	result := Values{
		Map:  allValues,
		Name: fmt.Sprintf("%s_%s", m.title, "all_values"),
	}
	return result
}

func (m *Model) FlattenModelValues() {
	for key, value := range m.ModelValues.Map {
		if value != nil {
			m.ModelValues.Map[key] = new(string)
			*m.ModelValues.Map[key] = *value
		} else {
			m.ModelValues.Map[key] = new(string)
		}
	}
}

func (m *Model) GetValuesByFormName(formName string) Values {
	for _, form := range m.forms {
		if form.Name == formName {
			values := form.GetValues()
			values.Name = fmt.Sprintf("%s_%s", m.title, formName)
			return values
		}
	}
	return Values{
		Map:  make(map[string]*string),
		Name: fmt.Sprintf("%s_%s", m.title, formName),
	}
}

func (m *Model) GetValuesByFormIndex(index int) Values {
	if index >= 0 && index < len(m.forms) {
		values := m.forms[index].GetValues()
		values.Name = fmt.Sprintf("%s_form_%d", m.title, index)
		return values
	}
	return Values{
		Map:  make(map[string]*string),
		Name: fmt.Sprintf("%s_form_%d", m.title, index),
	}
}

func (m *Model) GetValue(formName, key string) (string, bool) {
	formValues := m.GetValuesByFormName(formName)
	value, exists := formValues.Map[key]
	return *value, exists
}

func (m *Model) GetValueFlat(key string) (string, bool) {
	for _, form := range m.forms {
		formValues := form.GetValues()
		if value, exists := formValues.Map[key]; exists {
			return *value, true
		}
	}
	return "", false
}

func (m *Model) findNextValidForm(startIndex int) int {
	for i := startIndex; i < len(m.forms); i++ {
		if !m.forms[i].ShouldSkip(m.ModelValues, m.forms, i) {
			return i
		}
		m.forms[i].ExecuteSkipCallback(m.ModelValues, m.forms, i)
	}
	return -1
}

func (m *Model) moveToNextValidForm() bool {
	nextIndex := m.findNextValidForm(m.currentFormPtr + 1)
	if nextIndex == -1 {
		return false
	}
	m.currentFormPtr = nextIndex
	return true
}

func (m *Model) setCurrentFormToValidForm(startIndex int) bool {
	validIndex := m.findNextValidForm(startIndex)
	if validIndex == -1 {
		return false
	}
	m.currentFormPtr = validIndex
	return true
}

func (m *Model) handleDirectCompletion() tea.Cmd {
	m.currentFormPtr = len(m.forms) - 1
	m.currentForm = m.forms[m.currentFormPtr].Form(m.forms[m.currentFormPtr].Values, m.ModelValues)
	if m.currentForm == nil {
		return tea.Quit
	}
	m.currentForm.WithShowHelp(false)
	m.state = stateDone
	return m.currentForm.Init()
}

func (m *Model) HandleError(err []error) {
	m.Error = err
	m.jumpToFinalStatus()
}

func (m *Model) jumpToFinalStatus() {
	m.currentFormPtr = len(m.forms) - 1
	m.createFinalStatusForm()
	m.state = stateDone
}

func FlattenValuesInto(targetValues Values, sourceValues Values) {
	for key, value := range sourceValues.Map {
		if existingValue, exists := targetValues.Map[key]; exists && existingValue != nil {
			*existingValue = *value
		} else {
			newValue := value
			targetValues.Map[key] = newValue
		}
	}
}

func MergeValues(targetValues Values, sourceValues Values) Values {
	for key, value := range sourceValues.Map {
		if existingValue, exists := targetValues.Map[key]; exists && existingValue != nil {
			suffix := 1
			for {
				if _, exists := targetValues.Map[fmt.Sprintf("%s_%d", key, suffix)]; !exists {
					newKey := fmt.Sprintf("%s_%d", key, suffix)
					newValue := value
					targetValues.Map[newKey] = newValue
					break
				}
				suffix++
			}

		} else {
			newValue := value
			targetValues.Map[key] = newValue
		}
	}
	return targetValues
}

func ResetFormValues(values Values) {
	for key := range values.Map {
		values.Map[key] = new(string)
	}
}

func PreserveFormValues(values Values) {}

func SetSpecificValues(values Values, updates Values) {
	for key, value := range updates.Map {
		if _, exists := values.Map[key]; exists {
			newValue := value
			values.Map[key] = newValue
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
