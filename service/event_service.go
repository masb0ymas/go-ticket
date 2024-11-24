package service

import (
	"errors"
	"go-ticket/models"
	"go-ticket/repository"

	"github.com/google/uuid"
)

type EventService struct {
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{
		repo: repo,
	}
}

type CreateEventRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	LocationID  uuid.UUID `json:"location_id" validate:"required"`
	ScheduleID  uuid.UUID `json:"schedule_id" validate:"required"`
}

type UpdateEventRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	LocationID  uuid.UUID `json:"location_id" validate:"required"`
	ScheduleID  uuid.UUID `json:"schedule_id" validate:"required"`
}

func (s *EventService) GetAllEvents() ([]models.Event, error) {
	return s.repo.FindAllWithRelations()
}

func (s *EventService) GetEventById(id uuid.UUID) (*models.Event, error) {
	event, err := s.repo.FindWithRelations(id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, errors.New("event not found")
	}
	return event, nil
}

func (s *EventService) CreateEvent(req *CreateEventRequest) (*models.Event, error) {
	event := &models.Event{
		Name:        req.Name,
		Description: req.Description,
		LocationID:  req.LocationID,
		ScheduleID:  req.ScheduleID,
	}

	err := s.repo.Create(event)
	if err != nil {
		return nil, err
	}

	return s.GetEventById(event.ID)
}

func (s *EventService) UpdateEvent(id uuid.UUID, req *UpdateEventRequest) (*models.Event, error) {
	event, err := s.GetEventById(id)
	if err != nil {
		return nil, err
	}

	event.Name = req.Name
	event.Description = req.Description
	event.LocationID = req.LocationID
	event.ScheduleID = req.ScheduleID

	err = s.repo.Update(id, event)
	if err != nil {
		return nil, err
	}

	return s.GetEventById(id)
}

func (s *EventService) DeleteEvent(id uuid.UUID) error {
	event, err := s.GetEventById(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(event.ID)
}
