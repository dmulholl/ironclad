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

// Reads a single line of input from stdin. Prints the prompt to stderr to avoid polluting stdout.
func input(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)

	input, err := stdinReader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return strings.TrimSpace(input), nil
		}
		return "", err
	}

	return strings.TrimSpace(input), nil
}

// Reads masked input from stdin.
func inputMasked(prompt string) (string, error) {
	input, err := gopass.GetPasswdPrompt(prompt, true, os.Stdin, os.Stderr)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(input)), nil
}

// Launch a text editor and capture its output.
func inputViaEditor(template string) (string, error) {
	file, err := os.CreateTemp("", "ironclad-input")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	defer os.Remove(file.Name())

	err = file.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close temporary file: %w", err)
	}

	err = os.Chmod(file.Name(), 0600)
	if err != nil {
		return "", fmt.Errorf("failed to set file permissions: %w", err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	_, err = exec.LookPath(editor)
	if err != nil {
		return "", fmt.Errorf("failed to check for editor binary: %w", err)
	}

	err = os.WriteFile(file.Name(), []byte(template), 0600)
	if err != nil {
		return "", fmt.Errorf("failed to write template: %w", err)
	}

	cmd := exec.Command(editor, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run editor: %w", err)
	}

	output, err := os.ReadFile(file.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read editor output: %w", err)
	}

	return string(output), nil
}
