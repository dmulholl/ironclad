package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pelletier/go-toml"
)

// Filepath of the TOML configuration file.
var ConfigFile string

// Initialize the path to the config file depending on the OS.
func init() {
	if runtime.GOOS == "windows" {
		ConfigFile = filepath.Join(os.Getenv("LOCALAPPDATA"), "Ironclad", "goconfig.toml")
	} else {
		ConfigFile = filepath.Join(os.Getenv("HOME"), ".config", "ironclad", "goconfig.toml")
	}
}

// Get reads a value from the configuration file.
func Get(key string) (value string, found bool, err error) {
	config, err := loadToml()
	if err != nil {
		return "", false, fmt.Errorf("failed to load config file: %w", err)
	}

	if config.Has(key) {
		if value, ok := config.Get(key).(string); ok {
			return value, true, nil
		}
	}

	return "", false, nil
}

// Set writes a value to the configuration file.
func Set(key, value string) error {
	config, err := loadToml()
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	config.Set(key, value)

	if err := saveToml(config); err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}

	return nil
}

// Delete deletes an entry from the configuration file.
func Delete(key string) error {
	config, err := loadToml()
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	if config.Has(key) {
		if err := config.Delete(key); err != nil {
			return fmt.Errorf("failed to delete entry from config file: %w", err)
		}
	}

	if err := saveToml(config); err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}

	return nil
}

// FileExists returns true if the configuration file exists.
func FileExists() bool {
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// Load a config file's TOML content.
func loadToml() (*toml.Tree, error) {
	if FileExists() {
		return toml.LoadFile(ConfigFile)
	}
	return toml.TreeFromMap(make(map[string]interface{}))
}

// Save a TOML tree to file.
func saveToml(tree *toml.Tree) error {
	err := os.MkdirAll(filepath.Dir(ConfigFile), 0777)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return os.WriteFile(ConfigFile, []byte(tree.String()), 0600)
}
