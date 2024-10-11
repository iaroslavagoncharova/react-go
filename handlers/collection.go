package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/iaroslavagoncharova/react-go/models"
	"github.com/iaroslavagoncharova/react-go/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handlers struct {
	CollectionsCollection *mongo.Collection
	UsersCollection       *mongo.Collection
	WordsCollection       *mongo.Collection
}

func (h *Handlers) GetCollections(c *fiber.Ctx) error {
	var collections []models.Collection

	cursor, err := h.CollectionsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var col models.Collection
		if err := cursor.Decode(&col); err != nil {
			return err
		}
		collections = append(collections, col)
	}
	return c.JSON(collections)
}

func (h *Handlers) CreateCollection(c *fiber.Ctx) error {
	var col models.Collection

	if err := c.BodyParser(&col); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// use userId from the token
	userId := c.Locals("user").(string)
	col.UserID = userId

	fmt.Println(col)
	// validate user input
	if err := validation.ValidateStruct(col); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	insertResult, err := h.CollectionsCollection.InsertOne(context.Background(), col)
	if err != nil {
		return err
	}

	col.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(col)
}

func (h *Handlers) UpdateCollection(c *fiber.Ctx) error {
	var col models.UpdateCollection

	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := c.BodyParser(&col); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	update := bson.M{}

	if col.Name != nil {
		update["name"] = *col.Name
	}
	if col.Description != nil {
		update["description"] = *col.Description
	}

	if len(update) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	// validate user input
	if err := validation.ValidateStruct(col); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// use userId from the token
	userId := c.Locals("user").(string)

	// check if the collection belongs to the user
	var collection models.Collection
	err = h.CollectionsCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&collection)
	if err != nil {
		return err
	}
	if collection.UserID != userId {
		return c.Status(403).JSON(fiber.Map{
			"error": "You are not authorized to update this collection",
		})
	}

	filter := bson.M{"_id": objectID}
	_, err = h.CollectionsCollection.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Collection updated successfully",
	})
}

func (h *Handlers) DeleteCollection(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	// use userId from the token
	userId := c.Locals("user").(string)

	// check if the collection belongs to the user
	var col models.Collection
	err = h.CollectionsCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&col)
	if err != nil {
		return err
	}
	if col.UserID != userId {
		return c.Status(403).JSON(fiber.Map{
			"error": "You are not authorized to delete this collection",
		})
	}

	filter := bson.M{"_id": objectID}
	_, err = h.CollectionsCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Collection deleted successfully",
	})
}
