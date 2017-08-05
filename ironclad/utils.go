package main


import "github.com/howeyc/gopass"
import "golang.org/x/crypto/ssh/terminal"


import (
    "fmt"
    "bufio"
    "strings"
    "os"
    "os/exec"
    "path/filepath"
    "io/ioutil"
)


import (
    "github.com/dmulholland/ironclad/ironconfig"
)


// Reader for reading user input from stdin.
var stdinReader = bufio.NewReader(os.Stdin)


// Read a single line of input from stdin. We print the prompt to stderr to
// avoid polluting stdout with the password prompt. This means that the output
// of the dump and export commands can be printed to stdout by default and
// cleanly piped to files when required.
func input(prompt string) string {
    fmt.Fprint(os.Stderr, prompt)
    input, err := stdinReader.ReadString('\n')
    if err != nil {
        exit(err)
    }
    return strings.TrimSpace(input)
}


// Read a masked password from stdin.
func inputPass(prompt string) string {
    fmt.Fprint(os.Stderr, prompt)
    bytes, err := gopass.GetPasswdMasked()
    if err != nil {
        exit(err)
    }
    return strings.TrimSpace(string(bytes))
}


// Read from stdin until EOF.
func inputViaStdin() string {
    input, err := ioutil.ReadAll(stdinReader)
    if err != nil {
        exit(err)
    }
    return string(input)
}


// Launch a text editor and capture its output.
func inputViaEditor(file, template string) string {

    // Set the file for the editor to open.
    file = filepath.Join(ironconfig.ConfigDir, file)

    // Create a file for the editor to open.
    os.MkdirAll(ironconfig.ConfigDir, 0777)
    err := ioutil.WriteFile(file, []byte(template), 0600)
    if err != nil {
        exit(err)
    }

    // Determine the editor to use.
    editor := os.Getenv("EDITOR")
    if editor == "" {
        editor = "vim"
    }
    _, err = exec.LookPath(editor)
    if err != nil {
        exit("cannot locate text editor:", editor)
    }

    // Launch the editor and wait for it to complete.
    cmd := exec.Command(editor, file)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err = cmd.Run()
    if err != nil {
        exit(err)
    }

    // Load the editor's output.
    input, err := ioutil.ReadFile(file)
    if err != nil {
        exit(err)
    }

    // Delete the input file.
    err = os.Remove(file)
    if err != nil {
        exit(err)
    }

    return string(input)
}


// Returns true if stdout is connected to a terminal.
func stdoutIsTerminal() bool {
    return terminal.IsTerminal(int(os.Stdout.Fd()))
}


// Exit with an error message and non-zero error code.
func exit(objects ...interface{}) {
    fmt.Fprint(os.Stderr, "Error:")
    for _, obj := range objects {
        fmt.Fprintf(os.Stderr, " %v", obj)
    }
    fmt.Fprint(os.Stderr, "\n")
    os.Exit(1)
}


// Returns the width of the terminal window.
func terminalWidth() int {
    if terminal.IsTerminal(int(os.Stdout.Fd())) {
        width, _, err := terminal.GetSize(int(os.Stdout.Fd()))
        if err == nil {
            return width
        }
    }
    return 80
}


// Print a line of characters.
func line(char string) {
    length := terminalWidth()
    for i := 0; i < length; i++ {
        fmt.Print(char)
    }
    fmt.Println()
}


// Print an indented line of characters.
func indentedLine(char string) {
    length := terminalWidth() - 4
    fmt.Print("  ")
    for i := 0; i < length; i++ {
        fmt.Print(char)
    }
    fmt.Println()
}


// Returns the set of strings that are in slice1 but not in slice2.
func diff(slice1, slice2 []string) []string {
    diff := make([]string, 0)

    for _, s1 := range slice1 {
        found := false
        for _, s2 := range slice2 {
            if s1 == s2 {
                found = true
                break
            }
        }
        if !found {
            diff = append(diff, s1)
        }
    }

    return diff
}


// Add spaces to a string after every fourth character.
func addSpaces(input string) string {
    words := make([]string, 0)
    runes := []rune(input)

    for i := 0; i < len(runes); i += 4 {
        if i + 4 > len(runes) {
            words = append(words, string(runes[i:]))
        } else {
            words = append(words, string(runes[i:i+4]))
        }
    }

    return strings.Join(words, "  ")
}


// Returns a string of the specified length.
func charstr(length int, char rune) string {
    runes := make([]rune, length)
    for i := range runes {
        runes[i] = char
    }
    return string(runes)
}
