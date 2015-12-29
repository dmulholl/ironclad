package main


import (
    "fmt"
    "bufio"
    "strings"
    "os"
    "os/exec"
    "path/filepath"
    "golang.org/x/crypto/ssh/terminal"
    "io/ioutil"
)


// Reader for reading user input from stdin.
var stdinReader = bufio.NewReader(os.Stdin)


// Read a single line of input from stdin.
// We print the prompt to stderr to avoid polluting stdout with the password
// prompt. This means that the output of the dump and export commands can be
// printed to stdout by default and cleanly piped to files when required.
func input(prompt string) string {
    fmt.Fprint(os.Stderr, prompt)

    input, err := stdinReader.ReadString('\n')
    if err != nil {
        exit("Error reading input from stdin.")
    }

    return strings.TrimSpace(input)
}


// Launch a text editor and capture its output.
func inputViaEditor(filename, text string) string {

    // Set the filename for the editor to open.
    filename = filepath.Join(ironpath, filename)

    // Create a file for the editor to open.
    os.MkdirAll(ironpath, 0777)
    err := ioutil.WriteFile(filename, []byte(text), 0600)
    if err != nil {
        exit("Error:", err)
    }

    // Determine the editor to use.
    editor := os.Getenv("EDITOR")
    if editor == "" {
        editor = "vim"
    }
    _, err = exec.LookPath(editor)
    if err != nil {
        exit("Error: cannot locate text editor '" + editor + "'.")
    }

    // Launch the editor and wait for it to complete.
    cmd := exec.Command(editor, filename)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err = cmd.Run()
    if err != nil {
        exit("Error:", err)
    }

    // Load the editor's output.
    input, err := ioutil.ReadFile(filename)
    if err != nil {
        exit("Error:", err)
    }

    // Delete the input file.
    err = os.Remove(filename)
    if err != nil {
        exit("Error:", err)
    }

    return string(input)
}


// Returns true if stdout is connected to a terminal.
func stdoutIsTerminal() bool {
    return terminal.IsTerminal(int(os.Stdout.Fd()))
}


// Exit with an optional error message.
func exit(msg ...interface{}) {
    if len(msg) > 0 {
        elements := make([]string, 0)
        for _, element := range msg {
            elements = append(elements, fmt.Sprint(element))
        }
        output := strings.Join(elements, " ")
        if !strings.HasSuffix(output, ".") {
            output = output + "."
        }
        fmt.Fprintln(os.Stderr, output)
        os.Exit(1)
    } else {
        os.Exit(0)
    }
}


// Print a line of characters.
func line(char string) {
    for i := 0; i < 80; i++ {
        fmt.Print(char)
    }
    fmt.Println()
}


// Print an indented line of characters.
func iline(char string, indent int) {
    for i := 0; i < indent; i++ {
        fmt.Print(" ")
    }
    for i := 0; i < 80 - indent * 2; i++ {
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
