package main

import (
	"awesomeProject/internal/Story"
	"awesomeProject/internal/handlers"
	"awesomeProject/internal/service"
	"fmt"
	"net/http"
)

func main() {
	db, err := Story.InitDb()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	todoService := &service.ToDoService{
		DB: db,
	}
	handlers.TodoService = todoService
	http.HandleFunc("/tasks", handlers.GetHandler)
	http.HandleFunc("/create", handlers.CreateTaskHandler)
	http.HandleFunc("/tasks/change", handlers.ChangeHandler)
	http.HandleFunc("/tasks/delete", handlers.DeleteTask)

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./web/index.html")
	})

	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

}
