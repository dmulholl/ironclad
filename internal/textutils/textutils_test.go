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

func TestAddSpacer(t *testing.T) {
	testcases := []struct {
		name   string
		input  string
		spacer string
		n      int
		want   string
	}{
		{
			name:   "empty input empty spacer",
			input:  "",
			spacer: "",
			n:      3,
			want:   "",
		},
		{
			name:   "empty input",
			input:  "",
			spacer: "--",
			n:      3,
			want:   "",
		},
		{
			name:   "input a",
			input:  "a",
			spacer: "--",
			n:      3,
			want:   "a",
		},
		{
			name:   "input ab",
			input:  "ab",
			spacer: "--",
			n:      3,
			want:   "ab",
		},
		{
			name:   "input abc",
			input:  "abc",
			spacer: "--",
			n:      3,
			want:   "abc",
		},
		{
			name:   "input abcd",
			input:  "abcd",
			spacer: "--",
			n:      3,
			want:   "abc--d",
		},
		{
			name:   "input abcde",
			input:  "abcde",
			spacer: "--",
			n:      3,
			want:   "abc--de",
		},
		{
			name:   "input abcdef",
			input:  "abcdef",
			spacer: "--",
			n:      3,
			want:   "abc--def",
		},
		{
			name:   "input abcdefg",
			input:  "abcdefg",
			spacer: "--",
			n:      3,
			want:   "abc--def--g",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := AddSpacer(tc.input, tc.spacer, tc.n)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRuneString(t *testing.T) {
	testcases := []struct {
		name   string
		length int
		char   rune
		want   string
	}{
		{
			name:   "zero length",
			length: 0,
			char:   '-',
			want:   "",
		},
		{
			name:   "single rune output",
			length: 1,
			char:   '-',
			want:   "-",
		},
		{
			name:   "multirune output",
			length: 2,
			char:   '-',
			want:   "--",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := RuneString(tc.length, tc.char)
			assert.Equal(t, tc.want, got)
		})
	}
}
