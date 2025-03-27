package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

const BaseURL = "https://www.drone.jp"

type Article struct {
	Title     string
	URL       string
	Summary   string
	Thumbnail string
}

func FetchNews(url string) ([]Article, error) {
	var articles []Article
	c := colly.NewCollector()

	c.OnHTML(".news-item", func(e *colly.HTMLElement) {
		article := Article{}

		// 記事タイトル
		article.Title = e.ChildText(".entry-title a")

		// 記事URL（相対パスなら絶対URLに変換）
		relativeURL := e.ChildAttr(".entry-title a", "href")
		if strings.HasPrefix(relativeURL, "/") {
			article.URL = BaseURL + relativeURL
		} else {
			article.URL = relativeURL
		}

		// 記事の概要（要約 or 最初の数行）
		article.Summary = e.ChildText(".entry-summary")

		// サムネイル画像URL
		article.Thumbnail = e.ChildAttr(".post-thumbnail img", "src")

		fmt.Println(article)
		articles = append(articles, article)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}
	return articles, nil
}
