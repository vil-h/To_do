package service

import "awesomeProject/internal/models"

type ToDoService struct {
	tobos []models.To_do
}

func (t *ToDoService) Add(name string) int {
	NewId := len(t.tobos) + 1
	task := models.To_do{
		ID:       NewId,
		Name:     name,
		IsActive: false,
	}
	t.tobos = append(t.tobos, task)
	return NewId
}

func (t *ToDoService) GetAll() []models.To_do {
	return t.tobos
}

func (t *ToDoService) Change(id int) bool {
	for i := 0; i < len(t.tobos); i++ {
		if t.tobos[i].ID == id {
			t.tobos[i].IsActive = !t.tobos[i].IsActive
			return true
		}
	}
	return false

}
