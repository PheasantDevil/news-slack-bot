package models

import "time"

// Article ニュース記事の構造体
type Article struct {
	Title     string
	URL       string
	Summary   string
	Thumbnail string
	PostDate  time.Time
}

// FormatMessage は記事をSlackメッセージ形式にフォーマットします
func (a *Article) FormatMessage() string {
	message := "*" + a.Title + "*\n\n"
	if a.Summary != "" {
		message += a.Summary + "\n\n"
	}
	message += "<" + a.URL + "|記事を読む>"
	return message
}
