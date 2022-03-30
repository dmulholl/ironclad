package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo"
)

const (
	PoolLower     = "abcdefghijklmnopqrstuvwxyz"
	PoolUpper     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	PoolDigits    = "0123456789"
	PoolSymbols   = `!"$%^&*()_+=-[]{};'#:@~,./<>?|`
	PoolSimilars  = "il1|oO0"
	DefaultLength = 24
)

var genHelp = fmt.Sprintf(`
Usage: %s gen [length]

  Generates a random ASCII password. The password is automatically copied to
  the system clipboard. The password can alternatively be printed to stdout.

  The default password length is 24 characters. The default character pool
  consists of uppercase letters, lowercase letters, digits, and symbols.

  The full list of possible symbols is:

    %s

  The --exclude-similar option excludes the following characters from the
  pool:

    %s

Arguments:
  [length]                  Length of the password. Default: 24 characters.

Character Flags:
  -d, --digits              Include digits [0-9].
  -l, --lowercase           Include lowercase letters [a-z].
  -s, --symbols             Include symbols.
  -u, --uppercase           Include uppercase letters [A-Z].

Flags:
  -x, --exclude-similar     Exclude similar characters.
  -h, --help                Print this command's help text and exit.
  -p, --print               Print the password to stdout.
  -r, --readable            Add spaces for readability.
`, filepath.Base(os.Args[0]), PoolSymbols, PoolSimilars)

func registerGenCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("gen")
	cmdParser.Helptext = genHelp
	cmdParser.Callback = genCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewFlag("digits d")
	cmdParser.NewFlag("exclude-similar x")
	cmdParser.NewFlag("lowercase l")
	cmdParser.NewFlag("symbols s")
	cmdParser.NewFlag("uppercase u")
	cmdParser.NewFlag("readable r")
	cmdParser.NewFlag("print p")
}

func genCallback(cmdName string, cmdParser *argo.ArgParser) {
	var length int
	if cmdParser.HasArgs() {
		length = int(cmdParser.ArgsAsInts()[0])
	} else {
		length = DefaultLength
	}

	password := genPassword(
		length,
		cmdParser.Found("uppercase"),
		cmdParser.Found("lowercase"),
		cmdParser.Found("symbols"),
		cmdParser.Found("digits"),
		cmdParser.Found("exclude-similar"))

	if cmdParser.Found("readable") {
		password = addSpaces(password)
	}

	if cmdParser.Found("print") {
		fmt.Print(password)
		if stdoutIsTerminal() {
			fmt.Println()
		}
		return
	}

	writeToClipboard(password)
}

// Generate a new password.
func genPassword(length int, upper, lower, symbols, digits, excludeSimilar bool) string {

	// Assemble the character pool.
	var pool string
	if digits {
		pool += PoolDigits
	}
	if lower {
		pool += PoolLower
	}
	if symbols {
		pool += PoolSymbols
	}
	if upper {
		pool += PoolUpper
	}

	// Use the default pool if no options were specified.
	if pool == "" {
		pool = PoolDigits + PoolLower + PoolUpper
	}

	// Are we excluding similar characters from the pool?
	if excludeSimilar {
		newpool := make([]rune, 0)
		for _, r := range pool {
			if !strings.ContainsRune(PoolSimilars, r) {
				newpool = append(newpool, r)
			}
		}
		pool = string(newpool)
	}

	// Assemble the password.
	passBytes := make([]byte, length)
	for i := range passBytes {
		passBytes[i] = pool[randInt(len(pool))]
	}

	return string(passBytes)
}

// Generate a random integer in the uniform range [0, max).
func randInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		exit(err)
	}
	return int(n.Int64())
}
