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

// Copyright (c) 2025 John Dewey
// SPDX-License-Identifier: MIT

package ansible_test

import (
	"testing"

	"github.com/kluctl/kluctl/lib/go-jinja2"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/voidspan/internal/ansible"
)

type RenderJinjaFieldsPublicTestSuite struct {
	suite.Suite

	renderer *jinja2.Jinja2
}

func (s *RenderJinjaFieldsPublicTestSuite) SetupSuite() {
	r, err := jinja2.NewJinja2("test-renderer", 1)
	s.Require().NoError(err)
	s.renderer = r
}

func (s *RenderJinjaFieldsPublicTestSuite) TearDownSuite() {
	s.renderer.Close()
}

func (s *RenderJinjaFieldsPublicTestSuite) TestRenderJinjaFields() {
	tests := []struct {
		name              string
		input             map[string]interface{}
		task              ansible.Task
		vars              map[string]interface{}
		expected          map[string]interface{}
		expectErr         bool
		expectErrContains string
	}{
		{
			name: "simple string interpolation",
			input: map[string]interface{}{
				"message": "hello {{ name }}",
			},
			vars: map[string]interface{}{
				"name": "world",
			},
			expected: map[string]interface{}{
				"message": "hello world",
			},
		},
		{
			name: "nested maps with interpolation",
			input: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "{{ foo }}",
				},
			},
			vars: map[string]interface{}{
				"foo": "bar",
			},
			expected: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "bar",
				},
			},
		},
		{
			name: "non-string values untouched",
			input: map[string]interface{}{
				"number": 42,
				"bool":   true,
			},
			vars: map[string]interface{}{},
			expected: map[string]interface{}{
				"number": 42,
				"bool":   true,
			},
		},
		{
			name: "missing_variable_renders_empty_string",
			task: ansible.Task{
				Name:   "test",
				Module: "ansible.builtin.debug",
				RawArgs: map[string]interface{}{
					"msg": "{{ not_set }}",
				},
			},
			vars:              map[string]interface{}{},
			expectErr:         true,
			expectErrContains: "UndefinedError",
		},
		{
			name: "invalid syntax returns error",
			input: map[string]interface{}{
				"oops": "{{ invalid",
			},
			vars:              map[string]interface{}{},
			expectErr:         true,
			expectErrContains: "unexpected end of template",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var input map[string]interface{}

			if tc.input != nil {
				input = tc.input
			} else if tc.task.RawArgs != nil {
				input = tc.task.RawArgs
			} else {
				s.FailNow("no input or task.RawArgs provided")
			}

			result, err := ansible.RenderJinjaFields(input, tc.vars, s.renderer)

			if tc.expectErr {
				s.Error(err)
				if tc.expectErrContains != "" {
					s.Contains(err.Error(), tc.expectErrContains)
				}
				return
			}

			s.Require().NoError(err)
			s.Equal(tc.expected, result)
		})
	}
}

func TestRenderJinjaFieldsPublicTestSuite(t *testing.T) {
	suite.Run(t, new(RenderJinjaFieldsPublicTestSuite))
}
