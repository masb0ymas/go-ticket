package service

import (
	"errors"
	"go-ticket/models"
	"go-ticket/repository"
	"time"

	"github.com/google/uuid"
)

type TicketTypeService struct {
	repo *repository.TicketTypeRepository
}

func NewTicketTypeService(repo *repository.TicketTypeRepository) *TicketTypeService {
	return &TicketTypeService{
		repo: repo,
	}
}

type CreateTicketTypeRequest struct {
	EventID     uuid.UUID `json:"event_id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" validate:"required,min=0"`
	Quota       int       `json:"quota" validate:"required,min=1"`
}

type UpdateTicketTypeRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       *float64 `json:"price" validate:"omitempty,min=0"`
	Quota       *int     `json:"quota" validate:"omitempty,min=1"`
}

func (s *TicketTypeService) GetAllTicketTypes() ([]models.TicketType, error) {
	return s.repo.FindAll()
}

func (s *TicketTypeService) GetTicketTypeById(id uuid.UUID) (*models.TicketType, error) {
	return s.repo.FindById(id)
}

func (s *TicketTypeService) GetTicketTypesByEventId(eventId uuid.UUID) ([]models.TicketType, error) {
	return s.repo.FindByEventId(eventId)
}

func (s *TicketTypeService) GetAvailableTicketTypes(eventId uuid.UUID) ([]models.TicketType, error) {
	return s.repo.FindAvailable(eventId)
}

func (s *TicketTypeService) CreateTicketType(req *CreateTicketTypeRequest) (*models.TicketType, error) {
	ticketType := &models.TicketType{
		BaseModel: models.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		EventID:        req.EventID,
		Name:           req.Name,
		Description:    req.Description,
		Price:          req.Price,
		Quota:          req.Quota,
		RemainingQuota: req.Quota,
	}

	err := s.repo.Create(ticketType)
	if err != nil {
		return nil, err
	}

	return ticketType, nil
}

func (s *TicketTypeService) UpdateTicketType(id uuid.UUID, req *UpdateTicketTypeRequest) (*models.TicketType, error) {
	ticketType, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		ticketType.Name = req.Name
	}

	if req.Description != "" {
		ticketType.Description = req.Description
	}

	if req.Price != nil {
		ticketType.Price = *req.Price
	}

	if req.Quota != nil {
		if *req.Quota < ticketType.Quota-ticketType.RemainingQuota {
			return nil, errors.New("new quota cannot be less than sold tickets")
		}
		quotaDiff := *req.Quota - ticketType.Quota
		ticketType.Quota = *req.Quota
		ticketType.RemainingQuota += quotaDiff
	}

	ticketType.UpdatedAt = time.Now()

	err = s.repo.Update(ticketType)
	if err != nil {
		return nil, err
	}

	return ticketType, nil
}

func (s *TicketTypeService) DeleteTicketType(id uuid.UUID) error {
	ticketType, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	if ticketType.Quota != ticketType.RemainingQuota {
		return errors.New("cannot delete ticket type with sold tickets")
	}

	ticketType.DeletedAt = &time.Time{}
	return s.repo.Update(ticketType)
}

func (s *TicketTypeService) UpdateQuota(id uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	return s.repo.UpdateQuota(id, quantity)
}
