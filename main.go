package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/iaroslavagoncharova/react-go/handlers"
	"github.com/iaroslavagoncharova/react-go/router"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/iaroslavagoncharova/react-go/config"
)

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
		WordsCollection:       wordsCollection,
	}

	// Setup routes
	router.SetupRoutes(app, handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Fatal(app.Listen(":" + port))
}