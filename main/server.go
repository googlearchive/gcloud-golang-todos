// Copyright 2015 Google Inc. All rights reserved.
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

// +build appenginevm
// This package implements a simple HTTP server providing a REST API to a todo
// handler.
//
// It provides six methods:
//
// 	GET	/todos		Retrieves all the todos.
// 	POST	/todos		Creates a new todo given a title.
//	DELETE	/todos		Deletes all completed todos.
// 	GET	/todos/{todoKey}	Retrieves the todo with the given key.
// 	PUT	/todos/{todoKey}	Updates the todo with the given key.
//	DELETE	/todos/{todoKey}	Deletes the todo with the given key.
//
// Every method below gives more information about every API call, its
// parameters, and its results.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/GoogleCloudPlatform/gcloud-golang-todos/todo"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const PathPrefix = "/api/todos"

// learnJSON contains marhsalled JSON served in place of TodoMVC's standard
// learn.json and represents a merged JSON object containing both learn.json
// and backend_learn.json.
// TodoMVC's client-side contains some logic which renders it prettily in
// the sidebar.
var learnJSON = mergeJSONOrDie("static/todomvc/learn.json", "static/backend_learn.json")

func init() {
	r := mux.NewRouter()
	r.HandleFunc(PathPrefix,
		errorHandler(ListTodos)).Methods("GET")
	r.HandleFunc(PathPrefix,
		errorHandler(NewTodo)).Methods("POST")
	r.HandleFunc(PathPrefix,
		errorHandler(DeleteCompletedTodos)).Methods("DELETE")
	r.HandleFunc(PathPrefix+"/{key}",
		errorHandler(GetTodo)).Methods("GET")
	r.HandleFunc(PathPrefix+"/{key}",
		errorHandler(UpdateTodo)).Methods("PUT")
	r.HandleFunc(PathPrefix+"/{key}",
		errorHandler(DeleteTodo)).Methods("DELETE")
	http.Handle("/", r)
	http.HandleFunc("/learn.json", WriteLearnJSON)
	http.HandleFunc("/api", IsApiEnabled)
}

// ListTodos handles GET requests on /todos.
// It requires no parameters and returns a list of todos.
//
// Example:
//
//   req: GET /todos/
//   res: 200 [
//          {"id": 1, "title": "Learn Go", "completed": false},
//          {"id": 2, "title": "Buy bread", "completed": true}
//        ]
func ListTodos(w http.ResponseWriter, r *http.Request, c context.Context) error {
	res, err := todo.All(c)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, res)
}

// NewTodo handles POST requests on /todos.
// The request body must contain a JSON object with a Title field.
// The status code of the response is used to indicate any error.
//
// Examples:
//
//   req: POST /todos/ {"title": ""}
//   res: 400 empty title
//
//   req: POST /todos/ {"title": "Buy bread"}
//   res: 201
func NewTodo(w http.ResponseWriter, r *http.Request, c context.Context) error {
	var req struct{ Title string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError{err}
	}
	t, err := todo.NewTodo(c, req.Title)
	if err != nil {
		return badRequestError{err}
	}
	if err = t.Save(c); err != nil {
		return err
	}
	return writeJSON(w, http.StatusCreated, t)
}

// parseKey obtains the key variable from the given request url,
// parses the obtained text, and returns a datastore key, along
// with any processing error.
func parseKey(r *http.Request) (*datastore.Key, error) {
	k, ok := mux.Vars(r)["key"]
	if !ok {
		return nil, fmt.Errorf("todo key not found")
	}
	return datastore.DecodeKey(k)
}

// GetTodo handles GET requsts to /todos/{todoKey}.
// There's no parameters and it returns a JSON encoded todo.
//
// Examples:
//
//   req: GET /todos/1
//   res: 200 {"id": 1, "title": "Buy bread", "completed": true}
//
//   req: GET /todos/42
//   res: 404 todo not found
func GetTodo(w http.ResponseWriter, r *http.Request, c context.Context) error {
	k, err := parseKey(r)
	if err != nil {
		return badRequestError{err}
	}
	t, err := todo.Get(c, k)
	if err != nil {
		return errTodoNotFound
	}
	return writeJSON(w, http.StatusOK, t)
}

