package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsFile(t *testing.T) {
	tmpDir := t.TempDir()

	filePath := filepath.Join(tmpDir, "sample.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	dirPath := filepath.Join(tmpDir, "nested")
	if err := os.Mkdir(dirPath, 0o755); err != nil {
		t.Fatalf("mkdir temp dir: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "regular file",
			path: filePath,
			want: true,
		},
		{
			name: "directory",
			path: dirPath,
			want: false,
		},
		{
			name: "nonexistent",
			path: filepath.Join(tmpDir, "missing.txt"),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsFile(tc.path)
			if got != tc.want {
				t.Fatalf("IsFile(%q) = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}
