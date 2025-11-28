package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	reviewsPage "server/handlers/reviewsPage"
)

// Подключение к MongoDB
func connectToDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/psychologyApp"))
	if err != nil {
		log.Fatal("Не удалось подключиться к MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB не отвечает:", err)
	}

	log.Println("Подключился к БД")

	db := client.Database("psychologyApp")
	
	initCollections(db)
}

// Инициализация обработчиков каждой коллекции
func initCollections(db *mongo.Database) {
	reviewsPage.InitReviews(db.Collection("reviews"))
}
