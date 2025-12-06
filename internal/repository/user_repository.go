package repository

import (
	"github.com/canhviet/go-clean-architecture/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) FindAll() ([]model.User, error) {
    var users []model.User
    return users, r.db.Find(&users).Error
}

func (r *UserRepository) FindByID(id uint) (model.User, error) {
    var user model.User
    return user, r.db.First(&user, id).Error
}

func (r *UserRepository) Create(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *model.User) error {
    return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
    return r.db.Delete(&model.User{}, id).Error
}