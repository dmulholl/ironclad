package main


import "github.com/dmulholland/clio/go/clio"


import (
    "os"
)


import (
    "github.com/dmulholland/ironclad/irondb"
    "github.com/dmulholland/ironclad/ironio"
)


// Load a database from an encrypted file.
func loadDB(args *clio.ArgParser) (filename, password string, db *irondb.DB) {

    // Determine the file to use. First check for a file specified on
    // the command line, next look for a cached filename from the
    // application's last run, if that fails prompt the user to enter a
    // filename.
    filename = args.GetStr("file")
    if filename == "" {
        var found bool
        if filename, found = getCachedFilename(); !found {
            filename = input("File: ")
        }
    }

    // Make sure the specified file exists.
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        exitfmt("'%v' does not exist", filename)
    }

    // Look for a cached password from the application's last run. This
    // password may be invalid for the current file so if it fails prompt the
    // user to enter a new one.
    if password, found := getCachedPassword(); found {
        data, err := ironio.Load(filename, password)
        if err != nil {
            password = inputPass("Password: ")
            data, err = ironio.Load(filename, password)
            if err != nil {
                exit(err)
            }
        }
        db, err := irondb.FromJSON(data)
        if err != nil {
            exit(err)
        }
        setCachedPassword(password)
        setCachedFilename(filename)
        return filename, password, db
    }

    // No command-line or cached password. Prompt the user to enter one.
    password = inputPass("Password: ")
    data, err := ironio.Load(filename, password)
    if err != nil {
        exit(err)
    }
    db, err = irondb.FromJSON(data)
    if err != nil {
        exit(err)
    }
    setCachedPassword(password)
    setCachedFilename(filename)
    return filename, password, db
}


// Encrypt and save a database file.
func saveDB(file, password string, db *irondb.DB) {

    // Serialize the database as a byte-slice of JSON.
    json, err := db.ToJSON()
    if err != nil {
        exit(err)
    }

    // Encrypt the serialized database and write it to disk.
    err = ironio.Save(file, password, json)
    if err != nil {
        exit(err)
    }
}
