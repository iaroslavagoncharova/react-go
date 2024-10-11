package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iaroslavagoncharova/react-go/handlers"
	"github.com/iaroslavagoncharova/react-go/middlewares"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "github.com/iaroslavagoncharova/react-go/docs"
)

func SetupRoutes(app *fiber.App, h *handlers.Handlers) {
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	// routes without authentication
	// show that api is working
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the language learning API!")
	})
	publicRoutes := app.Group("/api")
	publicRoutes.Get("/collections", h.GetCollections)
	publicRoutes.Get("/collections/:id", h.GetCollectionById)
	publicRoutes.Post("/login", h.Login)
	publicRoutes.Get("/users", h.GetUsers)
	publicRoutes.Get("/users/:id", h.GetUserById)
	publicRoutes.Post("/users", h.CreateUser)
	publicRoutes.Get("/collections/:collectionId/words", h.GetWordsByCollection)
	publicRoutes.Get("/words/:id", h.GetWordById)

	// routes with authentication
	privateRoutes := app.Group("/api", middlewares.AuthMiddleware)
	privateRoutes.Post("/collections", h.CreateCollection)
	privateRoutes.Patch("/collections/:id", h.UpdateCollection)
	privateRoutes.Delete("/collections/:id", h.DeleteCollection)
	privateRoutes.Patch("/users/:id", h.UpdateUser)
	privateRoutes.Delete("/users/:id", h.DeleteUser)
	privateRoutes.Post("/collections/:collectionId/words", h.CreateWord)
	privateRoutes.Patch("/words/:id", h.UpdateWord)
	privateRoutes.Delete("/words/:id", h.DeleteWord)

	// admin routes
	// adminRoutes := app.Group("/api", middlewares.AuthMiddleware, middlewares.Authorize("admin"))
}
