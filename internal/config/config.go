package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config represents the agent's configuration.
type Config struct {
	AgentName     string   `mapstructure:"agent_name"`
	HmacSecretEnv string   `mapstructure:"hmac_secret_env"`
	LogFile       string   `mapstructure:"log_file"`
	WebhookURL    string   `mapstructure:"webhook_url"`
	PathsToWatch  []string `mapstructure:"paths_to_watch"`
	FilesToWatch  []string `mapstructure:"files_to_watch"`
}

// LoadConfig reads the configuration from config.yaml or environment variables.
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Basic validation: ensure the HMAC secret env var is defined
	if cfg.HmacSecretEnv == "" {
		cfg.HmacSecretEnv = "LYNX_HMAC_SECRET"
	}

	return &cfg, nil
}

// InitConfig creates a default configuration file at the given path.
func InitConfig(path string) error {
	// 1. Check if the file already exists (we don't want to overwrite)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("configuration file already exists at: %s", path)
	}

	// 2. Define default configuration values
	// We'll use a direct YAML string to ensure comments and clear formatting.
	defaultYAML := `agent_name: "default-agent"
hmac_secret_env: "LYNX_HMAC_SECRET"
log_file: "/var/log/lynx.log"
webhook_url: ""

# Paths and files to monitor
paths_to_watch:
  - "/etc/ssh"
  - "/usr/local/bin"

files_to_watch:
  - "/etc/passwd"
  - "/etc/hosts"
`

	return os.WriteFile(path, []byte(defaultYAML), 0644)
}

// GetHmacSecret retrieves the secret from the configured environment variable.
func (c *Config) GetHmacSecret() (string, error) {
	secret := os.Getenv(c.HmacSecretEnv)
	if secret == "" {
		return "", fmt.Errorf("HMAC secret not found in environment variable: %s", c.HmacSecretEnv)
	}
	return secret, nil
}
