package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// 1. Run the initialization logic (this should fail initially as it's not implemented)
	err := InitConfig(configPath)
	if err != nil {
		t.Fatalf("InitConfig failed: %v", err)
	}

	// 2. Verify the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("InitConfig did not create the configuration file")
	}

	// 3. Load the created config and verify defaults
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed to read the initialized config: %v", err)
	}

	if cfg.AgentName != "default-agent" {
		t.Errorf("Expected AgentName 'default-agent', got '%s'", cfg.AgentName)
	}

	if len(cfg.PathsToWatch) == 0 {
		t.Error("InitConfig created an empty paths list")
	}
}
