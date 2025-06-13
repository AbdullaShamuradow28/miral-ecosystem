package models

import "time"

type Group struct {
	gorm.Model
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex"`
	CreatorID string    `json:"creator_id"` // ID пользователя, создавшего группу
	Creator   MCUser    `json:"creator" gorm:"foreignKey:CreatorID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Files     []File    `json:"files" gorm:"foreignKey:GroupID"` // Связь с файлами
}