package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	dashboardPage "server/handlers/dashboardPage"
	loginPage "server/handlers/loginPage"
	reviewsPage "server/handlers/reviewsPage"
	registrationPage "server/handlers/registrationPage"
	recommendationsPage "server/handlers/recommendationsPage"
	testsPage "server/handlers/testsPage"
	"server/internal/dashboard"
	"server/internal/recommendations"
	"server/internal/reviews"
	"server/internal/tests"
	"server/internal/user"
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

	recRepo := recommendations.NewMongoRepository(db)
	recService := recommendations.NewService(recRepo)
	recHandlers := recommendationsPage.NewHandlers(recService)

	userRepo := user.NewMongoRepository(db)
	userService := user.NewService(userRepo)
	authHandlers := loginPage.NewHandlers(userService)
	regHandlers := registrationPage.NewHandlers(userService)

	testsRepo := tests.NewMongoRepository(db)
	testsService := tests.NewService(testsRepo)
	testsHandlers := testsPage.NewHandlers(testsService)

	reviewsRepo := reviews.NewMongoRepository(db)
	reviewsService := reviews.NewService(reviewsRepo)
	reviewsHandlers := reviewsPage.NewHandlers(reviewsService)

	dashboardRepo := dashboard.NewMongoRepository(db)
	dashboardService := dashboard.NewService(dashboardRepo)
	dashboardHandlers := dashboardPage.NewHandlers(dashboardService)

	api := router.Group("/api")
	
	// /api/login/... (Обработчики событий авторизации)
	login := api.Group("/login")
	{
		login.POST("/google", func(c *gin.Context) {
			authHandlers.LoginWithGoogleHandler(c)
		})
		login.POST("/yandex", func(c *gin.Context) {
			authHandlers.LoginWithYandexHandler(c)
		})
		login.POST("/password", func(c *gin.Context) {
			authHandlers.LoginWithPasswordHandler(c)
		})
		login.POST("/lostPassword", func(c *gin.Context) {
			authHandlers.LostPasswordHandler(c)
		})
	}

	// /api/createAccount (Обработчик регистрации аккаунта)
	api.POST("/createAccount", func(c *gin.Context) {
		regHandlers.RegistrationPageHandler(c)
	})

	dashboard := api.Group("/dashboard")
	{
		dashboard.POST("/completed-tests", func(c *gin.Context) {
			dashboardHandlers.GetCompletedTestsHandler(c)
		})
		dashboard.POST("/users", func(c *gin.Context) {
			dashboardHandlers.GetUsersDataHandler(c)
		})
		dashboard.POST("/user-answers", func(c *gin.Context) {
			dashboardHandlers.GetUserAnswersHandler(c)
		})
		dashboard.POST("/block-user", func(c *gin.Context) {
			dashboardHandlers.BlockUserHandler(c)
		})
		dashboard.POST("/delete-user", func(c *gin.Context) {
			dashboardHandlers.DeleteUserHandler(c)
		})
		dashboard.POST("/delete-account", func(c *gin.Context) {
			dashboardHandlers.DeleteAccountHandler(c)
		})
		dashboard.POST("/change-user-data", func(c *gin.Context) {
			dashboardHandlers.ChangeUserDataHandler(c)
		})
		dashboard.POST("/terminal", func(c *gin.Context) {
			dashboardHandlers.TerminalCommandsHandler(c)
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
			reviewsHandlers.GetReviewsHandler(c)
		})
		// /api/reviews/createReview
		reviews.POST("/createReview", func(c *gin.Context) {
			reviewsHandlers.CreateReviewHandler(c)
		})
		// /api/reviews/updateReview
		reviews.POST("/updateReview", func(c *gin.Context) {
			reviewsHandlers.UpdateReviewHandler(c)
		})
		// /api/reviews/deleteReview
		reviews.POST("/deleteReview", func(c *gin.Context) {
			reviewsHandlers.DeleteReviewHandler(c)
		})
		// /api/reviews/approveOrDeny
		reviews.POST("/approveOrDeny", func(c *gin.Context) {
			reviewsHandlers.ApproveOrDenyReviewHandler(c)
		})
	}

	tests := api.Group("/tests")
	{
		tests.POST("/getTests", func(c *gin.Context) {
			testsHandlers.GetTestsHandler(c)
		})
		tests.POST("/getQuestions", func(c *gin.Context) {
			testsHandlers.GetQuestionsHandler(c)
		})
		tests.POST("/attemptTest", func(c *gin.Context) {
			testsHandlers.AttemptTestHandler(c)
		})
		tests.POST("/deleteTest", func(c *gin.Context) {
			testsHandlers.DeleteTestHandler(c)
		})
		tests.POST("/changeTest", func(c *gin.Context) {
			testsHandlers.ChangeTestHandler(c)
		})
		tests.POST("/addTest", func(c *gin.Context) {
			testsHandlers.AddTestHandler(c)
		})
	}

	recommendations := api.Group("/recommendations")
	{
		recommendations.GET("/list", func(c *gin.Context) {
			recHandlers.GetRecommendationsHandler(c)
		})
		recommendations.POST("/addBlock", func(c *gin.Context) {
			recHandlers.AddBlockHandler(c)
		})
		recommendations.POST("/updateBlock", func(c *gin.Context) {
			recHandlers.UpdateBlockHandler(c)
		})
		recommendations.POST("/deleteBlock", func(c *gin.Context) {
			recHandlers.DeleteBlockHandler(c)
		})
		recommendations.POST("/addSection", func(c *gin.Context) {
			recHandlers.AddSectionHandler(c)
		})
		recommendations.POST("/deleteSection", func(c *gin.Context) {
			recHandlers.DeleteSectionHandler(c)
		})
	}

	return router
}
