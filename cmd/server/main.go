package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/qluan1/go-todo-micro/internal/handlers"
)

func main() {
	logger := log.New(os.Stdout, "go-todo-micro ", log.LstdFlags)
	
	todoHandler := handlers.NewTodos(logger)

	sm := mux.NewRouter()
	
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/todos", todoHandler.GetTodos)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/todos", todoHandler.PostTodos)
	postRouter.Use(todoHandler.MiddlewareValidateTodo)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/todos/{id:[0-9]+}", todoHandler.PutTodo)
	putRouter.Use(todoHandler.MiddlewareValidateTodo)

	server := &http.Server{
		Addr:    ":8080",
		Handler: sm,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout: 1*time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	sig := <- sigChan
	logger.Println("Received terminate, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.Shutdown(tc)
}
