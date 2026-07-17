package models

import "time"

type To_do struct {
	ID          int        `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name"`
	IsActive    bool       `json:"is_active"`
	CreatedTime time.Time  `gorm:"autoCreateTime" json:"created_time"`
	Deadline    *time.Time `json:"deadline_time"`
	IsArchived  bool       `gorm:"default:false" json:"is_archived"`
}
