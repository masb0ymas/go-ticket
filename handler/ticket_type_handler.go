package handler

import (
	"go-ticket/service"
	"go-ticket/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TicketTypeHandler struct {
	service *service.TicketTypeService
}

func NewTicketTypeHandler(service *service.TicketTypeService) *TicketTypeHandler {
	return &TicketTypeHandler{
		service: service,
	}
}

func (h *TicketTypeHandler) RegisterRoutes(app *fiber.App) {
	ticketTypes := app.Group("/v1/ticket-types")
	ticketTypes.Get("/", h.GetAllTicketTypes)
	ticketTypes.Get("/:id", h.GetTicketTypeById)
	ticketTypes.Get("/event/:eventId", h.GetTicketTypesByEventId)
	ticketTypes.Get("/event/:eventId/available", h.GetAvailableTicketTypes)
	ticketTypes.Post("/", h.CreateTicketType)
	ticketTypes.Put("/:id", h.UpdateTicketType)
	ticketTypes.Delete("/:id", h.DeleteTicketType)
}

func (h *TicketTypeHandler) GetAllTicketTypes(c *fiber.Ctx) error {
	ticketTypes, err := h.service.GetAllTicketTypes()
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Ticket types retrieved successfully", ticketTypes)
}

func (h *TicketTypeHandler) GetTicketTypeById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid ticket type ID")
	}

	ticketType, err := h.service.GetTicketTypeById(id)
	if err != nil {
		return utils.SendNotFoundResponse(c, "Ticket type not found")
	}

	return utils.SendSuccessResponse(c, "Ticket type retrieved successfully", ticketType)
}

func (h *TicketTypeHandler) GetTicketTypesByEventId(c *fiber.Ctx) error {
	eventId, err := uuid.Parse(c.Params("eventId"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid event ID")
	}

	ticketTypes, err := h.service.GetTicketTypesByEventId(eventId)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Ticket types retrieved successfully", ticketTypes)
}

func (h *TicketTypeHandler) GetAvailableTicketTypes(c *fiber.Ctx) error {
	eventId, err := uuid.Parse(c.Params("eventId"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid event ID")
	}

	ticketTypes, err := h.service.GetAvailableTicketTypes(eventId)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Available ticket types retrieved successfully", ticketTypes)
}

func (h *TicketTypeHandler) CreateTicketType(c *fiber.Ctx) error {
	var req service.CreateTicketTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	ticketType, err := h.service.CreateTicketType(&req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendCreatedResponse(c, "Ticket type created successfully", ticketType)
}

func (h *TicketTypeHandler) UpdateTicketType(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid ticket type ID")
	}

	var req service.UpdateTicketTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	ticketType, err := h.service.UpdateTicketType(id, &req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Ticket type updated successfully", ticketType)
}

func (h *TicketTypeHandler) DeleteTicketType(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid ticket type ID")
	}

	err = h.service.DeleteTicketType(id)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Ticket type deleted successfully", nil)
}
