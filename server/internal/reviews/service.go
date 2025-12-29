package reviews

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidID    = errors.New("invalid id")
	ErrNotFound     = errors.New("not found")
	ErrForbidden    = errors.New("forbidden")
	ErrDatabase     = errors.New("database error")
)

// Service реализует бизнес-логику управления отзывами.
type Service struct {
	repo Repository
}

// NewService создает сервис отзывов.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ReviewItem объединяет отзыв и имя автора.
type ReviewItem struct {
	Review Review
	AuthorName string
}

// GetReviews возвращает все отзывы, кроме удаленных, с именем автора.
func (s *Service) GetReviews() ([]ReviewItem, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	reviews, err := s.repo.FindReviews(ctx, bson.M{"status": bson.M{"$ne": StatusDeleted}})
	if err != nil {
		return nil, ErrDatabase
	}

	userIDs := make([]primitive.ObjectID, 0)
	seen := make(map[primitive.ObjectID]struct{})
	for _, review := range reviews {
		if review.UserID.IsZero() {
			continue
		}
		if _, ok := seen[review.UserID]; !ok {
			seen[review.UserID] = struct{}{}
			userIDs = append(userIDs, review.UserID)
		}
	}

	userNames := make(map[primitive.ObjectID]string)
	if len(userIDs) > 0 {
		users, err := s.repo.FindUsersByIDs(ctx, userIDs)
		if err != nil {
			return nil, ErrDatabase
		}
		for _, user := range users {
			userNames[user.ID] = strings.TrimSpace(user.FirstName)
		}
	}

	result := make([]ReviewItem, 0, len(reviews))
	for _, review := range reviews {
		name := strings.TrimSpace(userNames[review.UserID])
		if name == "" {
			name = "Неизвестный автор"
		}
		result = append(result, ReviewItem{
			Review:     review,
			AuthorName: name,
		})
	}

	return result, nil
}

// CreateReviewInput описывает данные для создания отзыва.
type CreateReviewInput struct {
	UserID string
	ReviewBody string
}

// CreateReview создает новый отзыв со статусом модерации.
func (s *Service) CreateReview(input CreateReviewInput) (ReviewItem, error) {
	userID := strings.TrimSpace(input.UserID)
	body := strings.TrimSpace(input.ReviewBody)
	if userID == "" || body == "" {
		return ReviewItem{}, ErrInvalidInput
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ReviewItem{}, ErrInvalidID
	}

	ctx, cancel := withTimeout()
	defer cancel()

	author, err := s.repo.FindUserByID(ctx, userObjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ReviewItem{}, ErrNotFound
		}
		return ReviewItem{}, ErrDatabase
	}

	existing, err := s.repo.FindReviewByUser(ctx, userObjectID)
	if err == nil && existing.ID != primitive.NilObjectID {
		return ReviewItem{}, ErrForbidden
	}
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return ReviewItem{}, ErrDatabase
	}

	review := Review{
		UserID:     userObjectID,
		ReviewBody: body,
		Date:       time.Now().Format("02.01.2006"),
		Status:     StatusModeration,
	}

	insertedID, err := s.repo.InsertReview(ctx, review)
	if err != nil {
		return ReviewItem{}, ErrDatabase
	}
	review.ID = insertedID

	return ReviewItem{
		Review:     review,
		AuthorName: strings.TrimSpace(author.FirstName),
	}, nil
}

// UpdateReviewInput описывает данные для обновления отзыва.
type UpdateReviewInput struct {
	ReviewID  string
	UserID    string
	ReviewBody string
}

