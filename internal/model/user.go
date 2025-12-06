package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    Name      string         `gorm:"size:100;not null" json:"name"`
    Email     string         `gorm:"size:100;unique;not null" json:"email"`
    Age       int            `json:"age"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}