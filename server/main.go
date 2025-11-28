package main

import (
	"log"
)

func main() {

	// Подключение к MongoDB
	connectToDatabase()
	// Настраиваем роутер (без знания о MongoDB) в файле router.go
	router := setupRouter()

	//Запуск сервера
	log.Println("Сервер запущен на http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}

}
