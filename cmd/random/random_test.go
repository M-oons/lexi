package random

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestNewCmd_RunE_LoadWordsFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "words.txt")
	if err := os.WriteFile(filePath, []byte("Alpha\nBeta\n"), 0o644); err != nil {
		t.Fatalf("write temp words file: %v", err)
	}

	stdout, _, err := executeRandomCmd(t,
		"--words", filePath,
		"--count", "5",
		"--seed", "11",
	)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	lines := splitOutputLines(stdout)
	if len(lines) != 5 {
		t.Fatalf("line count = %d, want 5", len(lines))
	}

	valid := map[string]bool{
		"Alpha": true,
		"Beta":  true,
	}
	for i, line := range lines {
		if !valid[line] {
			t.Fatalf("line[%d] = %q, want one of Alpha/Beta", i, line)
		}
	}
}

func TestNewCmd_RunE_UsesDefaultWordsWhenWordsFlagEmpty(t *testing.T) {
	stdout, _, err := executeRandomCmd(t,
		"--count", "1",
		"--seed", "7",
	)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	lines := splitOutputLines(stdout)
	if len(lines) != 1 {
		t.Fatalf("line count = %d, want 1", len(lines))
	}
	if strings.TrimSpace(lines[0]) == "" {
		t.Fatal("expected non-empty output line")
	}
}

func TestNewCmd_RunE_CountDefaultsToOne(t *testing.T) {
	tests := []struct {
		name  string
		count string
	}{
		{name: "zero", count: "0"},
		{name: "negative", count: "-3"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _, err := executeRandomCmd(t,
				"--words", "only",
				"--count", tc.count,
			)
			if err != nil {
				t.Fatalf("execute command: %v", err)
			}

			lines := splitOutputLines(stdout)
			if len(lines) != 1 {
				t.Fatalf("line count = %d, want 1", len(lines))
			}
			if lines[0] != "only" {
				t.Fatalf("line[0] = %q, want %q", lines[0], "only")
			}
		})
	}
}

func TestNewCmd_RunE_ReplacesAllPlaceholderTypes(t *testing.T) {
	stdout, _, err := executeRandomCmd(t,
		"--words", "MiXeD",
		"--count", "2",
		"--format", "x-{}-{l}-{u}",
	)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	got := splitOutputLines(stdout)
	want := []string{"x-MiXeD-mixed-MIXED", "x-MiXeD-mixed-MIXED"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("output = %#v, want %#v", got, want)
	}
}

func TestNewCmd_RunE_LiteralFormatWithoutPlaceholders(t *testing.T) {
	stdout, _, err := executeRandomCmd(t,
		"--words", "alpha,beta",
		"--count", "3",
		"--format", "literal",
	)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	got := splitOutputLines(stdout)
	want := []string{"literal", "literal", "literal"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("output = %#v, want %#v", got, want)
	}
}

func TestNewCmd_RunE_InvalidRegexReturnsError(t *testing.T) {
	stdout, _, err := executeRandomCmd(t,
		"--words", "a,b",
		"--regex", "(",
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
}

func TestNewCmd_RunE_OrdersAscending(t *testing.T) {
	stdout, _, err := executeRandomCmd(t,
		"--words", "delta,alpha,charlie,bravo",
		"--count", "25",
		"--seed", "17",
		"--order", "asc",
	)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	lines := splitOutputLines(stdout)
	if len(lines) != 25 {
		t.Fatalf("line count = %d, want 25", len(lines))
	}

	sorted := append([]string(nil), lines...)
	sort.Strings(sorted)
	if !reflect.DeepEqual(lines, sorted) {
		t.Fatalf("output is not sorted ascending: %#v", lines)
	}
}

func TestNewCmd_RunE_OrdersDescending(t *testing.T) {
	stdout, _, err := executeRandomCmd(t,
		"--words", "delta,alpha,charlie,bravo",
		"--count", "25",
		"--seed", "17",
		"--order", "descending",
	)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}

	lines := splitOutputLines(stdout)
	if len(lines) != 25 {
		t.Fatalf("line count = %d, want 25", len(lines))
	}

	sorted := append([]string(nil), lines...)
	sort.Strings(sorted)
	sort.Sort(sort.Reverse(sort.StringSlice(sorted)))
	if !reflect.DeepEqual(lines, sorted) {
		t.Fatalf("output is not sorted descending: %#v", lines)
	}
}

func TestNewCmd_RunE_InvalidOrderReturnsError(t *testing.T) {
	stdout, _, err := executeRandomCmd(t,
		"--words", "a,b",
		"--order", "sideways",
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid order") {
		t.Fatalf("error = %q, want substring %q", err.Error(), "invalid order")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
}

func executeRandomCmd(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	cmd := NewCmd()
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return stdout.String(), stderr.String(), err
}

func splitOutputLines(stdout string) []string {
	stdout = strings.TrimSuffix(stdout, "\n")
	if stdout == "" {
		return []string{}
	}

	return strings.Split(stdout, "\n")
}
