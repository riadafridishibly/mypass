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
	"os"

	"github.com/riadafridishibly/mypass/backend"
	"github.com/riadafridishibly/mypass/config"
	"github.com/riadafridishibly/mypass/vkeys"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"golang.design/x/clipboard"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mypass",
	Short: "A dead simple password manager",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetBool("verbose") {
			jww.SetStdoutThreshold(jww.LevelDebug)
		}
		// If command is not init then load the database
		if cmd.CalledAs() != "init" {
			b, err := backend.Get()
			if err != nil {
				return err
			}
			pubKeys, err := b.PublicKeys()
			if err != nil {
				jww.ERROR.Println("failed to load public keys")
				return err
			}
			jww.INFO.Println("loaded public keys: ", pubKeys)
			viper.Set(vkeys.PublicKeys, pubKeys)
		}
		// TODO: we may not need to load the clipboard,
		// if the command is not select
		err := clipboard.Init()
		if err != nil {
			// jww.FATAL.Fatal("Failed to initialize keyboard: ", err)
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mypass.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "Verbose mode")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

var DefaultConfigPath = config.ExpandWithHome("~/.mypass.yaml")

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mypass" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mypass")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		jww.INFO.Println("Using config file:", viper.ConfigFileUsed())
	}
}
