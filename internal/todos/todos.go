package todos

import (
	"encoding/json"
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

var SampleTodos = &Todos{
	{
		ID: 1,
		Title: "Learn Go",
		Completed: false,
		CreatedAt: time.Now().UTC().String(),
		UpdatedAt: time.Now().UTC().String(),
	},
	{
		ID: 2,
		Title: "Learn Microservices",
		Completed: false,
		CreatedAt: time.Now().UTC().String(),
		UpdatedAt: time.Now().UTC().String(),
	},
}

func GetTodos() *Todos {
	return SampleTodos
}