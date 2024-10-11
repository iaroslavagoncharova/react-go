package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/iaroslavagoncharova/react-go/models"
	"github.com/iaroslavagoncharova/react-go/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Produce json
// @Success 200 {array} models.UserWithoutPassword "List of users"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /users [get]
func (h *Handlers) GetUsers(c *fiber.Ctx) error {
	var users []models.User

	cursor, err := h.UsersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "Error decoding user"})
		}
		users = append(users, user)
	}
	// return users without password
	var usersWithoutPassword []models.UserWithoutPassword
	for _, user := range users {
		usersWithoutPassword = append(usersWithoutPassword, models.UserWithoutPassword{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		})
	}
	return c.JSON(usersWithoutPassword)
}

// @Summary Get user by ID
// @Description Get a user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.UserWithoutPassword "User object"
// @Failure 400 {object} models.ErrorResponse "Invalid ID"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /users/{id} [get]
func (h *Handlers) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid ID"})
	}

	var user models.User
	err = h.UsersCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	return c.JSON(models.UserWithoutPassword{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	})
}

// @Summary Create a user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 201 {object} models.UserResponse "User created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input: Username, email, and password are required"
// @Failure 500 {object} models.ErrorResponse "Internal server error: Error hashing password"
// @Router /users [post]
func (h *Handlers) CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return err
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Username, email, and password are required"})
	}

	// validate user input
	err := validation.ValidateStruct(user)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: err.Error()})
	}

	// password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Error hashing password"})
	}
	user.Password = string(hashedPassword)

	// user default role is "user"
	user.Role = "user"

	insertResult, err := h.UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	user.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(models.UserResponse{
		Message: "User created successfully",
		User: models.UserWithoutPassword{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	})
}

// @Summary Update a user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UpdateUser true "User object with fields to update"
// @Success 200 {object} models.UserResponse "User updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 403 {object} models.ErrorResponse "Unauthorized: You are not authorized to update this user"
// @Failure 500 {object} models.ErrorResponse "Internal server error: Error updating user"
// @Router /users [patch]
func (h *Handlers) UpdateUser(c *fiber.Ctx) error {
	var updateUser models.UpdateUser
	role := c.Locals("role").(string)

	// user id from params
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid ID"})
	}

	// get the current id from token
	tokenUserId := c.Locals("user_id").(string)

	if err := c.BodyParser(&updateUser); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	update := bson.M{}

	if updateUser.Username != nil {
		update["username"] = *updateUser.Username
	}
	if updateUser.Email != nil {
		update["email"] = *updateUser.Email
	}
	if updateUser.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*updateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "Error hashing password"})
		}
		update["password"] = string(hashedPassword)
	}

	if len(update) == 0 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "No fields to update"})
	}

	// check privileges
	if role != "admin" {
		// if not admin, check if the user is updating their own profile
		if tokenUserId != id {
			return c.Status(403).JSON(models.ErrorResponse{Error: "You are not authorized to update this user"})
		}
	}

	// get the user by id
	var user models.User
	err = h.UsersCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	filter := bson.M{"_id": objectID}
	_, err = h.UsersCollection.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	return c.Status(200).JSON(models.UserResponse{
		Message: "User updated successfully",
		User: models.UserWithoutPassword{
			ID:       objectID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	})
}

// @Summary Delete a user
// @Description Delete a user profile
// @Tags users
// @Produce json
// @Success 200 {object} models.MessageResponse "User deleted successfully"
// @Failure 403 {object} models.ErrorResponse "Unauthorized: You are not authorized to delete this user"
// @Failure 500 {object} models.ErrorResponse "Internal server error: Error deleting user"
// @Router /users [delete]
func (h *Handlers) DeleteUser(c *fiber.Ctx) error {
	// user id from params
	id := c.Params("id")
	role := c.Locals("role").(string)

	// use the current user id from token
	userId := c.Locals("user_id").(string)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid ID"})
	}

	// check privileges
	if role != "admin" {
		// if not admin, check if the user is deleting their own profile
		if userId != id {
			return c.Status(403).JSON(models.ErrorResponse{Error: "You are not authorized to delete this user"})
		}
	}

	filter := bson.M{"_id": objectID}
	_, err = h.UsersCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	return c.Status(200).JSON(models.MessageResponse{Message: "User deleted successfully"})
}
