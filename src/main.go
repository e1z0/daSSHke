package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"os"
)

var (
	SHOW_MENU = false
	version   = "0.0"
)

func main() {
	// read settings
	err := readSettings()
	if err != nil {
		fmt.Printf("There are no settings saved, please read the manual for the application customization!\n")
	}
	options = GetHosts()
	// command line options
	menu := flag.Bool("m", false, "Show menu")
	conn := flag.String("c", "", "Connect to host, syntax user@host:port or host, example.: root@example.com:2222")
	var push, pull *bool
	if Settings.Sync {
		push = flag.Bool("push", false, "Push server list changes to the github gist")
		pull = flag.Bool("pull", false, "Pull server list changes from the github gist")
	}
	ver := flag.Bool("v", false, "Show version")
	flag.Parse()

	if *ver {
		fmt.Printf("daSSHke %s\nCopyright (c) 2025 Justinas K. (e1z0@icloud.com)\n", version)
		os.Exit(0)
	}

	if *push {
		err = uploadToGist()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		os.Exit(0)
	}
	if *pull {
		err = downloadFromGist()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		os.Exit(0)
	}

	if *conn != "" {
		// add host to internal list if enabled
		if Settings.AutoAddHosts {
			AddHost(*conn)
		}
		sshHost(*conn)
		os.Exit(0)
	}

	if *menu {
		Settings.ShowMenu = true
	}

	if Settings.ShowMenu {
		p := tea.NewProgram(initialModel())

		finalModel, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			os.Exit(1)
		}

		// Extract selected hostname and print it
		if m, ok := finalModel.(model); ok && m.selected != "" {
			sshHost(m.selected)
		}

	} else {

		p := tea.NewProgram(initialAutocompleteModel())

		finalModel, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			os.Exit(1)
		}

		// Extract selected hostname and print it
		if m, ok := finalModel.(autocompletemodel); ok && m.selected != "" {
			sshHost(m.selected)
		}

	}
}
