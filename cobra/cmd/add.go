// Copyright © 2015 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"path"
	"unicode"

	"github.com/spf13/cobra"
)

var (
	packageName string
	parentName  string

	addCmd = &cobra.Command{
		Use:     "add [command name]",
		Aliases: []string{"command"},
		Short:   "Add a command to a Cobra Application",
		Long: `Add (cobra add) will create a new command, with a license and
the appropriate structure for a Cobra-based CLI application,
and register it to its parent (default rootCmd).

If you want your command to be public, pass in the command name
with an initial uppercase letter.

Example: cobra add server -> resulting in a new cmd/server.go`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				er("add needs a name for the command")
			}

			commandName := validateCmdName(args[0])

			if packageName == "" {
				// derive packageName
			}

			wd, err := os.Getwd()
			if err != nil {
				er(err)
			}

			command := &Command{
				CmdName:   commandName,
				CmdParent: parentName,
				Project: &Project{
					AbsolutePath: fmt.Sprintf("%s/cmd", wd),
					AppName:      path.Base(packageName),
					PkgName:      packageName,
					Legal:        getLicense(),
					Copyright:    copyrightLine(),
				},
			}

			err = command.Create()
			if err != nil {
				er(err)
			}

			fmt.Printf("%s created at %s", command.CmdName, command.Project.AbsolutePath)
		},
	}
)

func init() {
	addCmd.Flags().StringVarP(&packageName, "package", "t", "", "target package name (e.g. github.com/spf13/hugo)")
	addCmd.Flags().StringVarP(&parentName, "parent", "p", "rootCmd", "variable name of parent command for this command")
}

// validateCmdName returns source without any dashes and underscore.
// If there will be dash or underscore, next letter will be uppered.
// It supports only ASCII (1-byte character) strings.
// https://github.com/spf13/cobra/issues/269
func validateCmdName(source string) string {
	i := 0
	l := len(source)
	// The output is initialized on demand, then first dash or underscore
	// occurs.
	var output string

	for i < l {
		if source[i] == '-' || source[i] == '_' {
			if output == "" {
				output = source[:i]
			}

			// If it's last rune and it's dash or underscore,
			// don't add it output and break the loop.
			if i == l-1 {
				break
			}

			// If next character is dash or underscore,
			// just skip the current character.
			if source[i+1] == '-' || source[i+1] == '_' {
				i++
				continue
			}

			// If the current character is dash or underscore,
			// upper next letter and add to output.
			output += string(unicode.ToUpper(rune(source[i+1])))
			// We know, what source[i] is dash or underscore and source[i+1] is
			// uppered character, so make i = i+2.
			i += 2
			continue
		}

		// If the current character isn't dash or underscore,
		// just add it.
		if output != "" {
			output += string(source[i])
		}
		i++
	}

	if output == "" {
		return source // source is initially valid name.
	}
	return output
}
