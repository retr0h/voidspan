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
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// parseTasks resolves raw Ansible task maps (including include_tasks) into typed Task objects.
func parseTasks(
	rawTasks []map[string]interface{},
	sourcePath string,
	rolesPath string,
) ([]Task, error) {
	tasks := make([]Task, 0, len(rawTasks))
	baseDir := filepath.Dir(sourcePath)

	for _, taskMap := range rawTasks {
		task := Task{
			Name:    safeString(taskMap["name"]),
			Vars:    make(map[string]interface{}),
			RawArgs: make(map[string]interface{}),
			Source:  sourcePath,
		}

		for k, v := range taskMap {
			switch k {
			case "name":
				// already handled
			case "vars":
				if varMap, ok := v.(map[string]interface{}); ok {
					task.Vars = varMap
				}
			case "loop":
				if loopStr, ok := v.(string); ok {
					task.Loop = loopStr
				}
			case "include_tasks", "ansible.builtin.include_tasks":
				includePath, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("include_tasks path must be a string")
				}

				fullPath := filepath.Join(baseDir, includePath)
				data, err := os.ReadFile(fullPath)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to read included task file %s: %w",
						fullPath,
						err,
					)
				}

				var includedRaw []map[string]interface{}
				if err := yaml.Unmarshal(data, &includedRaw); err != nil {
					return nil, fmt.Errorf(
						"failed to parse included task file %s: %w",
						fullPath,
						err,
					)
				}

				includedTasks, err := parseTasks(includedRaw, fullPath, rolesPath)
				if err != nil {
					return nil, err
				}

				tasks = append(tasks, includedTasks...)
				goto nextTask // skip appending this include as a normal task
			default:
				task.Module = k
				switch val := v.(type) {
				case map[string]interface{}:
					task.RawArgs = val
				default:
					task.RawArgs = map[string]interface{}{"__value__": val}
				}
			}
		}

		tasks = append(tasks, task)

	nextTask:
	}

	return tasks, nil
}