// UpdateReview обновляет текст отзыва пользователя.
func (s *Service) UpdateReview(input UpdateReviewInput) (ReviewItem, error) {
	reviewID := strings.TrimSpace(input.ReviewID)
	userID := strings.TrimSpace(input.UserID)
	body := strings.TrimSpace(input.ReviewBody)
	if reviewID == "" || userID == "" || body == "" {
		return ReviewItem{}, ErrInvalidInput
	}

	reviewObjectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return ReviewItem{}, ErrInvalidID
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ReviewItem{}, ErrInvalidID
	}

	ctx, cancel := withTimeout()
	defer cancel()

	existing, err := s.repo.FindReviewByID(ctx, reviewObjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ReviewItem{}, ErrNotFound
		}
		return ReviewItem{}, ErrDatabase
	}

	if existing.Status == StatusDeleted {
		return ReviewItem{}, errors.New("Удаленный отзыв нельзя редактировать")
	}
	if existing.UserID != userObjectID {
		return ReviewItem{}, ErrForbidden
	}

	if err := s.repo.UpdateReviewBody(ctx, reviewObjectID, body); err != nil {
		return ReviewItem{}, ErrDatabase
	}

	author, _ := s.repo.FindUserByID(ctx, userObjectID)

	existing.ReviewBody = body
	return ReviewItem{
		Review:     existing,
		AuthorName: strings.TrimSpace(author.FirstName),
	}, nil
}

// DeleteReviewInput описывает параметры удаления отзыва.
type DeleteReviewInput struct {
	ReviewID string
	UserID   string
	IsAdmin  bool
}

// DeleteReview помечает отзыв как удаленный.
func (s *Service) DeleteReview(input DeleteReviewInput) error {
	reviewID := strings.TrimSpace(input.ReviewID)
	if reviewID == "" {
		return ErrInvalidInput
	}

	reviewObjectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return ErrInvalidID
	}

	var userObjectID primitive.ObjectID
	if !input.IsAdmin {
		userID := strings.TrimSpace(input.UserID)
		if userID == "" {
			return errors.New("Не указан пользователь")
		}
		userObjectID, err = primitive.ObjectIDFromHex(userID)
		if err != nil {
			return ErrInvalidID
		}
	}

	ctx, cancel := withTimeout()
	defer cancel()

	existing, err := s.repo.FindReviewByID(ctx, reviewObjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}
		return ErrDatabase
	}

	if !input.IsAdmin && existing.UserID != userObjectID {
		return ErrForbidden
	}

	if existing.Status == StatusDeleted {
		return nil
	}

	if err := s.repo.UpdateReviewStatus(ctx, reviewObjectID, StatusDeleted); err != nil {
		return ErrDatabase
	}
	return nil
}

// ApproveOrDenyInput описывает решение по отзыву.
type ApproveOrDenyInput struct {
	ReviewID string
	AdminID  string
	Decision string
}

// ApproveOrDeny меняет статус отзыва по решению администратора.
func (s *Service) ApproveOrDeny(input ApproveOrDenyInput) (ReviewItem, error) {
	reviewID := strings.TrimSpace(input.ReviewID)
	adminID := strings.TrimSpace(input.AdminID)
	decision := strings.TrimSpace(strings.ToLower(input.Decision))

	if reviewID == "" || adminID == "" {
		return ReviewItem{}, errors.New("Не указаны обязательные поля")
	}

	var newStatus string
	switch decision {
	case "approve":
		newStatus = StatusApproved
	case "deny":
		newStatus = StatusDenied
	default:
		return ReviewItem{}, errors.New("Некорректное действие")
	}

	reviewObjectID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return ReviewItem{}, ErrInvalidID
	}

	adminObjectID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return ReviewItem{}, ErrInvalidID
	}

	ctx, cancel := withTimeout()
	defer cancel()

	admin, err := s.repo.FindUserByID(ctx, adminObjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ReviewItem{}, ErrForbidden
		}
		return ReviewItem{}, ErrDatabase
	}

	if strings.TrimSpace(admin.Status) != "Администратор" {
		return ReviewItem{}, ErrForbidden
	}

	existing, err := s.repo.FindReviewByID(ctx, reviewObjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ReviewItem{}, ErrNotFound
		}
		return ReviewItem{}, ErrDatabase
	}

	if existing.Status == StatusDeleted {
		return ReviewItem{}, errors.New("Удаленный отзыв нельзя изменить")
	}

	if err := s.repo.UpdateReviewStatus(ctx, reviewObjectID, newStatus); err != nil {
		return ReviewItem{}, ErrDatabase
	}

	author, _ := s.repo.FindUserByID(ctx, existing.UserID)
	existing.Status = newStatus

	return ReviewItem{
		Review:     existing,
		AuthorName: strings.TrimSpace(author.FirstName),
	}, nil
}
