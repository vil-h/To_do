package handlers

import (
	"awesomeProject/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var TodoService *service.ToDoService

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get", http.StatusMethodNotAllowed)
		return
	}

	filter := r.URL.Query().Get("filter")
	fmt.Printf("=== GetHandler: filter = '%s' ===\n", filter)
	if filter == "" {
		filter = "all"
	}
	if filter != "all" && filter != "active" && filter != "completed" && filter != "archived" {
		http.Error(w, "Ошибка валидации", http.StatusBadRequest)
		return
	}

	zadachi, err := TodoService.GetTask(filter)
	if err != nil {
		http.Error(w, "Проеб в БД:"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(&zadachi)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST", http.StatusMethodNotAllowed)
		return
	}
	type newCreateTask struct {
		Name     string `json:"name"`
		Deadline string `json:"deadline_time"`
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

	var deadline *time.Time
	if req.Deadline != "" {
		s, err := time.Parse("2006-01-02 15:04:05", req.Deadline)
		if err != nil {
			http.Error(w, "Неверный формат даты. Ожидается: 2006-01-02 15:04:05", http.StatusBadRequest)
			return
		}
		deadline = &s
	}
	NewId, err := TodoService.Add(req.Name, deadline)
	if err != nil {
		http.Error(w, "Проеб в БД блять:"+err.Error(), http.StatusMethodNotAllowed)
		return
	}
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
	newStatus, err := TodoService.Change(newId.Id)
	if err != nil {
		http.Error(w, "Не удалось изменить статус задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}
	res := otvetStatus{Status: newStatus}
	json.NewEncoder(w).Encode(res)

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE", http.StatusMethodNotAllowed)
		return
	}
	type deleteRequest struct {
		Id int `json:"id"`
	}
	var req deleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Кривой JSON"+err.Error(), http.StatusNotFound)
		return
	}
	err = TodoService.Archiv(req.Id)
	if err != nil {
		http.Error(w, "Не удалось удалить задачу: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"result": "success"}`))
}
