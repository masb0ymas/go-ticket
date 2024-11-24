package handler

import (
	"go-ticket/service"
	"go-ticket/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{
		service: service,
	}
}

func (h *EventHandler) RegisterRoutes(app *fiber.App) {
	events := app.Group("/v1/events")
	events.Get("/", h.GetAllEvents)
	events.Get("/:id", h.GetEventById)
	events.Post("/", h.CreateEvent)
	events.Put("/:id", h.UpdateEvent)
	events.Delete("/:id", h.DeleteEvent)
}

func (h *EventHandler) GetAllEvents(c *fiber.Ctx) error {
	events, err := h.service.GetAllEvents()
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Events retrieved successfully", events)
}

func (h *EventHandler) GetEventById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid event ID")
	}

	event, err := h.service.GetEventById(id)
	if err != nil {
		return utils.SendNotFoundResponse(c, "Event not found")
	}

	return utils.SendSuccessResponse(c, "Event retrieved successfully", event)
}

func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	var req service.CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	event, err := h.service.CreateEvent(&req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendCreatedResponse(c, "Event created successfully", event)
}

func (h *EventHandler) UpdateEvent(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid event ID")
	}

	var req service.UpdateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	event, err := h.service.UpdateEvent(id, &req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Event updated successfully", event)
}

func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid event ID")
	}

	err = h.service.DeleteEvent(id)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Event deleted successfully", nil)
}
