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

// Playbook represents a full Ansible playbook — a list of plays.
type Playbook []Play

// Play defines a single play in the playbook.
type Play struct {
	// Name is the descriptive name of the play (e.g., "Converge")
	Name string
	// Hosts defines the target hosts for this play (e.g., "all")
	Hosts string
	// Tasks is the ordered list of tasks to run in this play
	Tasks []Task
}

// Task represents an individual Ansible task.
type Task struct {
	// Name is the descriptive name of the task (e.g., "Install NGINX")
	Name string
	// Module is the name of the Ansible module invoked (e.g., "ansible.builtin.copy")
	Module string
	// RawArgs contains arguments passed to the module (raw, pre-processed)
	RawArgs map[string]interface{}
	// Vars holds optional vars specific to this task
	Vars map[string]interface{}
	// Loop holds the raw loop expression from the task (e.g., "{{ some_list }}").
	// Loop contains the raw loop expression from the task (e.g., "{{ some_list }}").
	Loop string
	// Source is the absolute or relative file path where this task was defined.
	Source string
}
