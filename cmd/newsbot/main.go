package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"newsbot/internal/scraper"
	"newsbot/internal/slack"
)

func main() {
	// Slack Webhook URLを取得
	webhookURL := os.Getenv("SLACK_WEBHOOK_NEWS_DRONE_JP_CH_URL")
	if webhookURL == "" {
		log.Fatal("Error: SLACK_WEBHOOK_NEWS_DRONE_JP_CH_URL is not set")
	}

	// 現在の日付を取得
	now := time.Now()
	fmt.Printf("Fetching news articles for %s...\n", now.Format("2006-01-02"))

	// 記事を取得（当日の記事のみ）
	articles, err := scraper.FetchNews("https://www.drone.jp/", now)
	if err != nil {
		log.Fatalf("Failed to fetch news: %v", err)
	}

	// 記事数をログ出力
	if len(articles) == 0 {
		log.Println("No articles found for today!")
		return
	}
	log.Printf("Found %d articles for today\n", len(articles))

	// Slackに投稿
	if err := slack.PostArticlesToSlack(webhookURL, articles); err != nil {
		log.Fatalf("Failed to post to Slack: %v", err)
	}

	fmt.Println("News posted to Slack successfully!")
}
