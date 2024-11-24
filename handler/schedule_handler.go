package handler

import (
	"go-ticket/service"
	"go-ticket/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ScheduleHandler struct {
	service *service.ScheduleService
}

func NewScheduleHandler(service *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		service: service,
	}
}

func (h *ScheduleHandler) RegisterRoutes(app *fiber.App) {
	schedules := app.Group("/v1/schedules")
	schedules.Get("/", h.GetAllSchedules)
	schedules.Get("/:id", h.GetScheduleById)
	schedules.Post("/", h.CreateSchedule)
	schedules.Put("/:id", h.UpdateSchedule)
	schedules.Delete("/:id", h.DeleteSchedule)
	schedules.Post("/search", h.SearchSchedules)
}

func (h *ScheduleHandler) GetAllSchedules(c *fiber.Ctx) error {
	schedules, err := h.service.GetAllSchedules()
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Schedules retrieved successfully", schedules)
}

func (h *ScheduleHandler) GetScheduleById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid schedule ID")
	}

	schedule, err := h.service.GetScheduleById(id)
	if err != nil {
		return utils.SendNotFoundResponse(c, "Schedule not found")
	}

	return utils.SendSuccessResponse(c, "Schedule retrieved successfully", schedule)
}

func (h *ScheduleHandler) CreateSchedule(c *fiber.Ctx) error {
	var req service.CreateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	schedule, err := h.service.CreateSchedule(&req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendCreatedResponse(c, "Schedule created successfully", schedule)
}

func (h *ScheduleHandler) UpdateSchedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid schedule ID")
	}

	var req service.UpdateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	schedule, err := h.service.UpdateSchedule(id, &req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Schedule updated successfully", schedule)
}

func (h *ScheduleHandler) DeleteSchedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid schedule ID")
	}

	err = h.service.DeleteSchedule(id)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Schedule deleted successfully", nil)
}

func (h *ScheduleHandler) SearchSchedules(c *fiber.Ctx) error {
	var req service.SearchScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	schedules, err := h.service.SearchSchedules(&req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Schedules retrieved successfully", schedules)
}
