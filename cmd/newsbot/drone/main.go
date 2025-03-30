package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"newsbot/internal/scraper"
	"newsbot/internal/slack"
)

// RunDroneBot Droneのニュースを取得してSlackに投稿
func RunDroneBot() error {
	// Slack Webhook URLを取得
	webhookURL := os.Getenv("SLACK_WEBHOOK_NEWS_DRONE_JP_CH_URL")
	if webhookURL == "" {
		return fmt.Errorf("SLACK_WEBHOOK_NEWS_DRONE_JP_CH_URL is not set")
	}

	// 現在の日付を取得
	now := time.Now()
	log.Printf("Fetching Drone news articles for %s...\n", now.Format("2006-01-02"))

	// 記事を取得（当日の記事のみ）
	articles, err := scraper.FetchDroneNews("https://www.drone.jp/", now)
	if err != nil {
		return fmt.Errorf("failed to fetch news: %v", err)
	}

	// 記事数をログ出力
	if len(articles) == 0 {
		log.Println("No articles found for today!")
		return nil
	}
	log.Printf("Found %d articles for today\n", len(articles))

	// Slackに投稿
	if err := slack.PostDroneArticlesToSlack(webhookURL, articles); err != nil {
		return fmt.Errorf("failed to post to Slack: %v", err)
	}

	log.Println("News posted to Slack successfully!")
	return nil
}

func main() {
	if err := RunDroneBot(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
