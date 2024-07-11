package todos

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Todo struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type Todos []*Todo

func (t *Todos) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

func (t *Todo) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

func (t *Todo) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}

func getNextID() int {
	lastTodo := SampleTodos[len(SampleTodos)-1]
	return lastTodo.ID + 1
}

func AddTodo(t *Todo) {
	t.ID = getNextID()
	t.CreatedAt = time.Now().UTC().String()
	SampleTodos = append(SampleTodos, t)
}

var SampleTodos = []*Todo{
	{
		ID: 1,
		Title: "Learn Go",
		Completed: false,
		CreatedAt: time.Now().UTC().String(),
	},
	{
		ID: 2,
		Title: "Learn Microservices",
		Completed: false,
		CreatedAt: time.Now().UTC().String(),
	},
}

func GetTodos() *Todos {
	todos := Todos(SampleTodos)
	return &todos
}

var ErrTodoNotFound = fmt.Errorf("Todo not found")

func GetTodoById(id int) (*Todo, error) {
	for _, t := range SampleTodos {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, ErrTodoNotFound 
}

func UpdateTodoByID(id int, t *Todo) (*Todo, error) {
	todo, err := GetTodoById(id)
	if err != nil {
		return nil, err
	}
	todo.Title = t.Title
	todo.Completed = t.Completed
	todo.UpdatedAt = time.Now().UTC().String()
	return todo, nil
}