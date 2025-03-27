package slack

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type SlackMessage struct {
	Text string `json:"text"`
}

func SendToSlack(message string) {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		log.Println("Slack Webhook URL is not set.")
		return
	}

	payload, err := json.Marshal(SlackMessage{Text: message})
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error sending Slack message: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Println("Slack message sent successfully!")
}
