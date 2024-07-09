package handlers

import (
	"log"
	"net/http"

	"github.com/qluan1/go-todo-micro/internal/todos"
)

type Todos struct {
	logger *log.Logger
}

func NewTodos(logger *log.Logger) *Todos {
	return &Todos{logger}
}

func (t *Todos) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t.GetTodos(rw, r)
		return
	}
	http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
}

func (t *Todos) GetTodos(rw http.ResponseWriter, r *http.Request) {
	tds := todos.GetTodos()
	err := tds.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Error reading todos", http.StatusBadRequest)
		t.logger.Println("Error reading todos", err)
	}
	t.logger.Println("GET /todos")
}