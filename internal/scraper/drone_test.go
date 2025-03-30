package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchDroneNews(t *testing.T) {
	// テスト用のHTMLを準備
	html := `
		<html>
			<body>
				<article>
					<h2><a href="https://example.com/article1">テスト記事1</a></h2>
					<div class="summary">テスト記事1のサマリー</div>
					<img src="https://example.com/image1.jpg">
					<div class="date">2024-03-30</div>
				</article>
				<article>
					<h2><a href="https://example.com/article2">テスト記事2</a></h2>
					<div class="summary">テスト記事2のサマリー</div>
					<img src="https://example.com/image2.jpg">
					<div class="date">2024-03-29</div>
				</article>
			</body>
		</html>
	`

	// モックサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
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
			name:       "本日の記事を取得",
			targetDate: time.Date(2024, 3, 30, 0, 0, 0, 0, time.UTC),
			wantCount:  1,
			wantErr:    false,
		},
		{
			name:       "昨日の記事を取得",
			targetDate: time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
			wantCount:  1,
			wantErr:    false,
		},
		{
			name:       "記事なし",
			targetDate: time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC),
			wantCount:  0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			articles, err := FetchDroneNews(server.URL, tt.targetDate)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, articles, tt.wantCount)

			if tt.wantCount > 0 {
				article := articles[0]
				assert.NotEmpty(t, article.Title)
				assert.NotEmpty(t, article.URL)
				assert.NotEmpty(t, article.Summary)
				assert.NotEmpty(t, article.Thumbnail)
				assert.Equal(t, tt.targetDate.Format("2006-01-02"), article.PostDate.Format("2006-01-02"))
			}
		})
	}
}
