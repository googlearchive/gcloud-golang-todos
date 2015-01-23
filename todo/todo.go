// Copyright 2014 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

// The package provides a TodoManager implemented using TDD techniques.
// The tests were developed before the code was written.
package todo

import "fmt"

type Todo struct {
	ID    int64  `json:"id"`        // Unique identifier
	Title string `json:"title"`     // Description
	Done  bool   `json:"completed"` // Is this todo done?
}

// NewTodo creates a new todo given a title, that can't be empty.
func NewTodo(title string) (*Todo, error) {
	if title == "" {
		return nil, fmt.Errorf("empty title")
	}
	return &Todo{0, title, false}, nil
}

// TodoManager manages a list of todos in memory.
type TodoManager struct {
	todos  []*Todo
	lastID int64
}

// NewTodoManager returns an empty TodoManager.
func NewTodoManager() *TodoManager {
	return &TodoManager{}
}

// Save saves the given Todo in the TodoManager.
func (m *TodoManager) Save(todo *Todo) error {
	if todo.ID == 0 {
		m.lastID++
		todo.ID = m.lastID
		m.todos = append(m.todos, cloneTodo(todo))
		return nil
	}

	for i, t := range m.todos {
		if t.ID == todo.ID {
			m.todos[i] = cloneTodo(todo)
			return nil
		}
	}
	return fmt.Errorf("unknown todo")
}

// cloneTodo creates and returns a deep copy of the given Todo.
func cloneTodo(t *Todo) *Todo {
	c := *t
	return &c
}

// All returns the list of all the Todos in the TodoManager.
func (m *TodoManager) All() []*Todo {
	return m.todos
}

// Find returns the Todo with the given id in the TodoManager and a boolean
// indicating if the id was found.
func (m *TodoManager) Find(ID int64) (*Todo, bool) {
	for _, t := range m.todos {
		if t.ID == ID {
			return t, true
		}
	}
	return nil, false
}
