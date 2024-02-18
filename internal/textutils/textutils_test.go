package textutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndent(t *testing.T) {
	testcases := []struct {
		name   string
		input  string
		prefix string
		want   string
	}{
		{
			name:   "empty input empty prefix",
			input:  "",
			prefix: "",
			want:   "",
		},
		{
			name:   "empty input",
			input:  "",
			prefix: "---",
			want:   "",
		},
		{
			name:   "single line input",
			input:  "foo bar baz",
			prefix: "---",
			want:   "---foo bar baz",
		},
		{
			name:   "multiline input",
			input:  "foo\nbar\nbaz",
			prefix: "---",
			want:   "---foo\n---bar\n---baz",
		},
		{
			name:   "multiline input empty lines",
			input:  "foo\n\nbar\n\nbaz",
			prefix: "---",
			want:   "---foo\n\n---bar\n\n---baz",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := Indent(tc.input, tc.prefix)
			assert.Equal(t, tc.want, got)
		})
	}
}
