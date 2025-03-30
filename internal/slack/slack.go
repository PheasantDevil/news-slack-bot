package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"newsbot/internal/scraper"
)

type SlackMessage struct {
	Text string `json:"text"`
}

func SendToSlack(webhookURL string, message string) error {
	payload, err := json.Marshal(SlackMessage{Text: message})
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error sending Slack message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("error: received non-200 response code: %d", resp.StatusCode)
	}

	log.Println("Slack message sent successfully!")
	return nil
}

// PostArticlesToSlack 記事をSlackに投稿
func PostArticlesToSlack(webhookURL string, articles []scraper.Article) error {
	if len(articles) == 0 {
		return fmt.Errorf("no articles to post")
	}

	// 各記事を個別に投稿
	for _, article := range articles {
		message := fmt.Sprintf("*%s*\n%s\n%s\n%s",
			article.Title,
			article.URL,
			article.Summary,
			article.Thumbnail)

		if err := SendToSlack(webhookURL, message); err != nil {
			log.Printf("Warning: Failed to post article '%s': %v", article.Title, err)
			continue
		}
	}

	// 記事一覧のサマリーを作成
	var summary strings.Builder
	summary.WriteString("*本日の記事一覧*\n")
	for i, article := range articles {
		summary.WriteString(fmt.Sprintf("%d. <%s|%s>\n", i+1, article.URL, article.Title))
	}

	// サマリーを投稿
	if err := SendToSlack(webhookURL, summary.String()); err != nil {
		return fmt.Errorf("failed to post summary: %v", err)
	}

	return nil
}
