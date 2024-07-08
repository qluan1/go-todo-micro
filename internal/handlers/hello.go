package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	logger *log.Logger
}

func NewHello(logger *log.Logger) *Hello {
	return &Hello{logger}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.logger.Println("Hello, World!")

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Error reading body", http.StatusBadRequest)
	}

	fmt.Fprintf(rw, "Hello, %s", d)
}