package repository

import (
	"go-ticket/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	*Repository[models.Transaction]
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{
		Repository: NewRepository[models.Transaction](db, "transactions"),
	}
}

// Custom methods for TransactionRepository
func (r *TransactionRepository) FindByUserId(userId uuid.UUID) ([]models.Transaction, error) {
	query := `
		SELECT t.*, u.* FROM transactions t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.user_id = $1 
		AND t.deleted_at IS NULL
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.Queryx(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		var user models.User
		err := rows.Scan(
			&transaction.ID, &transaction.UserID, &transaction.Status,
			&transaction.TotalAmount, &transaction.PaymentMethod,
			&transaction.PaymentStatus, &transaction.CreatedAt,
			&transaction.UpdatedAt, &transaction.DeletedAt,
			&user.ID, &user.Fullname, &user.Email, &user.Phone,
			&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		transaction.User = &user
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) FindWithDetails(id uuid.UUID) (*models.Transaction, error) {
	query := `
		SELECT t.*, u.*, td.*, tt.* FROM transactions t
		LEFT JOIN users u ON t.user_id = u.id
		LEFT JOIN transaction_details td ON td.transaction_id = t.id
		LEFT JOIN ticket_types tt ON td.ticket_type_id = tt.id
		WHERE t.id = $1 
		AND t.deleted_at IS NULL
	`

	rows, err := r.db.Queryx(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transaction *models.Transaction
	var details []models.TransactionDetail

	for rows.Next() {
		var user models.User
		var detail models.TransactionDetail
		var ticketType models.TicketType

		if transaction == nil {
			transaction = &models.Transaction{}
			err := rows.Scan(
				&transaction.ID, &transaction.UserID, &transaction.Status,
				&transaction.TotalAmount, &transaction.PaymentMethod,
				&transaction.PaymentStatus, &transaction.CreatedAt,
				&transaction.UpdatedAt, &transaction.DeletedAt,
				&user.ID, &user.Fullname, &user.Email, &user.Phone,
				&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
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
			transaction.User = &user
		}

		detail.TicketType = &ticketType
		details = append(details, detail)
	}

	if transaction != nil {
		transaction.Details = details
	}

	return transaction, nil
}

func (r *TransactionRepository) UpdateStatus(id uuid.UUID, status string) error {
	query := `
		UPDATE transactions 
		SET status = $1, 
			updated_at = NOW()
		WHERE id = $2 
		AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *TransactionRepository) UpdatePaymentStatus(id uuid.UUID, status string) error {
	query := `
		UPDATE transactions 
		SET payment_status = $1, 
			updated_at = NOW()
		WHERE id = $2 
		AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, status, id)
	return err
}
