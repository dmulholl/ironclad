package main

import (
	"os"
	"path/filepath"

	"github.com/dmulholl/argo/v4"
	"github.com/dmulholl/ironclad/internal/database"
	"github.com/dmulholl/ironclad/internal/fileio"
)

// Load a database from an encrypted file.
func loadDB(args *argo.ArgParser) (filename, masterpass string, db *database.DB) {
	// Determine the file to use.
	// 1. Has a filename been specified on the command line?
	// 2. Look for a cached filename.
	// 3. Prompt the user to enter a filename.
	filename = args.StringValue("file")
	if filename == "" {
		var found bool
		if filename, found = getCachedFilename(); !found {
			filename = input("Database File: ")
		}
	}
	filename, err := filepath.Abs(filename)
	if err != nil {
		exit(err)
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exit(err)
	}

	// Look for a cached master password from the application's last run. This
	// password may be invalid for the current file so if it fails prompt
	// the user to enter a new one.
	if masterpass, success := getCachedPassword(filename); success {
		data, err := fileio.Load(filename, masterpass)
		if err != nil {
			println("Error: the cached password was invalid for the database.")
			masterpass = inputPass("Master Password: ")
			data, err = fileio.Load(filename, masterpass)
			if err != nil {
				exit(err)
			}
		}
		db, err := database.FromJSON(data)
		if err != nil {
			exit(err)
		}
		setCachedPassword(filename, masterpass, db.CachePass)
		setCachedFilename(filename)
		return filename, masterpass, db
	}

	// No cached password. Prompt the user to enter one.
	masterpass = inputPass("Master Password: ")
	data, err := fileio.Load(filename, masterpass)
	if err != nil {
		exit(err)
	}
	db, err = database.FromJSON(data)
	if err != nil {
		exit(err)
	}
	setCachedPassword(filename, masterpass, db.CachePass)
	setCachedFilename(filename)
	return filename, masterpass, db
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
