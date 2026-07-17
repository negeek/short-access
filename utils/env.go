package utils

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnvFile reads a .env file and sets any variables it defines that aren't
// already present in the environment. A missing file is not an error, so it's
// safe to call where the variables come from elsewhere (like Docker).
//
// It's a small stand-in for a dotenv library and understands the common cases:
// KEY=VALUE lines, # comments, blank lines, optional surrounding quotes, and an
// optional leading "export".
func LoadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		value = trimQuotes(strings.TrimSpace(value))

		// Real environment variables win, so we never clobber them.
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
	return scanner.Err()
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		first, last := s[0], s[len(s)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
