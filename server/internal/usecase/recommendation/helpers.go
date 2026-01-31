package recommendation

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"server/internal/domain/entity"
)

// Константы для рекомендаций
const (
	DefaultTextMode  = "base"
	DefaultBlockText = "Новый текстовый блок — добавьте конкретное действие или мысль поддержки."
)

var pageNumberRegex = regexp.MustCompile(`\d+`)

// NormalizeRecommendationType приводит тип раздела к корректному виду
func NormalizeRecommendationType(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "Страница 1"
	}
	return trimmed
}

// SanitizeTextMode нормализует режим отображения текста
func SanitizeTextMode(mode string) entity.TextMode {
	switch strings.TrimSpace(mode) {
	case "bold":
		return entity.TextMode("bold")
	case "line":
		return entity.TextMode("line")
	case "bold-italics-line":
		return entity.TextMode("bold-italics-line")
	default:
		return entity.TextMode(DefaultTextMode)
	}
}

// ExtractPageNumber извлекает номер страницы из строки типа "Страница N"
func ExtractPageNumber(value string) int {
	match := pageNumberRegex.FindString(value)
	if match == "" {
		return 0
	}
	var number int
	fmt.Sscanf(match, "%d", &number)
	return number
}

// SortTypes сортирует типы разделов по номеру страницы
func SortTypes(types []string) {
	sort.Slice(types, func(i, j int) bool {
		left := ExtractPageNumber(types[i])
		right := ExtractPageNumber(types[j])
		if left == right {
			return types[i] < types[j]
		}
		return left < right
	})
}

// SortRecommendations сортирует рекомендации по номеру страницы и тексту
func SortRecommendations(recs []entity.Recommendation) {
	sort.Slice(recs, func(i, j int) bool {
		left := ExtractPageNumber(recs[i].RecommendationType)
		right := ExtractPageNumber(recs[j].RecommendationType)
		if left == right {
			return recs[i].RecommendationText < recs[j].RecommendationText
		}
		return left < right
	})
}
