package main

import (
	"awesomeProject/internal/handlers"
	"awesomeProject/internal/service"
	"fmt"
	"net/http"
)

func main() {
	todoService := &service.ToDoService{}
	handlers.TodoService = todoService
	http.HandleFunc("/tasks", handlers.HelloHandler)
	http.HandleFunc("/create", handlers.CreateTaskHandler)
	http.HandleFunc("/tasks/change", handlers.ChangeHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

}
