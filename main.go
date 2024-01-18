package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CS559-CSD-IITBH/order-service/routes"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Internal server error: Unable to load the env file")
	}

	// Replace this with your MongoDB Atlas connection string
	connectionString := os.Getenv("MONGO_URL")

	// Set MongoDB connection options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Internal server error: Unable to connect to Mongo")
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Internal server error: Unable to talk to Mongo")
	}

	fmt.Println("Connected to MongoDB!")

	// You can now use the "client" variable to interact with your MongoDB database.
	// For example, you can access a collection:
	orderCollection := client.Database(os.Getenv("MONGO_DB_NAME")).Collection(os.Getenv("MONGO_COLLECTION_ORDER"))
	cartCollection := client.Database(os.Getenv("MONGO_DB_NAME")).Collection(os.Getenv("MONGO_COLLECTION_CART"))

	// Session store in  NewFilesystemStore
	store := sessions.NewFilesystemStore("sessions/", []byte("secret-key"))

	// Set max age for cookie
	store.Options = &sessions.Options{
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	if err != nil {
		log.Fatalln("Internal server error: Unable to connect to the DB")
	}

	r := routes.SetupRouter(orderCollection, cartCollection, store)
	r.Run(":" + os.Getenv("PORT"))
}
