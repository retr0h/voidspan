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

// LoadRoleTasks loads tasks from a given role's tasks/main.yml file
func LoadRoleTasks(
	roleName string,
	rolesPath string,
) ([]Task, error) {
	roleDir := filepath.Join(rolesPath, roleName)
	tasksPath := filepath.Join(roleDir, "tasks", "main.yml")

	data, err := os.ReadFile(tasksPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read role tasks: %w", err)
	}

	var rawTasks []map[string]interface{}
	if err := yaml.Unmarshal(data, &rawTasks); err != nil {
		return nil, fmt.Errorf("failed to parse tasks YAML: %w", err)
	}

	return parseTasks(rawTasks, tasksPath, rolesPath)
}
