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

// Package todo provides a type and several functions for interacting
// with todos backed by the Appengine Datastore. Todo entities are stored
// as descendants of the same parent, such that the lookups can be made
// strongly consistent.
package todo

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// getTodoGroupKey returns a key representing the parent of all Todo entities.
func getTodoGroupKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "TodoGroup", "", 1, nil)
}

type Todo struct {
	ID    int64  `datastore:"-" json:"id"` // Unique identifier
	Title string `json:"title"`            // Description
	Done  bool   `json:"completed"`        // Is this todo done?
}

// Save stores the todo entity in the datastore with the proper parent key.
// It returns an error if one was generated during storage. Otherwise, the
// todo's ID field is populated as a side effect.
func (t *Todo) Save(c context.Context) error {
	k := datastore.NewKey(c, "Todo", "", t.ID, getTodoGroupKey(c))
	k, err := datastore.Put(c, k, t)
	if err == nil {
		t.ID = k.IntID()
	}
	return err
}

// All returns a slice containing all todos in the datastore and an error, if one
// occurred. If there was no error, all todos will have their ID field
// initialized.
//
// The underlying query is strongly consistent.
func All(c context.Context) ([]*Todo, error) {
	var todos []*Todo
	q := datastore.NewQuery("Todo").Ancestor(getTodoGroupKey(c))
	keys, err := q.GetAll(c, &todos)
	if err == nil {
		for i := 0; i < len(keys); i++ {
			todos[i].ID = keys[i].IntID()
		}
	}
	return todos, err
}

// Get returns the todo with the given id and a boolean indicating if it
// was found.
//
// The underlying query is strongly consistent.
func Get(c context.Context, id int64) (*Todo, bool) {
	k := datastore.NewKey(c, "Todo", "", id, getTodoGroupKey(c))
	todo := &Todo{}
	err := datastore.Get(c, k, todo)
	if err != nil {
		todo.ID = k.IntID()
	}
	return todo, err == nil
}

// Delete removes the todo with the given id from the datastore. It returns
// a boolean indicating if we were successful.
func Delete(c context.Context, id int64) bool {
	k := datastore.NewKey(c, "Todo", "", id, getTodoGroupKey(c))
	return datastore.Delete(c, k) == nil
}

// DeleteCompleted removes all finished todos from the datastore.
func DeleteCompleted(c context.Context) error {
	q := datastore.NewQuery("Todo").Ancestor(
		getTodoGroupKey(c)).KeysOnly().Filter("Done =", true)
	keys, err := q.GetAll(c, nil)
	if err != nil {
		return err
	}
	return datastore.DeleteMulti(c, keys)
}

// NewTodo creates a new todo given the given title, which can't be the empty
// string.
func NewTodo(title string) (*Todo, error) {
	if title == "" {
		return nil, fmt.Errorf("empty title")
	}
	return &Todo{0, title, false}, nil
}
