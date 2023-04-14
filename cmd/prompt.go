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
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
)

type FieldWithValue map[string]string

type FieldConfig struct {
	Mask       bool
	Default    string
	ValidateFn promptui.ValidateFunc
}

type FieldsWithConfig map[string]*FieldConfig

func Prompt(mp FieldsWithConfig, detailsTpl string) (FieldWithValue, error) {
	out := FieldWithValue{}
	type filedValue struct {
		Name  string
		Value FieldWithValue
	}
	fields := make([]filedValue, 0, len(mp))
	// Init out map
	for k, v := range mp {
		out[k] = v.Default
		fields = append(fields, filedValue{Name: k, Value: out})
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	const saveAndExit = "Save And Exit"
	fields = append(fields, filedValue{Name: saveAndExit, Value: out})
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "> {{ .Name | cyan }}",
		Inactive: "  {{ .Name  }}",
		Selected: "{{ .Name | red | cyan }}",
		Details:  detailsTpl,
	}
	searcher := func(input string, index int) bool {
		field := fields[index]
		name := strings.Replace(strings.ToLower(field.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
mainLoop:
	for {
		selectPrompt := promptui.Select{
			Label:        "Select to edit",
			Items:        fields,
			Templates:    templates,
			CursorPos:    len(fields) - 1,
			Size:         8,
			Searcher:     searcher,
			HideSelected: true,
			HideHelp:     false,
		}
		i, _, err := selectPrompt.Run()
		if fields[i].Name == saveAndExit {
			// validate all the fields except save and exit
			validateAllSuccess := true
			for idx, f := range fields {
				if f.Name == saveAndExit {
					continue
				}
				if mp[f.Name].ValidateFn(out[f.Name]) != nil {
					validateAllSuccess = false
					i = idx
					break
				}
			}
			if validateAllSuccess {
				break mainLoop
			}
		}
		if err != nil {
			return nil, fmt.Errorf("prompt failed: %w", err)
		}
		mi, ok := mp[fields[i].Name]
		if !ok {
			fmt.Println(">>> ", fields[i].Name, "not found")
			continue
		}
		validate := mi.ValidateFn

		templates := &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . | green }} ",
			Invalid: "{{ . | red }} ",
			Success: "{{ . | bold }} ",
		}
		fieldName := fields[i].Name
		mask := rune(0)
		if mp[fieldName].Mask {
			mask = '*'
		}
		prompt := promptui.Prompt{
			Label:       fieldName + ":",
			Templates:   templates,
			Validate:    validate,
			HideEntered: true,
			AllowEdit:   true,
			Default:     out[fieldName],
			Mask:        mask,
		}
		result, err := prompt.Run()
		if err != nil {
			return nil, fmt.Errorf("prompt failed: %w", err)
		}
		out[fieldName] = result
	}
	return out, nil
}
