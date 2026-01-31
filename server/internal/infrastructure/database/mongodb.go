package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"server/internal/infrastructure/config"
)

func NewMongoDatabase(cfg config.DatabaseConfig) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		log.Fatal("Не удалось подключиться к MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB не отвечает:", err)
	}

	log.Println("✓ Подключение к БД установлено")
	return client.Database(cfg.Database)
}
