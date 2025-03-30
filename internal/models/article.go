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
