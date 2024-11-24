package repository

import (
	"go-ticket/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	*Repository[models.Event]
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{
		Repository: NewRepository[models.Event](db, "events"),
	}
}

// Custom methods for EventRepository
func (r *EventRepository) FindWithRelations(id uuid.UUID) (*models.Event, error) {
	query := `
		SELECT e.*, l.*, s.*
		FROM events e
		LEFT JOIN locations l ON e.location_id = l.id
		LEFT JOIN schedules s ON e.schedule_id = s.id
		WHERE e.id = $1 AND e.deleted_at IS NULL
	`

	rows, err := r.db.Queryx(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var event models.Event
	var location models.Location
	var schedule models.Schedule

	if rows.Next() {
		err = rows.Scan(
			&event.ID, &event.Name, &event.Description, &event.LocationID, &event.ScheduleID,
			&event.CreatedAt, &event.UpdatedAt, &event.DeletedAt,
			&location.ID, &location.Name, &location.Address, &location.City,
			&location.State, &location.Country, &location.PostalCode,
			&location.CreatedAt, &location.UpdatedAt, &location.DeletedAt,
			&schedule.ID, &schedule.Title, &schedule.Description, &schedule.StartDate, &schedule.EndDate,
			&schedule.CreatedAt, &schedule.UpdatedAt, &schedule.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		event.Location = &location
		event.Schedule = &schedule
	}

	return &event, nil
}

func (r *EventRepository) FindAllWithRelations() ([]models.Event, error) {
	query := `
		SELECT e.*, l.*, s.*
		FROM events e
		LEFT JOIN locations l ON e.location_id = l.id
		LEFT JOIN schedules s ON e.schedule_id = s.id
		WHERE e.deleted_at IS NULL
	`

	rows, err := r.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var event models.Event
		var location models.Location
		var schedule models.Schedule

		err = rows.Scan(
			&event.ID, &event.Name, &event.Description, &event.LocationID, &event.ScheduleID,
			&event.CreatedAt, &event.UpdatedAt, &event.DeletedAt,
			&location.ID, &location.Name, &location.Address, &location.City,
			&location.State, &location.Country, &location.PostalCode,
			&location.CreatedAt, &location.UpdatedAt, &location.DeletedAt,
			&schedule.ID, &schedule.Title, &schedule.Description, &schedule.StartDate, &schedule.EndDate,
			&schedule.CreatedAt, &schedule.UpdatedAt, &schedule.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		event.Location = &location
		event.Schedule = &schedule
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) Create(event *models.Event) error {
	query := `
		INSERT INTO events (
			id, name, description, location_id, schedule_id,
			created_at, updated_at
		) VALUES (
			:id, :name, :description, :location_id, :schedule_id,
			:created_at, :updated_at
		)
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":          event.ID,
		"name":        event.Name,
		"description": event.Description,
		"location_id": event.LocationID,
		"schedule_id": event.ScheduleID,
		"created_at":  event.CreatedAt,
		"updated_at":  event.UpdatedAt,
	})
	return err
}

func (r *EventRepository) Update(event *models.Event) error {
	query := `
		UPDATE events SET
			name = :name,
			description = :description,
			location_id = :location_id,
			schedule_id = :schedule_id,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":          event.ID,
		"name":        event.Name,
		"description": event.Description,
		"location_id": event.LocationID,
		"schedule_id": event.ScheduleID,
		"updated_at":  event.UpdatedAt,
	})
	return err
}
