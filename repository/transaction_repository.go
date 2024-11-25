package repository

import (
	"go-ticket/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	*Repository[models.Transaction]
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{
		Repository: NewRepository[models.Transaction](db, "transactions"),
		db:         db,
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
			&transaction.ID, &transaction.UserID, &transaction.EventID, &transaction.TotalAmount, &transaction.Status,
			&transaction.PaymentMethod, &transaction.PaymentStatus, &transaction.PaymentUrl, &transaction.PaymentCallback,
			&transaction.CreatedAt, &transaction.UpdatedAt, &transaction.DeletedAt,
			&user.ID, &user.Fullname, &user.Email, &user.Phone, &user.Password,
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
				&transaction.ID, &transaction.UserID, &transaction.EventID, &transaction.TotalAmount, &transaction.Status,
				&transaction.PaymentMethod, &transaction.PaymentStatus, &transaction.PaymentUrl, &transaction.PaymentCallback,
				&transaction.CreatedAt, &transaction.UpdatedAt, &transaction.DeletedAt,
				&user.ID, &user.Fullname, &user.Email, &user.Phone, &user.Password,
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
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *TransactionRepository) UpdatePaymentStatus(id uuid.UUID, status string) error {
	query := `
		UPDATE transactions 
		SET payment_status = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, user_id, event_id, total_amount,
			status, payment_method, payment_status,
			payment_url, payment_callback,
			created_at, updated_at
		) VALUES (
			:id, :user_id, :event_id, :total_amount,
			:status, :payment_method, :payment_status,
			:payment_url, :payment_callback,
			:created_at, :updated_at
		)
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":               transaction.ID,
		"user_id":          transaction.UserID,
		"event_id":         transaction.EventID,
		"total_amount":     transaction.TotalAmount,
		"status":           transaction.Status,
		"payment_method":   transaction.PaymentMethod,
		"payment_status":   transaction.PaymentStatus,
		"payment_url":      transaction.PaymentUrl,
		"payment_callback": transaction.PaymentCallback,
		"created_at":       transaction.CreatedAt,
		"updated_at":       transaction.UpdatedAt,
	})
	return err
}

func (r *TransactionRepository) Update(transaction *models.Transaction) error {
	query := `
		UPDATE transactions SET
			total_amount = :total_amount,
			status = :status,
			payment_method = :payment_method,
			payment_status = :payment_status,
			payment_url = :payment_url,
			payment_callback = :payment_callback,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":               transaction.ID,
		"total_amount":     transaction.TotalAmount,
		"status":           transaction.Status,
		"payment_method":   transaction.PaymentMethod,
		"payment_status":   transaction.PaymentStatus,
		"payment_url":      transaction.PaymentUrl,
		"payment_callback": transaction.PaymentCallback,
		"updated_at":       transaction.UpdatedAt,
	})
	return err
}
