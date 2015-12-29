/*
    Package ironconfig handles reading from and writing to the application's
    configuration file.
*/
package ironconfig


import (
    "github.com/pelletier/go-toml"
    "os"
    "io/ioutil"
    "path/filepath"
)


// Location of the configuration file.
var Configfile string


// Load a config file's TOML content.
func load() (*toml.TomlTree, error) {
    if _, err := os.Stat(Configfile); err == nil {
        tree, err := toml.LoadFile(Configfile)
        if err != nil {
            return nil, err
        }
        return tree, nil
    } else {
        return toml.TreeFromMap(make(map[string]interface{})), nil
    }
}


// Save a TOML tree to file.
func save(tree *toml.TomlTree) error {
    err := os.MkdirAll(filepath.Dir(Configfile), 0777)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(Configfile, []byte(tree.ToString()), 0600)
}


// Get reads a value from the configuration file.
func Get(key string) (found bool, value string, err error) {
    config, err := load()
    if err != nil {
        return false, "", err
    }
    if config.Has(key) {
        if value, ok := config.Get(key).(string); ok {
            return true, value, nil
        }
    }
    return false, "", nil
}


// Set writes a value to the configuration file.
func Set(key, value string) error {
    config, err := load()
    if err != nil {
        return err
    }
    config.Set(key, value)
    return save(config)
}
