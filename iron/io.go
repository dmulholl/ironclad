package main


import (
    "os"
    "github.com/dmulholland/clio/go/clio"
    "github.com/dmulholland/ironclad/irondb"
    "github.com/dmulholland/ironclad/ironio"
)


// Load and decrypt a database.
func loadDB(parser *clio.ArgParser) (password, filename string, db *irondb.DB) {

    // Determine the filename to use. First check for a filename specified on
    // the command line, next look for a cached filename from the application's
    // last run, if that fails prompt the user to enter a filename.
    filename = parser.GetStr("file")
    if filename == "" {
        var found bool
        if filename, found = getCachedFilename(); !found {
            filename = input("Filename: ")
        }
    }

    // Make sure the specified file exists.
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        exitf("'%v' does not exist", filename)
    }

    // If a password has been specified on the command line, try it.
    password = parser.GetStr("db-password")
    if password != "" {
        data, err := ironio.Load(password, filename)
        if err != nil {
            exit(err)
        }
        db, err := irondb.FromJSON(data)
        if err != nil {
            exit(err)
        }
        setCachedPassword(password)
        setCachedFilename(filename)
        return password, filename, db
    }

    // Look for a cached password from the application's last run. This
    // password may be invalid for the current file so if it fails prompt the
    // user to enter a new one.
    if password, found := getCachedPassword(); found {
        data, err := ironio.Load(password, filename)
        if err != nil {
            password = input("Password: ")
            data, err = ironio.Load(password, filename)
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
        return password, filename, db
    }

    // No command-line or cached password. Prompt the user to enter one.
    password = input("Password: ")
    data, err := ironio.Load(password, filename)
    if err != nil {
        exit(err)
    }
    db, err = irondb.FromJSON(data)
    if err != nil {
        exit(err)
    }
    setCachedPassword(password)
    setCachedFilename(filename)
    return password, filename, db
}


// Encrypt and save a database.
func saveDB(password, filename string, db *irondb.DB) {

    // Serialize the database as a byte-slice of JSON.
    json, err := db.ToJSON()
    if err != nil {
        exit(err)
    }

    // Encrypt the serialized database and write it to disk.
    err = ironio.Save(password, filename, json)
    if err != nil {
        exit(err)
    }
}
