package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArticle_FormatMessage(t *testing.T) {
	// テストケース
	tests := []struct {
		name     string
		article  Article
		expected string
	}{
		{
			name: "通常の記事",
			article: Article{
				Title:     "テスト記事",
				URL:       "https://example.com/article",
				Summary:   "テスト記事のサマリー",
				Thumbnail: "https://example.com/image.jpg",
				PostDate:  time.Date(2024, 3, 30, 0, 0, 0, 0, time.UTC),
			},
			expected: "*テスト記事*\n\nテスト記事のサマリー\n\n<https://example.com/article|記事を読む>",
		},
		{
			name: "サマリーなしの記事",
			article: Article{
				Title:     "テスト記事",
				URL:       "https://example.com/article",
				Summary:   "",
				Thumbnail: "https://example.com/image.jpg",
				PostDate:  time.Date(2024, 3, 30, 0, 0, 0, 0, time.UTC),
			},
			expected: "*テスト記事*\n\n<https://example.com/article|記事を読む>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := tt.article.FormatMessage()
			assert.Equal(t, tt.expected, message)
		})
	}
}
