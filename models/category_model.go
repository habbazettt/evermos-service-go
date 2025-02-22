package models

import "time"

type Category struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	NamaCategory string    `json:"nama_category"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
