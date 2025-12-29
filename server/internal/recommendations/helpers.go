package recommendations

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var pageNumberRegex = regexp.MustCompile(`\d+`)

const dbTimeout = 10 * time.Second

func withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), dbTimeout)
}

// NormalizeRecommendationType приводит тип раздела к корректному виду.
func NormalizeRecommendationType(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "Страница 1"
	}
	return trimmed
}

// SanitizeTextMode нормализует режим отображения текста.
func SanitizeTextMode(mode string) string {
	switch strings.TrimSpace(mode) {
	case "bold":
		return "bold"
	case "line":
		return "line"
	case "bold-italics-line":
		return "bold-italics-line"
	default:
		return DefaultTextMode
	}
}

func extractPageNumber(value string) int {
	match := pageNumberRegex.FindString(value)
	if match == "" {
		return 0
	}
	var number int
	fmt.Sscanf(match, "%d", &number)
	return number
}

func sortTypes(types []string) {
	sort.Slice(types, func(i, j int) bool {
		left := extractPageNumber(types[i])
		right := extractPageNumber(types[j])
		if left == right {
			return types[i] < types[j]
		}
		return left < right
	})
}

// sortRecommendations сортирует рекомендации по номеру страницы и тексту внутри страницы.
func sortRecommendations(recs []Recommendation) {
	sort.Slice(recs, func(i, j int) bool {
		left := extractPageNumber(recs[i].RecommendationType)
		right := extractPageNumber(recs[j].RecommendationType)
		if left == right {
			return recs[i].RecommendationText < recs[j].RecommendationText
		}
		return left < right
	})
}
