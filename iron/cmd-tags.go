package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "sort"
)


// Help text for the 'tags' command.
var tagsHelptext = fmt.Sprintf(`
Usage: %s tags [FLAGS] [OPTIONS]

  List the tags in a database.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'tags' command.
func tagsCallback(parser *clio.ArgParser) {

    // Load the database.
    _, _, db := loadDB(parser)

    // Assemble a map of tags.
    tagmap := db.TagMap()

    // Extract a sorted slice of tag strings.
    tags := make([]string, 0)
    for tag := range tagmap {
        tags = append(tags, tag)
    }
    sort.Strings(tags)

    // Print the tag list.
    line("-")
    fmt.Println("  Tags")
    line("-")
    for _, tag := range tags {
        fmt.Printf("  %s [%d]\n", tag, len(tagmap[tag]))
    }
    line("-")
}
