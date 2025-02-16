package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type autocompletemodel struct {
	input         textinput.Model
	suggestions   []string
	filtered      []string
	index         int
	baseMatch     string
	autocompleted bool
	selected      string
}

var (
	AutotitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	AutoinputStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	AutohintStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	AutosuggestStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

func initialAutocompleteModel() autocompletemodel {
	input := textinput.New()
	input.Focus()

	return autocompletemodel{
		input:         input,
		suggestions:   options,
		filtered:      options,
		index:         -1,
		baseMatch:     "",
		autocompleted: false,
		selected:      "",
	}
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, str := range strs[1:] {
		for !strings.HasPrefix(str, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}
	return prefix
}

func AutofilterSuggestions(input string, options []string) []string {
	var filtered []string
	for _, option := range options {
		if strings.HasPrefix(option, input) {
			filtered = append(filtered, option)
		}
	}
	return filtered
}

func isValidHostname(input string, options []string) bool {
	for _, option := range options {
		if input == option {
			return true
		}
	}
	return false
}

func (m autocompletemodel) Init() tea.Cmd {
	return textinput.Blink
}

func (m autocompletemodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if len(m.filtered) == 0 {
				m.filtered = AutofilterSuggestions(m.input.Value(), m.suggestions)
				m.autocompleted = false
			}
			if len(m.filtered) > 0 {
				if !m.autocompleted {
					commonPrefix := longestCommonPrefix(m.filtered)
					m.input.SetValue(commonPrefix)
					m.autocompleted = true
					m.index = -1
				} else {
					m.index = (m.index + 1) % len(m.filtered)
					m.input.SetValue(m.filtered[m.index])
				}
				m.input.CursorEnd()
			}
			return m, nil
		case "esc":
			return m, tea.Quit
		case "enter":
			if isValidHostname(m.input.Value(), m.suggestions) {
				m.selected = m.input.Value()
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.filtered = AutofilterSuggestions(m.input.Value(), m.suggestions)
	m.autocompleted = false

	return m, cmd
}

func (m autocompletemodel) View() string {
	return fmt.Sprintf("%s\n%s %s\n%s",
		AutotitleStyle.Render("SSH Hostname Selector"),
		"Enter hostname:", inputStyle.Render(m.input.View()),
		AutohintStyle.Render("(Press TAB to autocomplete, ENTER to submit a valid hostname)")+"\n"+AutosuggestStyle.Render(fmt.Sprintf("Suggestions: %v", m.filtered)))
}
