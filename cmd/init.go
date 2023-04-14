/*
Copyright © 2023 Riad Afridi Shibly <riadafridishibly@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"filippo.io/age"
	"github.com/riadafridishibly/mypass/backend"
	"github.com/riadafridishibly/mypass/config"
	"github.com/riadafridishibly/mypass/models"
	"github.com/riadafridishibly/mypass/vkeys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

func initPrivateKeys(file string) (pubKeys []string, err error) {
	// If err == nil; file exists
	if _, err := os.Stat(file); err == nil {
		return nil, errors.New("private key file exists")
	}

	fmt.Print("Enter your master password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to read master password: %w", err)
		// jww.FATAL.Fatal("Failed to read password")
	}
	fmt.Println()
	fmt.Print("Enter your master password (again): ")
	password2, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to read master password: %w", err)
		// jww.FATAL.Fatal("Failed to read password")
	}
	fmt.Println()

	if !bytes.Equal(password, password2) {
		// jww.ERROR.Fatal("Password didn't match")
		return nil, fmt.Errorf("password didn't match")
	}

	viper.Set(vkeys.Password, string(password))

	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, fmt.Errorf("failed to create X25519 Key pairs: %v", err)
	}
	privKeys := models.PrivateKeys{
		Meta: models.Meta{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Keys: []models.SymSecretStr{
			models.SymSecretStr(identity.String()),
		},
	}
	// OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(privKeys)
	if err != nil {
		return nil, err
	}
	return []string{identity.Recipient().String()}, err
}

var ErrConfigExits = errors.New("config already exits")

func initDefaultConfig() error {
	file := viper.ConfigFileUsed()
	if file != "" {
		return ErrConfigExits
	}
	data, err := config.DefaultConfig.ToYaml()
	if err != nil {
		return err
	}
	// Populate default config
	viper.SetConfigType("yaml")
	err = viper.MergeConfig(bytes.NewReader(data))
	if err != nil {
		return err
	}
	return os.WriteFile(DefaultConfigPath, data, 0644)
}

func initPrivateKeysAndDb() error {
	// Maybe define config root?
	dir := filepath.Dir(viper.GetString(vkeys.PrivateKeysPath))
	err := os.MkdirAll(dir, 0o0700)
	if err != nil {
		// jww.FATAL.Fatalf("Failed to create dir: %v: err: %v", exHome("~/.mypass"), err)
		return err
	}

	pubKeys, err := initPrivateKeys(viper.GetString(vkeys.PrivateKeysPath))
	if err != nil {
		return err
	}

	b, err := backend.Get()
	if err != nil {
		return err
	}

	return b.AddPublicKeys(pubKeys...)
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize mypass application",
	Long: `Initialize the application, create necessary directories and files. 

If the application is already initialized it'll do nothing.
In case of error it'll report them in stdout.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Try initialize default config!
		err := initDefaultConfig()
		if err != nil && !errors.Is(err, ErrConfigExits) {
			return err
		}

		err = initPrivateKeysAndDb()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
