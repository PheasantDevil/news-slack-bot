package main

import (
	"fmt"
	"log"

	"newsbot/internal/scraper"
	"newsbot/internal/slack"
)

func main() {
	fmt.Println("Fetching news articles...")
	articles, err := scraper.FetchNews("https://www.drone.jp/")
	if err != nil {
		log.Fatal(err)
	}

	for _, article := range articles {
		message := fmt.Sprintf("*%s*\n%s\n%s\n%s", article.Title, article.URL, article.Summary, article.Thumbnail)
		slack.SendToSlack(message)
	}

	log.Println("News posted to Slack successfully!")
}
