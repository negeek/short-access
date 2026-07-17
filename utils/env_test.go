package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := `# a comment
export EXPORTED=yes
QUOTED="a value"
SINGLE='another'
PLAIN=plain

PREEXISTING=from-file
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	// A variable already in the environment must not be overwritten.
	t.Setenv("PREEXISTING", "from-env")

	if err := LoadEnvFile(path); err != nil {
		t.Fatalf("LoadEnvFile: %v", err)
	}

	cases := map[string]string{
		"EXPORTED":    "yes",
		"QUOTED":      "a value",
		"SINGLE":      "another",
		"PLAIN":       "plain",
		"PREEXISTING": "from-env",
	}
	for key, want := range cases {
		if got := os.Getenv(key); got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}
}

func TestLoadEnvFileMissingIsNoError(t *testing.T) {
	if err := LoadEnvFile(filepath.Join(t.TempDir(), "nope.env")); err != nil {
		t.Fatalf("missing file should be fine, got %v", err)
	}
}
