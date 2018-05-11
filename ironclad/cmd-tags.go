package main


import "github.com/dmulholland/args"


import (
    "fmt"
    "os"
    "path/filepath"
    "sort"
)


var tagsHelp = fmt.Sprintf(`
Usage: %s tags [FLAGS] [OPTIONS]

  List the tags in a database.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerTags(parser *args.ArgParser) {
    cmd := parser.NewCmd("tags", tagsHelp, tagsCallback)
    cmd.NewString("file f")
}


func tagsCallback(parser *args.ArgParser) {

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
    if len(tags) > 0 {
        line("─")
        fmt.Println("  Tags")
        line("─")
        for _, tag := range tags {
            fmt.Printf("  %s [%d]\n", tag, len(tagmap[tag]))
        }
        line("─")
    } else {
        line("─")
        fmt.Println("  No Tags")
        line("─")
    }
}
