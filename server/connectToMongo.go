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
func connectToMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Не удалось подключиться к MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB не отвечает:", err)
	}

	db := client.Database("psychologyApp")

	// Инициализируем пакет с обработчиками отзывов коллекцией MongoDB
	reviewsPage.InitReviews(db.Collection("reviews"))
}
