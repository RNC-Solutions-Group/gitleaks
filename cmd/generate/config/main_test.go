package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestGeneratedConfigMatchesTracked verifies that the checked-in configuration
// is exactly what the generator produces, without rewriting the tracked file.
// Line endings are normalized because Windows checkouts commonly use CRLF while
// the generator and Git blobs use LF.
func TestGeneratedConfigMatchesTracked(t *testing.T) {
	generatedPath := filepath.Join(t.TempDir(), "gitleaks.toml")
	if err := generate(generatedPath); err != nil {
		t.Fatalf("generate configuration: %v", err)
	}

	generated, err := os.ReadFile(generatedPath)
	if err != nil {
		t.Fatalf("read generated configuration: %v", err)
	}
	trackedPath := filepath.Join("..", "..", "..", "config", "gitleaks.toml")
	tracked, err := os.ReadFile(trackedPath)
	if err != nil {
		t.Fatalf("read tracked configuration %s: %v", trackedPath, err)
	}

	if !bytes.Equal(normalizeLineEndings(generated), normalizeLineEndings(tracked)) {
		t.Fatal("generated configuration differs from config/gitleaks.toml; run `make generate` and review the diff")
	}
}

func normalizeLineEndings(content []byte) []byte {
	return bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
}
