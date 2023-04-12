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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"filippo.io/age"
	"github.com/riadafridishibly/mypass/models"
	"github.com/riadafridishibly/mypass/vkeys"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

func exHome(p string) string {
	v, err := expandedHome(p)
	if err != nil {
		jww.FATAL.Fatalf("Failed to get home dir: %v", err)
	}
	return v
}

func initPrivateKeysAndDb() {
	// Create key pair
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		jww.ERROR.Fatalf("Failed to create X25519 Key pairs: %v", err)
	}

	fmt.Print("Enter your master password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		jww.FATAL.Fatal("Failed to read password")
	}
	fmt.Println()
	fmt.Print("Enter your master password (again): ")
	password2, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		jww.FATAL.Fatal("Failed to read password")
	}
	fmt.Println()

	if !bytes.Equal(password, password2) {
		jww.ERROR.Fatal("Password didn't match")
	}

	viper.Set(vkeys.Password, string(password))

	privKeys := models.PrivateKeys{
		Meta: models.Meta{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Keys: []models.SymSecretStr{
			models.SymSecretStr(identity.String()),
		},
	}

	db := models.Database{
		PublicKeys: []string{
			identity.Recipient().String(),
		},
	}

	err = os.MkdirAll(exHome("~/.mypass"), 0o0700)
	if err != nil {
		jww.FATAL.Fatalf("Failed to create dir: %v: err: %v", exHome("~/.mypass"), err)
	}

	jww.INFO.Printf("Writing database to: %s", viper.GetString(vkeys.DatabasePath))
	err = writeJSON(viper.GetString(vkeys.DatabasePath), db)
	if err != nil {
		jww.FATAL.Fatalf("Failed to write database: %v", err)
	}
	jww.INFO.Printf("Writing private keys to: %s", viper.GetString(vkeys.PrivateKeysPath))
	err = writeJSON(viper.GetString(vkeys.PrivateKeysPath), privKeys)
	if err != nil {
		jww.FATAL.Fatalf("Failed to write private_keys: %v", err)
	}
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func isDBExists() bool {
	return fileExists(viper.GetString(vkeys.DatabasePath))
}

func isPrivKeysExists() bool {
	return fileExists(viper.GetString(vkeys.PrivateKeysPath))
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize mypass application",
	Long: `Initialize the application, create necessary directories and files. 

If the application is already initialized it'll do nothing.
In case of error it'll report them in stdout.`,
	Run: func(cmd *cobra.Command, args []string) {
		file := viper.ConfigFileUsed()
		// FIXME: if config exists, but db and priv keys doesn't then create them
		if file != "" {
			jww.WARN.Printf("Config file exists: %v", file)
			dbExists := isDBExists()
			privKeyExists := isPrivKeysExists()
			if dbExists {
				jww.FATAL.Fatal("Database file already exists in:", viper.GetString(vkeys.DatabasePath))
			}
			if privKeyExists {
				jww.FATAL.Fatal("Private keys file already exists in:", viper.GetString(vkeys.PrivateKeysPath))
			}

			jww.INFO.Print("Initializing private keys and database files...")
			initPrivateKeysAndDb()
			os.Exit(0)
		}

		// Set config
		viper.Set(vkeys.PrivateKeysPath, exHome("~/.mypass/private_keys"))
		viper.Set(vkeys.DatabasePath, exHome("~/.mypass/db"))
		initPrivateKeysAndDb()
		vp := viper.New()
		vp.Set(vkeys.DatabasePath, viper.GetString(vkeys.DatabasePath))
		vp.Set(vkeys.PrivateKeysPath, viper.GetString(vkeys.PrivateKeysPath))

		jww.INFO.Printf("Writing config file to: %s", "~/.mypass.yaml")
		err := vp.WriteConfigAs(exHome("~/.mypass.yaml"))
		if err != nil {
			jww.ERROR.Fatalf("Failed to write config file: %v", err)
		}
	},
}

func expandedHome(p string) (string, error) {
	p = strings.TrimPrefix(p, "~/")
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, p), nil
}

func writeJSON(filepath string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		// TODO: handle failure case gracefully
		fmt.Println(string(data))
	}
	_, err = f.Write([]byte("\n"))
	return err
}

func init() {
	rootCmd.AddCommand(initCmd)
}
