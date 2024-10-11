package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/iaroslavagoncharova/react-go/middlewares"
	"github.com/iaroslavagoncharova/react-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handlers) Login(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	loginData := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	if err := c.BodyParser(loginData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// find user by email
	var dbUser models.User
	err := h.UsersCollection.FindOne(context.Background(), bson.M{"email": loginData.Email}).Decode(&dbUser)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginData.Password))
	if err != nil {
		// if passwords don't match
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// jwt token generation
	token, err := middlewares.GenerateJWT(dbUser.ID.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": token})
}
