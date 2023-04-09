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
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/riadafridishibly/mypass/app"
	"github.com/riadafridishibly/mypass/models"
	"golang.design/x/clipboard"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cfg struct {
	ShowPassword bool
}

type itemWithConfig struct {
	Cfg *cfg
	*models.Item
}

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:     "select",
	Short:   "Select a password from a the list",
	Long:    `Search password or interactively select items here`,
	Aliases: []string{"list"},
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := app.NewApp()
		if err != nil {
			return err
		}
		c := &cfg{
			ShowPassword: viper.GetBool("select.show-password"),
		}
		itemsRaw := a.ListAllItems()
		var items []itemWithConfig
		for _, i := range itemsRaw {
			items = append(items, itemWithConfig{Cfg: c, Item: i})
		}
		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "> {{ .String | underline }}",
			Inactive: "  {{ .String | faint }}",
			Selected: "\U00002714 {{ .String | green}}",
			Details: `
{{ "ID:" | faint }}	{{ .ID }}
{{ "Title:" | faint }}	{{ .Title }}
{{- if .Password}}
{{ "Type:" | faint }}	{{ "password" }}
 {{"Username:" | faint}}  {{.Password.Username}}
 {{"SiteName:" | faint}}  {{.Password.SiteName}}
 {{"URL:" | faint}}       {{.Password.URL}}
{{- if .Cfg.ShowPassword}}
 {{"Password:" | faint}}  {{.Password.Password}}
{{end}}
{{end}}
{{- if .SSH}}
{{ "Type:" | faint }}	{{ "ssh" }}
 {{"Host:" | faint}}      {{.SSH.Host}}
 {{"Port:" | faint}}      {{.SSH.Port}}
 {{"Username:" | faint}}  {{.SSH.Username}}
{{- if .Cfg.ShowPassword}}
 {{"Password:" | faint}}  {{.SSH.Password}}
{{end}}
{{end}}`,
		}

		// TODO: Fork the modify search allow to rank results
		searcher := func(input string, index int) bool {
			singleItem := items[index]
			name := strings.ToLower(singleItem.Title)
			input = strings.ToLower(input)

			return strings.Contains(name, input)
		}
		prompt := promptui.Select{
			Label:        "Select one",
			Items:        items,
			Templates:    templates,
			Size:         10,
			Searcher:     searcher,
			HideSelected: true,
		}
		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return err
		}
		v, err := items[i].GetPassword()
		if err != nil {
			return err
		}
		clipboard.Write(clipboard.FmtText, []byte(v))
		fmt.Printf("Password copied for %q to clipboard.\n", items[i].InnerItemString())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)

	selectCmd.Flags().Bool("show-password", false, "show-password")
	viper.BindPFlag("select.show-password", selectCmd.Flags().Lookup("show-password"))
}
