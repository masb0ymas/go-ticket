package service

import (
	"errors"
	"go-ticket/models"
	"go-ticket/repository"
	"time"

	"github.com/google/uuid"
)

type ScheduleService struct {
	repo *repository.ScheduleRepository
}

func NewScheduleService(repo *repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{
		repo: repo,
	}
}

type CreateScheduleRequest struct {
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

type UpdateScheduleRequest struct {
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

type SearchScheduleRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (s *ScheduleService) GetAllSchedules() ([]models.Schedule, error) {
	return s.repo.FindAll()
}

func (s *ScheduleService) GetScheduleById(id uuid.UUID) (*models.Schedule, error) {
	schedule, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, errors.New("schedule not found")
	}
	return schedule, nil
}

func (s *ScheduleService) CreateSchedule(req *CreateScheduleRequest) (*models.Schedule, error) {
	schedule := &models.Schedule{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	err := s.repo.Create(schedule)
	if err != nil {
		return nil, err
	}

	return s.GetScheduleById(schedule.ID)
}

func (s *ScheduleService) UpdateSchedule(id uuid.UUID, req *UpdateScheduleRequest) (*models.Schedule, error) {
	schedule, err := s.GetScheduleById(id)
	if err != nil {
		return nil, err
	}

	schedule.StartDate = req.StartDate
	schedule.EndDate = req.EndDate

	err = s.repo.Update(schedule)
	if err != nil {
		return nil, err
	}

	return s.GetScheduleById(id)
}

func (s *ScheduleService) DeleteSchedule(id uuid.UUID) error {
	schedule, err := s.GetScheduleById(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(schedule.ID)
}

func (s *ScheduleService) SearchSchedules(req *SearchScheduleRequest) ([]models.Schedule, error) {
	if req.StartDate == "" || req.EndDate == "" {
		return s.GetAllSchedules()
	}

	return s.repo.FindByDateRange(req.StartDate, req.EndDate)
}
