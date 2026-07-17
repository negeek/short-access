package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// apiKeyPrefix marks our keys so they are easy to recognise in logs and code.
const apiKeyPrefix = "sa_"

// GenerateAPIKey returns a new random API key. The caller shows this to the user
// once and never stores it; only its hash is kept.
func GenerateAPIKey() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return apiKeyPrefix + base64.RawURLEncoding.EncodeToString(buf), nil
}

// HashAPIKey returns the SHA-256 hash of a key. API keys are long and random, so
// a fast hash is enough and lets us look a key up directly by its hash.
func HashAPIKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}
