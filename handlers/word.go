package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/iaroslavagoncharova/react-go/models"
	"github.com/iaroslavagoncharova/react-go/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Get words by collection
// @Description Get a list of all words in a collection
// @Tags words
// @Produce json
// @Param collectionId path string true "Collection ID"
// @Success 200 {array} models.Word "List of words"
// @Failure 400 {object} models.ErrorResponse "Invalid collection ID"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /collections/{collectionId}/words [get]
func (h *Handlers) GetWordsByCollection(c *fiber.Ctx) error {
	collectionID := c.Params("collectionId")
	objectID, err := primitive.ObjectIDFromHex(collectionID)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid collection ID"})
	}

	var words []models.Word
	filter := bson.M{"collectionId": objectID}
	cursor, err := h.WordsCollection.Find(context.Background(), filter)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var word models.Word
		if err := cursor.Decode(&word); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Error: "Error decoding word"})
		}
		words = append(words, word)
	}
	return c.JSON(words)
}

// @Summary Create a word
// @Description Create a new word
// @Tags words
// @Accept json
// @Produce json
// @Param collectionId path string true "Collection ID"
// @Param word body models.Word true "Word object"
// @Success 201 {object} models.WordResponse "Word created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input: Word and translation are required"
// @Failure 403 {object} models.ErrorResponse "You are not authorized to add words to this collection"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /collections/{collectionId}/words [post]
func (h *Handlers) CreateWord(c *fiber.Ctx) error {
	collectionID := c.Params("collectionId")
	objectID, err := primitive.ObjectIDFromHex(collectionID)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid collection ID"})
	}

	word := new(models.Word)
	if err := c.BodyParser(word); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	// validate user input
	if err := validation.ValidateStruct(word); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: err.Error()})
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
		return c.Status(403).JSON(models.ErrorResponse{Error: "You are not authorized to add words to this collection"})
	}

	word.CollectionID = objectID
	insertResult, err := h.WordsCollection.InsertOne(context.Background(), word)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	word.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(models.WordResponse{
		Message: "Word created successfully",
		Word:    *word,
	})
}

// @Summary Update a word
// @Description Update a word by ID
// @Tags words
// @Accept json
// @Produce json
// @Param id path string true "Word ID"
// @Param word body models.UpdateWord true "Word object with fields to update"
// @Success 200 {object} models.WordResponse "Word updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid word ID or request body"
// @Failure 403 {object} models.ErrorResponse "You are not authorized to update this word"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /words/{id} [patch]
func (h *Handlers) UpdateWord(c *fiber.Ctx) error {
	var updateWord models.UpdateWord

	wordID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(wordID)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid word ID"})
	}

	if err := c.BodyParser(&updateWord); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	update := bson.M{}

	if updateWord.CollectionID != nil {
		update["collectionId"] = *updateWord.CollectionID
	}
	if updateWord.Difficulty != nil {
		update["difficulty"] = *updateWord.Difficulty
	}
	if updateWord.Translation != nil {
		update["translation"] = *updateWord.Translation
	}
	if updateWord.Word != nil {
		update["word"] = *updateWord.Word
	}

	if len(update) == 0 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "No fields to update"})
	}

	// use userId from the token
	userId := c.Locals("user").(string)

	// check if the collection belongs to the user
	var wordObj models.Word
	err = h.WordsCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&wordObj)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	var col models.Collection
	err = h.CollectionsCollection.FindOne(context.Background(), bson.M{"_id": wordObj.CollectionID}).Decode(&col)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}
	if col.UserID != userId {
		return c.Status(403).JSON(models.ErrorResponse{Error: "You are not authorized to update this word"})
	}

	filter := bson.M{"_id": objectID}

	_, err = h.WordsCollection.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}
	return c.Status(200).JSON(models.WordResponse{
		Message: "Word updated successfully",
		Word:    models.Word{ID: objectID, CollectionID: wordObj.CollectionID, Word: *updateWord.Word, Translation: *updateWord.Translation, Difficulty: *updateWord.Difficulty},
	})
}

// @Summary Delete a word
// @Description Delete a word by ID
// @Tags words
// @Param id path string true "Word ID"
// @Success 200 {object} models.MessageResponse "Word deleted successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid word ID"
// @Failure 403 {object} models.ErrorResponse "You are not authorized to delete this word"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /words/{id} [delete]
func (h *Handlers) DeleteWord(c *fiber.Ctx) error {
	wordID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(wordID)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid word ID"})
	}

	// use userId from the token
	userId := c.Locals("user").(string)

	// check if the collection belongs to the user
	var word models.Word
	err = h.WordsCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&word)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	var col models.Collection
	err = h.CollectionsCollection.FindOne(context.Background(), bson.M{"_id": word.CollectionID}).Decode(&col)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}
	if col.UserID != userId {
		return c.Status(403).JSON(models.ErrorResponse{Error: "You are not authorized to delete this word"})
	}

	filter := bson.M{"_id": objectID}
	_, err = h.WordsCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	return c.Status(200).JSON(models.MessageResponse{Message: "Word deleted successfully"})
}
