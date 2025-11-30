package main

import (
	"log"
)

func main() {

	// Подключение к MongoDB
	db := connectToDatabase()
	// Настраиваем роутер, передавая подключение к БД во все обработчики
	router := setupRouter(db)

	//Запуск сервера
	log.Println("Сервер запущен на http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}

}
