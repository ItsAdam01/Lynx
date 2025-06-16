package fs

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ItsAdam01/Lynx/internal/crypto"
)

// BaselineMetadata contains information about when and how the baseline was created.
type BaselineMetadata struct {
	GeneratedAt time.Time `json:"generated_at"`
	TotalFiles  int       `json:"total_files"`
	ConfigHash  string    `json:"config_hash"` // Protects the ignore list and paths
}

// Baseline represents the full state of the monitored files.
type Baseline struct {
	Metadata  BaselineMetadata  `json:"metadata"`
	Hashes    map[string]string `json:"hashes"`
	Signature string            `json:"signature"`
}

// NewBaseline initializes a Baseline with the given hashes and configuration fingerprint.
func NewBaseline(hashes map[string]string, configHash string) *Baseline {
	return &Baseline{
		Metadata: BaselineMetadata{
			GeneratedAt: time.Now(),
			TotalFiles:  len(hashes),
			ConfigHash:  configHash,
		},
		Hashes: hashes,
	}
}

// Save marshals the baseline to JSON, signs it with HMAC, and writes it to disk.
func (b *Baseline) Save(path, secret string) error {
	// Clear signature before signing
	b.Signature = ""

	payload, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("failed to marshal baseline: %w", err)
	}

	b.Signature = crypto.SignPayload(payload, secret)

	finalJSON, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal signed baseline: %w", err)
	}

	return os.WriteFile(path, finalJSON, 0600)
}

// LoadBaseline reads the baseline from disk and verifies its HMAC signature.
func LoadBaseline(path, secret string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read baseline file: %w", err)
	}

	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("failed to unmarshal baseline: %w", err)
	}

	// To verify, we must extract the signature and verify the rest of the payload
	storedSignature := b.Signature
	b.Signature = "" // Reset to match state during original signing

	payload, _ := json.Marshal(b) // This is the payload that was signed

	if !crypto.VerifySignature(payload, storedSignature, secret) {
		return nil, fmt.Errorf("baseline signature verification failed - file may be tampered")
	}

	b.Signature = storedSignature // Restore signature
	return &b, nil
}
