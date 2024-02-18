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
