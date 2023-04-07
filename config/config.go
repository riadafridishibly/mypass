package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/riadafridishibly/mypass/models"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// Config location ~/.config/mypass/config.yaml
type Config struct {
	PrivateKeyPaths     []string
	CachedPasswordPaths []string
	DatabasePath        string
	EncryptionLevel     int
}

var DefaultConfig = &Config{
	PrivateKeyPaths: []string{
		"~/.mypass/private_keys",
	},
	CachedPasswordPaths: []string{
		"~/.mypass/cached_pass",
	},
	DatabasePath: "~/.mypass/db",
}

func Init() error {
	// Handle `init`, when the app is not yet initialized!
	if err := LoadCachedPassword(); err != nil {
		return err
	}

	if err := LoadPrivateKeys(); err != nil {
		return err
	}

	return nil
}

func LoadCachedPassword() error {
	// Try loading password from file
	// Else prompt for password
	// TODO: Add some compution here!
	p, err := expandedPath(DefaultConfig.CachedPasswordPaths[0])
	if err != nil {
		return err
	}
	data, err := os.ReadFile(p)
	if err == nil {
		viper.Set("password", string(data))
		return nil
	}

	fmt.Printf("Enter your master password: ")
	data, err = term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to read password from stdin: %w", err)
	}
	viper.Set("password", string(data))
	fmt.Println()
	return nil
}

func expandedPath(p string) (string, error) {
	p = strings.TrimPrefix(p, "~/")
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, p), nil
}

func LoadPrivateKeys() error {
	// Look for cached password
	// If not prompt for password
	// Unlock the private keys
	// Set private_keys in viper
	p, err := expandedPath(DefaultConfig.PrivateKeyPaths[0])
	if err != nil {
		return err
	}
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	var privKeys models.PrivateKeys
	err = json.NewDecoder(f).Decode(&privKeys)
	if err != nil {
		return err
	}
	var pkSlice []string
	for _, key := range privKeys.Keys {
		pkSlice = append(pkSlice, string(key))
	}
	viper.Set("private_keys", pkSlice)
	return nil
}

func LoadDatabase() (*models.Database, error) {
	var db models.Database
	data, err := os.ReadFile(DefaultConfig.DatabasePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &db)
	if err != nil {
		return nil, err
	}
	return &db, nil
}
