package main

import (
	"fmt"
	"sort"

	"github.com/dmulholl/argo/v4"
)

var tagsCmdHelptext = `
Usage: ironclad tags

  Lists the tags in a database.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerTagsCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("tags")
	cmdParser.Helptext = tagsCmdHelptext
	cmdParser.Callback = tagsCmdCallback
	cmdParser.NewStringOption("file f", "")
}

func tagsCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	_, _, db := loadDB(cmdParser)

	tagmap := db.TagMap()

	tags := make([]string, 0)
	for tag := range tagmap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	if len(tags) > 0 {
		printLineOfChar("─")
		fmt.Println("  Tags")
		printLineOfChar("─")
		for _, tag := range tags {
			fmt.Printf("  %s [%d]\n", tag, len(tagmap[tag]))
		}
		printLineOfChar("─")
	} else {
		printLineOfChar("─")
		fmt.Println("  No Tags")
		printLineOfChar("─")
	}

	return nil
}
