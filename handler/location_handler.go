package handler

import (
	"go-ticket/service"
	"go-ticket/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LocationHandler struct {
	service *service.LocationService
}

func NewLocationHandler(service *service.LocationService) *LocationHandler {
	return &LocationHandler{
		service: service,
	}
}

func (h *LocationHandler) RegisterRoutes(app *fiber.App) {
	locations := app.Group("/v1/locations")
	locations.Get("/", h.GetAllLocations)
	locations.Get("/:id", h.GetLocationById)
	locations.Post("/", h.CreateLocation)
	locations.Put("/:id", h.UpdateLocation)
	locations.Delete("/:id", h.DeleteLocation)
	locations.Post("/search", h.SearchLocations)
}

func (h *LocationHandler) GetAllLocations(c *fiber.Ctx) error {
	locations, err := h.service.GetAllLocations()
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Locations retrieved successfully", locations)
}

func (h *LocationHandler) GetLocationById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid location ID")
	}

	location, err := h.service.GetLocationById(id)
	if err != nil {
		return utils.SendNotFoundResponse(c, "Location not found")
	}

	return utils.SendSuccessResponse(c, "Location retrieved successfully", location)
}

func (h *LocationHandler) CreateLocation(c *fiber.Ctx) error {
	var req service.CreateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	location, err := h.service.CreateLocation(&req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendCreatedResponse(c, "Location created successfully", location)
}

func (h *LocationHandler) UpdateLocation(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid location ID")
	}

	var req service.UpdateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	location, err := h.service.UpdateLocation(id, &req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Location updated successfully", location)
}

func (h *LocationHandler) DeleteLocation(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid location ID")
	}

	err = h.service.DeleteLocation(id)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Location deleted successfully", nil)
}

func (h *LocationHandler) SearchLocations(c *fiber.Ctx) error {
	var req service.SearchLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	locations, err := h.service.SearchLocations(&req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Locations retrieved successfully", locations)
}
