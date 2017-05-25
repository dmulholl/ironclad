package main


import "github.com/dmulholland/clio/go/clio"


import (
    "fmt"
    "os"
    "strings"
    "path/filepath"
    "crypto/rand"
    "math/big"
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
var genHelp = fmt.Sprintf(`
Usage: %s gen [FLAGS] ARGUMENTS

  Generate a random ASCII password. The password is automatically copied to
  the system clipboard. The password can alternatively be printed to stdout.

  The default password length is 24 characters. The default character pool
  consists of uppercase letters, lowercase letters, and digits.

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
      --help                Print this command's help text and exit.
  -p, --print               Print the password to stdout.
  -r, --readable            Add spaces for readability.
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

    // Set the default pool if no options were specified.
    if pool == "" {
        pool = PoolDigits + PoolLower + PoolUpper
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
    if parser.GetFlag("print") {
        fmt.Print(password)
        if stdoutIsTerminal() {
            fmt.Println()
        }
        return
    }

    // Copy the password to the clipboard.
    writeToClipboard(password)
}


// Generate a random integer in the uniform range [0, max).
func randInt(max int) int {
    n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
    if err != nil {
        exit(err)
    }
    return int(n.Int64())
}
