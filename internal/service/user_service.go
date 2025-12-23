package service

import (
	"time"

	"github.com/canhviet/go-clean-architecture/internal/dto"
	"github.com/canhviet/go-clean-architecture/internal/model"
	"github.com/canhviet/go-clean-architecture/internal/repository"
)

type UserService struct {
    repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetAll() ([]dto.UserResponse, error) {
    users, err := s.repo.FindAll()
    return s.toResponseList(users), err
}

func (s* UserService) GetList(limt int, page int)([]dto.UserResponse, error) {
    users, err := s.repo.GetPagedAndFiltered(limt, page)

    return  s.toResponseList(users), err
}

func (s *UserService) GetByID(id uint) (dto.UserResponse, error) {
    user, err := s.repo.FindByID(id)
    return s.toResponse(user), err
}

func (s *UserService) Create(req dto.CreateUserRequest) (dto.UserResponse, error) {
    user := model.User{
        Name:  req.Name,
        Email: req.Email,
        Age:   req.Age,
    }
    err := s.repo.Create(&user)
    return s.toResponse(user), err
}

func (s *UserService) Update(id uint, req dto.UpdateUserRequest) (dto.UserResponse, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return dto.UserResponse{}, err
    }

    if req.Name != "" {
        user.Name = req.Name
    }
    if req.Email != "" {
        user.Email = req.Email
    }
    if req.Age > 0 {
        user.Age = req.Age
    }

    err = s.repo.Update(&user)
    return s.toResponse(user), err
}

func (s *UserService) Delete(id uint) error {
    return s.repo.Delete(id)
}

// Helper: convert model → dto
func (s *UserService) toResponse(u model.User) dto.UserResponse {
    return dto.UserResponse{
        ID:        u.ID,
        Name:      u.Name,
        Email:     u.Email,
        Age:       u.Age,
        CreatedAt: u.CreatedAt.Format(time.RFC3339),
        UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
    }
}

func (s *UserService) toResponseList(users []model.User) []dto.UserResponse {
    result := make([]dto.UserResponse, len(users))
    for i, u := range users {
        result[i] = s.toResponse(u)
    }
    return result
}