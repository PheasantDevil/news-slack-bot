package scraper

import (
	"fmt"
	"log"
	"net/http"
	"newsbot/internal/models"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// FetchDroneNews Droneのニュースを取得
func FetchDroneNews(url string, targetDate time.Time) ([]models.Article, error) {
	// HTTPリクエストを作成
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// User-Agentを設定
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// HTTPクライアントを作成してリクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスのステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// HTMLをパース
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var articles []models.Article

	// 記事一覧を取得
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		// タイトルとURLを取得
		titleElem := s.Find("h2 a")
		title := titleElem.Text()
		url, exists := titleElem.Attr("href")
		if !exists {
			log.Printf("Warning: URL not found for article '%s'", title)
			return
		}

		// サマリーを取得
		summary := s.Find(".summary").Text()

		// サムネイル画像を取得
		thumbnail := ""
		imgElem := s.Find("img")
		if imgElem.Length() > 0 {
			thumbnail, _ = imgElem.Attr("src")
		}

		// 投稿日を取得
		dateStr := s.Find(".date").Text()
		postDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("Warning: Failed to parse date '%s' for article '%s': %v", dateStr, title, err)
			return
		}

		// 本日の記事のみを追加
		if postDate.Equal(targetDate) {
			articles = append(articles, models.Article{
				Title:     title,
				URL:       url,
				Summary:   summary,
				Thumbnail: thumbnail,
				PostDate:  postDate,
			})
		}
	})

	return articles, nil
}
