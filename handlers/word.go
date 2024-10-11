package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/iaroslavagoncharova/react-go/models"
	"github.com/iaroslavagoncharova/react-go/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handlers) GetWordsByCollection(c *fiber.Ctx) error {
	collectionID := c.Params("collectionId")
	objectID, err := primitive.ObjectIDFromHex(collectionID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid collection ID",
		})
	}

	var words []models.Word
	filter := bson.M{"collectionId": objectID}
	cursor, err := h.WordsCollection.Find(context.Background(), filter)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var word models.Word
		if err := cursor.Decode(&word); err != nil {
			return err
		}
		words = append(words, word)
	}
	return c.JSON(words)
}

func (h *Handlers) CreateWord(c *fiber.Ctx) error {
	collectionID := c.Params("collectionId")
	objectID, err := primitive.ObjectIDFromHex(collectionID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid collection ID",
		})
	}

	word := new(models.Word)
	if err := c.BodyParser(word); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// validate user input
	if err := validation.ValidateStruct(word); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
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
			"error": "You are not authorized to add words to this collection",
		})
	}

	word.CollectionID = objectID
	insertResult, err := h.WordsCollection.InsertOne(context.Background(), word)
	if err != nil {
		return err
	}

	word.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(word)
}

func (h *Handlers) UpdateWord(c *fiber.Ctx) error {
	var updateWord models.UpdateWord

	wordID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(wordID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid word ID",
		})
	}

	if err := c.BodyParser(&updateWord); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
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
		return c.Status(400).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	// use userId from the token
	userId := c.Locals("user").(string)

	// check if the collection belongs to the user
	var wordObj models.Word
	err = h.WordsCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&wordObj)
	if err != nil {
		return err
	}

	var col models.Collection
	err = h.CollectionsCollection.FindOne(context.Background(), bson.M{"_id": wordObj.CollectionID}).Decode(&col)
	if err != nil {
		return err
	}
	if col.UserID != userId {
		return c.Status(403).JSON(fiber.Map{
			"error": "You are not authorized to update this word",
		})
	}

	filter := bson.M{"_id": objectID}

	_, err = h.WordsCollection.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Word updated successfully",
	})
}

func (h *Handlers) DeleteWord(c *fiber.Ctx) error {
	wordID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(wordID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid word ID",
		})
	}

	// use userId from the token
	userId := c.Locals("user").(string)

	// check if the collection belongs to the user
	var word models.Word
	err = h.WordsCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&word)
	if err != nil {
		return err
	}

	var col models.Collection
	err = h.CollectionsCollection.FindOne(context.Background(), bson.M{"_id": word.CollectionID}).Decode(&col)
	if err != nil {
		return err
	}
	if col.UserID != userId {
		return c.Status(403).JSON(fiber.Map{
			"error": "You are not authorized to delete this word",
		})
	}

	filter := bson.M{"_id": objectID}
	_, err = h.WordsCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Word deleted successfully",
	})
}