package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// テストデータを構造体で定義
type economicTimesTestArticle struct {
	title     string
	url       string
	summary   string
	thumbnail string
	date      string
}

// テストケースのデータを生成
func generateEconomicTimesTestHTML(articles []economicTimesTestArticle) string {
	html := `<html><body>`
	for _, article := range articles {
		html += `
			<div class="news-list-item">
				<h3><a href="` + article.url + `">` + article.title + `</a></h3>
				<div class="news-list-item__summary">` + article.summary + `</div>
				<div class="news-list-item__image"><img src="` + article.thumbnail + `"></div>
				<div class="news-list-item__date">` + article.date + `</div>
			</div>
		`
	}
	html += `</body></html>`
	return html
}

func TestFetchEconomicTimesNews(t *testing.T) {
	// 現在の日付を基準にテストデータを生成
	now := time.Now()
	today := now.Format("2006.01.02")
	yesterday := now.AddDate(0, 0, -1).Format("2006.01.02")

	// テストデータを定義
	testArticles := []economicTimesTestArticle{
		{
			title:     "今日の記事1",
			url:       "https://economictimes.indiatimes.com/news/article1",
			summary:   "今日の記事1のサマリー",
			thumbnail: "https://example.com/image1.jpg",
			date:      today,
		},
		{
			title:     "今日の記事2",
			url:       "https://economictimes.indiatimes.com/news/article2",
			summary:   "今日の記事2のサマリー",
			thumbnail: "https://example.com/image2.jpg",
			date:      today,
		},
		{
			title:     "昨日の記事",
			url:       "https://economictimes.indiatimes.com/news/article3",
			summary:   "昨日の記事のサマリー",
			thumbnail: "https://example.com/image3.jpg",
			date:      yesterday,
		},
	}

	// テスト用のHTMLを生成
	html := generateEconomicTimesTestHTML(testArticles)

	// モックサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(html)); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// テストケース
	tests := []struct {
		name       string
		targetDate time.Time
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "今日の記事を取得",
			targetDate: now,
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:       "昨日の記事を取得",
			targetDate: now.AddDate(0, 0, -1),
			wantCount:  1,
			wantErr:    false,
		},
		{
			name:       "明日の記事を取得（記事なし）",
			targetDate: now.AddDate(0, 0, 1),
			wantCount:  0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			articles, err := FetchEconomicTimesNews(server.URL, tt.targetDate)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, articles, tt.wantCount)

			if tt.wantCount > 0 {
				// 期待される記事の日付を取得
				expectedDate := tt.targetDate.Format("2006.01.02")

				// 各記事の内容を検証
				for _, article := range articles {
					assert.NotEmpty(t, article.Title, "タイトルが空です")
					assert.NotEmpty(t, article.URL, "URLが空です")
					assert.NotEmpty(t, article.Summary, "サマリーが空です")
					assert.NotEmpty(t, article.Thumbnail, "サムネイルが空です")
					assert.Equal(t, expectedDate, article.PostDate.Format("2006.01.02"), "日付が一致しません")
				}
			}
		})
	}
}
