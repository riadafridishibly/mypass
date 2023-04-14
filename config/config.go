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
	"gopkg.in/yaml.v3"

	jww "github.com/spf13/jwalterweatherman"
)

// Config location ~/.config/mypass/config.yaml
type Config struct {
	// Stored private keys
	PrivateKeys string `yaml:"private_keys"`
	// Cached password path
	CachedPassword  string `yaml:"cached_password"`
	DatabasePath    string `yaml:"database"`
	EncryptionLevel int    `yaml:"encryption_level"`
	Backend         string `yaml:"backend"`
}

func (cfg *Config) ToYaml() ([]byte, error) {
	return yaml.Marshal(cfg)
}

var DefaultConfig = &Config{
	PrivateKeys:    ExpandWithHome("~/.mypass/private_keys"),
	CachedPassword: ExpandWithHome("~/.mypass/cached_pass"),
	DatabasePath:   ExpandWithHome("~/.mypass/db.sqlite"),
	Backend:        "sqlite3",
}

func LoadCachedPassword() error {
	// Already in viper
	if viper.GetString(vkeys.Password) != "" {
		return nil
	}
	p := viper.GetString(vkeys.CachedPassword)
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
		jww.DEBUG.Printf("Failed to create random file: %v, err: %v", randomFile, err)
		return err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.CopyN(io.MultiWriter(buf, f), rand.Reader, 1<<20)
	if err != nil {
		jww.DEBUG.Printf("Failed to copy random data to random file: %v, err: %v", randomFile, err)
		return err
	}
	data, err := encryption.EncryptWithPassword(pass, buf.String())
	if err != nil {
		jww.DEBUG.Println("Failed to encrypt password with random file, err:", err)
		return err
	}
	err = os.WriteFile(cacheFile, data, 0600)
	if err != nil {
		jww.DEBUG.Printf("Failed to write cache file: %v, err: %v", cacheFile, err)
		return err
	}
	return nil
}

func expandedHome(p string) (string, error) {
	p = strings.TrimPrefix(p, "~/")
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, p), nil
}

func ExpandWithHome(p string) string {
	v, err := expandedHome(p)
	if err != nil {
		jww.FATAL.Fatalf("Failed to get home dir: %v", err)
	}
	return v
}

func LoadPrivateKeys() error {
	if len(viper.GetStringSlice(vkeys.PrivateKeys)) > 0 {
		return nil
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
	var privKeySlice []string
	for _, key := range privKeys.Keys {
		privKeySlice = append(privKeySlice, string(key))
	}
	viper.Set(vkeys.PrivateKeys, privKeySlice)
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
