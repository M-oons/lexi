package random

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadWords(t *testing.T) {
	tmpDir := t.TempDir()

	filePath := filepath.Join(tmpDir, "words.txt")
	if err := os.WriteFile(filePath, []byte("alpha\nbeta\n"), 0o644); err != nil {
		t.Fatalf("write temp words file: %v", err)
	}

	tests := []struct {
		name      string
		input     string
		separator string
		want      []string
	}{
		{
			name:      "inline defaults to comma",
			input:     "a,b",
			separator: "",
			want:      []string{"a", "b"},
		},
		{
			name:      "file defaults to newline",
			input:     filePath,
			separator: "",
			want:      []string{"alpha", "beta"},
		},
		{
			name:      "explicit escaped newline separator",
			input:     "a\nb",
			separator: "\\\\n",
			want:      []string{"a", "b"},
		},
		{
			name:      "explicit escaped tab separator",
			input:     "a\tb",
			separator: "\\\\t",
			want:      []string{"a", "b"},
		},
		{
			name:      "nonexistent path treated as inline",
			input:     filepath.Join(tmpDir, "no_such_file.txt"),
			separator: "",
			want:      []string{filepath.Join(tmpDir, "no_such_file.txt")},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := LoadWords(tc.input, tc.separator)
			if err != nil {
				t.Fatalf("LoadWords(%q, %q) returned error: %v", tc.input, tc.separator, err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("LoadWords(%q, %q) = %#v, want %#v", tc.input, tc.separator, got, tc.want)
			}
		})
	}
}

func TestLoadDefaultWords_Sanity(t *testing.T) {
	words := LoadDefaultWords("")
	if len(words) == 0 {
		t.Fatal("LoadDefaultWords returned empty slice")
	}

	for i, word := range words {
		if word == "" {
			t.Fatalf("word[%d] is empty", i)
		}
		if strings.TrimSpace(word) != word {
			t.Fatalf("word[%d] is not trimmed: %q", i, word)
		}
	}
}

func TestFilterWords(t *testing.T) {
	words := []string{"", "a", "ab", "abc", "abcd"}

	tests := []struct {
		name      string
		minLength int
		maxLength int
		regex     string
		want      []string
		wantErr   bool
	}{
		{
			name:      "min and max default values",
			minLength: 0,
			maxLength: 0,
			regex:     "",
			want:      []string{"a", "ab", "abc", "abcd"},
		},
		{
			name:      "bounded lengths",
			minLength: 2,
			maxLength: 3,
			regex:     "",
			want:      []string{"ab", "abc"},
		},
		{
			name:      "max default when non-positive",
			minLength: 3,
			maxLength: 0,
			regex:     "",
			want:      []string{"abc", "abcd"},
		},
		{
			name:      "regex filtering",
			minLength: 1,
			maxLength: 0,
			regex:     "^ab",
			want:      []string{"ab", "abc", "abcd"},
		},
		{
			name:      "regex and bounds together",
			minLength: 1,
			maxLength: 2,
			regex:     "^a.$",
			want:      []string{"ab"},
		},
		{
			name:      "min greater than max",
			minLength: 4,
			maxLength: 2,
			regex:     "",
			want:      []string{},
		},
		{
			name:      "invalid regex returns error",
			minLength: 1,
			maxLength: 0,
			regex:     "(",
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FilterWords(words, tc.minLength, tc.maxLength, tc.regex)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("FilterWords returned unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("FilterWords(...) = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestParseLines(t *testing.T) {
	tests := []struct {
		name      string
		lines     string
		separator string
		want      []string
	}{
		{
			name:      "trim and drop empties",
			lines:     " a , , b ",
			separator: ",",
			want:      []string{"a", "b"},
		},
		{
			name:      "newline with trailing newline",
			lines:     "a\nb\n",
			separator: "\n",
			want:      []string{"a", "b"},
		},
		{
			name:      "crlf input",
			lines:     "a\r\nb\r\n",
			separator: "\n",
			want:      []string{"a", "b"},
		},
		{
			name:      "separator not found",
			lines:     "  lone  ",
			separator: ",",
			want:      []string{"lone"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := parseLines(tc.lines, tc.separator)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("parseLines(%q, %q) = %#v, want %#v", tc.lines, tc.separator, got, tc.want)
			}
		})
	}
}

func TestNormalizeSeparator(t *testing.T) {
	tests := []struct {
		name      string
		separator string
		want      string
	}{
		{
			name:      "empty",
			separator: "",
			want:      "",
		},
		{
			name:      "newline",
			separator: "\\\\n",
			want:      "\n",
		},
		{
			name:      "carriage return",
			separator: "\\\\r",
			want:      "\r",
		},
		{
			name:      "tab",
			separator: "\\\\t",
			want:      "\t",
		},
		{
			name:      "mixed",
			separator: "x\\\\n,y",
			want:      "x\n,y",
		},
		{
			name:      "unknown escape unchanged",
			separator: "\\x",
			want:      "\\x",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeSeparator(tc.separator)
			if got != tc.want {
				t.Fatalf("normalizeSeparator(%q) = %q, want %q", tc.separator, got, tc.want)
			}
		})
	}
}
