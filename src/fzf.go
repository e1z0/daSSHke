package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
	"strings"
)

var NoBellStdout = &noBellStdout{}

type noBellStdout struct{}

func (n *noBellStdout) Write(p []byte) (int, error) {
	if len(p) == 1 && p[0] == readline.CharBell {
		return 0, nil
	}
	return readline.Stdout.Write(p)
}

func (n *noBellStdout) Close() error {
	return readline.Stdout.Close()
}

// Function to filter SSH hosts using fuzzy search
func searchHosts(input string) []string {
	if input == "" {
		return options // Show all hosts if no input
	}

	var results []string
	for _, host := range options {
		if fuzzy.Match(input, host) {
			results = append(results, host)
		}
	}

	// If no matches, return something helpful
	if len(results) == 0 {
		return []string{"No matches found"}
	}

	return results
}

func fz() {
	fmt.Println("🔍 Search SSH Host (Type to filter, use arrow keys to select)")

	templates := &promptui.SelectTemplates{
		Label: "{{ . }}?",
		//Active:   "\U0001F336  ({{ . | cyan }})",
		Active:   "🔥 ({{ . | cyan }})",
		Inactive: "  {{ . | cyan }}",
		Selected: "🔥 {{ . | red | cyan }}",
	}

	// Get filtered results based on user input
	prompt := promptui.Select{
		Stdout:    NoBellStdout,
		Label:     "Select Host",
		Size:      20,
		Templates: templates,
		Items:     options, // Start with full list
		Searcher: func(input string, index int) bool {
			item := options[index]
			return strings.Contains(item, input) || fuzzy.Match(input, item)
		},
	}

	_, selectedHost, err := prompt.Run()
	if err != nil {
		fmt.Println("❌ Selection canceled")
		return
	}

	fmt.Println("✅ Selected Host:", selectedHost)
	sshHost(selectedHost)

}
