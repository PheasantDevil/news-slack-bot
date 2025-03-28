package scraper

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Article struct {
	Title     string
	URL       string
	Summary   string
	Thumbnail string
	PostDate  time.Time
}

// FetchNews スクレイピング処理
func FetchNews(url string, targetDate time.Time) ([]Article, error) {
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

		// 日付を取得してパース
		dateStr := s.Find("time.date.published").AttrOr("datetime", "")
		postDate, err := time.Parse("2006-01-02T15:04:05-07:00", dateStr)
		if err != nil {
			fmt.Printf("Warning: Failed to parse date for article %d: %v\n", i, err)
			return
		}

		// 指定された日付と同じ日の記事のみを追加
		if isSameDate(postDate, targetDate) {
			fmt.Printf("Found article %d:\nTitle: %s\nURL: %s\nSummary: %s\nThumbnail: %s\nDate: %s\n\n",
				i, title, url, summary, thumbnail, postDate.Format("2006-01-02"))

			if title != "" {
				articles = append(articles, Article{
					Title:     strings.TrimSpace(title),
					URL:       url,
					Summary:   strings.TrimSpace(summary),
					Thumbnail: thumbnail,
					PostDate:  postDate,
				})
			}
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

// isSameDate 2つの時刻が同じ日付かどうかを判定
func isSameDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
