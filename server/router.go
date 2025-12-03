package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	dashboardPage "server/handlers/dashboardPage"
	loginPage "server/handlers/loginPage"
	reviewsPage "server/handlers/reviewsPage"
	registrationPage "server/handlers/registrationPage"
)

// setupRouter добавляет маршруты, прокидывая подключение к БД в каждый обработчик.
func setupRouter(db *mongo.Database) *gin.Engine {
	router := gin.Default()

	// CORS (для разработки с React)
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	api := router.Group("/api")
	
	// /api/login/... (Обработчики событий авторизации)
	login := api.Group("/login")
	{
		login.POST("/google", func(c *gin.Context) {
			loginPage.LoginWithGoogleHandler(db, c)
		})
		login.POST("/yandex", func(c *gin.Context) {
			loginPage.LoginWithYandexHandler(db, c)
		})
		login.POST("/password", func(c *gin.Context) {
			loginPage.LoginWithPasswordHandler(db, c)
		})
		login.POST("/lostPassword", func(c *gin.Context) {
			loginPage.LostPasswordHandler(db, c)
		})
	}

	// /api/createAccount (Обработчик регистрации аккаунта)
	api.POST("/createAccount", func(c *gin.Context) {
		registrationPage.RegistrationPageHandler(db, c)
	})

	dashboard := api.Group("/dashboard")
	{
		dashboard.POST("/completed-tests", func(c *gin.Context) {
			dashboardPage.GetCompletedTestsHandler(db, c)
		})
		dashboard.POST("/users", func(c *gin.Context) {
			dashboardPage.GetUsersDataHandler(db, c)
		})
		dashboard.POST("/user-answers", func(c *gin.Context) {
			dashboardPage.GetUserAnswersHandler(db, c)
		})
		dashboard.POST("/block-user", func(c *gin.Context) {
			dashboardPage.BlockUserHandler(db, c)
		})
		dashboard.POST("/delete-user", func(c *gin.Context) {
			dashboardPage.DeleteUserHandler(db, c)
		})
		dashboard.POST("/delete-account", func(c *gin.Context) {
			dashboardPage.DeleteAccountHandler(db, c)
		})
		dashboard.POST("/change-user-data", func(c *gin.Context) {
			dashboardPage.ChangeUserDataHandler(db, c)
		})
	}


	// /api/account/...
	// account := api.Group("/account")
	// {
	// 	account.POST("/createAccount", accountPage.CreateAccountHandler)
	// }

	// Группа маршрутов для отзывов: /api/reviews/...
	reviews := api.Group("/reviews")
	{
		// /api/reviews/getReviews
		reviews.GET("/getReviews", func(c *gin.Context) {
			reviewsPage.GetReviewsHandler(db, c)
		})
		// /api/reviews/createReview
		reviews.POST("/createReview", func(c *gin.Context) {
			reviewsPage.CreateReviewHandler(db, c)
		})
		// /api/reviews/updateReview
		reviews.POST("/updateReview", func(c *gin.Context) {
			reviewsPage.UpdateReviewHandler(db, c)
		})
		// /api/reviews/deleteReview
		reviews.POST("/deleteReview", func(c *gin.Context) {
			reviewsPage.DeleteReviewHandler(db, c)
		})
		// /api/reviews/approveOrDeny
		reviews.POST("/approveOrDeny", func(c *gin.Context) {
			reviewsPage.ApproveOrDenyReviewHandler(db, c)
		})
	}

	return router
}
