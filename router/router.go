package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iaroslavagoncharova/react-go/handlers"
	"github.com/iaroslavagoncharova/react-go/middlewares"
)

func SetupRoutes(app *fiber.App, h *handlers.Handlers) {
	// collections routes
	app.Get("/api/collections", h.GetCollections)
	app.Post("/api/login", h.Login)
	
	app.Use(middlewares.AuthMiddleware)

	app.Post("/api/collections", middlewares.AuthMiddleware, h.CreateCollection)
	app.Patch("/api/collections/:id", middlewares.AuthMiddleware, h.UpdateCollection)
	app.Delete("/api/collections/:id", middlewares.AuthMiddleware, h.DeleteCollection)

	// // words routes
	// app.Get("/api/collections/:collectionId/words", getWordsByCollection)
	// app.Post("/api/collections/:collectionId/words", AuthMiddleware, createWord)
	// app.Patch("/api/words/:id", AuthMiddleware, updateWord)
	// app.Delete("/api/words/:id", AuthMiddleware, deleteWord)

	// // users routes
	// app.Get("/api/users", getUsers)
	// app.Post("/api/users", createUser)
	// app.Patch("/api/users/:id", AuthMiddleware, updateUser)
	// app.Delete("/api/users/:id", AuthMiddleware, deleteUser)
}
