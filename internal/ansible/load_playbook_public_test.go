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

type LoadPlaybookPublicTestSuite struct {
	suite.Suite

	tmpDir string
}

func (s *LoadPlaybookPublicTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "voidspan-test-*")
	s.Require().NoError(err)
	s.tmpDir = dir
}

func (s *LoadPlaybookPublicTestSuite) TearDownTest() {
	_ = os.RemoveAll(s.tmpDir)
}

func (s *LoadPlaybookPublicTestSuite) TestLoadPlaybook() {
	tests := []struct {
		name              string
		playbookYAML      string
		expected          []ansible.Play
		expectErr         bool
		expectErrContains string
		prepare           func(dir string)
	}{
		{
			name: "basic inline debug",
			playbookYAML: `
---
- name: test play
  hosts: all
  tasks:
    - name: inline | debug hello
      ansible.builtin.debug:
        msg: "hello world"
`,
			expected: []ansible.Play{{
				Name:  "test play",
				Hosts: "all",
				Tasks: []ansible.Task{{
					Name:   "inline | debug hello",
					Module: "ansible.builtin.debug",
					RawArgs: map[string]interface{}{
						"msg": "hello world",
					},
					Vars: map[string]interface{}{},
					Loop: "",
				}},
			}},
		},
		{
			name: "inline with loop",
			playbookYAML: `
---
- name: loop play
  hosts: all
  tasks:
    - name: inline | debug loop
      ansible.builtin.debug:
        msg: "looping {{ item }}"
      loop: "{{ ['one', 'two'] }}"
`,
			expected: []ansible.Play{{
				Name:  "loop play",
				Hosts: "all",
				Tasks: []ansible.Task{{
					Name:   "inline | debug loop",
					Module: "ansible.builtin.debug",
					RawArgs: map[string]interface{}{
						"msg": "looping {{ item }}",
					},
					Vars: map[string]interface{}{},
					Loop: "{{ ['one', 'two'] }}",
				}},
			}},
		},
		{
			name: "invalid YAML",
			playbookYAML: `
---
- name: bad play
  hosts: all
  tasks:
    - name: broken
      ansible.builtin.debug
        msg: "missing colon above"
`,
			expectErr:         true,
			expectErrContains: "failed to parse YAML",
		},
		{
			name: "invalid tasks type",
			playbookYAML: `
---
- name: bad task play
  hosts: all
  tasks: this is not a list
`,
			expected: []ansible.Play{},
		},
		{
			name: "parseTasks returns error (bad include_tasks)",
			playbookYAML: `
---
- name: test play
  hosts: all
  tasks:
    - name: bad task
      ansible.builtin.include_tasks:
        foo: bar
`,
			expectErr:         true,
			expectErrContains: "include_tasks path must be a string",
		},
		{
			name: "include_role missing name",
			playbookYAML: `
---
- name: test play
  hosts: all
  tasks:
    - name: include role without name
      ansible.builtin.include_role: {}
`,
			expectErr:         true,
			expectErrContains: "include_role task is missing 'name'",
		},
		{
			name: "include_role invalid role path",
			playbookYAML: `
---
- name: test play
  hosts: all
  tasks:
    - name: include role that does not exist
      ansible.builtin.include_role:
        name: not_a_real_role
`,
			expectErr:         true,
			expectErrContains: "failed to load role",
		},
		{
			name: "include_role loads tasks",
			playbookYAML: `
---
- name: test play
  hosts: all
  tasks:
    - name: include real role
      ansible.builtin.include_role:
        name: myrole
`,
			expected: []ansible.Play{{
				Name:  "test play",
				Hosts: "all",
				Tasks: []ansible.Task{{
					Name:   "role | main | test",
					Module: "ansible.builtin.debug",
					RawArgs: map[string]interface{}{
						"msg": "from role",
					},
					Vars: map[string]interface{}{},
					Loop: "",
				}},
			}},
			prepare: func(dir string) {
				roleDir := filepath.Join(dir, "roles", "myrole", "tasks")
				_ = os.MkdirAll(roleDir, 0o755)
				_ = os.WriteFile(filepath.Join(roleDir, "main.yml"), []byte(`
---
- name: role | main | test
  ansible.builtin.debug:
    msg: "from role"
`), 0o644)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			if tc.prepare != nil {
				tc.prepare(s.tmpDir)
			}

			playbookPath := filepath.Join(s.tmpDir, "playbook.yml")
			err := os.WriteFile(playbookPath, []byte(tc.playbookYAML), 0o644)
			s.Require().NoError(err)

			data, err := os.ReadFile(playbookPath)
			s.Require().NoError(err)

			actual, err := ansible.LoadPlaybook(
				data,
				playbookPath,
				filepath.Join(s.tmpDir, "roles"),
			)

			if tc.expectErr {
				s.Error(err)
				if tc.expectErrContains != "" {
					s.Contains(err.Error(), tc.expectErrContains)
				}
				return
			}

			s.Require().NoError(err)
			s.Equal(len(tc.expected), len(actual))

			for i := range actual {
				s.Equal(tc.expected[i].Name, actual[i].Name)
				s.Equal(tc.expected[i].Hosts, actual[i].Hosts)
				s.Equal(len(tc.expected[i].Tasks), len(actual[i].Tasks))

				for j := range actual[i].Tasks {
					exp := tc.expected[i].Tasks[j]
					act := actual[i].Tasks[j]

					s.Equal(exp.Name, act.Name)
					s.Equal(exp.Module, act.Module)
					s.Equal(exp.Loop, act.Loop)
					s.Equal(exp.RawArgs, act.RawArgs)
					s.Equal(exp.Vars, act.Vars)
					s.NotEmpty(act.Source)
				}
			}
		})
	}
}

func TestLoadPlaybookPublicTestSuite(t *testing.T) {
	suite.Run(t, new(LoadPlaybookPublicTestSuite))
}
