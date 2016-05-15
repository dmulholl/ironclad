package main


import (
    "fmt"
    "os"
    "strings"
    "path/filepath"
    "crypto/rand"
    "math/big"
    "github.com/dmulholland/clio/go/clio"
)


const (
    PoolLower = "abcdefghijklmnopqrstuvwxyz"
    PoolUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    PoolDigits = "0123456789"
    PoolSymbols = `!"$%^&*()_+=-[]{};'#:@~,./<>?|`
    PoolSimilars = "il1|oO0"
    DefaultLength = 24
)


// Help text for the 'gen' command.
var genHelptext = fmt.Sprintf(`
Usage: %s gen [FLAGS] ARGUMENTS

  Generate a random ASCII password.

  The default password length is 24 characters. The default character pool
  consists of the full range of uppercase letters, lowercase letters, digits,
  and symbols.

  The full list of possible symbols is:

    %s

  The --exclude-similar option excludes the following characters from the
  pool:

    %s

Arguments:
  [length]                  Length of the password. Default: 24 characters.

Flags:
  -d, --digits              Include digits [0-9].
  -e, --exclude-similar     Exclude similar characters.
      --help                Print this command's help text and exit.
  -l, --lowercase           Include lowercase letters [a-z].
  -r, --readable            Add a space after every fourth character.
  -s, --symbols             Include symbols.
  -u, --uppercase           Include uppercase letters [A-Z].
`, filepath.Base(os.Args[0]), PoolSymbols, PoolSimilars)


// Callback for the 'gen' command.
func genCallback(parser *clio.ArgParser) {

    var length int
    var pool string

    // Has a password length been specified?
    if parser.HasArgs() {
        length = parser.GetArgsAsInts()[0]
    } else {
        length = DefaultLength
    }

    // Assemble the character pool.
    if parser.GetFlag("digits") {
        pool += PoolDigits
    }
    if parser.GetFlag("lowercase") {
        pool += PoolLower
    }
    if parser.GetFlag("symbols") {
        pool += PoolSymbols
    }
    if parser.GetFlag("uppercase") {
        pool += PoolUpper
    }

    // Default to using all character classes if no options were specified.
    if pool == "" {
        pool = PoolDigits + PoolLower + PoolSymbols + PoolUpper
    }

    // Are we excluding similar characters from the pool?
    if parser.GetFlag("exclude-similar") {
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
    password := string(passBytes)

    // Add spaces if required.
    if parser.GetFlag("readable") {
        password = addSpaces(password)
    }

    // Print the password to stdout.
    fmt.Print(password)

    // Only print a newline if we're connected to a terminal.
    if stdoutIsTerminal() {
        fmt.Println()
    }
}


// Generate a random integer in the uniform range [0, max).
func randInt(max int) int {
    n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
    if err != nil {
        exit("Error:", err)
    }
    return int(n.Int64())
}
