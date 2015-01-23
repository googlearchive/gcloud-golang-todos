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

package todo

import "testing"

func newTodoOrFatal(t *testing.T, title string) *Todo {
	todo, err := NewTodo(title)
	if err != nil {
		t.Fatalf("new todo: %v", err)
	}
	return todo
}

func TestNewTodo(t *testing.T) {
	title := "learn Go"
	todo := newTodoOrFatal(t, title)
	if todo.Title != title {
		t.Errorf("expected title %q, got %q", title, todo.Title)
	}
	if todo.Done {
		t.Errorf("new todo is done")
	}
}

func TestNewTodoEmptyTitle(t *testing.T) {
	_, err := NewTodo("")
	if err == nil {
		t.Errorf("expected 'empty title' error, got nil")
	}
}

func TestSaveTodoAndRetrieve(t *testing.T) {
	todo := newTodoOrFatal(t, "learn Go")

	m := NewTodoManager()
	m.Save(todo)

	all := m.All()
	if len(all) != 1 {
		t.Errorf("expected 1 todo, got %v", len(all))
	}
	if *all[0] != *todo {
		t.Errorf("expected %v, got %v", todo, all[0])
	}
}

func TestSaveAndRetrieveTwoTodos(t *testing.T) {
	learnGo := newTodoOrFatal(t, "learn Go")
	learnTDD := newTodoOrFatal(t, "learn TDD")

	m := NewTodoManager()
	m.Save(learnGo)
	m.Save(learnTDD)

	all := m.All()
	if len(all) != 2 {
		t.Errorf("expected 2 todos, got %v", len(all))
	}
	if *all[0] != *learnGo && *all[1] != *learnGo {
		t.Errorf("missing todo: %v", learnGo)
	}
	if *all[0] != *learnTDD && *all[1] != *learnTDD {
		t.Errorf("missing todo: %v", learnTDD)
	}
}

func TestSaveModifyAndRetrieve(t *testing.T) {
	todo := newTodoOrFatal(t, "learn Go")
	m := NewTodoManager()
	m.Save(todo)

	todo.Done = true
	if m.All()[0].Done {
		t.Errorf("saved todo wasn't done")
	}
}

func TestSaveTwiceAndRetrieve(t *testing.T) {
	todo := newTodoOrFatal(t, "learn Go")
	m := NewTodoManager()
	m.Save(todo)
	m.Save(todo)

	all := m.All()
	if len(all) != 1 {
		t.Errorf("expected 1 todo, got %v", len(all))
	}
	if *all[0] != *todo {
		t.Errorf("expected todo %v, got %v", todo, all[0])
	}
}

func TestSaveAndFind(t *testing.T) {
	todo := newTodoOrFatal(t, "learn Go")
	m := NewTodoManager()
	m.Save(todo)

	nt, ok := m.Find(todo.ID)
	if !ok {
		t.Errorf("didn't find todo")
	}
	if *todo != *nt {
		t.Errorf("expected %v, got %v", todo, nt)
	}
}
