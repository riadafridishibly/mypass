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
	"github.com/riadafridishibly/mypass/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new items to the database",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if fn := rootCmd.PersistentPreRunE; fn != nil {
			err := fn(cmd, args)
			if err != nil {
				return err
			}
		}
		// FIXME: We don't need to load private keys when we add items,
		// currently encryption and decryption is done in MarshalJSON and UnmarshalJSON methods.
		// In sqlite backend, we call GetItemByID after inserting it, that causes a call to
		// UnmarshalJSON, which will return error if PrivateKeys are not already loaded in viper.
		// We may introduce Pre and Post hook, That will solve the issue.
		err := config.LoadCachedPassword()
		if err != nil {
			return err
		}
		err = config.LoadPrivateKeys()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().String("title", "", "Item title")
	viper.BindPFlag("add.title", addCmd.PersistentFlags().Lookup("title"))

	addCmd.PersistentFlags().String("namespace", "default", "Namespace")
	viper.BindPFlag("add.namespace", addCmd.PersistentFlags().Lookup("namespace"))

	addCmd.PersistentFlags().BoolP("interactive", "i", false, "Interactive mode")
	viper.BindPFlag("add.interactive", addCmd.PersistentFlags().Lookup("interactive"))
}
