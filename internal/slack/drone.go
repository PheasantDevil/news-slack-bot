package slack

import (
	"fmt"
	"log"
	"strings"

	"newsbot/internal/models"
)

// PostDroneArticlesToSlack Droneの記事をSlackに投稿
func PostDroneArticlesToSlack(webhookURL string, articles []models.Article) error {
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
	summary.WriteString("*Drone - 本日の記事一覧*\n")
	for i, article := range articles {
		summary.WriteString(fmt.Sprintf("%d. <%s|%s>\n", i+1, article.URL, article.Title))
	}

	// サマリーを投稿
	if err := SendToSlack(webhookURL, summary.String()); err != nil {
		return fmt.Errorf("failed to post summary: %v", err)
	}

	return nil
}
