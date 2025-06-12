package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// HashFile calculates the SHA-256 hash of a file at the given path.
func HashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("unable to open file for hashing: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("error reading file during hashing: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// SignPayload generates an HMAC-SHA256 signature for a given payload using a secret key.
func SignPayload(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

// VerifySignature compares a payload and a provided signature using the secret key.
func VerifySignature(payload []byte, signature, secret string) bool {
	expectedSignature := SignPayload(payload, secret)
	// Use hmac.Equal for constant-time comparison to prevent timing attacks.
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
