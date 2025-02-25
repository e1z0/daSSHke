package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

var Settings Config
var IniServers []IniHost

type Config struct {
	LoadetcHosts bool
	LoadsshHosts bool
	AutoAddHosts bool
	Sync         bool
	GistID       string
	GistSecret   string
}

type IniHost struct {
	Hostname    string
	HostDetails string
}

func returnConfPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "daSSHke", "settings.ini")
}

func FindIniHost(host string) (error, IniHost) {
	hst := IniHost{}
	for _, key := range IniServers {
		if key.Hostname == host {
			return nil, key
		}
	}
	return fmt.Errorf("Unable to find server in conf"), hst
}

func IniWrServers(servers []IniHost) error {
	f := returnConfPath()
	cfg, err := ini.LooseLoad(f) // Load file, create if missing
	if err != nil {
		return err
	}

	// Get or create the [servers] section
	section, err := cfg.GetSection("servers")
	if err != nil {
		section, _ = cfg.NewSection("servers")
	}

	for _, key := range section.Keys() {
		section.DeleteKey(key.Name())
	}
	for _, server := range servers {
		section.Key(server.Hostname).SetValue(server.HostDetails)
	}
	// Save changes back to the file
	err = cfg.SaveTo(f)
	if err != nil {
		return err
	}
	return nil
}

func AddHost(hostname string) {
	_, host, _ := parseSSHHost(hostname)

	if host != "" {
		f := returnConfPath()
		cfg, err := ini.LooseLoad(f) // Load file, create if missing
		if err != nil {
			fmt.Println("Error loading conf file: ", err)
			return
		}

		// Get or create the [servers] section
		section, err := cfg.GetSection("servers")
		if err != nil {
			section, _ = cfg.NewSection("servers")
		}

		// Set a new key-value pair
		section.Key(host).SetValue(hostname)
		// Save changes back to the file
		err = cfg.SaveTo(f)
		if err != nil {
			fmt.Println("Error saving conf file: ", err)
			return
		}
	}
}

func readSettings() error {
	f := returnConfPath()
	cfg, err := ini.Load(f)
	if err != nil {
		return err
	}
	Settings = Config{
		LoadetcHosts: cfg.Section("General").Key("loadetchosts").MustBool(false),
		LoadsshHosts: cfg.Section("General").Key("loadsshhosts").MustBool(true),
		AutoAddHosts: cfg.Section("General").Key("autoaddhosts").MustBool(true),
		Sync:         cfg.Section("General").Key("sync").MustBool(false),
		GistID:       cfg.Section("General").Key("gistid").MustString(""),
		GistSecret:   cfg.Section("General").Key("gistsecret").MustString(""),
	}
	section := cfg.Section("servers")
	for _, key := range section.Keys() {
		iniServ := IniHost{Hostname: key.Name(), HostDetails: key.Value()}
		IniServers = append(IniServers, iniServ)
	}

	return nil
}
