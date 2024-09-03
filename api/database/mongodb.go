package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectMongoDB() {

	env := godotenv.Load()
	if env != nil {
		log.Fatalf("Erro loading .env: %v", env)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not defined")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error verifying connection to MongoDB:", err)
	}

	log.Println("Connected to MongoDB!")
}

func GetClient() *mongo.Client {
	return client
}

func GetCollection(databaseName, collectionName string) *mongo.Collection {
	return GetClient().Database(databaseName).Collection(collectionName)
}
