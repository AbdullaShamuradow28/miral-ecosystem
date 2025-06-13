package models

import "time"

type Project struct {
	gorm.Model
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex"`
	OwnerID   string    `json:"owner_id"`
	Owner     MCUser    `json:"owner" gorm:"foreignKey:OwnerID"`
	StorageUsed int64     `json:"storage_used" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Collections []Collection `json:"collections" gorm:"foreignKey:ProjectID"`
}