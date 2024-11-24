package handler

import (
	"go-ticket/service"
	"go-ticket/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

func (h *TransactionHandler) RegisterRoutes(app *fiber.App) {
	transactions := app.Group("/v1/transactions")
	transactions.Get("/", h.GetAllTransactions)
	transactions.Get("/:id", h.GetTransactionById)
	transactions.Get("/user/:userId", h.GetTransactionsByUserId)
	transactions.Post("/", h.CreateTransaction)
	transactions.Put("/:id/status", h.UpdateTransactionStatus)
	transactions.Put("/:id/payment-status", h.UpdatePaymentStatus)
}

func (h *TransactionHandler) GetAllTransactions(c *fiber.Ctx) error {
	transactions, err := h.service.GetAllTransactions()
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Transactions retrieved successfully", transactions)
}

func (h *TransactionHandler) GetTransactionById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid transaction ID")
	}

	transaction, err := h.service.GetTransactionById(id)
	if err != nil {
		return utils.SendNotFoundResponse(c, "Transaction not found")
	}

	return utils.SendSuccessResponse(c, "Transaction retrieved successfully", transaction)
}

func (h *TransactionHandler) GetTransactionsByUserId(c *fiber.Ctx) error {
	userId, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid user ID")
	}

	transactions, err := h.service.GetTransactionsByUserId(userId)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Transactions retrieved successfully", transactions)
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var req service.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	transaction, err := h.service.CreateTransaction(&req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendCreatedResponse(c, "Transaction created successfully", transaction)
}

func (h *TransactionHandler) UpdateTransactionStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid transaction ID")
	}

	var req service.UpdateTransactionStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	err = h.service.UpdateTransactionStatus(id, &req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Transaction status updated successfully", nil)
}

func (h *TransactionHandler) UpdatePaymentStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid transaction ID")
	}

	var req service.UpdatePaymentStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	err = h.service.UpdatePaymentStatus(id, &req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Payment status updated successfully", nil)
}
