package service

import (
	"errors"
	"go-ticket/models"
	"go-ticket/repository"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"email"`
	Phone string `json:"phone"`
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetUserById(id uuid.UUID) (*models.User, error) {
	return s.repo.FindById(id)
}

func (s *UserService) CreateUser(req *CreateUserRequest) (*models.User, error) {
	// Check if email already exists
	existingUser, err := s.repo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Check if phone already exists
	existingUser, err = s.repo.FindByPhone(req.Phone)
	if err == nil && existingUser != nil {
		return nil, errors.New("phone number already registered")
	}

	user := &models.User{
		BaseModel: models.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Fullname: req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, req *UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if req.Email != "" && req.Email != user.Email {
		existingUser, err := s.repo.FindByEmail(req.Email)
		if err == nil && existingUser != nil {
			return nil, errors.New("email already registered")
		}
		user.Email = req.Email
	}

	if req.Phone != "" && req.Phone != user.Phone {
		existingUser, err := s.repo.FindByPhone(req.Phone)
		if err == nil && existingUser != nil {
			return nil, errors.New("phone number already registered")
		}
		user.Phone = req.Phone
	}

	if req.Name != "" {
		user.Fullname = req.Name
	}

	user.UpdatedAt = time.Now()

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	user, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	user.DeletedAt = &time.Time{}
	return s.repo.Update(user)
}
