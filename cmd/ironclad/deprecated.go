package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/howeyc/gopass"
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
	input, err := io.ReadAll(stdinReader)
	if err != nil {
		exit(err)
	}
	return string(input)
}

// Launch a text editor and capture its output.
func inputViaEditor(template string) string {
	file, err := os.CreateTemp("", "ironclad-input")
	if err != nil {
		exit(err)
	}

	defer os.Remove(file.Name())

	err = file.Close()
	if err != nil {
		exit(err)
	}

	err = os.Chmod(file.Name(), 0600)
	if err != nil {
		exit(err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	_, err = exec.LookPath(editor)
	if err != nil {
		exit(err)
	}

	// Write the template to the temporary file.
	err = os.WriteFile(file.Name(), []byte(template), 0600)
	if err != nil {
		exit(err)
	}

	// Launch the editor and wait for it to complete.
	cmd := exec.Command(editor, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		exit(err)
	}

	// Load the editor's output.
	output, err := os.ReadFile(file.Name())
	if err != nil {
		exit(err)
	}

	return string(output)
}

// Exit with an error message and non-zero error code.
func exit(message any) {
	fmt.Fprintf(os.Stderr, "error: %s\n", message)
	os.Exit(1)
}
