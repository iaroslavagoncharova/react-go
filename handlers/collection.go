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

// @Summary Get collections
// @Description Get a list of all collections
// @Tags collections
// @Produce json
// @Success 200 {array} models.Collection "List of collections"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /collections [get]
func (h *Handlers) GetCollections(c *fiber.Ctx) error {
	var collections []models.Collection

	cursor, err := h.CollectionsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var col models.Collection
		if err := cursor.Decode(&col); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "Error decoding collection"})
		}
		collections = append(collections, col)
	}
	return c.JSON(collections)
}

// @Summary Get collection by ID
// @Description Get a collection by ID
// @Tags collections
// @Produce json
// @Param id path string true "Collection ID"
// @Success 200 {object} models.Collection "Collection object"
// @Failure 400 {object} models.ErrorResponse "Invalid collection ID"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /collections/{id} [get]
func (h *Handlers) GetCollectionById(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid collection ID"})
	}

	var col models.Collection
	err = h.CollectionsCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&col)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}
	return c.JSON(col)
}

// @Summary Create a collection
// @Description Create a new collection
// @Tags collections
// @Accept json
// @Produce json
// @Param collection body models.Collection true "Collection object"
// @Success 201 {object} models.CollectionResponse "Collection created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input: Name and description are required"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /collections [post]
func (h *Handlers) CreateCollection(c *fiber.Ctx) error {
	var col models.Collection

	if err := c.BodyParser(&col); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	// use userId from the token
	userId := c.Locals("user").(string)
	col.UserID = userId

	fmt.Println(col)
	// validate user input
	if err := validation.ValidateStruct(col); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: err.Error()})
	}

	insertResult, err := h.CollectionsCollection.InsertOne(context.Background(), col)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	col.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(models.CollectionResponse{
		Message:    "Collection created successfully",
		Collection: col,
	})
}

// @Summary Update a collection
// @Description Update a collection by ID
// @Tags collections
// @Accept json
// @Produce json
// @Param id path string true "Collection ID"
// @Param collection body models.UpdateCollection true "Collection object"
// @Success 200 {object} models.CollectionResponse "Collection updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid ID or request body"
// @Failure 403 {object} models.ErrorResponse "You are not authorized to update this collection"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /collections/{id} [patch]
func (h *Handlers) UpdateCollection(c *fiber.Ctx) error {
	var col models.UpdateCollection

	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid ID"})
	}

	if err := c.BodyParser(&col); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	update := bson.M{}

	if col.Name != nil {
		update["name"] = *col.Name
	}
	if col.Description != nil {
		update["description"] = *col.Description
	}

	if len(update) == 0 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "No fields to update"})
	}

	// validate user input
	if err := validation.ValidateStruct(col); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: err.Error()})
	}

	// use userId from the token
	userId := c.Locals("user").(string)

	// check if the collection belongs to the user
	var collection models.Collection
	err = h.CollectionsCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&collection)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}
	if collection.UserID != userId {
		return c.Status(403).JSON(models.ErrorResponse{Error: "You are not authorized to update this collection"})
	}

	filter := bson.M{"_id": objectID}
	_, err = h.CollectionsCollection.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}
	return c.Status(200).JSON(models.CollectionResponse{
		Message:    "Collection updated successfully",
		Collection: collection,
	})
}

// @Summary Delete a collection
// @Description Delete a collection by ID
// @Tags collections
// @Produce json
// @Param id path string true "Collection ID"
// @Success 200 {object} models.MessageResponse "Collection deleted successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid ID"
// @Failure 403 {object} models.ErrorResponse "You are not authorized to delete this collection"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /collections/{id} [delete]
func (h *Handlers) DeleteCollection(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid ID"})
	}

	// use userId from the token
	userId := c.Locals("user").(string)

	// check if the collection belongs to the user
	var col models.Collection
	err = h.CollectionsCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&col)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}
	if col.UserID != userId {
		return c.Status(403).JSON(models.ErrorResponse{Error: "You are not authorized to delete this collection"})
	}

	filter := bson.M{"_id": objectID}
	_, err = h.CollectionsCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	return c.Status(200).JSON(models.MessageResponse{Message: "Collection deleted successfully"})
}
