package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    Name      string         `gorm:"size:100;not null" json:"name"`
    Username string           `gorm:"size:100;not null" json:"username"`
    Email     string         `gorm:"size:100;unique;not null" json:"email"`
    Age       int            `json:"age"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    Password string `gorm:"size:255;not null" json:"-"`
}

// Hook: Tự động hash password trước khi tạo user mới
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.Password != "" {
        return u.hashPassword()
    }
    return nil
}

// Hook: Tự động hash lại nếu password bị thay đổi
func (u *User) BeforeSave(tx *gorm.DB) error {
    // Kiểm tra xem field Password có trong danh sách các field bị thay đổi không
    if tx.Statement.Changed("Password") {
        if u.Password != "" {
            return u.hashPassword()
        }
    }
    return nil
}

// Hàm helper để hash password
func (u *User) hashPassword() error {
    hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashed)
    return nil
}

// Hàm để kiểm tra password khi login (dùng ở usecase/service)
func (u *User) ComparePassword(password string) error {
    return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}