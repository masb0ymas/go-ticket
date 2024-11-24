package service

import (
	"errors"
	"go-ticket/models"
	"go-ticket/repository"
	"time"

	"github.com/google/uuid"
)

type LocationService struct {
	repo *repository.LocationRepository
}

func NewLocationService(repo *repository.LocationRepository) *LocationService {
	return &LocationService{
		repo: repo,
	}
}

type CreateLocationRequest struct {
	Name       string `json:"name" validate:"required"`
	Address    string `json:"address" validate:"required"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state"`
	Country    string `json:"country" validate:"required"`
	PostalCode string `json:"postal_code"`
}

type UpdateLocationRequest struct {
	Name       string `json:"name" validate:"required"`
	Address    string `json:"address" validate:"required"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state"`
	Country    string `json:"country" validate:"required"`
	PostalCode string `json:"postal_code"`
}

type SearchLocationRequest struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

func (s *LocationService) GetAllLocations() ([]models.Location, error) {
	return s.repo.FindAll()
}

func (s *LocationService) GetLocationById(id uuid.UUID) (*models.Location, error) {
	location, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, errors.New("location not found")
	}
	return location, nil
}

func (s *LocationService) CreateLocation(req *CreateLocationRequest) (*models.Location, error) {
	location := &models.Location{
		BaseModel: models.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:       req.Name,
		Address:    req.Address,
		City:       req.City,
		State:      req.State,
		Country:    req.Country,
		PostalCode: req.PostalCode,
	}

	err := s.repo.Create(location)
	if err != nil {
		return nil, err
	}

	return s.GetLocationById(location.ID)
}

func (s *LocationService) UpdateLocation(id uuid.UUID, req *UpdateLocationRequest) (*models.Location, error) {
	location, err := s.GetLocationById(id)
	if err != nil {
		return nil, err
	}

	location.Name = req.Name
	location.Address = req.Address
	location.City = req.City
	location.State = req.State
	location.Country = req.Country
	location.PostalCode = req.PostalCode

	err = s.repo.Update(location)
	if err != nil {
		return nil, err
	}

	return s.GetLocationById(id)
}

func (s *LocationService) DeleteLocation(id uuid.UUID) error {
	location, err := s.GetLocationById(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(location.ID)
}

func (s *LocationService) SearchLocations(req *SearchLocationRequest) ([]models.Location, error) {
	if req.City != "" {
		return s.repo.FindByCity(req.City)
	}
	if req.Country != "" {
		return s.repo.FindByCountry(req.Country)
	}
	return s.GetAllLocations()
}
