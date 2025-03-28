package scraper

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Article struct {
	Title     string
	URL       string
	Summary   string
	Thumbnail string
}

// FetchNews スクレイピング処理
func FetchNews(url string) ([]Article, error) {
	// 🔥 ステップ1: HTTPリクエストの送信
	fmt.Println("Fetching URL:", url)

	// User-Agent を設定してリクエスト
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	// 🔥 ステップ2: HTTP ステータスコードを確認
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Final URL after request:", resp.Request.URL.String())

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch page, status code: %d", resp.StatusCode)
	}

	// 🔥 ステップ3: HTML を goquery に渡す前に出力
	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// 取得したHTMLを表示
	htmlText := string(htmlBytes)
	fmt.Println("Fetched HTML content (first 500 chars):", htmlText[:min(500, len(htmlText))])

	// goquery に渡すために新しい Reader を作成
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}
	fmt.Println("HTML parsed successfully!")

	// 🔥 ステップ4: 記事を取得
	articles := []Article{}
	doc.Find(".p-wrap").Each(func(i int, s *goquery.Selection) {
		// デバッグ: 各記事のHTML構造を出力
		html, _ := s.Html()
		fmt.Printf("Article HTML structure:\n%s\n", html)

		title := s.Find("h3.entry-title a").Text()
		url, exists := s.Find("h3.entry-title a").Attr("href")
		if !exists {
			fmt.Printf("Warning: URL not found for article %d\n", i)
		}
		summary := s.Find(".entry-summary").Text()
		thumbnail, _ := s.Find("img.featured-img").Attr("src")

		fmt.Printf("Found article %d:\nTitle: %s\nURL: %s\nSummary: %s\nThumbnail: %s\n\n",
			i, title, url, summary, thumbnail)

		if title != "" {
			articles = append(articles, Article{
				Title:     strings.TrimSpace(title),
				URL:       url,
				Summary:   strings.TrimSpace(summary),
				Thumbnail: thumbnail,
			})
		}
	})

	fmt.Printf("Total articles found: %d\n", len(articles))
	if len(articles) == 0 {
		// デバッグ: 全体のHTML構造を確認
		fmt.Println("\nFull HTML structure:")
		doc.Find("body").Each(func(i int, s *goquery.Selection) {
			html, _ := s.Html()
			fmt.Printf("%s\n", html)
		})
	}

	return articles, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
