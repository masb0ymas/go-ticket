package handler

import (
	"go-ticket/service"
	"go-ticket/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	users := app.Group("/v1/users")
	users.Get("/", h.GetAllUsers)
	users.Get("/:id", h.GetUserById)
	users.Post("/", h.CreateUser)
	users.Put("/:id", h.UpdateUser)
	users.Delete("/:id", h.DeleteUser)
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "Users retrieved successfully", users)
}

func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid user ID")
	}

	user, err := h.service.GetUserById(id)
	if err != nil {
		return utils.SendNotFoundResponse(c, "User not found")
	}

	return utils.SendSuccessResponse(c, "User retrieved successfully", user)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req service.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendCreatedResponse(c, "User created successfully", user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid user ID")
	}

	var req service.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendBadRequestResponse(c, "Invalid request body")
	}

	user, err := h.service.UpdateUser(id, &req)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "User updated successfully", user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendBadRequestResponse(c, "Invalid user ID")
	}

	err = h.service.DeleteUser(id)
	if err != nil {
		return utils.SendInternalServerErrorResponse(c, err)
	}

	return utils.SendSuccessResponse(c, "User deleted successfully", nil)
}
