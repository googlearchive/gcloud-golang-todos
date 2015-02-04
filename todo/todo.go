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

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Todo struct {
	ID    int64  `datastore:"-" json:"id"` // Unique identifier
	Title string `json:"title"`            // Description
	Done  bool   `json:"completed"`        // Is this todo done?
}

func (t *Todo) Save(c context.Context) error {
	k := datastore.NewKey(c, "Todo", "", t.ID, nil)
	k, err := datastore.Put(c, k, t)
	if err == nil {
		t.ID = k.IntID()
	}
	return err
}

func All(c context.Context) ([]*Todo, error) {
	var todos []*Todo
	keys, err := datastore.NewQuery("Todo").GetAll(c, &todos)
	if err == nil {
		for i := 0; i < len(keys); i++ {
			todos[i].ID = keys[i].IntID()
		}
	}
	return todos, err
}

// Get returns the Todo with the given id and a boolean indicating if it
// was found.
func Get(c context.Context, id int64) (*Todo, bool) {
	k := datastore.NewKey(c, "Todo", "", id, nil)
	todo := &Todo{}
	err := datastore.Get(c, k, todo)
	if err != nil {
		todo.ID = k.IntID()
	}
	return todo, err == nil
}

func Delete(c context.Context, id int64) bool {
	k := datastore.NewKey(c, "Todo", "", id, nil)
	return datastore.Delete(c, k) == nil
}

// NewTodo creates a new todo given the given title, which can't be the empty
// string.
func NewTodo(title string) (*Todo, error) {
	if title == "" {
		return nil, fmt.Errorf("empty title")
	}
	return &Todo{0, title, false}, nil
}
