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

package ansible

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// LoadPlaybook parses Ansible-style playbook YAML data and resolves tasks and roles.
func LoadPlaybook(
	data []byte,
	playbookPath string,
	rolesPath string,
) ([]Play, error) {
	var rawPlays []map[string]interface{}
	if err := yaml.Unmarshal(data, &rawPlays); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	parsedPlays := make([]Play, 0, len(rawPlays))
	for _, rawPlay := range rawPlays {
		play := Play{
			Name:  safeString(rawPlay["name"]),
			Hosts: safeString(rawPlay["hosts"]),
		}

		rawTasks, ok := rawPlay["tasks"].([]interface{})
		if !ok {
			continue
		}

		var rawTaskMaps []map[string]interface{}
		for _, t := range rawTasks {
			if tm, ok := t.(map[string]interface{}); ok {
				rawTaskMaps = append(rawTaskMaps, tm)
			}
		}

		parsedTasks, err := parseTasks(rawTaskMaps, playbookPath, rolesPath)
		if err != nil {
			return nil, err
		}

		for _, task := range parsedTasks {
			if task.Module == "include_role" || task.Module == "ansible.builtin.include_role" {
				roleName := safeString(task.RawArgs["name"])
				if roleName == "" {
					return nil, fmt.Errorf("include_role task is missing 'name': %+v", task.RawArgs)
				}

				roleTasks, err := LoadRoleTasks(roleName, rolesPath)
				if err != nil {
					return nil, fmt.Errorf("failed to load role %q: %w", roleName, err)
				}

				play.Tasks = append(play.Tasks, roleTasks...)
				continue
			}

			play.Tasks = append(play.Tasks, task)
		}

		parsedPlays = append(parsedPlays, play)
	}

	return parsedPlays, nil
}
