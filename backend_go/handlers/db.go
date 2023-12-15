// handlers/db.go
package handlers

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
)

// ConnectDB connects to MongoDB
func ConnectDB() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	cl, err := mongo.Connect(context.TODO(), clientOptions)
	client = cl
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}

// CloseDB disconnects from MongoDB
func CloseDB() {
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
