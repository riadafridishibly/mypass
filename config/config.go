package config

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/riadafridishibly/mypass/encryption"
	"github.com/riadafridishibly/mypass/models"
	"github.com/riadafridishibly/mypass/vkeys"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// Config location ~/.config/mypass/config.yaml
type Config struct {
	// Stored private keys
	PrivateKeys string `yaml:"private_keys"`
	// Cached password path
	CachedPassword  string `yaml:"cached_password"`
	DatabasePath    string `yaml:"database"`
	EncryptionLevel int    `yaml:"encryption_level"`
}

var DefaultConfig = &Config{
	PrivateKeys:    "~/.mypass/private_keys",
	CachedPassword: "~/.mypass/cached_pass",
	DatabasePath:   "~/.mypass/db",
}

func LoadCachedPassword() error {
	// Already in viper
	if viper.GetString(vkeys.Password) != "" {
		return nil
	}
	p, err := expandedPath(DefaultConfig.CachedPassword)
	if err != nil {
		return err
	}
	// Try load from cache (TODO: encrypt the cache)
	data, err := decryptCache(p+".rnd", p)
	if err == nil {
		viper.Set(vkeys.Password, string(data))
		// We need to test if the password is correct!
		err = LoadPrivateKeys()
		if err == nil {
			return nil
		}
	}

	fmt.Printf("Enter your master password: ")
	data, err = term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to read password from stdin: %w", err)
	}
	viper.Set(vkeys.Password, string(data))
	fmt.Println()
	// Save to cache
	_ = encryptCache(p+".rnd", p, data)
	return nil
}

func decryptCache(randomFile, cacheFile string) ([]byte, error) {
	data, err := os.ReadFile(randomFile)
	if err != nil {
		return nil, err
	}
	cache, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}
	// TODO: set lower work factor
	return encryption.DecryptWithPassword(cache, string(data))
}

func encryptCache(randomFile, cacheFile string, pass []byte) error {
	// Create a file of 1MB size
	f, err := os.Create(randomFile)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.CopyN(io.MultiWriter(buf, f), rand.Reader, 1<<20)
	if err != nil {
		return err
	}
	data, err := encryption.EncryptWithPassword(pass, buf.String())
	if err != nil {
		return err
	}
	return os.WriteFile(cacheFile, data, 0600)
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
	if len(viper.GetStringSlice(vkeys.PrivateKeys)) > 0 {
		return nil
	}
	if err := LoadCachedPassword(); err != nil {
		return err
	}
	privKeyPath := viper.GetString(vkeys.PrivateKeysPath)
	if privKeyPath == "" {
		return fmt.Errorf("private keys not found, path: %q", privKeyPath)
	}
	f, err := os.Open(privKeyPath)
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
	viper.Set(vkeys.PrivateKeys, pkSlice)
	return nil
}

func LoadDatabase() (*models.Database, error) {
	if err := LoadPrivateKeys(); err != nil {
		return nil, err
	}

	if db := viper.Get("db"); db != nil {
		v, ok := db.(*models.Database)
		if ok {
			return v, nil
		}
	}
	dbPath := viper.GetString("database")
	if dbPath == "" {
		return nil, errors.New("database path is not defined")
	}
	var db models.Database
	data, err := os.ReadFile(dbPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &db)
	if err != nil {
		return nil, err
	}
	viper.Set(vkeys.PublicKeys, db.PublicKeys)
	viper.Set("db", &db)
	return &db, nil
}
