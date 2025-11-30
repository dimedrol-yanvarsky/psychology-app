package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Подключение к MongoDB
func connectToDatabase() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal("Не удалось подключиться к MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB не отвечает:", err)
	}

	log.Println("Подключился к БД")

	return client.Database("psychologyApp")
}
