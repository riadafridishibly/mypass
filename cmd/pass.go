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
	"fmt"
	"os"

	"github.com/riadafridishibly/mypass/app"
	"github.com/riadafridishibly/mypass/models"
	"github.com/riadafridishibly/mypass/vkeys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// passCmd represents the pass command
var passCmd = &cobra.Command{
	Use:   "pass",
	Short: "Add password item",
	RunE: func(cmd *cobra.Command, args []string) error {
		pass := viper.GetString("pass.password")
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
			Username: viper.GetString("pass.username"),
			SiteName: viper.GetString("pass.site"),
			URL:      viper.GetString("pass.url"),
			Password: models.AsymSecretStr(pass),
		}

		a, err := app.NewApp()
		if err != nil {
			return err
		}
		err = a.AddItem(&models.Item{
			Title:     viper.GetString("add.title"),
			Namespace: viper.GetString("add.namespace"),
			Password:  p,
		})
		if err != nil {
			return err
		}
		return writeJSON(viper.GetString(vkeys.DatabasePath), a.DB)
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

	cobra.MarkFlagRequired(passCmd.Flags(), "username")
	cobra.MarkFlagRequired(passCmd.Flags(), "site")
}
