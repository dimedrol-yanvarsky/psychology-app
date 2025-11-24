package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	loginPage "server/handlers/loginPage"
	reviewsPage "server/handlers/reviewsPage" // замени server на имя модуля
)

// setupRouter настраивает только маршруты, не работая с MongoDB/коллекциями.
func setupRouter() *gin.Engine {
	router := gin.Default()

	// CORS (для разработки с React)
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	api := router.Group("/api")

	// Группа маршрутов для отзывов: /api/reviews/...
	reviews := api.Group("/reviews")
	{
		// /api/reviews/getReviews
		reviews.GET("/getReviews", reviewsPage.GetReviewsHandler)

		// /api/reviews/createReview
		reviews.POST("/createReview", reviewsPage.CreateReviewHandler)
	}

	// /api/login/...
	login := api.Group("/login")
	{
		login.POST("/google", loginPage.LoginWithGoogleHandler)
		login.POST("/yandex", loginPage.LoginWithYandexHandler)
		login.POST("/password", loginPage.LoginWithPasswordHandler)
	}

	// Здесь же можешь добавлять другие группы:
	// account := api.Group("/account")
	// account.POST("/createAccount", accountPage.CreateAccountHandler)
	// ...

	return router
}
