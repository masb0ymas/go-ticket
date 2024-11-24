package main

import (
	"log"
	"os"

	"go-ticket/database"
	"go-ticket/handler"
	"go-ticket/repository"
	"go-ticket/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.DB.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Initialize repositories
	eventRepo := repository.NewEventRepository(database.DB)
	scheduleRepo := repository.NewScheduleRepository(database.DB)
	locationRepo := repository.NewLocationRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	ticketTypeRepo := repository.NewTicketTypeRepository(database.DB)
	transactionRepo := repository.NewTransactionRepository(database.DB)
	transactionDetailRepo := repository.NewTransactionDetailRepository(database.DB)

	// Initialize services
	eventService := service.NewEventService(eventRepo)
	scheduleService := service.NewScheduleService(scheduleRepo)
	locationService := service.NewLocationService(locationRepo)
	userService := service.NewUserService(userRepo)
	ticketTypeService := service.NewTicketTypeService(ticketTypeRepo)
	transactionService := service.NewTransactionService(transactionRepo, transactionDetailRepo, ticketTypeRepo)

	// Initialize handlers
	eventHandler := handler.NewEventHandler(eventService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	locationHandler := handler.NewLocationHandler(locationService)
	userHandler := handler.NewUserHandler(userService)
	ticketTypeHandler := handler.NewTicketTypeHandler(ticketTypeService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Register routes
	eventHandler.RegisterRoutes(app)
	scheduleHandler.RegisterRoutes(app)
	locationHandler.RegisterRoutes(app)
	userHandler.RegisterRoutes(app)
	ticketTypeHandler.RegisterRoutes(app)
	transactionHandler.RegisterRoutes(app)

	// Get port from environment variable or use default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
