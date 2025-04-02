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

package ansible_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/retr0h/voidspan/internal/ansible"
)

type LoadRoleTasksPublicTestSuite struct {
	suite.Suite

	tmpDir string
}

func (s *LoadRoleTasksPublicTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "voidspan-role-test-*")
	s.Require().NoError(err)
	s.tmpDir = dir
}

func (s *LoadRoleTasksPublicTestSuite) TearDownTest() {
	_ = os.RemoveAll(s.tmpDir)
}

func (s *LoadRoleTasksPublicTestSuite) TestLoadRoleTasks() {
	tests := []struct {
		name              string
		roleName          string
		expected          []ansible.Task
		expectErr         bool
		expectErrContains string
		prepare           func(roleDir string)
	}{
		{
			name:     "valid role task",
			roleName: "role1",
			expected: []ansible.Task{{
				Name:   "role1 | hello",
				Module: "ansible.builtin.debug",
				RawArgs: map[string]interface{}{
					"msg": "hi from role1",
				},
				Vars: map[string]interface{}{},
				Loop: "",
			}},
			prepare: func(roleDir string) {
				tasksDir := filepath.Join(roleDir, "tasks")
				_ = os.MkdirAll(tasksDir, 0o755)
				_ = os.WriteFile(filepath.Join(tasksDir, "main.yml"), []byte(`
---
- name: role1 | hello
  ansible.builtin.debug:
    msg: hi from role1
`), 0o644)
			},
		},
		{
			name:              "missing main.yml",
			roleName:          "missing",
			expectErr:         true,
			expectErrContains: "failed to read role tasks",
		},
		{
			name:              "invalid YAML in main.yml",
			roleName:          "badyaml",
			expectErr:         true,
			expectErrContains: "failed to parse tasks YAML",
			prepare: func(roleDir string) {
				tasksDir := filepath.Join(roleDir, "tasks")
				_ = os.MkdirAll(tasksDir, 0o755)
				_ = os.WriteFile(filepath.Join(tasksDir, "main.yml"), []byte(`
---
- name: bad
  ansible.builtin.debug
    msg: bad indentation
`), 0o644)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			roleDir := filepath.Join(s.tmpDir, tc.roleName)
			if tc.prepare != nil {
				tc.prepare(roleDir)
			}

			tasks, err := ansible.LoadRoleTasks(tc.roleName, s.tmpDir)

			if tc.expectErr {
				s.Error(err)
				if tc.expectErrContains != "" {
					s.Contains(err.Error(), tc.expectErrContains)
				}
				return
			}

			s.Require().NoError(err)
			s.Equal(len(tc.expected), len(tasks))

			for i := range tasks {
				exp := tc.expected[i]
				act := tasks[i]

				s.Equal(exp.Name, act.Name)
				s.Equal(exp.Module, act.Module)
				s.Equal(exp.Loop, act.Loop)
				s.Equal(exp.RawArgs, act.RawArgs)
				s.Equal(exp.Vars, act.Vars)
				s.NotEmpty(act.Source)
			}
		})
	}
}

func TestLoadRoleTasksPublicTestSuite(t *testing.T) {
	suite.Run(t, new(LoadRoleTasksPublicTestSuite))
}
