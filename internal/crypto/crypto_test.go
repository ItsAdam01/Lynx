package crypto

import (
	"os"
	"testing"
)

func TestHashFile(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := []byte("hello, cybersecurity world!")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Expected hash for "hello, cybersecurity world!"
	// Generated using: echo -n "hello, cybersecurity world!" | sha256sum
	expectedHash := "dbbe7ee53692a509a86f33fde86bd6438d9986c657e99674c68d1b185785bf65"

	hash, err := HashFile(tmpfile.Name())
	if err != nil {
		t.Errorf("HashFile returned error: %v", err)
	}

	if hash != expectedHash {
		t.Errorf("HashFile mismatch. Expected %s, got %s", expectedHash, hash)
	}
}

func TestHMACSignVerify(t *testing.T) {
	payload := []byte("this is a test baseline payload")
	secret := "my-super-secret-key"

	signature := SignPayload(payload, secret)
	if signature == "" {
		t.Fatal("SignPayload returned empty signature")
	}

	if !VerifySignature(payload, signature, secret) {
		t.Error("VerifySignature failed to validate a correct signature")
	}

	// Test with wrong secret
	if VerifySignature(payload, signature, "wrong-secret") {
		t.Error("VerifySignature should have failed with the wrong secret")
	}

	// Test with tampered payload
	tamperedPayload := []byte("this is a tampered baseline payload")
	if VerifySignature(tamperedPayload, signature, secret) {
		t.Error("VerifySignature should have failed with a tampered payload")
	}
}
