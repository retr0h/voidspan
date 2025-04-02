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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type ParseTasksTestSuite struct {
	suite.Suite

	tmpDir string
}

func (s *ParseTasksTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "voidspan-parse-tasks-*")
	s.Require().NoError(err)
	s.tmpDir = dir
}

func (s *ParseTasksTestSuite) TearDownTest() {
	_ = os.RemoveAll(s.tmpDir)
}

func (s *ParseTasksTestSuite) TestParseTasks() {
	tests := []struct {
		name              string
		taskYAML          string
		expected          []Task
		expectErr         bool
		expectErrContains string
		prepare           func(roleDir string)
	}{
		{
			name: "basic task with debug",
			taskYAML: `
- name: say hello
  ansible.builtin.debug:
    msg: "hello"
`,
			expected: []Task{{
				Name:   "say hello",
				Module: "ansible.builtin.debug",
				RawArgs: map[string]interface{}{
					"msg": "hello",
				},
				Vars: map[string]interface{}{},
				Loop: "",
			}},
		},
		{
			name: "task with vars",
			taskYAML: `
- name: task with vars
  ansible.builtin.shell:
    cmd: "echo {{ message }}"
  vars:
    message: hello
`,
			expected: []Task{{
				Name:   "task with vars",
				Module: "ansible.builtin.shell",
				RawArgs: map[string]interface{}{
					"cmd": "echo {{ message }}",
				},
				Vars: map[string]interface{}{
					"message": "hello",
				},
				Loop: "",
			}},
		},
		{
			name: "task with loop",
			taskYAML: `
- name: looping
  ansible.builtin.debug:
    msg: "{{ item }}"
  loop: "{{ ['a', 'b'] }}"
`,
			expected: []Task{{
				Name:   "looping",
				Module: "ansible.builtin.debug",
				RawArgs: map[string]interface{}{
					"msg": "{{ item }}",
				},
				Vars: map[string]interface{}{},
				Loop: "{{ ['a', 'b'] }}",
			}},
		},
		{
			name: "include_tasks with bad value",
			taskYAML: `
- name: bad include
  ansible.builtin.include_tasks:
    foo: bar
`,
			expectErr:         true,
			expectErrContains: "include_tasks path must be a string",
		},
		{
			name: "include_tasks missing file",
			taskYAML: `
- name: include file that doesn't exist
  ansible.builtin.include_tasks: not_here.yml
`,
			expectErr:         true,
			expectErrContains: "failed to read included task file",
		},
		{
			name: "include_tasks with valid file",
			taskYAML: `
- name: include real tasks
  ansible.builtin.include_tasks: included.yml
`,
			expected: []Task{{
				Name:   "included | say hi",
				Module: "ansible.builtin.debug",
				RawArgs: map[string]interface{}{
					"msg": "hi from included",
				},
				Vars: map[string]interface{}{},
				Loop: "",
			}},
			prepare: func(dir string) {
				err := os.WriteFile(filepath.Join(dir, "included.yml"), []byte(`
- name: included | say hi
  ansible.builtin.debug:
    msg: hi from included
`), 0o644)
				if err != nil {
					panic(err)
				}
			},
		},
		{
			name: "include_tasks with invalid YAML",
			taskYAML: `
- name: include broken file
  ansible.builtin.include_tasks: included.yml
`,
			expectErr:         true,
			expectErrContains: "failed to parse included task file",
			prepare: func(dir string) {
				err := os.WriteFile(filepath.Join(dir, "included.yml"), []byte(`
- name: bad
  ansible.builtin.debug
    msg: "oops
`), 0o644)
				if err != nil {
					panic(err)
				}
			},
		},
		{
			name: "include_tasks nested failure",
			taskYAML: `
- name: include bad nested file
  ansible.builtin.include_tasks: included.yml
`,
			expectErr:         true,
			expectErrContains: "include_tasks path must be a string",
			prepare: func(dir string) {
				err := os.WriteFile(filepath.Join(dir, "included.yml"), []byte(`
- name: invalid nested include
  ansible.builtin.include_tasks:
    foo: bar
`), 0o644)
				if err != nil {
					panic(err)
				}
			},
		},
		{
			name: "module with short-form string arg",
			taskYAML: `
- name: run shell
  ansible.builtin.shell: echo hello
`,
			expected: []Task{{
				Name:   "run shell",
				Module: "ansible.builtin.shell",
				RawArgs: map[string]interface{}{
					"__value__": "echo hello",
				},
				Vars: map[string]interface{}{},
				Loop: "",
			}},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			if tc.prepare != nil {
				tc.prepare(s.tmpDir)
			}

			var rawTasks []map[string]interface{}
			err := yaml.Unmarshal([]byte(tc.taskYAML), &rawTasks)
			s.Require().NoError(err)

			sourcePath := filepath.Join(s.tmpDir, "source.yml")
			tasks, err := parseTasks(rawTasks, sourcePath, filepath.Join(s.tmpDir, "roles"))

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
				s.Equal(exp.RawArgs, act.RawArgs)
				s.Equal(exp.Vars, act.Vars)
				s.Equal(exp.Loop, act.Loop)
				s.NotEmpty(act.Source)
			}
		})
	}
}

func TestParseTasksSuite(t *testing.T) {
	suite.Run(t, new(ParseTasksTestSuite))
}
