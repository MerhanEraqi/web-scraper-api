package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient is a global MongoDB client instance.
var MongoClient *mongo.Client

// ConnectMongoDB connects to MongoDB, creates a database and collection if they don't exist.
func ConnectMongoDB() error {
	// Load MongoDB environment variables.
	mongoURI := os.Getenv("MONGO_URI")
	mongoDB := os.Getenv("MONGO_DB")
	collectionName := os.Getenv("MONGO_COLLECTION")

	// Set a timeout context for MongoDB operations.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the MongoDB server to ensure the connection works.
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB successfully.")
	MongoClient = client

	// Access the database.
	database := MongoClient.Database(mongoDB)

	// Check and create the collection if it doesn't exist.
	log.Printf("Checking existence of collection '%s'...", collectionName)
	filter := bson.D{{Key: "name", Value: collectionName}}
	collections, err := database.ListCollectionNames(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	if len(collections) == 0 {
		if err := database.CreateCollection(ctx, collectionName); err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
		log.Printf("Collection '%s' created successfully.", collectionName)
	} else {
		log.Printf("Collection '%s' already exists.", collectionName)
	}

	return nil
}
