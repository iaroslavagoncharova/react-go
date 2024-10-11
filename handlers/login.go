package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/iaroslavagoncharova/react-go/middlewares"
	"github.com/iaroslavagoncharova/react-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Login
// @Description Login a user
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginCredentials true "Login credentials"
// @Success 200 {object} models.TokenResponse "Token generated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 401 {object} models.ErrorResponse "Invalid email or password"
// @Failure 500 {object} models.ErrorResponse "Failed to generate token"
// @Router /login [post]
func (h *Handlers) Login(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	loginData := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	if err := c.BodyParser(loginData); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	// find user by email
	var dbUser models.User
	err := h.UsersCollection.FindOne(context.Background(), bson.M{"email": loginData.Email}).Decode(&dbUser)
	if err != nil {
		return c.Status(403).JSON(models.ErrorResponse{Error: "Invalid email or password"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginData.Password))
	if err != nil {
		// if passwords don't match
		return c.Status(403).JSON(models.ErrorResponse{Error: "Invalid email or password"})
	}

	// jwt token generation
	token, err := middlewares.GenerateJWT(dbUser.ID.Hex(), dbUser.Role)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Failed to generate token"})
	}

	return c.JSON(models.TokenResponse{
		Message: "Token generated successfully", 
		Token: token})
}
