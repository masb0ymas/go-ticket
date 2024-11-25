package repository

import (
	"go-ticket/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionDetailRepository struct {
	*Repository[models.TransactionDetail]
}

func NewTransactionDetailRepository(db *sqlx.DB) *TransactionDetailRepository {
	return &TransactionDetailRepository{
		Repository: NewRepository[models.TransactionDetail](db, "transaction_details"),
	}
}

// Custom methods for TransactionDetailRepository
func (r *TransactionDetailRepository) FindByTransactionId(transactionId uuid.UUID) ([]models.TransactionDetail, error) {
	query := `
		SELECT td.*, tt.* FROM transaction_details td
		LEFT JOIN ticket_types tt ON td.ticket_type_id = tt.id
		WHERE td.transaction_id = $1 
		AND td.deleted_at IS NULL
	`

	rows, err := r.db.Queryx(query, transactionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []models.TransactionDetail
	for rows.Next() {
		var detail models.TransactionDetail
		var ticketType models.TicketType
		err := rows.Scan(
			&detail.ID, &detail.TransactionID, &detail.TicketTypeID,
			&detail.Quantity, &detail.PricePerTicket, &detail.Subtotal,
			&detail.CreatedAt, &detail.UpdatedAt, &detail.DeletedAt,
			&ticketType.ID, &ticketType.EventID, &ticketType.Name,
			&ticketType.Description, &ticketType.Price, &ticketType.Quota,
			&ticketType.RemainingQuota, &ticketType.CreatedAt,
			&ticketType.UpdatedAt, &ticketType.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		detail.TicketType = &ticketType
		details = append(details, detail)
	}

	return details, nil
}

func (r *TransactionDetailRepository) BulkCreate(details []models.TransactionDetail) error {
	query := `
		INSERT INTO transaction_details (
			id, transaction_id, ticket_type_id,
			quantity, price_per_ticket, subtotal,
			created_at, updated_at
		) VALUES (
			:id, :transaction_id, :ticket_type_id,
			:quantity, :price_per_ticket, :subtotal,
			:created_at, :updated_at
		)
	`

	_, err := r.db.NamedExec(query, details)
	return err
}