// UpdateTodo handles PUT requests to /todos/{todoKey}.
// The request body must contain a JSON encoded todo.
//
// Example:
//
//   req: PUT /todos/1 {"title": "Learn Go", "completed": true}
//   res: 200
func UpdateTodo(w http.ResponseWriter, r *http.Request, c context.Context) error {
	k, err := parseKey(r)
	if err != nil {
		return badRequestError{err}
	}
	var t todo.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return badRequestError{err}
	}
	t.Key = k
	if _, err := todo.Get(c, k); err != nil {
		return errTodoNotFound
	}
	if err = t.Save(c); err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, t)
}

// DeleteTodo handles DELETE requests to /todos/{todoKey}.
// Returns a badRequestError error if the key cannot be parsed, and
// errTodoNotFound if no corresponding todo can be found.
//
// Example:
//
//   req: DELETE /todos/1
//   res: 204
//
//   req: DELETE /todos/asdf
//   res: 400
//
//   req: DELETE /todos/42
//   res: 404 todo not found
func DeleteTodo(w http.ResponseWriter, r *http.Request, c context.Context) error {
	k, err := parseKey(r)
	if err != nil {
		return badRequestError{err}
	}
	if err := todo.Delete(c, k); err != nil {
		return errTodoNotFound
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// DeleteCompletedTodos handles DELETE requests to /todos.
// It attempts to delete all todos which have been marked completed and returns
// an error if one occurred.
//
// Example:
//
//   req: DELETE /todos
//   res: 204
//
//   req: DELETE /todos (Some error happens)
//   res: 500
func DeleteCompletedTodos(w http.ResponseWriter, r *http.Request, c context.Context) error {
	if err := todo.DeleteCompleted(c); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// mergeJSONOrDie merges the k/v pairs in all indicated files with the
// base map derived from the first file. It returns a slice of bytes
// representing the final marshaled JSON, and panics if any file cannot be
// read, have its data unmarshaled, or if the resulting map cannot be
// marshaled.
// Importantly, this function only works with files corresponding to JSON
// objects! Files that correspond to JSON arrays will fail!
// Also, if more than one file shares the same key, only the value from
// the final file will be present in the output.
func mergeJSONOrDie(files ...string) []byte {
	if len(files) == 0 {
		panic(fmt.Errorf("Called with no files!"))
	}
	b, err := ioutil.ReadFile(files[0])
	if err != nil {
		panic(err)
	}
	var base map[string]interface{}
	if err = json.Unmarshal(b, &base); err != nil {
		panic(err)
	}
	for _, file := range files[1:] {
		b, err = ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		var next map[string]interface{}
		if err = json.Unmarshal(b, &next); err != nil {
			panic(err)
		}
		for k, v := range next {
			base[k] = v
		}
	}
	b, err = json.Marshal(base)
	if err != nil {
		panic(err)
	}
	return b
}

// errorHandler wraps a function returning an error by handling the error and returning a http.Handler.
// If the error is of the one of the types defined above, it is handled as described for every type.
// If the error is of another type, it is considered as an internal error and its message is logged.
func errorHandler(f func(w http.ResponseWriter, r *http.Request, c context.Context) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		err := f(w, r, c)
		if err == nil {
			return
		}
		if err == errTodoNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		switch err.(type) {
		case badRequestError:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Errorf(c, "internal error %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// writeJSON marshals v into its JSON representation and writes that to w
// along with the given status code. It also takes care to set the HTTP
// Content-Type header. Any error encountered during marshalling is returned.
func writeJSON(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

// WriteLearnJSON writes a sidebar JSON file consumable by the TodoMVC
// JS which is specific to our Go backend implementation.
func WriteLearnJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(learnJSON)
}

// IsApiEnabled writes an HTTP 200 indicating that the TodoMvc API is enabled for this app.
func IsApiEnabled(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// badRequestError is a type of error returned by a handler function to indicate a malformed request.
type badRequestError struct{ error }

// errTodoNotFound should be returned by a handler function to indicate that the requested data could not be found.
var errTodoNotFound = fmt.Errorf("todo not found")
