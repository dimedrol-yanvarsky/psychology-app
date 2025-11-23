package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// глобальная коллекция, доступная из обработчиков
var reviewsCollection *mongo.Collection

func main() {
	// Подключение к MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Не удалось подключиться к MongoDB:", err)
	}

	db := client.Database("reviewsdb")
	reviewsCollection = db.Collection("reviews")

	router := gin.Default()

	// CORS (для разработки с React на http://localhost:3000)
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	api := router.Group("/api")
	{
		api.POST("/reviews", createReviewHandler) // маршрут получения данных от клиента
		api.GET("/reviews", listReviewsHandler)   // маршрут для получения списка отзывов
	}

	log.Println("Сервер запущен на http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
