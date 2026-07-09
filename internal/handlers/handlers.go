package handlers

import (
	"awesomeProject/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

var TodoService *service.ToDoService

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get", http.StatusMethodNotAllowed)
		return
	}

	zadachi := TodoService.GetAll()
	var result string
	for _, task := range zadachi {
		result += fmt.Sprintf("[%d] %s (Выполнено: %t)\n", task.ID, task.Name, task.IsActive)
	}
	w.Write([]byte(result))
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST", http.StatusMethodNotAllowed)
		return
	}
	type newCreateTask struct {
		Name string `json:"name"`
	}
	type otvetAddId struct {
		ID int `json:"id"`
	}
	var req newCreateTask
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Кривой JSON"+err.Error(), http.StatusNotFound)
		return
	}
	NewId := TodoService.Add(req.Name)
	res := otvetAddId{ID: NewId}
	json.NewEncoder(w).Encode(res)

}

func ChangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Only PATCH", http.StatusMethodNotAllowed)
		return
	}
	type IdStatus struct {
		Id int `json:"id"`
	}
	type otvetStatus struct {
		Status bool `json:"Status"`
	}

	var newId IdStatus
	err := json.NewDecoder(r.Body).Decode(&newId)
	if err != nil {
		http.Error(w, "Кривой JSON"+err.Error(), http.StatusNotFound)
		return
	}
	newStatus := TodoService.Change(newId.Id)
	res := otvetStatus{Status: newStatus}
	json.NewEncoder(w).Encode(res)

}
