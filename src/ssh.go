package main

import (
	"fmt"
	hostsfile "github.com/kevinburke/hostsfile/lib"
	"github.com/kevinburke/ssh_config"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	options []string // ssh hosts selection options
)

func GetHosts() []string {
	var hosts []string
	if len(IniServers) > 0 {
		for _, key := range IniServers {
			hosts = append(hosts, key.Hostname)
		}
	}
	if Settings.LoadsshHosts {
		hosts = append(hosts, GetsshConfigHosts()...)
	}
	if Settings.LoadetcHosts {
		err, etchosts := GetetcHosts()
		if err == nil && len(etchosts) > 0 {
			hosts = append(hosts, etchosts...)
		}
	}
	return hosts
}

func GetsshConfigHosts() []string {
	var hosts []string
	f, _ := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	cfg, _ := ssh_config.Decode(f)
	for _, host := range cfg.Hosts {
		// i know this is ugly hack, but sadly the library itself does not have a clear method to list all hosts from ssh config file...
		if !strings.Contains(fmt.Sprintf("%s", host.Patterns), "*") {
			re := regexp.MustCompile(`\[(.*?)\]`)
			matches := re.FindStringSubmatch(fmt.Sprintf("%s", host.Patterns))
			if len(matches) > 1 {
				hosts = append(hosts, matches[1])
			}
		}
	}
	return hosts
}

func GetetcHosts() (error, []string) {
	var hosts []string
	f, err := os.Open("/etc/hosts")
	if err != nil {
		fmt.Printf("Unable to open /etc/hosts file: %s\n", err)
		return err, hosts
	}
	h, err := hostsfile.Decode(f)
	if err != nil {
		fmt.Printf("Unable to decode /etc/hosts file: %s\n", err)
		return err, hosts
	}
	// the different library for host extraction but from the same developer and almost the same problems, again there are no normal way to return all elements, just ip addresses
	// but no dns name, maybe because the one ip can have more than one, ok... when we will return the first one..
	for _, host := range h.Records() {
		var hostname string
		for key := range host.Hostnames {
			hostname = key
			break // get only first one
		}
		// skip local and invalid hosts
		matched, err := regexp.MatchString(`localhost|broadcasthost`, hostname)
		if err != nil {
			continue
		}
		if matched {
			continue
		}
		if hostname != "" {
			hosts = append(hosts, hostname)
		}

	}
	return nil, hosts
}

// parseSSHHost extracts user, host, and port from the input string
func parseSSHHost(input string) (string, string, string) {
	// Default values
	user := "" // Default SSH user
	host := ""
	port := "" // Default SSH port (empty means use system default)

	// Regex pattern to match [user@]host[:port]
	re := regexp.MustCompile(`^(?:(\w+)@)?([\w\.-]+)(?::(\d+))?$`)
	matches := re.FindStringSubmatch(input)

	if len(matches) == 0 {
		fmt.Println("Invalid SSH host format")
		return "", "", ""
	}

	// Extract user, host, and port
	if matches[1] != "" {
		user = matches[1]
	}
	host = matches[2]
	if matches[3] != "" {
		port = matches[3]
	}

	return user, host, port
}

func sshHost(hostname string) error {
	// if host is already in the local conf file then load it's value...
	err, hostval := FindIniHost(hostname)
	if err == nil && hostval.HostDetails != "" {
		hostname = hostval.HostDetails
	}

        
        sshConfig :=filepath.Join(os.Getenv("HOME"), ".ssh", "config")
        


	user, host, port := parseSSHHost(hostname)
	// Construct the SSH command arguments
	sshArgs := []string{}
	if port != "" {
		sshArgs = append(sshArgs, "-p", port)
	}
        // check if user's ssh config exist and then pass paramater to ssh to load it
        if _, err := os.Stat(sshConfig); err == nil {
                 sshArgs = append(sshArgs, "-F",sshConfig)
        }
        //sshArgs = append(sshArgs, "-v")
        if user != "" {
                 sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", user, host))
        } else {
                 sshArgs = append(sshArgs, host)
        }
        if port != "" {
                 sshArgs = append(sshArgs, "-p",port)
        }
	fmt.Printf("\nEstablishing connection to: %s\n\n", hostname)
	cmd := exec.Command("ssh", sshArgs...)
	// Attach standard input/output
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Run SSH interactively
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("SSH error: %s", err)
	}

	return nil
}
