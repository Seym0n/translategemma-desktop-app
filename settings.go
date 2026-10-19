package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Settings are the user preferences persisted between app launches.
type Settings struct {
	Model      string `json:"model"`
	SourceCode string `json:"sourceCode"`
	TargetCode string `json:"targetCode"`
}

func settingsPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "TranslateGemma", "settings.json"), nil
}

// loadSettings reads the saved settings; a missing or unreadable file yields zero values.
func loadSettings() Settings {
	var s Settings
	path, err := settingsPath()
	if err != nil {
		return s
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return s
	}
	_ = json.Unmarshal(data, &s)
	return s
}

func saveSettings(s Settings) error {
	path, err := settingsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
