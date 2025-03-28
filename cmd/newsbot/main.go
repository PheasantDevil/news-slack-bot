package main

import (
	"fmt"
	"log"

	"newsbot/internal/scraper"
	"newsbot/internal/slack"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	fmt.Println("Fetching news articles...")

	articles, err := scraper.FetchNews("https://www.drone.jp/")
	if err != nil {
		log.Fatalf("Failed to fetch news: %v", err)
	}

	// 🔥 ここで取得した記事のデバッグ出力を追加
	if len(articles) == 0 {
		log.Println("No articles found!")
	} else {
		for _, article := range articles {
			fmt.Println("Title:", article.Title)
			fmt.Println("URL:", article.URL)
			fmt.Println("Summary:", article.Summary)
			fmt.Println("Thumbnail:", article.Thumbnail)
		}
	}

	// Slack へ投稿
	for _, article := range articles {
		message := fmt.Sprintf("*%s*\n%s\n%s\n%s", article.Title, article.URL, article.Summary, article.Thumbnail)
		slack.SendToSlack(message)
	}
	fmt.Println("News posted to Slack successfully!")
}
