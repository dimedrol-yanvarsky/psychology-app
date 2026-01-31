package entity

// ReviewID представляет уникальный идентификатор отзыва
type ReviewID string

func (id ReviewID) String() string { return string(id) }
func (id ReviewID) IsEmpty() bool  { return id == "" }

// ReviewStatus описывает статус отзыва
type ReviewStatus string

const (
	ReviewStatusDeleted    ReviewStatus = "Удален"
	ReviewStatusModeration ReviewStatus = "Модерируется"
	ReviewStatusApproved   ReviewStatus = "Добавлен"
	ReviewStatusDenied     ReviewStatus = "Отклонен"
)

// Review - доменная сущность отзыва
type Review struct {
	ID         ReviewID
	UserID     UserID
	ReviewBody string
	Date       string
	Status     ReviewStatus
}

// ReviewWithAuthor - отзыв с информацией об авторе
type ReviewWithAuthor struct {
	Review     Review
	AuthorName string
}

// IsDeleted проверяет, удален ли отзыв
func (r *Review) IsDeleted() bool {
	return r.Status == ReviewStatusDeleted
}

// IsApproved проверяет, одобрен ли отзыв
func (r *Review) IsApproved() bool {
	return r.Status == ReviewStatusApproved
}

// IsModeration проверяет, на модерации ли отзыв
func (r *Review) IsModeration() bool {
	return r.Status == ReviewStatusModeration
}
