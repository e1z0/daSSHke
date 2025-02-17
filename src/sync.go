package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var ()

const (
	GistAPI = "https://api.github.com/gists"
)

// GistRequest represents the JSON payload for updating a Gist
type GistRequest struct {
	Description string              `json:"description"`
	Files       map[string]GistFile `json:"files"`
}

type GistFile struct {
	Content string `json:"content"`
}

// Upload servers list to GitHub Gist
func uploadToGist() error {
	gistID := Settings.GistID
	token := Settings.GistSecret

	if gistID == "" || token == "" {
		return fmt.Errorf("gistid or gistsecret not set")
	}

	serversList, err := json.Marshal(IniServers)
	if err != nil {
		return err
	}

	// Prepare JSON payload
	payload := GistRequest{
		Description: "daSSHke servers list",
		Files: map[string]GistFile{
			"servers.json": {Content: string(serversList)},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/%s", GistAPI, gistID)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API error: %s", string(body))
	}

	fmt.Println("✅ Servers list pushed to GitHub Gist successfully!")
	return nil
}

func downloadFromGist() error {
	gistID := Settings.GistID
	token := Settings.GistSecret

	if gistID == "" || token == "" {
		return fmt.Errorf("gistid or gistsecret not set")
	}

	// Request Gist
	url := fmt.Sprintf("%s/%s", GistAPI, gistID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API error: %s", string(body))
	}

	// Parse JSON
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// Extract content
	files := result["files"].(map[string]interface{})
	serversData := files["servers.json"].(map[string]interface{})["content"].(string)
	var servers []IniHost
	err = json.Unmarshal([]byte(serversData), &servers)
	if err != nil {
		return err
	}
	err = IniWrServers(servers)
	if err != nil {
		return err
	}

	fmt.Println("✅ Servers list pulled from GitHub gist successfully!")
	return nil
}
