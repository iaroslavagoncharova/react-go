package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/iaroslavagoncharova/react-go/handlers"
	"github.com/iaroslavagoncharova/react-go/router"
	// "github.com/iaroslavagoncharova/react-go/middlewares"
	"github.com/iaroslavagoncharova/react-go/config"
)

// type Word struct {
// 	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
// 	CollectionID primitive.ObjectID `json:"collectionId" bson:"collectionId" validate:"required"`
// 	Word         string             `json:"word" bson:"word" validate:"required,min=1,max=50"`
// 	Translation  string             `json:"translation" bson:"translation" validate:"required,min=1,max=100"`
// 	Difficulty   int                `json:"difficulty" bson:"difficulty" validate:"required,gte=1,lte=5"`
// }

// type User struct {
// 	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
// 	Username string             `json:"username" bson:"username" validate:"required,min=3,max=20"`
// 	Email    string             `json:"email" bson:"email" validate:"required,email"`
// 	Password string             `json:"password" bson:"password" validate:"required,min=6"`
// }

var collectionsCollection *mongo.Collection
var wordsCollection *mongo.Collection
var usersCollection *mongo.Collection

func main() {
	fmt.Println("Starting the server...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.LoadConfig()

	MONGO_DB := os.Getenv("MONGODB_URI")

	clientOptions := options.Client().ApplyURI(MONGO_DB)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	fmt.Println("Connected to MongoDB!")

	db := client.Database("react_go")
	collectionsCollection = db.Collection("collections")
	wordsCollection = db.Collection("words")
	usersCollection = db.Collection("users")

	app := fiber.New()

	handler := &handlers.Handlers{
		CollectionsCollection: collectionsCollection,
		UsersCollection:       usersCollection,
	}

	// Setup routes
	router.SetupRoutes(app, handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Fatal(app.Listen(":" + port))
}

// func getWordsByCollection(c *fiber.Ctx) error {
// 	collectionID := c.Params("collectionId")
// 	objectID, err := primitive.ObjectIDFromHex(collectionID)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error": "Invalid collection ID",
// 		})
// 	}

// 	var words []Word
// 	filter := bson.M{"collectionId": objectID}
// 	cursor, err := wordsCollection.Find(context.Background(), filter)
// 	if err != nil {
// 		return err
// 	}
// 	defer cursor.Close(context.Background())

// 	for cursor.Next(context.Background()) {
// 		var word Word
// 		if err := cursor.Decode(&word); err != nil {
// 			return err
// 		}
// 		words = append(words, word)
// 	}
// 	return c.JSON(words)
// }

// func createWord(c *fiber.Ctx) error {
// 	collectionID := c.Params("collectionId")
// 	objectID, err := primitive.ObjectIDFromHex(collectionID)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error": "Invalid collection ID",
// 		})
// 	}

// 	word := new(Word)
// 	if err := c.BodyParser(word); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
// 	}

// 	// validate user input
// 	if err := validate.Struct(word); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	// use userId from the token
// 	userId := c.Locals("user").(string)

// 	// check if the collection belongs to the user
// 	var col Collection
// 	err = collectionsCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&col)
// 	if err != nil {
// 		return err
// 	}
// 	if col.UserID != userId {
// 		return c.Status(403).JSON(fiber.Map{
// 			"error": "You are not authorized to add words to this collection",
// 		})
// 	}

// 	word.CollectionID = objectID
// 	insertResult, err := wordsCollection.InsertOne(context.Background(), word)
// 	if err != nil {
// 		return err
// 	}

// 	word.ID = insertResult.InsertedID.(primitive.ObjectID)

// 	return c.Status(201).JSON(word)
// }

// func updateWord(c *fiber.Ctx) error {
// 	wordID := c.Params("id")
// 	objectID, err := primitive.ObjectIDFromHex(wordID)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error": "Invalid word ID",
// 		})
// 	}

// 	word := new(Word)
// 	if err := c.BodyParser(word); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
// 	}

// 	// validate user input
// 	if err := validate.Struct(word); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	// use userId from the token
// 	userId := c.Locals("user").(string)

// 	// check if the collection belongs to the user
// 	var wordObj Word
// 	err = wordsCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&wordObj)
// 	if err != nil {
// 		return err
// 	}

// 	var col Collection
// 	err = collectionsCollection.FindOne(context.Background(), bson.M{"_id": wordObj.CollectionID}).Decode(&col)
// 	if err != nil {
// 		return err
// 	}
// 	if col.UserID != userId {
// 		return c.Status(403).JSON(fiber.Map{
// 			"error": "You are not authorized to update this word",
// 		})
// 	}

// 	filter := bson.M{"_id": objectID}
// 	update := bson.M{"$set": bson.M{
// 		"word":        c.FormValue("word"),
// 		"translation": c.FormValue("translation"),
// 		"difficulty":  c.FormValue("difficulty"),
// 	}}

// 	_, err = wordsCollection.UpdateOne(context.Background(), filter, update)
// 	if err != nil {
// 		return err
// 	}
// 	return c.Status(200).JSON(fiber.Map{
// 		"message": "Word updated successfully",
// 	})
// }

// func deleteWord(c *fiber.Ctx) error {
// 	wordID := c.Params("id")
// 	objectID, err := primitive.ObjectIDFromHex(wordID)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error": "Invalid word ID",
// 		})
// 	}

// 	// use userId from the token
// 	userId := c.Locals("user").(string)

// 	// check if the collection belongs to the user
// 	var word Word
// 	err = wordsCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&word)
// 	if err != nil {
// 		return err
// 	}

// 	var col Collection
// 	err = collectionsCollection.FindOne(context.Background(), bson.M{"_id": word.CollectionID}).Decode(&col)
// 	if err != nil {
// 		return err
// 	}
// 	if col.UserID != userId {
// 		return c.Status(403).JSON(fiber.Map{
// 			"error": "You are not authorized to delete this word",
// 		})
// 	}

// 	filter := bson.M{"_id": objectID}
// 	_, err = wordsCollection.DeleteOne(context.Background(), filter)
// 	if err != nil {
// 		return err
// 	}

// 	return c.Status(200).JSON(fiber.Map{
// 		"message": "Word deleted successfully",
// 	})
// }

// func getUsers(c *fiber.Ctx) error {
// 	var users []User

// 	cursor, err := usersCollection.Find(context.Background(), bson.M{})
// 	if err != nil {
// 		return err
// 	}
// 	defer cursor.Close(context.Background())

// 	for cursor.Next(context.Background()) {
// 		var user User
// 		if err := cursor.Decode(&user); err != nil {
// 			return err
// 		}
// 		users = append(users, user)
// 	}
// 	return c.JSON(users)
// }

// func createUser(c *fiber.Ctx) error {
// 	user := new(User)
// 	if err := c.BodyParser(user); err != nil {
// 		return err
// 	}

// 	if user.Username == "" || user.Email == "" || user.Password == "" {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error": "Username, email, and password are required",
// 		})
// 	}

// 	// validate user input
// 	err := validate.Struct(user)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	// password hashing
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
// 	}
// 	user.Password = string(hashedPassword)

// 	insertResult, err := usersCollection.InsertOne(context.Background(), user)
// 	if err != nil {
// 		return err
// 	}

// 	user.ID = insertResult.InsertedID.(primitive.ObjectID)

// 	return c.Status(201).JSON(user)
// }

// func updateUser(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error": "Invalid ID",
// 		})
// 	}

// 	update := bson.M{"$set": bson.M{
// 		"username": c.FormValue("username"),
// 		"email":    c.FormValue("email"),
// 		"password": c.FormValue("password"),
// 	}}

// 	// use userId from the token
// 	userId := c.Locals("user").(string)

// 	// check if the user is updating their own profile
// 	var user User
// 	err = usersCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
// 	if err != nil {
// 		return err
// 	}
// 	if user.ID.Hex() != userId {
// 		return c.Status(403).JSON(fiber.Map{
// 			"error": "You are not authorized to update this user",
// 		})
// 	}

// 	filter := bson.M{"_id": objectID}
// 	_, err = usersCollection.UpdateOne(context.Background(), filter, update)
// 	if err != nil {
// 		return err
// 	}
// 	return c.Status(200).JSON(fiber.Map{
// 		"message": "User updated successfully",
// 	})
// }

// func deleteUser(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error": "Invalid ID",
// 		})
// 	}

// 	// use userId from the token
// 	userId := c.Locals("user").(string)

// 	// check if the user is deleting their own profile
// 	var user User
// 	err = usersCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
// 	if err != nil {
// 		return err
// 	}
// 	if user.ID.Hex() != userId {
// 		return c.Status(403).JSON(fiber.Map{
// 			"error": "You are not authorized to delete this user",
// 		})
// 	}

// 	filter := bson.M{"_id": objectID}
// 	_, err = usersCollection.DeleteOne(context.Background(), filter)
// 	if err != nil {
// 		return err
// 	}

// 	return c.Status(200).JSON(fiber.Map{
// 		"message": "User deleted successfully",
// 	})
// }

// func loginHandler(c *fiber.Ctx) error {
// 	loginData := new(struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	})

// 	if err := c.BodyParser(loginData); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
// 	}

// 	// find user by email
// 	var user User
// 	err := usersCollection.FindOne(context.Background(), bson.M{"email": loginData.Email}).Decode(&user)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
// 	if err != nil {
// 		// if passwords don't match
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
// 	}

// 	// jwt token generation
// 	token, err := GenerateJWT(user.ID.Hex())
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
// 	}

// 	return c.JSON(fiber.Map{"token": token})
// }