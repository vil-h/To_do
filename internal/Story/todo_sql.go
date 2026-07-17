package Story

import (
	"awesomeProject/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.To_do{})
	if err != nil {
		return nil, err
	}
	return db, nil

}
