package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmulholl/argo/v4"
	"github.com/dmulholl/ironclad/internal/config"
	"github.com/dmulholl/ironclad/internal/database"
	"github.com/dmulholl/ironclad/internal/fileio"
)

// Determine the database filename.
//  1. If a filename has been specified on the command line, use it.
//  2. Look for a cached filename in the config file.
//  3. Prompt the user to enter a filename.
func getDatabaseFilename(argParser *argo.ArgParser) (string, error) {
	filename := argParser.StringValue("file")

	if filename == "" {
		cached, found, err := config.Get("file")
		if err != nil {
			return "", fmt.Errorf("failed to check config file for cached filename: %w", err)
		}

		filename = cached
		if !found {
			filename = input("Database file: ")
		}
	}

	abspath, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for filename: %w", err)
	}

	if _, err := os.Stat(abspath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", abspath)
	}

	return abspath, nil
}

// Load a database from an encrypted file. Returns (masterpass, database, error).
func loadDB(filename string) (string, *database.DB, error) {
	if masterpass, success := getCachedPassword(filename); success {
		data, err := fileio.Load(filename, masterpass)
		if err == nil {
			db, err := database.FromJSON(data)
			if err != nil {
				return "", nil, fmt.Errorf("failed to unmarshall database: %w", err)
			}

			err = config.Set("file", filename)
			if err != nil {
				return "", nil, fmt.Errorf("failed to cache filename: %w", err)
			}

			setCachedPassword(filename, masterpass, db.CachePass)
			return masterpass, db, nil
		}
	}

	masterpass := inputPass("Master Password: ")

	data, err := fileio.Load(filename, masterpass)
	if err != nil {
		return "", nil, fmt.Errorf("failed to load database: %w", err)
	}

	db, err := database.FromJSON(data)
	if err != nil {
		return "", nil, fmt.Errorf("failed to unmarshall database: %w", err)
	}

	err = config.Set("file", filename)
	if err != nil {
		return "", nil, fmt.Errorf("failed to cache filename: %w", err)
	}

	setCachedPassword(filename, masterpass, db.CachePass)
	return masterpass, db, nil
}

// Encrypt and save a database file.
func saveDB(filename, password string, db *database.DB) {
	// Serialize the database as a byte-slice of JSON.
	json, err := db.ToJSON()
	if err != nil {
		exit(err)
	}

	// Encrypt the serialized database and write it to disk.
	err = fileio.Save(filename, password, json)
	if err != nil {
		exit(err)
	}
}
