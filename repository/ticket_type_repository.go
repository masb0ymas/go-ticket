package repository

import (
	"errors"
	"go-ticket/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TicketTypeRepository struct {
	*Repository[models.TicketType]
}

func NewTicketTypeRepository(db *sqlx.DB) *TicketTypeRepository {
	return &TicketTypeRepository{
		Repository: NewRepository[models.TicketType](db, "ticket_types"),
	}
}

// Custom methods for TicketTypeRepository
func (r *TicketTypeRepository) FindByEventId(eventId uuid.UUID) ([]models.TicketType, error) {
	query := `
		SELECT t.*, e.* FROM ticket_types t
		LEFT JOIN events e ON t.event_id = e.id
		WHERE t.event_id = $1 
		AND t.deleted_at IS NULL
	`

	rows, err := r.db.Queryx(query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ticketTypes []models.TicketType
	for rows.Next() {
		var ticketType models.TicketType
		var event models.Event
		err := rows.Scan(
			&ticketType.ID, &ticketType.EventID, &ticketType.Name,
			&ticketType.Description, &ticketType.Price, &ticketType.Quota,
			&ticketType.RemainingQuota, &ticketType.CreatedAt,
			&ticketType.UpdatedAt, &ticketType.DeletedAt,
			&event.ID, &event.Name, &event.Description,
			&event.LocationID, &event.ScheduleID,
			&event.CreatedAt, &event.UpdatedAt, &event.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		ticketType.Event = &event
		ticketTypes = append(ticketTypes, ticketType)
	}

	return ticketTypes, nil
}

func (r *TicketTypeRepository) FindAvailable(eventId uuid.UUID) ([]models.TicketType, error) {
	query := `
		SELECT t.*, e.* FROM ticket_types t
		LEFT JOIN events e ON t.event_id = e.id
		WHERE t.event_id = $1 
		AND t.remaining_quota > 0
		AND t.deleted_at IS NULL
	`

	rows, err := r.db.Queryx(query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ticketTypes []models.TicketType
	for rows.Next() {
		var ticketType models.TicketType
		var event models.Event
		err := rows.Scan(
			&ticketType.ID, &ticketType.EventID, &ticketType.Name,
			&ticketType.Description, &ticketType.Price, &ticketType.Quota,
			&ticketType.RemainingQuota, &ticketType.CreatedAt,
			&ticketType.UpdatedAt, &ticketType.DeletedAt,
			&event.ID, &event.Name, &event.Description,
			&event.LocationID, &event.ScheduleID,
			&event.CreatedAt, &event.UpdatedAt, &event.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		ticketType.Event = &event
		ticketTypes = append(ticketTypes, ticketType)
	}

	return ticketTypes, nil
}

func (r *TicketTypeRepository) UpdateQuota(id uuid.UUID, quantity int) error {
	query := `
		UPDATE ticket_types 
		SET remaining_quota = remaining_quota - $1
		WHERE id = $2 
		AND deleted_at IS NULL
		AND remaining_quota >= $1
	`

	result, err := r.db.Exec(query, quantity, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("insufficient ticket quota")
	}

	return nil
}
