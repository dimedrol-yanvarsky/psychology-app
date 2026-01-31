package entity

// RecommendationID представляет уникальный идентификатор рекомендации
type RecommendationID string

func (id RecommendationID) String() string { return string(id) }
func (id RecommendationID) IsEmpty() bool  { return id == "" }

// TextMode описывает режим текста рекомендации
type TextMode string

const (
	TextModeBase   TextMode = "base"
	TextModeNormal TextMode = "обычный режим"
)

// Recommendation - доменная сущность рекомендации
type Recommendation struct {
	ID                 RecommendationID
	RecommendationText string
	TextMode           TextMode
	RecommendationType string // Страница 1, Страница 2, и т.д.
}

// IsBaseMode проверяет, использует ли рекомендация базовый режим текста
func (r *Recommendation) IsBaseMode() bool {
	return r.TextMode == TextModeBase
}
