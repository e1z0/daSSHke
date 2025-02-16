package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type model struct {
	input       textinput.Model
	suggestions []string
	filtered    []string
	index       int
	selected    string
}

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	inputStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	hintStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	suggestStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
)

func initialModel() model {
	input := textinput.New()
	input.Focus()

	return model{
		input:       input,
		suggestions: options,
		filtered:    options,
		index:       -1,
		selected:    "",
	}
}

func filterSuggestions(input string, options []string) []string {
	var filtered []string
	for _, option := range options {
		if strings.HasPrefix(option, input) {
			filtered = append(filtered, option)
		}
	}
	return filtered
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.filtered = filterSuggestions(m.input.Value(), m.suggestions)
			if len(m.filtered) > 0 {
				m.index = 0
			} else {
				m.index = -1
			}
			return m, nil
		case "up":
			if m.index > 0 {
				m.index--
			}
			return m, nil
		case "down":
			if m.index < len(m.filtered)-1 {
				m.index++
			}
			return m, nil
		case "enter":
			// Allow selection only if a valid item is highlighted
			if m.index != -1 {
				m.selected = m.filtered[m.index]
				return m, tea.Quit
			}
			return m, nil
		case "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.filtered = filterSuggestions(m.input.Value(), m.suggestions)

	return m, cmd
}

func (m model) View() string {
	view := fmt.Sprintf("%s\n%s %s\n%s",
		titleStyle.Render("SSH Hostname Selector"),
		"Enter hostname:", inputStyle.Render(m.input.View()),
		hintStyle.Render("(Press TAB to filter, ARROW UP/DOWN to select, ENTER to confirm)"))

	for i, suggestion := range m.filtered {
		if i == m.index {
			view += fmt.Sprintf("\n> %s", selectedStyle.Render(suggestion)) // Highlight selected option
		} else {
			view += fmt.Sprintf("\n  %s", suggestion)
		}
	}

	return view
}
