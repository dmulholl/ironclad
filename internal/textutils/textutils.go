package textutils

import "strings"

// Inserts the prefix string at the beginning of each non-empty line.
func Indent(text, prefix string) string {
	var builder strings.Builder
	isBeginningOfLine := true

	for _, r := range text {
		if isBeginningOfLine && r != '\n' {
			builder.WriteString(prefix)
		}
		builder.WriteRune(r)
		isBeginningOfLine = (r == '\n')
	}

	return builder.String()
}

// Adds a spacer to a string after every nth character.
func AddSpacer(text, spacer string, n int) string {
	var builder strings.Builder
	runes := []rune(text)

	for i := 0; i < len(runes); i += n {
		if i > 0 {
			builder.WriteString(spacer)
		}

		if i+n > len(runes) {
			builder.WriteString(string(runes[i:]))
		} else {
			builder.WriteString(string(runes[i : i+n]))
		}
	}

	return builder.String()
}
