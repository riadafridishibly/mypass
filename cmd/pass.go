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
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/riadafridishibly/mypass/backend"
	"github.com/riadafridishibly/mypass/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

func newDefaultPasswordItem() *models.Item {
	return &models.Item{
		Title:     "Password Item",
		Namespace: "default",
		Type:      "password",
		Meta: models.Meta{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Password: &models.PasswordItem{
			Username: "",
			SiteName: "",
			URL:      "",
			Password: "",
		},
	}
}

func newPassFieldsWithConfig(i *models.Item) FieldsWithConfig {
	return FieldsWithConfig{
		"Title": &FieldConfig{
			Default:    i.Title,
			ValidateFn: func(string) error { return nil },
		},
		"Namespace": &FieldConfig{
			Default: i.Namespace,
			ValidateFn: func(s string) error {
				if s == "" {
					return errors.New("namespace can't be empty")
				}
				return nil
			},
		},
		"Username": &FieldConfig{
			Default: i.Password.Username,
			ValidateFn: func(s string) error {
				if s == "" {
					return errors.New("username can't be empty")
				}
				return nil
			},
		},
		"SiteName": &FieldConfig{
			Default: i.Password.SiteName,
			ValidateFn: func(s string) error {
				if s == "" {
					return errors.New("sitename can't be empty")
				}
				return nil
			},
		},
		"URL": &FieldConfig{
			Default: i.Password.URL,
			ValidateFn: func(s string) error {
				return nil
			},
		},
		"Password": &FieldConfig{
			Mask:    true,
			Default: string(i.Password.Password),
			ValidateFn: func(s string) error {
				return nil
			},
		},
	}
}

const passDetailsTpl = `
--------- Password Credential ----------
{{ "Title:" | faint }}	{{ .Value.Title }}
{{ "Namespace:" | faint }}	{{ .Value.Namespace }}
{{ "Username:" | faint }}	{{ .Value.Username }}
{{ "SiteName:" | faint }}	{{ .Value.SiteName }}
{{ "URL:" | faint }}	{{ .Value.URL }}`

// passCmd represents the pass command
var passCmd = &cobra.Command{
	Use:   "pass",
	Short: "Add password item",
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := Prompt(
			newPassFieldsWithConfig(newDefaultPasswordItem()),
			passDetailsTpl,
		)
		if err != nil {
			return err
		}
		pass := v["Password"]
		if pass == "" {
			// try read password from stdin
			fmt.Print("Enter password:")
			data, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			fmt.Println()
			pass = string(data)
		}

		p := &models.PasswordItem{
			Username: v["Username"],
			SiteName: v["SiteName"],
			URL:      v["URL"],
			Password: models.AsymSecretStr(pass),
		}

		a, err := backend.New()
		if err != nil {
			return err
		}
		_, err = a.CreateItem(&models.Item{
			Title:     v["Title"],
			Namespace: v["Namespace"],
			Password:  p,
		})
		if err != nil {
			return err
		}
		return a.Flush()
	},
}

func init() {
	addCmd.AddCommand(passCmd)
	passCmd.Flags().String("username", "", "Username")
	viper.BindPFlag("pass.username", passCmd.Flags().Lookup("username"))

	passCmd.Flags().String("password", "", "Password (not recommended, use stdin)")
	viper.BindPFlag("pass.password", passCmd.Flags().Lookup("password"))

	passCmd.Flags().String("site", "", "Site host name. eg. gmail.com, github.com")
	viper.BindPFlag("pass.site", passCmd.Flags().Lookup("site"))

	passCmd.Flags().String("url", "", "Site login url")
	viper.BindPFlag("pass.url", passCmd.Flags().Lookup("url"))

	// cobra.MarkFlagRequired(passCmd.Flags(), "username")
	// cobra.MarkFlagRequired(passCmd.Flags(), "site")
}
