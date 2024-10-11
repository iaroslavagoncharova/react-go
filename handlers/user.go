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

func (h *Handlers) GetUsers(c *fiber.Ctx) error {
	var users []models.User

	cursor, err := h.UsersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return err
		}
		users = append(users, user)
	}
	return c.JSON(users)
}

func (h *Handlers) CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return err
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Username, email, and password are required",
		})
	}

	// validate user input
	err := validation.ValidateStruct(user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}
	user.Password = string(hashedPassword)

	insertResult, err := h.UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	user.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(user)
}

func (h *Handlers) UpdateUser(c *fiber.Ctx) error {
	var updateUser models.UpdateUser
	userId := c.Locals("user").(string)

	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := c.BodyParser(&updateUser); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
		}
		update["password"] = string(hashedPassword)
	}

	if len(update) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	// check if the user is updating their own profile
	var user models.User
	err = h.UsersCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return err
	}
	if user.ID.Hex() != userId {
		return c.Status(403).JSON(fiber.Map{
			"error": "You are not authorized to update this user",
		})
	}

	filter := bson.M{"_id": objectID}
	_, err = h.UsersCollection.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

func (h *Handlers) DeleteUser(c *fiber.Ctx) error {
	// use userId from the token
	userId := c.Locals("user").(string)
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	// check if the user is deleting their own profile
	var user models.User
	err = h.UsersCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return err
	}
	if user.ID.Hex() != userId {
		return c.Status(403).JSON(fiber.Map{
			"error": "You are not authorized to delete this user",
		})
	}

	filter := bson.M{"_id": objectID}
	_, err = h.UsersCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
