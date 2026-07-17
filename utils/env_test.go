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
  SPACED  =  trimmed
EMPTY=
DB_URL=postgres://u:p@host:5432/db?sslmode=disable
IGNORED_NO_EQUALS

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
		"SPACED":      "trimmed",                                     // key and value are trimmed
		"EMPTY":       "",                                            // empty value is allowed
		"DB_URL":      "postgres://u:p@host:5432/db?sslmode=disable", // only the first '=' splits
		"PREEXISTING": "from-env",                                    // env wins over the file
	}
	for key, want := range cases {
		if got := os.Getenv(key); got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}

	// A line without '=' is not a variable and must be skipped.
	if _, ok := os.LookupEnv("IGNORED_NO_EQUALS"); ok {
		t.Error("a line without '=' should be ignored")
	}
}

func TestLoadEnvFileMissingIsNoError(t *testing.T) {
	if err := LoadEnvFile(filepath.Join(t.TempDir(), "nope.env")); err != nil {
		t.Fatalf("missing file should be fine, got %v", err)
	}
}
