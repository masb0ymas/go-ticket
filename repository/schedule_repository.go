package repository

import (
	"go-ticket/models"

	"github.com/jmoiron/sqlx"
)

type ScheduleRepository struct {
	*Repository[models.Schedule]
}

func NewScheduleRepository(db *sqlx.DB) *ScheduleRepository {
	return &ScheduleRepository{
		Repository: NewRepository[models.Schedule](db, "schedules"),
	}
}

// Custom methods for ScheduleRepository
func (r *ScheduleRepository) FindByDateRange(startDate, endDate string) ([]models.Schedule, error) {
	query := `
		SELECT * FROM schedules 
		WHERE start_date >= $1 
		AND end_date <= $2 
		AND deleted_at IS NULL
	`

	var schedules []models.Schedule
	err := r.db.Select(&schedules, query, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *ScheduleRepository) Create(schedule *models.Schedule) error {
	query := `
		INSERT INTO schedules (
			id, title, description, start_date, end_date,
			created_at, updated_at
		) VALUES (
			:id, :title, :description, :start_date, :end_date,
			:created_at, :updated_at
		)
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":          schedule.ID,
		"title":       schedule.Title,
		"description": schedule.Description,
		"start_date":  schedule.StartDate,
		"end_date":    schedule.EndDate,
		"created_at":  schedule.CreatedAt,
		"updated_at":  schedule.UpdatedAt,
	})
	return err
}

func (r *ScheduleRepository) Update(schedule *models.Schedule) error {
	query := `
		UPDATE schedules SET
			title = :title,
			description = :description,
			start_date = :start_date,
			end_date = :end_date,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":          schedule.ID,
		"title":       schedule.Title,
		"description": schedule.Description,
		"start_date":  schedule.StartDate,
		"end_date":    schedule.EndDate,
		"updated_at":  schedule.UpdatedAt,
	})
	return err
}
