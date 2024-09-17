package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qluan1/go-todo-micro/internal/todos"
)

type Todos struct {
	logger *log.Logger
}

func NewTodos(logger *log.Logger) *Todos {
	return &Todos{logger}
}

func (t *Todos) GetTodos(rw http.ResponseWriter, _ *http.Request) {
	t.logger.Println("GET /todos")

	tds := todos.GetTodos()

	err := tds.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Error reading todos", http.StatusBadRequest)
		t.logger.Println("Error reading todos", err)
	}
}

func (t *Todos) PostTodos(rw http.ResponseWriter, r *http.Request) {
	t.logger.Println("POST /todos")

	todo := r.Context().Value(TodoKey{}).(*todos.Todo)
	
	todos.AddTodo(todo)
	err := todo.ToJSON(rw) // return the todo with the ID
	if err != nil {
		http.Error(rw, "Error reading todo", http.StatusBadRequest)
		t.logger.Println("Error reading todo", err)
	}
	t.logger.Println("Todo added")
}

func (t *Todos) PutTodo(rw http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	t.logger.Printf("PUT /todos/%s\n", s)

	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}

	todo := r.Context().Value(TodoKey{}).(*todos.Todo)

	// update the todo with id
	todo, err = todos.UpdateTodoByID(id, todo)
	if err == todos.ErrTodoNotFound {
		http.Error(rw, "Todo not found", http.StatusNotFound)
		return
	}
	
	if err != nil {
		http.Error(rw, "Error updating todo", http.StatusInternalServerError)
		return
	}

	err = todo.ToJSON(rw) // return the todo with the ID
	if err != nil {
		http.Error(rw, "Error reading todo", http.StatusBadRequest)
		t.logger.Println("Error reading todo", err)
	}
	t.logger.Printf("Todo %d updated", id)
}

type TodoKey struct{}

func (t *Todos) MiddlewareValidateTodo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// validate the todo
		todo := &todos.Todo{}

		err := todo.FromJSON(r.Body)
		if err != nil {
			t.logger.Println("[Error] deserializing todo", err)
			http.Error(rw, "Fail to unmarshal todo", http.StatusBadRequest)
			return
		}

		err = todo.Validate()
		if err != nil {
			t.logger.Println("[Error] validating todo", err)
			http.Error(rw, "Invalid todo", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), TodoKey{}, todo)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}