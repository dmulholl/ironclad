/*
    Package ironconfig provides read and write access to the application's
    TOML configuration file.
*/
package ironconfig


import "github.com/pelletier/go-toml"


import (
    "os"
    "io/ioutil"
    "path/filepath"
)


// Location of the configuration file.
var ConfigDir = filepath.Join(os.Getenv("HOME"), ".config", "ironclad")
var ConfigFile = filepath.Join(ConfigDir, "goconfig.toml")


// Get reads a value from the configuration file.
func Get(key string) (value string, found bool, err error) {
    config, err := loadToml()
    if err != nil {
        return "", false, err
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
        return err
    }
    config.Set(key, value)
    return saveToml(config)
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
        tree, err := toml.LoadFile(ConfigFile)
        if err != nil {
            return nil, err
        }
        return tree, nil
    } else {
        tree, err := toml.TreeFromMap(make(map[string]interface{}))
        return tree, err
    }
}


// Save a TOML tree to file.
func saveToml(tree *toml.Tree) error {
    err := os.MkdirAll(filepath.Dir(ConfigFile), 0777)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(ConfigFile, []byte(tree.String()), 0600)
}
