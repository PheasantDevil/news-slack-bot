package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SlackMessage Slackに送信するメッセージの構造体
type SlackMessage struct {
	Text string `json:"text"`
}

// File: internal/slack/common.go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SendToSlack メッセージをSlackに送信
func SendToSlack(webhookURL string, message string) error {
	// Validate inputs
	if webhookURL == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}

	// メッセージを作成
	msg := SlackMessage{
		Text: message,
	}

	// JSONにエンコード
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// HTTPリクエストを作成
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// リクエストを送信
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスを読み取る
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// ステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
