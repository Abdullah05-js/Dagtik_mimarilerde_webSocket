package database

import (
	"ChatAPP/v2/model"
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *mongo.Database

func InitMongo() error {
	url := os.Getenv("MONGO_URL")
	if url == "" {
		log.Fatal("mongo url required")
	}
	client, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		return err
	}

	DB = client.Database("Cluster0")

	InitModels()

	log.Println("monogDB connected successfully")

	return nil
}

func InitModels() {
	// models initialization here
	model.Message = DB.Collection("Message")
}

func CloseMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := DB.Client().Disconnect(ctx); err != nil {
		return err
	}
	return nil
}
