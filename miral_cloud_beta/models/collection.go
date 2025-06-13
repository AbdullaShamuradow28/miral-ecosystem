package models

import "time"

// Collection - аналог Collection из твоего прототипа
type Collection struct {
	gorm.Model
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"` // Имя коллекции
	ProjectID uint      `json:"project_id"`
	Project   Project   `json:"project" gorm:"foreignKey:ProjectID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Document struct {
	gorm.Model
	ID           string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` // UUID для документа
	CollectionID uint   `json:"collection_id"`
	Collection   Collection `json:"collection" gorm:"foreignKey:CollectionID"`
	Data         string `json:"data" gorm:"type:text"` // Храним как JSON-строку
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}