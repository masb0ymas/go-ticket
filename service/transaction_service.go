package service

import (
	"errors"
	"go-ticket/models"
	"go-ticket/repository"
	"time"

	"github.com/google/uuid"
)

type TransactionService struct {
	repo           *repository.TransactionRepository
	detailRepo     *repository.TransactionDetailRepository
	ticketTypeRepo *repository.TicketTypeRepository
}

func NewTransactionService(
	repo *repository.TransactionRepository,
	detailRepo *repository.TransactionDetailRepository,
	ticketTypeRepo *repository.TicketTypeRepository,
) *TransactionService {
	return &TransactionService{
		repo:           repo,
		detailRepo:     detailRepo,
		ticketTypeRepo: ticketTypeRepo,
	}
}

type TransactionDetailRequest struct {
	TicketTypeID uuid.UUID `json:"ticket_type_id" validate:"required"`
	Quantity     int       `json:"quantity" validate:"required,min=1"`
}

type CreateTransactionRequest struct {
	UserID        uuid.UUID                  `json:"user_id" validate:"required"`
	EventID       uuid.UUID                  `json:"event_id" validate:"required"`
	PaymentMethod string                     `json:"payment_method" validate:"required"`
	PaymentUrl    string                     `json:"payment_url" validate:"required"`
	Details       []TransactionDetailRequest `json:"details" validate:"required,min=1"`
}

type UpdateTransactionStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

type UpdatePaymentStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

func (s *TransactionService) GetAllTransactions() ([]models.Transaction, error) {
	return s.repo.FindAll()
}

func (s *TransactionService) GetTransactionById(id uuid.UUID) (*models.Transaction, error) {
	return s.repo.FindWithDetails(id)
}

func (s *TransactionService) GetTransactionsByUserId(userId uuid.UUID) ([]models.Transaction, error) {
	return s.repo.FindByUserId(userId)
}

func (s *TransactionService) CreateTransaction(req *CreateTransactionRequest) (*models.Transaction, error) {
	// Calculate total amount and validate ticket availability
	var totalAmount float64
	var details []models.TransactionDetail

	for _, detail := range req.Details {
		ticketType, err := s.ticketTypeRepo.FindById(detail.TicketTypeID)
		if err != nil {
			return nil, err
		}

		if ticketType.RemainingQuota < detail.Quantity {
			return nil, errors.New("insufficient ticket quota")
		}

		subTotal := ticketType.Price * float64(detail.Quantity)
		totalAmount += subTotal

		details = append(details, models.TransactionDetail{
			BaseModel: models.BaseModel{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			TicketTypeID:   detail.TicketTypeID,
			Quantity:       detail.Quantity,
			PricePerTicket: ticketType.Price,
			Subtotal:       subTotal,
		})
	}

	// Create transaction
	transaction := &models.Transaction{
		BaseModel: models.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:        req.UserID,
		EventID:       req.EventID,
		Status:        "pending",
		TotalAmount:   totalAmount,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: "pending",
		PaymentUrl:    req.PaymentUrl,
	}

	err := s.repo.Create(transaction)
	if err != nil {
		return nil, err
	}

	// Create transaction details and update ticket quotas
	for i := range details {
		details[i].TransactionID = transaction.ID
		err = s.ticketTypeRepo.UpdateQuota(details[i].TicketTypeID, details[i].Quantity)
		if err != nil {
			// TODO: Implement rollback
			return nil, err
		}
	}

	err = s.detailRepo.BulkCreate(details)
	if err != nil {
		// TODO: Implement rollback
		return nil, err
	}

	transaction.Details = details
	return transaction, nil
}

func (s *TransactionService) UpdateTransactionStatus(id uuid.UUID, req *UpdateTransactionStatusRequest) error {
	_, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	validStatuses := map[string]bool{
		"pending":   true,
		"confirmed": true,
		"cancelled": true,
		"completed": true,
	}

	if !validStatuses[req.Status] {
		return errors.New("invalid status")
	}

	return s.repo.UpdateStatus(id, req.Status)
}

func (s *TransactionService) UpdatePaymentStatus(id uuid.UUID, req *UpdatePaymentStatusRequest) error {
	_, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	validStatuses := map[string]bool{
		"pending":  true,
		"paid":     true,
		"failed":   true,
		"refunded": true,
	}

	if !validStatuses[req.Status] {
		return errors.New("invalid payment status")
	}

	return s.repo.UpdatePaymentStatus(id, req.Status)
}
