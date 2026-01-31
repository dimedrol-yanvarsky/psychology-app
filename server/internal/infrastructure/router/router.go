package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	httpController "server/internal/adapter/controller/http"
)

type Controllers struct {
	Auth           *httpController.AuthController
	Test           *httpController.TestController
	Review         *httpController.ReviewController
	Recommendation *httpController.RecommendationController
	Dashboard      *httpController.DashboardController
}

func NewRouter(controllers Controllers) *gin.Engine {
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	api := router.Group("/api")

	// Auth routes
	login := api.Group("/login")
	{
		login.POST("/password", controllers.Auth.LoginWithPassword)
		login.POST("/google", controllers.Auth.LoginWithGoogle)
		login.POST("/yandex", controllers.Auth.LoginWithYandex)
		login.POST("/lostPassword", controllers.Auth.LostPassword)
	}
	api.POST("/createAccount", controllers.Auth.Register)

	// Tests routes
	tests := api.Group("/tests")
	{
		tests.POST("/getTests", controllers.Test.GetTests)
		tests.POST("/getQuestions", controllers.Test.GetQuestions)
		tests.POST("/attemptTest", controllers.Test.AttemptTest)
		tests.POST("/deleteTest", controllers.Test.DeleteTest)
		tests.POST("/changeTest", controllers.Test.ChangeTest)
		tests.POST("/addTest", controllers.Test.AddTest)
	}

	// Reviews routes
	reviews := api.Group("/reviews")
	{
		reviews.GET("/getReviews", controllers.Review.GetReviews)
		reviews.POST("/createReview", controllers.Review.CreateReview)
		reviews.POST("/updateReview", controllers.Review.UpdateReview)
		reviews.POST("/deleteReview", controllers.Review.DeleteReview)
		reviews.POST("/approveOrDeny", controllers.Review.ApproveOrDeny)
	}

	// Recommendations routes
	recommendations := api.Group("/recommendations")
	{
		recommendations.GET("/list", controllers.Recommendation.List)
		recommendations.POST("/addBlock", controllers.Recommendation.AddBlock)
		recommendations.POST("/updateBlock", controllers.Recommendation.UpdateBlock)
		recommendations.POST("/deleteBlock", controllers.Recommendation.DeleteBlock)
		recommendations.POST("/addSection", controllers.Recommendation.AddSection)
		recommendations.POST("/deleteSection", controllers.Recommendation.DeleteSection)
	}

	// Dashboard routes
	dashboard := api.Group("/dashboard")
	{
		dashboard.POST("/completed-tests", controllers.Dashboard.GetCompletedTests)
		dashboard.POST("/users", controllers.Dashboard.GetUsersData)
		dashboard.POST("/user-answers", controllers.Dashboard.GetUserAnswers)
		dashboard.POST("/block-user", controllers.Dashboard.BlockUser)
		dashboard.POST("/delete-user", controllers.Dashboard.DeleteUser)
		dashboard.POST("/delete-account", controllers.Dashboard.DeleteAccount)
		dashboard.POST("/change-user-data", controllers.Dashboard.ChangeUserData)
		dashboard.POST("/terminal", controllers.Dashboard.TerminalCommands)
	}

	return router
}
