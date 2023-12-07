package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// cfCLIConfig is the CF CLI configuration.
type cfCLIConfig struct {
	AccessToken           string
	RefreshToken          string
	Target                string
	AuthorizationEndpoint string
	UaaEndpoint           string
	UAAOAuthClient        string
	UAAOAuthClientSecret  string
	UAAGrantType          string
	SSHOAuthClient        string
	SSLDisabled           bool
}

// findCFHomeDir finds the CF Home directory.
func findCFHomeDir() (string, error) {
	cfHomeDir := os.Getenv("CF_HOME")
	if cfHomeDir != "" {
		return cfHomeDir, nil
	}
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine user's home directory: %w", err)
	}
	return userHomeDir, nil
}

// createConfigFromCFCLIConfig reads the CF Home configuration from the specified directory.
func loadCFCLIConfig(cfHomeDir string) (*cfCLIConfig, error) {
	configFile := filepath.Join(filepath.Join(cfHomeDir, ".cf"), "config.json")
	cfJSON, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", configFile, err)
	}
	var cfgHome cfCLIConfig
	if err = json.Unmarshal(cfJSON, &cfgHome); err != nil {
		return nil, fmt.Errorf("error while unmarshalling CF CLI config: %w", err)
	}
	return &cfgHome, nil
}
