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

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pass := viper.GetString("ssh.password")
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

		p := &models.SSHItem{
			Host:     viper.GetString("ssh.host"),
			Port:     viper.GetUint16("ssh.port"),
			Username: viper.GetString("ssh.username"),
			Password: models.AsymSecretStr(pass),
		}

		a, err := app.NewApp()
		if err != nil {
			return err
		}
		err = a.AddItem(&models.Item{
			Title:     viper.GetString("add.title"),
			Namespace: viper.GetString("add.namespace"),
			SSH:       p,
		})
		if err != nil {
			return err
		}
		return writeJSON(viper.GetString(vkeys.DatabasePath), a.DB)
	},
}

func init() {
	addCmd.AddCommand(sshCmd)

	sshCmd.Flags().String("username", "", "Username")
	viper.BindPFlag("ssh.username", sshCmd.Flags().Lookup("username"))

	sshCmd.Flags().String("password", "", "Password (not recommended, use stdin)")
	viper.BindPFlag("ssh.password", sshCmd.Flags().Lookup("password"))

	sshCmd.Flags().String("host", "", "Site host name. eg. example.com")
	viper.BindPFlag("ssh.host", sshCmd.Flags().Lookup("host"))

	sshCmd.Flags().Uint16("port", 22, "Port")
	viper.BindPFlag("ssh.port", sshCmd.Flags().Lookup("port"))

	cobra.MarkFlagRequired(sshCmd.Flags(), "username")
	cobra.MarkFlagRequired(sshCmd.Flags(), "host")
}
