package main

import (
	"encoding/json"
	// "fmt"
	// "fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
type todo struct {
	ID   string `json:"ID"`
	Item string `json:"Item"`
}
var todos = []todo{
	{ID: "1", Item: "mobile"},
	{ID: "2", Item: "PC"},
}
func main() {
	// {ID: "3",Item: "Keyboard"}

	chi.RegisterMethod("request")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", start)
	r.Get("/todos", getTodos)
	r.Post("/add", addTodo)
	r.Delete("/delete/{id}", deleteTodo)

	http.ListenAndServe(":3030", r)
}
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	for i, it := range todos {
		if it.ID == id {
			newTodo:=it
			todos = append(todos[0:i], todos[i+1:]...)
			json.NewEncoder(w).Encode(newTodo)
			break
		}
	}

}
func addTodo(w http.ResponseWriter, r *http.Request) {
	var temp todo
	json.NewDecoder(r.Body).Decode(&temp)
	todos = append(todos, temp)
	json.NewEncoder(w).Encode(temp)

}
func getTodos(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(todos)
}
func start(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
