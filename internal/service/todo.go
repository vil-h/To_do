package service

import (
	"awesomeProject/internal/models"
	"time"

	"gorm.io/gorm"
)

type ToDoService struct {
	DB *gorm.DB
}

func (t *ToDoService) Add(name string, time *time.Time) (int, error) {
	task := &models.To_do{
		Name:     name,
		IsActive: false,
		Deadline: time,
	}

	err := t.DB.Create(task).Error
	if err != nil {
		return 0, err
	}
	return task.ID, nil
}

func (t *ToDoService) GetTask(filter string) ([]models.To_do, error) {
	var tobos []models.To_do

	query := t.DB.Model(&models.To_do{})

	if filter == "archived" {
		query = query.Where("is_archived = ?", true)
	} else {
		query = query.Where("is_archived = ?", false)
	}

	switch filter {
	case "active":
		query = query.Where("is_active = ?", false)
	case "completed":
		query = query.Where("is_active = ?", true)
	case "archived":
		query = query.Where("is_archived = ?", true)
	case "all":
	default:
	}
	err := query.Find(&tobos).Error
	if err != nil {
		return nil, err
	}
	return tobos, nil
}

func (t *ToDoService) Change(id int) (bool, error) {
	var todo models.To_do
	err := t.DB.First(&todo, id).Error
	if err != nil {
		return false, err
	}
	newStatus := !todo.IsActive

	err = t.DB.Model(&todo).Update("is_active", newStatus).Error
	if err != nil {
		return false, err
	}
	return newStatus, nil
}

func (t *ToDoService) Archiv(id int) error {
	err := t.DB.Model(models.To_do{}).Where("id == ?", id).Update("is_archived", true).Error
	return err
}
