package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iaroslavagoncharova/react-go/handlers"
	"github.com/iaroslavagoncharova/react-go/middlewares"
)

func SetupRoutes(app *fiber.App, h *handlers.Handlers) {
	// routes without authentication
	publicRoutes := app.Group("/api")
	publicRoutes.Get("/collections", h.GetCollections)
	publicRoutes.Post("/login", h.Login)
	publicRoutes.Get("/users", h.GetUsers)
	publicRoutes.Post("/users", h.CreateUser)
	publicRoutes.Get("/collections/:collectionId/words", h.GetWordsByCollection)

	// routes with authentication
	privateRoutes := app.Group("/api", middlewares.AuthMiddleware)
	privateRoutes.Post("/collections", h.CreateCollection)
	privateRoutes.Patch("/collections/:id", h.UpdateCollection)
	privateRoutes.Delete("/collections/:id", h.DeleteCollection)
	privateRoutes.Patch("/users", h.UpdateUser)
	privateRoutes.Delete("/users", h.DeleteUser)
	privateRoutes.Post("/collections/:collectionId/words", h.CreateWord)
	privateRoutes.Patch("/words/:id", h.UpdateWord)
	privateRoutes.Delete("/words/:id", h.DeleteWord)
}
