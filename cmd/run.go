// Copyright (c) 2025 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/retr0h/voidspan/internal/ansible"
)

var runCmd = &cobra.Command{
	Use:   "run [playbook.yml]",
	Short: "Run an Ansible-style playbook with Voidspan",
	Run: func(_ *cobra.Command, _ []string) {
		playbookPath := viper.GetString("playbook")
		rolesPath := viper.GetString("roles-path")

		data, err := os.ReadFile(playbookPath)
		if err != nil {
			log.Fatalf("failed to read playbook: %v", err)
		}

		plays, err := ansible.LoadPlaybook(data, playbookPath, rolesPath)
		if err != nil {
			log.Fatalf("failed to parse playbook: %v", err)
		}

		for _, play := range plays {
			fmt.Printf("▶ Play: %s (hosts: %s)\n", play.Name, play.Hosts)
			for _, task := range play.Tasks {
				fmt.Printf("  ▸ Task: %s\n", task.Name)
				fmt.Printf("    Module: %s\n", task.Module)
				fmt.Printf("    Args: %+v\n", task.RawArgs)
				fmt.Printf("    Vars: %+v\n", task.Vars)
				if task.Loop != "" {
					fmt.Printf("    Loop: %s\n", task.Loop)
				}

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().
		StringP("roles-path", "r", "roles", "Path to the base directory containing Ansible roles")
	runCmd.PersistentFlags().
		StringP("playbook", "p", "playbook.yaml", "Path to the Ansible playbook file to parse and run")

	_ = viper.BindPFlag("playbook", runCmd.PersistentFlags().Lookup("playbook"))
	_ = viper.BindPFlag("roles-path", runCmd.PersistentFlags().Lookup("roles-path"))

	_ = runCmd.MarkPersistentFlagRequired("playbook")
	_ = runCmd.MarkPersistentFlagRequired("roles-path")
}
