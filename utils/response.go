package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func SendResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func SendSuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return SendResponse(c, fiber.StatusOK, message, data)
}

func SendCreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return SendResponse(c, fiber.StatusCreated, message, data)
}

func SendErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return SendResponse(c, statusCode, message, nil)
}

func SendBadRequestResponse(c *fiber.Ctx, message string) error {
	return SendErrorResponse(c, fiber.StatusBadRequest, message)
}

func SendNotFoundResponse(c *fiber.Ctx, message string) error {
	return SendErrorResponse(c, fiber.StatusNotFound, message)
}

func SendInternalServerErrorResponse(c *fiber.Ctx, err error) error {
	return SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
}
