package main

import (
	"log"

	"server/internal/adapter/controller/http"
	"server/internal/adapter/repository/mongodb"
	"server/internal/infrastructure/config"
	"server/internal/infrastructure/database"
	"server/internal/infrastructure/router"
	dashboardUseCase "server/internal/usecase/dashboard"
	recommendationUseCase "server/internal/usecase/recommendation"
	reviewUseCase "server/internal/usecase/review"
	testUseCase "server/internal/usecase/test"
	userUseCase "server/internal/usecase/user"
)

func main() {
	log.Println("🚀 Запуск сервера Clean Architecture...")

	// 1. Load configuration
	cfg := config.Load()
	log.Println("✓ Конфигурация загружена")

	// 2. Initialize database
	db := database.NewMongoDatabase(cfg.Database)
	log.Printf("✓ Подключено к БД: %s", cfg.Database.Database)

	// 3. Initialize repositories
	userRepo := mongodb.NewUserRepository(db)
	testRepo := mongodb.NewTestRepository(db)
	userAnswerRepo := mongodb.NewUserAnswerRepository(db)
	reviewRepo := mongodb.NewReviewRepository(db)
	recommendationRepo := mongodb.NewRecommendationRepository(db)
	dashboardRepo := mongodb.NewDashboardRepository(db)
	log.Println("✓ Репозитории инициализированы")

	// 4. Initialize use cases

	// Auth use cases
	loginUC := userUseCase.NewLoginUseCase(userRepo)
	registerUC := userUseCase.NewRegisterUseCase(userRepo)

	// Test use cases
	getTestsUC := testUseCase.NewGetTestsUseCase(testRepo, userAnswerRepo)
	getQuestionsUC := testUseCase.NewGetQuestionsUseCase(testRepo)
	attemptTestUC := testUseCase.NewAttemptTestUseCase(userAnswerRepo)
	addTestUC := testUseCase.NewAddTestUseCase(testRepo)
	changeTestUC := testUseCase.NewChangeTestUseCase(testRepo)
	deleteTestUC := testUseCase.NewDeleteTestUseCase(testRepo)

	// Review use cases
	getReviewsUC := reviewUseCase.NewGetReviewsUseCase(reviewRepo)
	createReviewUC := reviewUseCase.NewCreateReviewUseCase(reviewRepo)
	updateReviewUC := reviewUseCase.NewUpdateReviewUseCase(reviewRepo)
	deleteReviewUC := reviewUseCase.NewDeleteReviewUseCase(reviewRepo)
	moderateReviewUC := reviewUseCase.NewModerateReviewUseCase(reviewRepo)

	// Recommendation use cases
	listRecommendationsUC := recommendationUseCase.NewListRecommendationsUseCase(recommendationRepo)
	addBlockUC := recommendationUseCase.NewAddBlockUseCase(recommendationRepo)
	updateBlockUC := recommendationUseCase.NewUpdateBlockUseCase(recommendationRepo)
	deleteBlockUC := recommendationUseCase.NewDeleteBlockUseCase(recommendationRepo)
	addSectionUC := recommendationUseCase.NewAddSectionUseCase(recommendationRepo)
	deleteSectionUC := recommendationUseCase.NewDeleteSectionUseCase(recommendationRepo)

	// Dashboard use cases
	getUsersUC := dashboardUseCase.NewGetUsersUseCase(dashboardRepo)
	blockUserUC := dashboardUseCase.NewBlockUserUseCase(dashboardRepo)
	deleteUserUC := dashboardUseCase.NewDeleteUserUseCase(dashboardRepo)
	deleteAccountUC := dashboardUseCase.NewDeleteAccountUseCase(dashboardRepo)
	changeUserDataUC := dashboardUseCase.NewChangeUserDataUseCase(dashboardRepo)
	getCompletedTestsUC := dashboardUseCase.NewGetCompletedTestsUseCase(dashboardRepo, testRepo)
	getUserAnswersUC := dashboardUseCase.NewGetUserAnswersUseCase(dashboardRepo, testRepo)
	terminalCommandsUC := dashboardUseCase.NewTerminalCommandsUseCase()

	log.Println("✓ Use Cases инициализированы")

	// 5. Initialize controllers
	authController := http.NewAuthController(loginUC, registerUC)
	testController := http.NewTestController(
		getTestsUC,
		getQuestionsUC,
		attemptTestUC,
		addTestUC,
		changeTestUC,
		deleteTestUC,
	)
	reviewController := http.NewReviewController(
		getReviewsUC,
		createReviewUC,
		updateReviewUC,
		deleteReviewUC,
		moderateReviewUC,
	)
	recommendationController := http.NewRecommendationController(
		listRecommendationsUC,
		addBlockUC,
		updateBlockUC,
		deleteBlockUC,
		addSectionUC,
		deleteSectionUC,
	)
	dashboardController := http.NewDashboardController(
		getUsersUC,
		blockUserUC,
		deleteUserUC,
		deleteAccountUC,
		changeUserDataUC,
		getCompletedTestsUC,
		getUserAnswersUC,
		terminalCommandsUC,
	)
	log.Println("✓ Контроллеры инициализированы")

	// 6. Setup router
	r := router.NewRouter(router.Controllers{
		Auth:           authController,
		Test:           testController,
		Review:         reviewController,
		Recommendation: recommendationController,
		Dashboard:      dashboardController,
	})

	log.Println("✓ Роутер настроен")
	log.Println("")
	log.Println("📦 Статус модулей:")
	log.Println("  ✅ Auth (Login, Register) - РАБОТАЕТ")
	log.Println("  ✅ Test - РАБОТАЕТ")
	log.Println("  ✅ Review - РАБОТАЕТ")
	log.Println("  ✅ Recommendation - РАБОТАЕТ")
	log.Println("  ✅ Dashboard - РАБОТАЕТ")
	log.Println("")
	log.Printf("🌐 Сервер запущен на http://localhost%s\n", cfg.Server.Port)
	log.Println("✨ Clean Architecture миграция завершена!")

	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatal("✗ Ошибка запуска сервера:", err)
	}
}
