package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/qluan1/go-todo-micro/internal/todos"
)

type Todos struct {
	logger *log.Logger
}

func NewTodos(logger *log.Logger) *Todos {
	return &Todos{logger}
}

func (t *Todos) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// TODO implement url check for /todos/:id and /todos
	if r.Method == http.MethodGet {
		t.get(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		t.post(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		t.put(rw, r)
		return
	}
	http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
}

func (t *Todos) get(rw http.ResponseWriter, _ *http.Request) {
	t.logger.Println("GET /todos")

	tds := todos.GetTodos()

	err := tds.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Error reading todos", http.StatusBadRequest)
		t.logger.Println("Error reading todos", err)
	}
}

func (t *Todos) post(rw http.ResponseWriter, r *http.Request) {
	t.logger.Println("POST /todos")

	todo := &todos.Todo{}

	err := todo.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Fail to unmarshal todo", http.StatusBadRequest)
	}
	
	todos.AddTodo(todo)
	err = todo.ToJSON(rw) // return the todo with the ID
	if err != nil {
		http.Error(rw, "Error reading todo", http.StatusBadRequest)
		t.logger.Println("Error reading todo", err)
	}
}

func (t *Todos) put(rw http.ResponseWriter, r *http.Request) {
	t.logger.Println("PUT /todos/:id")

	// get the id from the URL
	s := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(s)

	if err != nil {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}

	// get the todo from the body
	todo := &todos.Todo{}
	err = todo.FromJSON(r.Body)
	
	if err != nil {
		http.Error(rw, "Fail to unmarshal todo", http.StatusBadRequest)
		return
	}

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