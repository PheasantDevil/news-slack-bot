package main

import (
	"fmt"
	"log"
	"newsbot/internal/scraper"
	"newsbot/internal/slack"
	"os"
	"time"
)

// RunEconomicTimesBot The Economic Timesのニュースを取得してSlackに投稿
func RunEconomicTimesBot() error {
	// Slack Webhook URLを環境変数から取得
	webhookURL := os.Getenv("SLACK_WEBHOOK_ECONOMIC_TIMES_URL")
	if webhookURL == "" {
		return fmt.Errorf("SLACK_WEBHOOK_ECONOMIC_TIMES_URL is not set")
	}

	// 本日の記事を取得
	targetDate := time.Now()
	articles, err := scraper.FetchEconomicTimesNews("https://economictimes.indiatimes.com/news", targetDate)
	if err != nil {
		return fmt.Errorf("failed to fetch news: %v", err)
	}

	if len(articles) == 0 {
		log.Println("No articles found for today")
		return nil
	}

	log.Printf("Found %d articles for today", len(articles))

	// Slackに投稿
	if err := slack.PostEconomicTimesArticlesToSlack(webhookURL, articles); err != nil {
		return fmt.Errorf("failed to post articles to Slack: %v", err)
	}

	log.Println("Successfully posted articles to Slack")
	return nil
}

func main() {
	if err := RunEconomicTimesBot(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
